package main

import (
	"fmt"
	"net"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var (
	pgIdentifierPattern   = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]{0,62}$`)
	pgDatabaseNamePattern = regexp.MustCompile(`^[A-Za-z0-9_.-]{1,63}$`)
	pgUserNamePattern     = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_.@-]{0,62}$`)
	sshUserPattern        = regexp.MustCompile(`^[A-Za-z0-9_.-]{1,64}$`)
)

func missingConfigValue(value string) bool {
	trimmed := strings.TrimSpace(value)
	return trimmed == "" || strings.EqualFold(trimmed, "null") || strings.EqualFold(trimmed, "<nil>")
}

func containsUnsafeControl(value string) bool {
	return strings.ContainsAny(value, "\x00\r\n")
}

func validTunnelModeValue(value string) bool {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "", "auto", "on", "managed", "existing", "external", "local", "off", "direct", "disabled", "false", "0":
		return true
	default:
		return false
	}
}

func validateHostPort(name, value string) error {
	if missingConfigValue(value) {
		return fmt.Errorf("%s is required", name)
	}
	if containsUnsafeControl(value) || strings.ContainsAny(value, " \t'\"`;|&<>") {
		return fmt.Errorf("%s contains unsupported characters", name)
	}
	host, portText, err := net.SplitHostPort(value)
	if err != nil {
		return fmt.Errorf("%s must be in host:port format", name)
	}
	if missingConfigValue(host) {
		return fmt.Errorf("%s host is required", name)
	}
	port, err := strconv.Atoi(portText)
	if err != nil || port <= 0 || port > 65535 {
		return fmt.Errorf("%s port must be between 1 and 65535", name)
	}
	return nil
}

func validateSecretValue(name, value string) error {
	if missingConfigValue(value) {
		return fmt.Errorf("%s is required", name)
	}
	if containsUnsafeControl(value) {
		return fmt.Errorf("%s contains unsupported control characters", name)
	}
	return nil
}

func requiredConfigErrors() []string {
	var errs []string
	checkRequired := func(name, value string) {
		if missingConfigValue(value) {
			errs = append(errs, fmt.Sprintf("%s is required", name))
		}
	}
	checkPattern := func(name, value string, pattern *regexp.Regexp) {
		if missingConfigValue(value) {
			errs = append(errs, fmt.Sprintf("%s is required", name))
			return
		}
		if containsUnsafeControl(value) || !pattern.MatchString(value) {
			errs = append(errs, fmt.Sprintf("%s contains unsupported characters", name))
		}
	}

	if err := validateHostPort("SSH_HOST", sshHost); err != nil {
		errs = append(errs, err.Error())
	}
	checkPattern("SSH_USER", sshUser, sshUserPattern)
	if _, err := resolveValidatedKeyPath(); err != nil {
		errs = append(errs, err.Error())
	}
	checkPattern("DB_USER", dbUser, pgUserNamePattern)
	checkPattern("DB_NAME", dbName, pgDatabaseNamePattern)
	checkPattern("DB_SCHEMA", dbSchema, pgIdentifierPattern)
	if err := validateSecretValue("DB_PASS", dbPass); err != nil {
		errs = append(errs, err.Error())
	}
	if err := validateSecretValue("ADMIN_TOKEN", effectiveAdminToken()); err != nil {
		errs = append(errs, err.Error())
	}
	checkRequired("LISTEN_ADDR", listenAddr)

	if !validTunnelModeValue(sshTunnelMode) {
		errs = append(errs, "SSH_TUNNEL_MODE must be one of: auto, existing, off")
	}
	if sshTunnelHost != "" {
		if containsUnsafeControl(sshTunnelHost) || strings.ContainsAny(sshTunnelHost, " \t'\"`;|&<>") {
			errs = append(errs, "SSH_TUNNEL_LOCAL_HOST contains unsupported characters")
		}
	}
	if dbPort <= 0 || dbPort > 65535 {
		errs = append(errs, "DB_PORT must be between 1 and 65535")
	}
	if dbTunnelLocalPort < 0 || dbTunnelLocalPort > 65535 {
		errs = append(errs, "DB_TUNNEL_LOCAL_PORT must be between 0 and 65535")
	}
	return errs
}

func effectiveAdminToken() string {
	if !missingConfigValue(adminToken) {
		return adminToken
	}
	return os.Getenv("ADMIN_TOKEN")
}
