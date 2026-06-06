package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
)

func TestAuditMiddlewareCapturesBlockedHighRiskMutationsMissingReason(t *testing.T) {
	tests := []struct {
		method string
		path   string
	}{
		{http.MethodDelete, "/api/v1/auth/discord/player-links/discord-123"},
		{http.MethodPost, "/api/v1/reconnect"},
		{http.MethodPost, "/api/v1/battlegroup/exec"},
		{http.MethodPost, "/api/v1/players/give-item"},
		{http.MethodPost, "/api/v1/players/give-currency"},
		{http.MethodPost, "/api/v1/players/grant-live"},
		{http.MethodPost, "/api/v1/players/give-faction-rep"},
		{http.MethodPost, "/api/v1/players/give-scrip"},
		{http.MethodPost, "/api/v1/players/award-xp"},
		{http.MethodPost, "/api/v1/players/award-char-xp"},
		{http.MethodPost, "/api/v1/players/award-intel"},
		{http.MethodPost, "/api/v1/players/kick"},
		{http.MethodDelete, "/api/v1/players/item/123"},
		{http.MethodPost, "/api/v1/players/item/stack-size"},
		{http.MethodPost, "/api/v1/players/reset-spec"},
		{http.MethodPost, "/api/v1/players/set-faction-tier"},
		{http.MethodPost, "/api/v1/players/journey/complete"},
		{http.MethodPost, "/api/v1/players/journey/reset"},
		{http.MethodPost, "/api/v1/players/journey/wipe"},
		{http.MethodPost, "/api/v1/players/delete-tutorials"},
		{http.MethodPost, "/api/v1/players/wipe-codex"},
		{http.MethodPost, "/api/v1/players/set-spec-xp"},
		{http.MethodPost, "/api/v1/players/repair-item"},
		{http.MethodPost, "/api/v1/players/teleport"},
		{http.MethodPost, "/api/v1/storage/123/give-item"},
		{http.MethodPost, "/api/v1/database/sql"},
		{http.MethodPost, "/api/v1/logs/stream-ticket"},
		{http.MethodPost, "/api/v1/notify"},
		{http.MethodPost, "/api/v1/blueprints/import"},
	}

	for i, tc := range tests {
		t.Run(fmt.Sprintf("%s %s", tc.method, tc.path), func(t *testing.T) {
			t.Setenv("ADMIN_AUDIT_LOG", filepath.Join(t.TempDir(), "audit.jsonl"))
			t.Setenv("ADMIN_REQUIRE_REASON", "1")

			safety := mutationSafetyForPath(tc.method, tc.path)
			if !safety.RequiresReason || (safety.Risk != "high" && safety.Risk != "destructive") {
				t.Fatalf("test route must require reason and be high/destructive, got %#v", safety)
			}

			downstreamCalled := false
			handler := auditMiddleware(mutationSafetyMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				downstreamCalled = true
				w.WriteHeader(http.StatusAccepted)
			})))

			body := strings.NewReader(`{"player_id":123,"account_id":77,"actor_id":"actor-99","item_id":456,"storage_id":789,"pod":"pod-1","cmd":"restart"}`)
			req := httptest.NewRequest(tc.method, tc.path, body)
			req.Header.Set("X-Request-ID", fmt.Sprintf("blocked-route-%02d", i))
			res := httptest.NewRecorder()

			handler.ServeHTTP(res, req)
			if downstreamCalled {
				t.Fatal("blocked mutation reached downstream handler")
			}
			if res.Code != http.StatusBadRequest {
				t.Fatalf("expected blocked mutation status 400, got %d", res.Code)
			}

			events, err := readAdminAuditEvents(10)
			if err != nil {
				t.Fatalf("read audit events: %v", err)
			}
			if len(events) != 1 {
				t.Fatalf("expected 1 audit event, got %d", len(events))
			}
			event := events[0]
			if event.Method != tc.method || event.Path != tc.path {
				t.Fatalf("unexpected method/path: %#v", event)
			}
			if event.Action != safety.Action || event.Risk != safety.Risk || event.Destructive != safety.Destructive {
				t.Fatalf("unexpected safety metadata: got %#v want %#v", event, safety)
			}
			if !event.RequiresReason || !event.RequiresPreview {
				t.Fatalf("expected reason and preview flags on blocked mutation event: %#v", event)
			}
			if event.Status != http.StatusBadRequest || event.Result != "failure" {
				t.Fatalf("expected failure audit event with status 400, got %#v", event)
			}
			if event.Reason != "" {
				t.Fatalf("blocked missing-reason event should not invent a reason, got %q", event.Reason)
			}
			if event.Target["player_id"] != "123" || event.Target["account_id"] != "77" || event.Target["actor_id"] != "actor-99" {
				t.Fatalf("expected target metadata on blocked mutation, got %#v", event.Target)
			}
			if event.RequestID == "" {
				t.Fatalf("expected request id on blocked mutation event: %#v", event)
			}
		})
	}
}

func TestAuditMiddlewareCapturesBlockedOversizedReasonInspection(t *testing.T) {
	t.Setenv("ADMIN_AUDIT_LOG", filepath.Join(t.TempDir(), "audit.jsonl"))
	t.Setenv("ADMIN_REQUIRE_REASON", "1")

	downstreamCalled := false
	handler := auditMiddleware(mutationSafetyMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		downstreamCalled = true
		w.WriteHeader(http.StatusAccepted)
	})))

	oversizedBody := strings.NewReader(strings.Repeat("x", int(maxAuditInspectableBodyBytes)+1))
	req := httptest.NewRequest(http.MethodPost, "/api/v1/players/give-item", oversizedBody)
	req.Header.Set("X-Request-ID", "oversized-reason-body")
	res := httptest.NewRecorder()

	handler.ServeHTTP(res, req)
	if downstreamCalled {
		t.Fatal("oversized blocked mutation reached downstream handler")
	}
	if res.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400 for oversized reason-inspection body, got %d", res.Code)
	}

	events, err := readAdminAuditEvents(10)
	if err != nil {
		t.Fatalf("read audit events: %v", err)
	}
	if len(events) != 1 {
		t.Fatalf("expected 1 audit event, got %d", len(events))
	}
	event := events[0]
	if event.Action != "post:players.give-item" || event.Risk != "high" {
		t.Fatalf("unexpected audit classification: %#v", event)
	}
	if event.Status != http.StatusBadRequest || event.Result != "failure" {
		t.Fatalf("expected blocked oversized body failure audit, got %#v", event)
	}
	if event.Target != nil {
		t.Fatalf("oversized body should not be parsed into target metadata, got %#v", event.Target)
	}
}
