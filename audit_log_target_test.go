package main

import (
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
)

func TestExtractMutationAuditMetadataCapturesExpandedTargetFields(t *testing.T) {
	body := `{"reason":"target coverage","player_id":123,"account_id":77,"actor_id":"actor-1","controller_id":"controller-1","fls_id":"0xabc123","item_id":456,"item_template":"StaticItem'/Game/Items/Water'","item_template_id":"water-template","template_id":"template-alias","quantity":5,"amount":6,"quality":0,"vehicle_id":"sandbike-01","guild_id":"guild-1","rank":"officer","command_path":"rmq.add_item","command":"AddItemToInventory"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/players/give-item", strings.NewReader(body))
	metadata := extractMutationAuditMetadata(req)

	want := map[string]string{
		"player_id":        "123",
		"account_id":       "77",
		"actor_id":         "actor-1",
		"controller_id":    "controller-1",
		"fls_id":           "0xabc123",
		"item_id":          "456",
		"item_template":    "StaticItem'/Game/Items/Water'",
		"item_template_id": "water-template",
		"template_id":      "template-alias",
		"quantity":         "5",
		"amount":           "6",
		"quality":          "0",
		"vehicle_id":       "sandbike-01",
		"guild_id":         "guild-1",
		"rank":             "officer",
		"command_path":     "rmq.add_item",
		"command":          "AddItemToInventory",
	}
	for key, expected := range want {
		if got := metadata.Target[key]; got != expected {
			t.Fatalf("target[%q] = %q, want %q; target=%#v", key, got, expected, metadata.Target)
		}
	}
	if metadata.Reason != "target coverage" {
		t.Fatalf("reason = %q, want target coverage", metadata.Reason)
	}
}

func TestAuditMiddlewareCapturesRouteSpecificTargetMetadata(t *testing.T) {
	t.Setenv("ADMIN_AUDIT_LOG", filepath.Join(t.TempDir(), "audit.jsonl"))
	handler := auditMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusAccepted)
	}))

	body := `{"reason":"target coverage","player_id":123,"actor_id":"actor-1","fls_id":"0xabc123","item_template":"StaticItem'/Game/Items/Water'","quantity":5,"quality":0,"command_path":"rmq.add_item"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/players/give-item", strings.NewReader(body))
	req.Header.Set("X-Request-ID", "route-target-coverage")
	res := httptest.NewRecorder()
	handler.ServeHTTP(res, req)

	if res.Code != http.StatusAccepted {
		t.Fatalf("status = %d, want 202", res.Code)
	}
	events, err := readAdminAuditEvents(10)
	if err != nil {
		t.Fatalf("read audit events: %v", err)
	}
	if len(events) != 1 {
		t.Fatalf("expected one audit event, got %d", len(events))
	}
	event := events[0]
	want := map[string]string{
		"player_id":     "123",
		"actor_id":      "actor-1",
		"fls_id":        "0xabc123",
		"item_template": "StaticItem'/Game/Items/Water'",
		"quantity":      "5",
		"quality":       "0",
		"command_path":  "rmq.add_item",
	}
	for key, expected := range want {
		if got := event.Target[key]; got != expected {
			t.Fatalf("target[%q] = %q, want %q; target=%#v", key, got, expected, event.Target)
		}
	}
	if event.RequestID != "route-target-coverage" {
		t.Fatalf("request id = %q, want route-target-coverage", event.RequestID)
	}
}

func TestAuditTargetMetadataRedactsExpandedSensitiveFields(t *testing.T) {
	body := `{"reason":"target coverage","player_id":123,"fls_id":"token=secret-token-value","command_path":"rabbitmqctl eval password=supersecret","item_template":"safe-template"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/players/give-item", strings.NewReader(body))
	metadata := extractMutationAuditMetadata(req)

	for key, value := range metadata.Target {
		if strings.Contains(value, "secret-token-value") || strings.Contains(value, "supersecret") {
			t.Fatalf("target[%q] was not redacted: %q", key, value)
		}
	}
	if metadata.Target["item_template"] != "safe-template" {
		t.Fatalf("safe target metadata was unexpectedly changed: %#v", metadata.Target)
	}
}
