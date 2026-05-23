package main

import (
	"bufio"
	"encoding/json"
	"net/http"
	"os"
	"sort"
	"strings"
	"sync"
	"time"
)

const defaultAdminAuditPath = "admin-audit.jsonl"

type adminAuditEvent struct {
	Timestamp  string `json:"timestamp"`
	Method     string `json:"method"`
	Path       string `json:"path"`
	Action     string `json:"action"`
	Status     int    `json:"status"`
	DurationMS int64  `json:"duration_ms"`
	Result     string `json:"result"`
}

type statusCaptureWriter struct {
	http.ResponseWriter
	status int
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

		started := time.Now()
		capture := &statusCaptureWriter{ResponseWriter: w}
		next.ServeHTTP(capture, r)

		status := capture.status
		if status == 0 {
			status = http.StatusOK
		}
		_ = appendAdminAuditEvent(adminAuditEvent{
			Timestamp:  started.UTC().Format(time.RFC3339Nano),
			Method:     r.Method,
			Path:       r.URL.Path,
			Action:     auditActionName(r.Method, r.URL.Path),
			Status:     status,
			DurationMS: time.Since(started).Milliseconds(),
			Result:     auditResultForStatus(status),
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

func handleAdminAuditEvents(w http.ResponseWriter, r *http.Request) {
	events, err := readAdminAuditEvents(100)
	if err != nil {
		jsonErr(w, err, http.StatusInternalServerError)
		return
	}
	jsonOK(w, events)
}
