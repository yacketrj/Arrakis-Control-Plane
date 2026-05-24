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
	if endpoint, err := discoverKubernetesDBEndpoint(client); err == nil {
		return endpoint, nil
	}
	if endpoint, err := discoverDockerDBEndpoint(client); err == nil {
		return endpoint, nil
	}
	return dbEndpointDiscovery{}, fmt.Errorf("database endpoint not found via kubectl or docker")
}

func discoverKubernetesDBEndpoint(client *ssh.Client) (dbEndpointDiscovery, error) {
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
