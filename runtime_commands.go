package main

import (
	"fmt"
	"regexp"
	"strings"
)

const (
	runtimeModeAuto       = "auto"
	runtimeModeKubernetes = "kubernetes"
	runtimeModeDocker     = "docker"
)

var dockerNamePattern = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9_.-]*$`)
var dockerIDPattern = regexp.MustCompile(`^[a-fA-F0-9]{12,64}$`)

func normalizeRuntime(value string) string {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "", "auto":
		return runtimeModeAuto
	case "k8", "k8s", "kubectl", "kubernetes", "hyper-v", "hyperv", "hyper_v":
		return runtimeModeKubernetes
	case "docker", "compose", "docker-compose", "docker_compose", "docker compose":
		return runtimeModeDocker
	default:
		return strings.ToLower(strings.TrimSpace(value))
	}
}

func supportedRuntime(value string) bool {
	switch normalizeRuntime(value) {
	case runtimeModeAuto, runtimeModeKubernetes, runtimeModeDocker:
		return true
	default:
		return false
	}
}

func runtimeUsesKubernetesCommands() bool {
	return normalizeRuntime(serverRuntime) == runtimeModeKubernetes
}

func runtimeUsesDockerCommands() bool {
	return normalizeRuntime(serverRuntime) == runtimeModeDocker
}

func dockerContainerIDsCommand() string {
	return `docker ps -a -q 2>&1`
}

func dockerLogTargetsCommand() string {
	return `docker ps -a --format '{{.Names}}|{{.ID}}' 2>&1`
}

func dockerStatusCommand() string {
	return `docker ps -a --format 'table {{.Names}}\t{{.ID}}\t{{.Image}}\t{{.Status}}\t{{.Ports}}' 2>&1`
}

func dockerStatsCommand() string {
	return `docker stats --no-stream --format 'table {{.Name}}\t{{.CPUPerc}}\t{{.MemUsage}}\t{{.MemPerc}}\t{{.NetIO}}\t{{.BlockIO}}\t{{.PIDs}}' 2>&1`
}

func runtimeBGStatusCommand() string {
	if runtimeUsesDockerCommands() {
		return dockerStatusCommand()
	}
	return fmt.Sprintf("sudo kubectl get pods -n %s -o wide 2>&1", globalPodNS)
}

func runtimeBGPodsCommand() string {
	if runtimeUsesDockerCommands() {
		return dockerContainerIDsCommand()
	}
	return fmt.Sprintf("sudo kubectl get pods -n %s --no-headers 2>&1", globalPodNS)
}

func runtimeLogPodsCommand() string {
	if runtimeUsesDockerCommands() {
		return dockerLogTargetsCommand()
	}
	return fmt.Sprintf("sudo kubectl get pods -n %s --no-headers -o custom-columns=NAME:.metadata.name 2>&1", globalPodNS)
}

func runtimeLogStreamCommand(ns, target string) string {
	if runtimeUsesDockerCommands() {
		return "docker logs -f " + shellQuote(target) + " 2>&1"
	}
	return strings.Join([]string{"sudo", "kubectl", "logs", "-f", "-n", ns, target}, " ") + " 2>&1"
}

func runtimeHealthSpecs(namespace string) []battlegroupHealthSpec {
	if runtimeUsesDockerCommands() {
		return []battlegroupHealthSpec{
			{Name: "container_overview", Description: "Combined container view with name, ID, image, status, and ports.", Command: dockerStatusCommand()},
			{Name: "container_metrics", Description: "One-shot container resource snapshot with CPU, memory, network I/O, block I/O, and PID count.", Command: dockerStatsCommand()},
		}
	}
	return []battlegroupHealthSpec{
		{Name: "pods", Description: "All namespace pods with node placement, pod IPs, readiness, restarts, and age.", Command: fmt.Sprintf("sudo kubectl get pods -n %s -o wide 2>&1", namespace)},
		{Name: "services", Description: "Namespace services and NodePort/ClusterIP exposure.", Command: fmt.Sprintf("sudo kubectl get svc -n %s -o wide 2>&1", namespace)},
		{Name: "statefulsets", Description: "Stateful service readiness.", Command: fmt.Sprintf("sudo kubectl get statefulsets -n %s -o wide 2>&1", namespace)},
		{Name: "deployments", Description: "Deployment rollout readiness.", Command: fmt.Sprintf("sudo kubectl get deployments -n %s -o wide 2>&1", namespace)},
		{Name: "persistent_volumes", Description: "Persistent volume claims.", Command: fmt.Sprintf("sudo kubectl get pvc -n %s -o wide 2>&1", namespace)},
		{Name: "recent_events", Description: "Recent namespace events sorted by timestamp.", Command: fmt.Sprintf("sudo kubectl get events -n %s --sort-by=.lastTimestamp 2>&1", namespace)},
		{Name: "nodes", Description: "Cluster node readiness, version, and age.", Command: "sudo kubectl get nodes -o wide 2>&1"},
		{Name: "pod_metrics", Description: "Pod CPU and memory usage when metrics-server is available.", Command: fmt.Sprintf("sudo kubectl top pods -n %s 2>&1", namespace)},
	}
}

func isValidRuntimeLogTarget(ns, target string) bool {
	if runtimeUsesDockerCommands() {
		return dockerIDPattern.MatchString(target) || dockerNamePattern.MatchString(target)
	}
	return isAllowedLogTarget(ns, target)
}
