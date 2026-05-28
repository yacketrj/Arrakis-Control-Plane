package main

import (
	"regexp"
	"strings"
)

var (
	sensitiveAssignmentPattern = regexp.MustCompile(`(?i)\b((?:admin[_-]?token|funcom[_-]?token|db[_-]?pass(?:word)?|database[_-]?url|postgres(?:ql)?[_-]?url|ssh[_-]?key|client[_-]?secret|secret|password|token)\s*[:=]\s*)([^,\s\"']+)`)
	sensitiveJSONPattern       = regexp.MustCompile(`(?i)(\"(?:admin[_-]?token|funcom[_-]?token|db[_-]?pass(?:word)?|database[_-]?url|postgres(?:ql)?[_-]?url|ssh[_-]?key|client[_-]?secret|secret|password|token)\"\s*:\s*\")[^\"]+(\")`)
	bearerTokenPattern         = regexp.MustCompile(`(?i)\b(Bearer\s+)[A-Za-z0-9._~+/=-]{12,}`)
	postgresURLPasswordPattern = regexp.MustCompile(`(?i)\b((?:postgres|postgresql)://[^:@\s]+:)([^@\s]+)(@[^\s]+)`)
)

// RedactSensitiveText removes high-risk credentials from logs, audit records,
// diagnostic errors, and exported text. Keep this helper conservative and
// deterministic so tests can assert exact redaction behavior.
func RedactSensitiveText(value string) string {
	if value == "" {
		return value
	}

	redacted := value
	redacted = postgresURLPasswordPattern.ReplaceAllString(redacted, `${1}[REDACTED]${3}`)
	redacted = sensitiveJSONPattern.ReplaceAllString(redacted, `${1}[REDACTED]${2}`)
	redacted = sensitiveAssignmentPattern.ReplaceAllString(redacted, `${1}[REDACTED]`)
	redacted = bearerTokenPattern.ReplaceAllString(redacted, `${1}[REDACTED]`)
	return redacted
}

func containsRedaction(value string) bool {
	return strings.Contains(value, "[REDACTED]")
}
