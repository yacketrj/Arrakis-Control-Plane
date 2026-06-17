package main

import (
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/crypto/ssh"
)

func dockerDuneStackAvailable(client *ssh.Client) bool {
	out, err := sshCombined(client, `docker ps --format '{{.Names}}' 2>/dev/null`)
	if err != nil {
		return false
	}
	value := strings.ToLower(out)
	markers := []string{"dune-postgres", "dune-server-", "dune-director", "dune-orchestrator", "dune-server-gateway"}
	for _, marker := range markers {
		if strings.Contains(value, marker) {
			return true
		}
	}
	return false
}

func discoverDockerDBEndpoint(client *ssh.Client) (dbEndpointDiscovery, error) {
	out, err := sshCombined(client, `docker ps --format '{{.ID}}|{{.Names}}|{{.Image}}' 2>/dev/null`)
	if err != nil {
		return dbEndpointDiscovery{}, fmt.Errorf("docker ps: %w", err)
	}
	for _, line := range strings.Split(out, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.Split(line, "|")
		if len(parts) < 3 {
			continue
		}
		id := strings.TrimSpace(parts[0])
		name := strings.TrimSpace(parts[1])
		image := strings.TrimSpace(parts[2])
		if !looksLikeDBContainer(name, image) {
			continue
		}
		endpoint, err := inspectDockerDBEndpoint(client, id, name)
		if err == nil {
			return endpoint, nil
		}
	}
	return dbEndpointDiscovery{}, fmt.Errorf("database container not found")
}

func looksLikeDBContainer(name, image string) bool {
	value := strings.ToLower(name + " " + image)
	markers := []string{"postgres", "postgresql", "db-dbdepl", "database", "_db", "-db"}
	for _, marker := range markers {
		if strings.Contains(value, marker) {
			return true
		}
	}
	return false
}

func inspectDockerDBEndpoint(client *ssh.Client, containerID, fallbackName string) (dbEndpointDiscovery, error) {
	publishedPort := dockerDBPort(client, containerID)
	if publishedPort > 0 {
		return dbEndpointDiscovery{
			Runtime:   runtimeDocker,
			Namespace: string(runtimeDocker),
			Name:      fallbackName,
			Host:      "127.0.0.1",
			Port:      publishedPort,
		}, nil
	}

	host, _ := sshCombined(client, `docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' `+shellQuote(containerID)+` 2>/dev/null`)
	host = strings.TrimSpace(host)
	if host == "" {
		host = "127.0.0.1"
	}
	return dbEndpointDiscovery{
		Runtime:   runtimeDocker,
		Namespace: string(runtimeDocker),
		Name:      fallbackName,
		Host:      host,
		Port:      5432,
	}, nil
}

func dockerDBPort(client *ssh.Client, containerID string) int {
	for _, containerPort := range []string{fmt.Sprintf("%d/tcp", dbPort), "5432/tcp"} {
		out, _ := sshCombined(client, `docker port `+shellQuote(containerID)+` `+shellQuote(containerPort)+` 2>/dev/null | head -1`)
		out = strings.TrimSpace(out)
		if out == "" {
			continue
		}
		if idx := strings.LastIndex(out, ":"); idx != -1 {
			if port, err := strconv.Atoi(strings.TrimSpace(out[idx+1:])); err == nil && port > 0 {
				return port
			}
		}
	}
	return 0
}
