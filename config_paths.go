package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var percentEnvPattern = regexp.MustCompile(`%([A-Za-z_][A-Za-z0-9_]*)%`)

func expandSupportedConfigPath(raw string) (string, error) {
	path := strings.TrimSpace(raw)
	path = strings.Trim(path, `"'`)
	if path == "" {
		return "", nil
	}
	if strings.Contains(strings.ToLower(path), "$env:") {
		return "", fmt.Errorf("PowerShell-style paths are not supported in .env; use %%USERPROFILE%%\\.ssh\\id_rsa instead of %s", raw)
	}

	var expandErr error
	path = percentEnvPattern.ReplaceAllStringFunc(path, func(match string) string {
		name := strings.Trim(match, "%")
		value := os.Getenv(name)
		if value == "" {
			expandErr = fmt.Errorf("environment variable %s referenced by path %s is not set", name, raw)
			return match
		}
		return value
	})
	if expandErr != nil {
		return "", expandErr
	}

	if path == "~" || strings.HasPrefix(path, `~\`) || strings.HasPrefix(path, "~/") {
		home, err := os.UserHomeDir()
		if err != nil || home == "" {
			return "", fmt.Errorf("cannot expand ~ in path %s: user home directory not available", raw)
		}
		if path == "~" {
			path = home
		} else {
			path = filepath.Join(home, path[2:])
		}
	}

	return filepath.Clean(path), nil
}

func validateReadableFilePath(label, raw string) (string, error) {
	path, err := expandSupportedConfigPath(raw)
	if err != nil {
		return "", fmt.Errorf("%s: %w", label, err)
	}
	if missingConfigValue(path) {
		return "", fmt.Errorf("%s is required", label)
	}
	info, err := os.Stat(path)
	if err != nil {
		return "", fmt.Errorf("%s does not exist or is not readable: %s", label, path)
	}
	if info.IsDir() {
		return "", fmt.Errorf("%s points to a directory, not a file: %s", label, path)
	}
	return path, nil
}

func resolveValidatedKeyPath() (string, error) {
	if strings.TrimSpace(sshKeyPath) != "" {
		return validateReadableFilePath("SSH_KEY", sshKeyPath)
	}
	return validateReadableFilePath("SSH_KEY", resolveKeyPath())
}
