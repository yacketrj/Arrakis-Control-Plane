package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"strings"
	"sync"
	"time"
)

const defaultAdminAuditPath = "admin-audit.jsonl"
const maxAuditInspectableBodyBytes int64 = 1 << 20

type adminAuditEvent struct {
	Timestamp        string            `json:"timestamp"`
	Method           string            `json:"method"`
	Path             string            `json:"path"`
	Action           string            `json:"action"`
	Risk             string            `json:"risk"`
	Reason           string            `json:"reason,omitempty"`
	Target           map[string]string `json:"target,omitempty"`
	Status           int               `json:"status"`
	DurationMS       int64             `json:"duration_ms"`
	Result           string            `json:"result"`
	RequiresReason   bool              `json:"requires_reason,omitempty"`
	RequiresPreview  bool              `json:"requires_preview,omitempty"`
	Destructive      bool              `json:"destructive,omitempty"`
	RollbackHint     string            `json:"rollback_hint,omitempty"`
	OperatorWarnings []string          `json:"operator_warnings,omitempty"`
	RecommendedPath  string            `json:"recommended_path,omitempty"`
}

type statusCaptureWriter struct {
	http.ResponseWriter
	status int
}

type mutationAuditMetadata struct {
	Reason string
	Target map[string]string
}

var adminAuditMu sync.Mutex

func (w *statusCaptureWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func (w *statusCaptureWriter) Write(body []byte) (int, error) {
	if w.status == 0 {
		w.status = http.StatusOK
	}
	return w.ResponseWriter.Write(body)
}

func auditMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !isAuditableRequest(r) {
			next.ServeHTTP(w, r)
			return
		}

		metadata := extractMutationAuditMetadata(r)
		safety := mutationSafetyForPath(r.Method, r.URL.Path)
		started := time.Now()
		capture := &statusCaptureWriter{ResponseWriter: w}
		next.ServeHTTP(capture, r)

		status := capture.status
		if status == 0 {
			status = http.StatusOK
		}
		_ = appendAdminAuditEvent(adminAuditEvent{
			Timestamp:        started.UTC().Format(time.RFC3339Nano),
			Method:           r.Method,
			Path:             r.URL.Path,
			Action:           safety.Action,
			Risk:             safety.Risk,
			Reason:           metadata.Reason,
			Target:           metadata.Target,
			Status:           status,
			DurationMS:       time.Since(started).Milliseconds(),
			Result:           auditResultForStatus(status),
			RequiresReason:   safety.RequiresReason,
			RequiresPreview:  safety.RequiresPreview,
			Destructive:      safety.Destructive,
			RollbackHint:     safety.RollbackHint,
			OperatorWarnings: safety.OperatorWarnings,
			RecommendedPath:  safety.RecommendedPath,
		})
	})
}

func isAuditableRequest(r *http.Request) bool {
	if isPublicPath(r.URL.Path) {
		return false
	}
	switch r.Method {
	case http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete:
		return true
	default:
		return false
	}
}

func auditActionName(method, path string) string {
	name := strings.TrimPrefix(path, "/api/v1/")
	name = strings.Trim(name, "/")
	name = strings.ReplaceAll(name, "/", ".")
	if name == "" {
		name = "root"
	}
	return strings.ToLower(method) + ":" + name
}

func mutationRiskForRequest(method, path string) string {
	if method == http.MethodDelete || strings.Contains(path, "/wipe") || strings.Contains(path, "/delete") || strings.Contains(path, "/blueprints/import") {
		return "destructive"
	}
	if strings.Contains(path, "/give-item") || strings.Contains(path, "/grant-live") || strings.Contains(path, "/teleport") || strings.Contains(path, "/journey/") || strings.Contains(path, "/set-faction") {
		return "high"
	}
	if method == http.MethodPost || method == http.MethodPut || method == http.MethodPatch {
		return "medium"
	}
	return "low"
}

func auditResultForStatus(status int) string {
	if status >= 200 && status < 400 {
		return "success"
	}
	return "failure"
}

func adminAuditPath() string {
	if path := strings.TrimSpace(os.Getenv("ADMIN_AUDIT_LOG")); path != "" {
		return path
	}
	return defaultAdminAuditPath
}

func appendAdminAuditEvent(event adminAuditEvent) error {
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	adminAuditMu.Lock()
	defer adminAuditMu.Unlock()

	file, err := os.OpenFile(adminAuditPath(), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o600)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(append(data, '\n'))
	return err
}

func readAdminAuditEvents(limit int) ([]adminAuditEvent, error) {
	if limit <= 0 || limit > 500 {
		limit = 100
	}
	file, err := os.Open(adminAuditPath())
	if os.IsNotExist(err) {
		return []adminAuditEvent{}, nil
	}
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var events []adminAuditEvent
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var event adminAuditEvent
		if err := json.Unmarshal(scanner.Bytes(), &event); err == nil {
			events = append(events, event)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	sort.SliceStable(events, func(i, j int) bool { return events[i].Timestamp > events[j].Timestamp })
	if len(events) > limit {
		events = events[:limit]
	}
	return events, nil
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
	for _, key := range []string{"player_id", "account_id", "actor_id", "controller_id", "item_id", "faction_id", "storage_id"} {
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
		trimmed := sanitizedAuditString(v, 128)
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
	value = strings.TrimSpace(value)
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

func handleAdminAuditEvents(w http.ResponseWriter, r *http.Request) {
	events, err := readAdminAuditEvents(100)
	if err != nil {
		jsonErr(w, err, http.StatusInternalServerError)
		return
	}
	jsonOK(w, events)
}
