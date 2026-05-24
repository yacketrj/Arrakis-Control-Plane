package main

import (
	"fmt"
	"strings"

	"golang.org/x/crypto/ssh"
)

type runtimeKind string

const (
	runtimeKubernetes    runtimeKind = "kubernetes"
	runtimeDockerCompose runtimeKind = "docker-compose"
	runtimeDocker        runtimeKind = "docker"
	runtimeUnknown       runtimeKind = "unknown"
)

type dbEndpointDiscovery struct {
	Runtime   runtimeKind
	Namespace string
	Name      string
	Host      string
	Port      int
}

func discoverDatabaseEndpoint(client *ssh.Client) (dbEndpointDiscovery, error) {
	requested := normalizeRuntime(serverRuntime)
	switch requested {
	case runtimeModeKubernetes, runtimeModeHyperV:
		endpoint, err := discoverKubernetesDBEndpoint(client)
		if err != nil {
			return dbEndpointDiscovery{}, err
		}
		return endpoint, nil
	case runtimeModeDocker, runtimeModeDockerCompose:
		endpoint, err := discoverDockerDBEndpoint(client)
		if err != nil {
			return dbEndpointDiscovery{}, err
		}
		return endpoint, nil
	case runtimeModeAuto:
		if endpoint, err := discoverKubernetesDBEndpoint(client); err == nil {
			return endpoint, nil
		}
		if endpoint, err := discoverDockerDBEndpoint(client); err == nil {
			return endpoint, nil
		}
		return dbEndpointDiscovery{}, fmt.Errorf("database endpoint not found through Kubernetes or Docker")
	default:
		return dbEndpointDiscovery{}, fmt.Errorf("unsupported SERVER_RUNTIME %q", serverRuntime)
	}
}

func applyDetectedRuntime(runtime runtimeKind) {
	if normalizeRuntime(serverRuntime) == runtimeModeAuto {
		switch runtime {
		case runtimeKubernetes:
			serverRuntime = runtimeModeKubernetes
		case runtimeDockerCompose:
			serverRuntime = runtimeModeDockerCompose
		case runtimeDocker:
			serverRuntime = runtimeModeDocker
		}
	}
}

func discoverKubernetesDBEndpoint(client *ssh.Client) (dbEndpointDiscovery, error) {
	if out, err := sshCombined(client, `command -v kubectl >/dev/null 2>&1 && echo ok`); err != nil || strings.TrimSpace(out) != "ok" {
		return dbEndpointDiscovery{}, fmt.Errorf("kubectl not available")
	}
	out, err := sshCombined(client, `sudo kubectl get pods -A -o jsonpath='{range .items[*]}{.metadata.namespace}{" "}{.metadata.name}{" "}{.status.podIP}{"\n"}{end}' 2>/dev/null | grep db-dbdepl-sts | head -1`)
	if err != nil {
		return dbEndpointDiscovery{}, fmt.Errorf("kubectl: %w", err)
	}
	parts := strings.Fields(strings.TrimSpace(out))
	if len(parts) < 3 {
		return dbEndpointDiscovery{}, fmt.Errorf("database pod not found in kubernetes")
	}
	return dbEndpointDiscovery{
		Runtime:   runtimeKubernetes,
		Namespace: parts[0],
		Name:      parts[1],
		Host:      parts[2],
		Port:      dbPort,
	}, nil
}

func sshCombined(client *ssh.Client, cmd string) (string, error) {
	sess, err := client.NewSession()
	if err != nil {
		return "", fmt.Errorf("SSH session: %w", err)
	}
	defer sess.Close()
	out, err := sess.CombinedOutput(cmd)
	return strings.TrimSpace(string(out)), err
}
