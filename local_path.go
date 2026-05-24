package main

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	powerShellEnvPathPattern = regexp.MustCompile(`(?i)\$env:([A-Z_][A-Z0-9_]*)`)
	windowsEnvPathPattern    = regexp.MustCompile(`%([^%]+)%`)
)

func expandLocalPath(raw string) string {
	path := strings.TrimSpace(raw)
	path = strings.Trim(path, "\"'")
	path = powerShellEnvPathPattern.ReplaceAllStringFunc(path, func(match string) string {
		parts := strings.SplitN(match, ":", 2)
		if len(parts) != 2 {
			return match
		}
		if value := os.Getenv(parts[1]); value != "" {
			return value
		}
		return match
	})
	path = windowsEnvPathPattern.ReplaceAllStringFunc(path, func(match string) string {
		key := strings.Trim(match, "%")
		if value := os.Getenv(key); value != "" {
			return value
		}
		return match
	})
	path = os.ExpandEnv(path)
	if strings.HasPrefix(path, "~") && (len(path) == 1 || path[1] == '/' || path[1] == '\\') {
		if home, err := os.UserHomeDir(); err == nil && home != "" {
			suffix := strings.TrimLeft(path[1:], `/\\`)
			if suffix == "" {
				path = home
			} else {
				path = filepath.Join(home, suffix)
			}
		}
	}
	return filepath.Clean(path)
}
