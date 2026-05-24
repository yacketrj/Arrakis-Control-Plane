package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

var bgCmdAllowlist = map[string]bool{
	"start": true, "stop": true, "restart": true,
	"update": true, "backup": true, "restore": true,
}

type battlegroupHealthSection struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Command     string `json:"command"`
	Output      string `json:"output"`
	Error       string `json:"error,omitempty"`
}

type battlegroupHealthResponse struct {
	Namespace string                     `json:"namespace"`
	CheckedAt string                     `json:"checked_at"`
	Sections  []battlegroupHealthSection `json:"sections"`
}

type battlegroupHealthSpec struct {
	Name        string
	Description string
	Command     string
}

func battlegroupHealthSpecs(namespace string) []battlegroupHealthSpec {
	return []battlegroupHealthSpec{
		{
			Name:        "pods",
			Description: "All namespace pods with node placement, pod IPs, readiness, restarts, and age.",
			Command:     fmt.Sprintf("sudo kubectl get pods -n %s -o wide 2>&1", namespace),
		},
		{
			Name:        "services",
			Description: "Namespace services and NodePort/ClusterIP exposure for BGD, RMQ, DB, and related components.",
			Command:     fmt.Sprintf("sudo kubectl get svc -n %s -o wide 2>&1", namespace),
		},
		{
			Name:        "statefulsets",
			Description: "Stateful service readiness, especially PostgreSQL and RabbitMQ stateful workloads.",
			Command:     fmt.Sprintf("sudo kubectl get statefulsets -n %s -o wide 2>&1", namespace),
		},
		{
			Name:        "deployments",
			Description: "Deployment rollout readiness for BGD, Text Router, gateway, and operator-managed services.",
			Command:     fmt.Sprintf("sudo kubectl get deployments -n %s -o wide 2>&1", namespace),
		},
		{
			Name:        "persistent_volumes",
			Description: "Persistent volume claims for database and stateful service storage health checks.",
			Command:     fmt.Sprintf("sudo kubectl get pvc -n %s -o wide 2>&1", namespace),
		},
		{
			Name:        "recent_events",
			Description: "Recent namespace events sorted by timestamp for crash, scheduling, image, storage, or readiness issues.",
			Command:     fmt.Sprintf("sudo kubectl get events -n %s --sort-by=.lastTimestamp 2>&1", namespace),
		},
		{
			Name:        "nodes",
			Description: "Cluster node readiness, version, and age. This is cluster-scope and intentionally read-only.",
			Command:     "sudo kubectl get nodes -o wide 2>&1",
		},
		{
			Name:        "pod_metrics",
			Description: "Pod CPU and memory usage when metrics-server is available. Errors here are informational.",
			Command:     fmt.Sprintf("sudo kubectl top pods -n %s 2>&1", namespace),
		},
	}
}

func handleBGStatus(w http.ResponseWriter, r *http.Request) {
	out, err := sshExec(fmt.Sprintf("sudo kubectl get pods -n %s -o wide 2>&1", globalPodNS))
	if err != nil {
		jsonErr(w, fmt.Errorf("kubectl: %w — output: %s", err, out), 500)
		return
	}
	jsonOK(w, map[string]string{"output": out})
}

func handleBGHealth(w http.ResponseWriter, r *http.Request) {
	if !isValidK8sName(globalPodNS) {
		jsonErr(w, fmt.Errorf("invalid namespace %q", globalPodNS), http.StatusBadRequest)
		return
	}

	response := battlegroupHealthResponse{
		Namespace: globalPodNS,
		CheckedAt: time.Now().UTC().Format(time.RFC3339),
	}
	for _, spec := range battlegroupHealthSpecs(globalPodNS) {
		section := battlegroupHealthSection{
			Name:        spec.Name,
			Description: spec.Description,
			Command:     spec.Command,
		}
		out, err := sshExec(spec.Command)
		section.Output = out
		if err != nil {
			section.Error = err.Error()
		}
		response.Sections = append(response.Sections, section)
	}
	jsonOK(w, response)
}

func handleBGExec(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Cmd string `json:"cmd"`
	}
	if err := decode(r, &req); err != nil {
		jsonErr(w, err, 400)
		return
	}
	if !bgCmdAllowlist[req.Cmd] {
		jsonErr(w, fmt.Errorf("unknown command %q", req.Cmd), 400)
		return
	}
	out, err := sshExec(fmt.Sprintf("sudo ~/.dune/download/scripts/battlegroup.sh %s 2>&1", req.Cmd))
	if err != nil {
		jsonErr(w, fmt.Errorf("exec: %w — output: %s", err, out), 500)
		return
	}
	jsonOK(w, map[string]string{"output": out})
}

func handleBGPods(w http.ResponseWriter, r *http.Request) {
	out, err := sshExec(fmt.Sprintf("sudo kubectl get pods -n %s --no-headers 2>&1", globalPodNS))
	if err != nil {
		jsonErr(w, fmt.Errorf("kubectl: %w", err), 500)
		return
	}
	lines := strings.Split(strings.TrimSpace(out), "\n")
	jsonOK(w, map[string]any{"pods": lines, "namespace": globalPodNS})
}
