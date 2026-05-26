package main

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var windowsEnvPathPattern = regexp.MustCompile(`%([^%]+)%`)

func expandLocalPath(raw string) string {
	path := strings.TrimSpace(raw)
	path = strings.Trim(path, "\"'")

	// Config files intentionally support Windows percent-variable expansion and
	// home-directory expansion only. PowerShell expressions are not expanded here;
	// validation rejects them with a clear operator-facing error.
	path = windowsEnvPathPattern.ReplaceAllStringFunc(path, func(match string) string {
		key := strings.Trim(match, "%")
		if value := os.Getenv(key); value != "" {
			return value
		}
		return match
	})
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
