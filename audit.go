package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

const defaultAuditLogFile = "admin-audit.jsonl"

var (
	auditMu sync.Mutex
	idSegmentPattern = regexp.MustCompile(`^\d+$`)
)

type adminAuditEvent struct {
	Timestamp  string            `json:"timestamp"`
	Method     string            `json:"method"`
	Path       string            `json:"path"`
	Action     string            `json:"action"`
	Risk       string            `json:"risk,omitempty"`
	Reason     string            `json:"reason,omitempty"`
	Target     map[string]string `json:"target,omitempty"`
	Status     int               `json:"status"`
	DurationMS int64             `json:"duration_ms"`
	Result     string            `json:"result"`
	RemoteAddr string            `json:"remote_addr,omitempty"`
	Error      string            `json:"error,omitempty"`
}

type auditResponseWriter struct {
	http.ResponseWriter
	status int
}

func (w *auditResponseWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func (w *auditResponseWriter) Write(b []byte) (int, error) {
	if w.status == 0 {
		w.status = http.StatusOK
	}
	return w.ResponseWriter.Write(b)
}

func auditMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !shouldAuditRequest(r) {
			next.ServeHTTP(w, r)
			return
		}

		started := time.Now()
		aw := &auditResponseWriter{ResponseWriter: w}
		next.ServeHTTP(aw, r)

		status := aw.status
		if status == 0 {
			status = http.StatusOK
		}
		event := buildAuditEvent(r, status, time.Since(started))
		if err := recordAdminAction(event); err != nil {
			fmt.Fprintf(os.Stderr, "audit: failed to record admin action: %v\n", err)
		}
	})
}

func shouldAuditRequest(r *http.Request) bool {
	if r.Method == http.MethodOptions || isPublicPath(r.URL.Path) {
		return false
	}
	if strings.HasPrefix(r.URL.Path, "/api/v1/audit/") {
		return false
	}
	switch r.Method {
	case http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete:
		return true
	default:
		return false
	}
}

func buildAuditEvent(r *http.Request, status int, duration time.Duration) adminAuditEvent {
	result := "success"
	if status >= 400 {
		result = "failure"
	}
	event := adminAuditEvent{
		Timestamp:  time.Now().UTC().Format(time.RFC3339Nano),
		Method:     r.Method,
		Path:       r.URL.Path,
		Action:     auditActionName(r),
		Risk:       auditRisk(r),
		Reason:     sanitizedAuditReason(r),
		Target:     auditTarget(r),
		Status:     status,
		DurationMS: duration.Milliseconds(),
		Result:     result,
		RemoteAddr: auditRemoteAddr(r.RemoteAddr),
	}
	if result == "failure" {
		event.Error = http.StatusText(status)
	}
	return event
}

func recordAdminAction(event adminAuditEvent) error {
	path := auditLogPath()
	if dir := filepath.Dir(path); dir != "." && dir != "" {
		if err := os.MkdirAll(dir, 0700); err != nil {
			return err
		}
	}
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}
	auditMu.Lock()
	defer auditMu.Unlock()
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(append(data, '\n'))
	return err
}

func handleAdminAuditEvents(w http.ResponseWriter, r *http.Request) {
	limit := 250
	if raw := strings.TrimSpace(r.URL.Query().Get("limit")); raw != "" {
		if parsed, err := strconv.Atoi(raw); err == nil && parsed > 0 && parsed <= 1000 {
			limit = parsed
		}
	}
	events, err := readRecentAuditEvents(limit)
	if err != nil {
		if os.IsNotExist(err) {
			jsonOK(w, []adminAuditEvent{})
			return
		}
		jsonErr(w, err, http.StatusInternalServerError)
		return
	}
	jsonOK(w, events)
}

func readRecentAuditEvents(limit int) ([]adminAuditEvent, error) {
	f, err := os.Open(auditLogPath())
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var events []adminAuditEvent
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		var event adminAuditEvent
		if err := json.Unmarshal([]byte(line), &event); err == nil {
			events = append(events, event)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	if len(events) > limit {
		events = events[len(events)-limit:]
	}
	for i, j := 0, len(events)-1; i < j; i, j = i+1, j-1 {
		events[i], events[j] = events[j], events[i]
	}
	return events, nil
}

func auditLogPath() string {
	if path := strings.TrimSpace(os.Getenv("AUDIT_LOG_FILE")); path != "" {
		return path
	}
	return defaultAuditLogFile
}

func auditActionName(r *http.Request) string {
	path := strings.TrimPrefix(r.URL.Path, "/api/v1/")
	path = strings.Trim(path, "/")
	if path == "" {
		path = "root"
	}
	parts := strings.Split(path, "/")
	for i, part := range parts {
		if idSegmentPattern.MatchString(part) || part == "{id}" {
			parts[i] = "id"
		}
	}
	return strings.ToLower(r.Method) + ":" + strings.Join(parts, ".")
}

func auditRisk(r *http.Request) string {
	path := strings.ToLower(r.URL.Path)
	switch {
	case r.Method == http.MethodDelete:
		return "destructive"
	case strings.Contains(path, "/wipe") || strings.Contains(path, "/delete") || strings.Contains(path, "/reset") || strings.Contains(path, "/kick"):
		return "destructive"
	case strings.Contains(path, "/teleport") || strings.Contains(path, "/give") || strings.Contains(path, "/award") || strings.Contains(path, "/grant") || strings.Contains(path, "/repair") || strings.Contains(path, "/import"):
		return "high"
	case strings.Contains(path, "/reconnect") || strings.Contains(path, "/notify") || strings.Contains(path, "/refresh"):
		return "medium"
	default:
		return "medium"
	}
}

func sanitizedAuditReason(r *http.Request) string {
	reason := strings.TrimSpace(r.Header.Get("X-Admin-Reason"))
	if reason == "" {
		reason = strings.TrimSpace(r.URL.Query().Get("reason"))
	}
	reason = redactAuditString(reason)
	if len(reason) > 256 {
		reason = reason[:256]
	}
	return reason
}

func auditTarget(r *http.Request) map[string]string {
	target := map[string]string{}
	if id := r.PathValue("id"); id != "" {
		target["id"] = id
	}
	path := strings.TrimPrefix(r.URL.Path, "/api/v1/")
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) > 0 && parts[0] != "" {
		target["domain"] = parts[0]
	}
	if len(parts) > 1 && idSegmentPattern.MatchString(parts[1]) {
		target[parts[0]+"_id"] = parts[1]
	}
	for _, key := range []string{"player_id", "account_id", "actor_id", "item_id", "container_id"} {
		if value := strings.TrimSpace(r.URL.Query().Get(key)); value != "" {
			target[key] = redactAuditString(value)
		}
	}
	if len(target) == 0 {
		return nil
	}
	return target
}

func auditRemoteAddr(addr string) string {
	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		return redactAuditString(addr)
	}
	return host
}

func redactAuditString(value string) string {
	if value == "" {
		return ""
	}
	redacted := value
	patterns := []string{"token", "password", "passwd", "secret", "key", "authorization"}
	lower := strings.ToLower(redacted)
	for _, pattern := range patterns {
		if strings.Contains(lower, pattern) {
			return "[redacted]"
		}
	}
	return strings.ReplaceAll(redacted, adminToken, "[redacted]")
}
