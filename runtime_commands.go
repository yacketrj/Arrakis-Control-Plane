package main

import (
	"fmt"
	"regexp"
	"strings"
)

const (
	runtimeModeAuto          = "auto"
	runtimeModeKubernetes    = "kubernetes"
	runtimeModeDockerCompose = "docker-compose"
	runtimeModeDocker        = "docker"
	runtimeModeHyperV        = "hyperv"
)

var dockerNamePattern = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9_.-]*$`)

func normalizeRuntime(value string) string {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "", "auto":
		return runtimeModeAuto
	case "k8s", "kubernetes":
		return runtimeModeKubernetes
	case "compose", "docker-compose", "docker_compose", "docker compose":
		return runtimeModeDockerCompose
	case "docker":
		return runtimeModeDocker
	case "hyper-v", "hyperv", "hyper_v":
		return runtimeModeHyperV
	default:
		return strings.ToLower(strings.TrimSpace(value))
	}
}

func supportedRuntime(value string) bool {
	switch normalizeRuntime(value) {
	case runtimeModeAuto, runtimeModeKubernetes, runtimeModeDockerCompose, runtimeModeDocker, runtimeModeHyperV:
		return true
	default:
		return false
	}
}

func runtimeUsesKubernetesCommands() bool {
	mode := normalizeRuntime(serverRuntime)
	return mode == runtimeModeKubernetes || mode == runtimeModeHyperV
}

func runtimeUsesDockerCommands() bool {
	mode := normalizeRuntime(serverRuntime)
	return mode == runtimeModeDocker || mode == runtimeModeDockerCompose
}

func dockerContainerNamesCommand() string {
	if normalizeRuntime(serverRuntime) == runtimeModeDockerCompose {
		return `docker ps --filter label=com.docker.compose.project --format '{{.Names}}' 2>&1`
	}
	return `docker ps --format '{{.Names}}' 2>&1`
}

func dockerStatusCommand() string {
	if normalizeRuntime(serverRuntime) == runtimeModeDockerCompose {
		return `docker ps --filter label=com.docker.compose.project --format 'table {{.Names}}\t{{.Image}}\t{{.Status}}\t{{.Ports}}' 2>&1`
	}
	return `docker ps --format 'table {{.Names}}\t{{.Image}}\t{{.Status}}\t{{.Ports}}' 2>&1`
}

func runtimeBGStatusCommand() string {
	if runtimeUsesDockerCommands() {
		return dockerStatusCommand()
	}
	return fmt.Sprintf("sudo kubectl get pods -n %s -o wide 2>&1", globalPodNS)
}

func runtimeBGPodsCommand() string {
	if runtimeUsesDockerCommands() {
		return dockerContainerNamesCommand()
	}
	return fmt.Sprintf("sudo kubectl get pods -n %s --no-headers 2>&1", globalPodNS)
}

func runtimeLogPodsCommand() string {
	if runtimeUsesDockerCommands() {
		return dockerContainerNamesCommand()
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
			{Name: "containers", Description: "Container names, images, status, and published ports.", Command: dockerStatusCommand()},
			{Name: "container_names", Description: "Detected Docker container names only.", Command: dockerContainerNamesCommand()},
			{Name: "compose_projects", Description: "Docker Compose project and service labels when available.", Command: `docker ps --filter label=com.docker.compose.project --format '{{.Names}}\t{{.Label "com.docker.compose.project"}}\t{{.Label "com.docker.compose.service"}}' 2>&1`},
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
		return dockerNamePattern.MatchString(target)
	}
	return isAllowedLogTarget(ns, target)
}
