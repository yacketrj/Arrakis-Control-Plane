package main

import (
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/crypto/ssh"
)

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
	runtime := runtimeDocker
	name := fallbackName
	project, _ := sshCombined(client, `docker inspect -f '{{index .Config.Labels "com.docker.compose.project"}}' `+shellQuote(containerID)+` 2>/dev/null`)
	service, _ := sshCombined(client, `docker inspect -f '{{index .Config.Labels "com.docker.compose.service"}}' `+shellQuote(containerID)+` 2>/dev/null`)
	project = strings.TrimSpace(project)
	service = strings.TrimSpace(service)
	if project != "" && project != "<no value>" {
		runtime = runtimeDockerCompose
		if service != "" && service != "<no value>" {
			name = project + "/" + service
		} else {
			name = project + "/" + fallbackName
		}
	}

	host, _ := sshCombined(client, `docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' `+shellQuote(containerID)+` 2>/dev/null`)
	host = strings.TrimSpace(host)
	port := dockerDBPort(client, containerID)
	if host == "" {
		host = "127.0.0.1"
	}
	if port == 0 {
		port = 5432
	}
	return dbEndpointDiscovery{
		Runtime:   runtime,
		Namespace: string(runtime),
		Name:      name,
		Host:      host,
		Port:      port,
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
	if dbPort > 0 {
		return dbPort
	}
	return 5432
}
