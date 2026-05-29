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

	piiAssignmentPattern = regexp.MustCompile(`(?i)\b((?:email|e-mail|mail|display[_-]?name|character[_-]?name|player[_-]?name|account[_-]?id|actor[_-]?id|controller[_-]?id|fls[_-]?id|steam[_-]?id|discord[_-]?id|ip[_-]?addr(?:ess)?|remote[_-]?addr)\s*[:=]\s*)([^,\s\"']+)`)
	piiJSONPattern       = regexp.MustCompile(`(?i)(\"(?:email|e-mail|mail|display[_-]?name|character[_-]?name|player[_-]?name|account[_-]?id|actor[_-]?id|controller[_-]?id|fls[_-]?id|steam[_-]?id|discord[_-]?id|ip[_-]?addr(?:ess)?|remote[_-]?addr)\"\s*:\s*\")[^\"]+(\")`)
	emailPattern         = regexp.MustCompile(`\b[A-Za-z0-9._%+\-]+@[A-Za-z0-9.\-]+\.[A-Za-z]{2,}\b`)
	ipv4Pattern          = regexp.MustCompile(`\b(?:(?:25[0-5]|2[0-4][0-9]|1?[0-9]{1,2})\.){3}(?:25[0-5]|2[0-4][0-9]|1?[0-9]{1,2})(?::[0-9]{1,5})?\b`)
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

// RedactPIIText removes user-identifying values from exported diagnostics,
// support bundles, client-facing errors, and any artifact that may be shared
// outside the server operator's trust boundary. It intentionally redacts more
// aggressively than RedactSensitiveText.
func RedactPIIText(value string) string {
	if value == "" {
		return value
	}

	redacted := RedactSensitiveText(value)
	redacted = piiJSONPattern.ReplaceAllString(redacted, `${1}[REDACTED_PII]${2}`)
	redacted = piiAssignmentPattern.ReplaceAllString(redacted, `${1}[REDACTED_PII]`)
	redacted = emailPattern.ReplaceAllString(redacted, `[REDACTED_EMAIL]`)
	redacted = ipv4Pattern.ReplaceAllString(redacted, `[REDACTED_IP]`)
	return redacted
}

func containsRedaction(value string) bool {
	return strings.Contains(value, "[REDACTED]") || strings.Contains(value, "[REDACTED_PII]")
}
