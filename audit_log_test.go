package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
)

func TestAuditMiddlewareCapturesMutatingRequests(t *testing.T) {
	t.Setenv("ADMIN_AUDIT_LOG", filepath.Join(t.TempDir(), "audit.jsonl"))

	handler := auditMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
	}))

	req := httptest.NewRequest(http.MethodPost, "/api/v1/players/give-item", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	events, err := readAdminAuditEvents(10)
	if err != nil {
		t.Fatalf("readAdminAuditEvents returned error: %v", err)
	}
	if len(events) != 1 {
		t.Fatalf("expected one audit event, got %d", len(events))
	}
	event := events[0]
	if event.Method != http.MethodPost {
		t.Fatalf("expected POST method, got %q", event.Method)
	}
	if event.Path != "/api/v1/players/give-item" {
		t.Fatalf("unexpected path: %q", event.Path)
	}
	if event.Action != "post:players.give-item" {
		t.Fatalf("unexpected action: %q", event.Action)
	}
	if event.Risk != "high" {
		t.Fatalf("expected high risk classification, got %q", event.Risk)
	}
	if event.Status != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", event.Status)
	}
	if event.Result != "success" {
		t.Fatalf("expected success result, got %q", event.Result)
	}
}

func TestAuditMiddlewareCapturesFailureStatus(t *testing.T) {
	t.Setenv("ADMIN_AUDIT_LOG", filepath.Join(t.TempDir(), "audit.jsonl"))

	handler := auditMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "bad request", http.StatusBadRequest)
	}))

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/players/item/123", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	events, err := readAdminAuditEvents(10)
	if err != nil {
		t.Fatalf("readAdminAuditEvents returned error: %v", err)
	}
	if len(events) != 1 {
		t.Fatalf("expected one audit event, got %d", len(events))
	}
	if events[0].Result != "failure" {
		t.Fatalf("expected failure result, got %q", events[0].Result)
	}
	if events[0].Risk != "destructive" {
		t.Fatalf("expected destructive risk, got %q", events[0].Risk)
	}
	if events[0].Status != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", events[0].Status)
	}
}

func TestAuditMiddlewareCapturesReasonAndTargetMetadata(t *testing.T) {
	t.Setenv("ADMIN_AUDIT_LOG", filepath.Join(t.TempDir(), "audit.jsonl"))

	var downstreamBody string
	handler := auditMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("downstream read body: %v", err)
		}
		downstreamBody = string(body)
		w.WriteHeader(http.StatusOK)
	}))

	body := `{"player_id":42,"account_id":77,"reason":" support grant\nverified by admin ","admin_token":"must-not-log"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/players/give-item", strings.NewReader(body))
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if downstreamBody != body {
		t.Fatalf("audit middleware did not restore request body, got %q", downstreamBody)
	}
	events, err := readAdminAuditEvents(10)
	if err != nil {
		t.Fatalf("readAdminAuditEvents returned error: %v", err)
	}
	if len(events) != 1 {
		t.Fatalf("expected one audit event, got %d", len(events))
	}
	event := events[0]
	if event.Reason != "support grant verified by admin" {
		t.Fatalf("unexpected sanitized reason: %q", event.Reason)
	}
	if event.Target["player_id"] != "42" || event.Target["account_id"] != "77" {
		t.Fatalf("unexpected target metadata: %#v", event.Target)
	}
	if _, exists := event.Target["admin_token"]; exists {
		t.Fatalf("secret-like field should not be included in target metadata: %#v", event.Target)
	}
}

func TestAuditMiddlewareSkipsReadOnlyAndPublicRequests(t *testing.T) {
	t.Setenv("ADMIN_AUDIT_LOG", filepath.Join(t.TempDir(), "audit.jsonl"))

	handler := auditMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	requests := []*http.Request{
		httptest.NewRequest(http.MethodGet, "/api/v1/players", nil),
		httptest.NewRequest(http.MethodPost, "/api/v1/public/status", nil),
	}
	for _, req := range requests {
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
	}

	events, err := readAdminAuditEvents(10)
	if err != nil {
		t.Fatalf("readAdminAuditEvents returned error: %v", err)
	}
	if len(events) != 0 {
		t.Fatalf("expected no audit events, got %#v", events)
	}
}

func TestReadAdminAuditEventsLimitAndSort(t *testing.T) {
	t.Setenv("ADMIN_AUDIT_LOG", filepath.Join(t.TempDir(), "audit.jsonl"))

	oldEvent := adminAuditEvent{Timestamp: "2026-05-23T10:00:00Z", Method: http.MethodPost, Path: "/api/v1/old", Action: "post:old", Status: 200, Result: "success"}
	newEvent := adminAuditEvent{Timestamp: "2026-05-23T11:00:00Z", Method: http.MethodPost, Path: "/api/v1/new", Action: "post:new", Status: 200, Result: "success"}
	if err := appendAdminAuditEvent(oldEvent); err != nil {
		t.Fatalf("append old event: %v", err)
	}
	if err := appendAdminAuditEvent(newEvent); err != nil {
		t.Fatalf("append new event: %v", err)
	}

	events, err := readAdminAuditEvents(1)
	if err != nil {
		t.Fatalf("readAdminAuditEvents returned error: %v", err)
	}
	if len(events) != 1 {
		t.Fatalf("expected one event after limit, got %d", len(events))
	}
	if events[0].Action != "post:new" {
		t.Fatalf("expected newest event first, got %#v", events[0])
	}
}

func TestMutationRiskForExpandedHighRiskRoutes(t *testing.T) {
	tests := []struct {
		method string
		path   string
		want   string
	}{
		{http.MethodPost, "/api/v1/battlegroup/exec", "high"},
		{http.MethodPost, "/api/v1/players/award-char-xp", "high"},
		{http.MethodPost, "/api/v1/players/set-spec-xp", "high"},
		{http.MethodPost, "/api/v1/players/repair-item", "high"},
		{http.MethodPost, "/api/v1/storage/42/give-item", "high"},
		{http.MethodPost, "/api/v1/players/reset-spec", "destructive"},
		{http.MethodPost, "/api/v1/blueprints/import", "destructive"},
		{http.MethodPost, "/api/v1/database/sql", "medium"},
	}

	for _, tt := range tests {
		if got := mutationRiskForRequest(tt.method, tt.path); got != tt.want {
			t.Fatalf("mutationRiskForRequest(%q, %q) = %q, want %q", tt.method, tt.path, got, tt.want)
		}
	}
}

func TestAuditMiddlewareWritesRedactedEvent(t *testing.T) {
	auditPath := t.TempDir() + "/admin-audit.jsonl"
	t.Setenv("ADMIN_AUDIT_LOG", auditPath)

	handler := auditMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))

	body := `{"player_id":123,"item_id":"water ADMIN_TOKEN=abcdefghijklmnopqrstuvwxyz0123456789ABCDEFG","reason":"support correction DB_PASS=supersecret"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/players/give-item", strings.NewReader(body))
	req.RemoteAddr = "192.0.2.10:54321"
	req.Header.Set("X-Admin-Token", "abcdefghijklmnopqrstuvwxyz0123456789ABCDEFG")
	req.Header.Set("X-Request-ID", "req-123")
	res := httptest.NewRecorder()

	handler.ServeHTTP(res, req)
	if res.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", res.Code)
	}

	events, err := readAdminAuditEvents(10)
	if err != nil {
		t.Fatalf("read audit events: %v", err)
	}
	if len(events) != 1 {
		t.Fatalf("expected 1 audit event, got %d", len(events))
	}
	event := events[0]
	if event.Risk != "high" || !event.RequiresReason || !event.RequiresPreview {
		t.Fatalf("unexpected safety classification in audit event: %#v", event)
	}
	if event.RequestID != "req-123" {
		t.Fatalf("expected request id, got %q", event.RequestID)
	}
	if event.RemoteAddr != "192.0.2.10" {
		t.Fatalf("expected remote host only, got %q", event.RemoteAddr)
	}
	if len(event.AdminTokenHash) != 16 {
		t.Fatalf("expected short token hash prefix, got %q", event.AdminTokenHash)
	}
	if strings.Contains(event.Reason, "supersecret") || !strings.Contains(event.Reason, "[REDACTED]") {
		t.Fatalf("expected redacted reason, got %q", event.Reason)
	}
	if itemID := event.Target["item_id"]; strings.Contains(itemID, "abcdefghijklmnopqrstuvwxyz0123456789ABCDEFG") || !strings.Contains(itemID, "[REDACTED]") {
		t.Fatalf("expected redacted item target, got %q", itemID)
	}
}

func TestAdminReasonEnforcementDefaultsOn(t *testing.T) {
	t.Setenv("ADMIN_REQUIRE_REASON", "")
	if !adminReasonEnforcementEnabled() {
		t.Fatal("expected admin reason enforcement to default on")
	}
}
