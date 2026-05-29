package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

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
