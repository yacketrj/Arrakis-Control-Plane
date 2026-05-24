package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func missingConfigValue(value string) bool {
	trimmed := strings.TrimSpace(value)
	return trimmed == "" || strings.EqualFold(trimmed, "null") || strings.EqualFold(trimmed, "<nil>")
}

func requiredConfigErrors() []string {
	var errs []string
	check := func(name, value string) {
		if missingConfigValue(value) {
			errs = append(errs, fmt.Sprintf("%s is required", name))
		}
	}

	check("SSH_HOST", sshHost)
	check("SSH_USER", sshUser)
	check("SSH_KEY", resolveKeyPath())
	check("DB_NAME", dbName)
	check("DB_USER", dbUser)
	check("DB_PASS", dbPass)
	check("DB_SCHEMA", dbSchema)
	check("ADMIN_TOKEN", effectiveAdminToken())
	check("LISTEN_ADDR", listenAddr)

	if dbPort <= 0 || dbPort > 65535 {
		errs = append(errs, "DB_PORT must be between 1 and 65535")
	}
	if dbTunnelLocalPort < 0 || dbTunnelLocalPort > 65535 {
		errs = append(errs, "DB_TUNNEL_LOCAL_PORT must be between 0 and 65535")
	}
	if _, err := strconv.Atoi(fmt.Sprintf("%d", dbPort)); err != nil {
		errs = append(errs, "DB_PORT must be a valid integer")
	}
	if keyPath := resolveKeyPath(); !missingConfigValue(keyPath) {
		if _, err := os.Stat(keyPath); err != nil {
			errs = append(errs, fmt.Sprintf("SSH_KEY does not exist or is not readable: %s", keyPath))
		}
	}
	return errs
}

func effectiveAdminToken() string {
	if !missingConfigValue(adminToken) {
		return adminToken
	}
	return os.Getenv("ADMIN_TOKEN")
}
