package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
)

type mutationAuditMetadata struct {
	Reason string
	Target map[string]string
}

var mutationAuditTargetKeys = []string{
	"player_id",
	"account_id",
	"actor_id",
	"controller_id",
	"fls_id",
	"item_id",
	"item_template",
	"item_template_id",
	"template_id",
	"quantity",
	"amount",
	"quality",
	"faction_id",
	"storage_id",
	"container_id",
	"vehicle_id",
	"guild_id",
	"rank",
	"pod",
	"service",
	"cmd",
	"command",
	"command_path",
}

func extractMutationAuditMetadata(r *http.Request) mutationAuditMetadata {
	metadata := mutationAuditMetadata{Target: map[string]string{}}
	metadata.Reason = sanitizedAuditString(r.Header.Get("X-Admin-Reason"), 256)
	if r.Body == nil || r.ContentLength > maxAuditInspectableBodyBytes {
		if len(metadata.Target) == 0 {
			metadata.Target = nil
		}
		return metadata
	}

	body, err := io.ReadAll(io.LimitReader(r.Body, maxAuditInspectableBodyBytes+1))
	if err != nil {
		r.Body = io.NopCloser(bytes.NewReader(body))
		return metadata
	}
	r.Body = io.NopCloser(bytes.NewReader(body))
	if int64(len(body)) > maxAuditInspectableBodyBytes || len(bytes.TrimSpace(body)) == 0 {
		return metadata
	}

	var payload map[string]any
	if err := json.Unmarshal(body, &payload); err != nil {
		return metadata
	}

	if metadata.Reason == "" {
		metadata.Reason = sanitizedAuditString(payloadString(payload, "reason"), 256)
	}
	for _, key := range mutationAuditTargetKeys {
		if value, ok := auditScalar(payload[key]); ok {
			metadata.Target[key] = value
		}
	}
	if len(metadata.Target) == 0 {
		metadata.Target = nil
	}
	return metadata
}

func payloadString(payload map[string]any, key string) string {
	if raw, ok := payload[key]; ok {
		if value, ok := raw.(string); ok {
			return value
		}
	}
	return ""
}

func auditScalar(value any) (string, bool) {
	switch v := value.(type) {
	case nil:
		return "", false
	case string:
		trimmed := sanitizedAuditString(RedactSensitiveText(v), 128)
		if trimmed == "" {
			return "", false
		}
		return trimmed, true
	case float64:
		return fmt.Sprintf("%.0f", v), true
	case bool:
		return fmt.Sprintf("%t", v), true
	default:
		return "", false
	}
}

func sanitizedAuditString(value string, maxLen int) string {
	value = RedactSensitiveText(strings.TrimSpace(value))
	if value == "" {
		return ""
	}
	value = strings.ReplaceAll(value, "\n", " ")
	value = strings.ReplaceAll(value, "\r", " ")
	value = strings.Join(strings.Fields(value), " ")
	if len(value) > maxLen {
		return value[:maxLen]
	}
	return value
}

func auditRemoteAddr(r *http.Request) string {
	if r == nil {
		return ""
	}
	remote := strings.TrimSpace(r.RemoteAddr)
	if remote == "" {
		return ""
	}
	host, _, err := net.SplitHostPort(remote)
	if err == nil {
		return sanitizedAuditString(host, 128)
	}
	return sanitizedAuditString(remote, 128)
}

func auditAdminTokenHash(r *http.Request) string {
	if r == nil {
		return ""
	}
	token := bearerToken(r.Header.Get("Authorization"))
	if token == "" {
		token = r.Header.Get("X-Admin-Token")
	}
	if token == "" {
		return ""
	}
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])[:16]
}
