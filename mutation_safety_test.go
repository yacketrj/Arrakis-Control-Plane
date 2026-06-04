package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestMutationSafetyGiveItem(t *testing.T) {
	got := mutationSafetyForPath(http.MethodPost, "/api/v1/players/give-item")
	if got.Risk != "high" || !got.RequiresReason || !got.RequiresPreview || got.Destructive {
		t.Fatalf("unexpected give-item classification: %#v", got)
	}
	if !strings.Contains(got.RecommendedPath, "Direct Inventory Write") {
		t.Fatalf("expected direct inventory write guidance, got %q", got.RecommendedPath)
	}
	if len(got.OperatorWarnings) == 0 {
		t.Fatalf("expected online-player warning")
	}
}

func TestMutationSafetyDeleteItem(t *testing.T) {
	got := mutationSafetyForPath(http.MethodDelete, "/api/v1/players/item/99")
	if got.Risk != "destructive" || !got.Destructive || !got.RequiresReason || !got.RequiresPreview {
		t.Fatalf("unexpected delete classification: %#v", got)
	}
	if got.RollbackHint == "" {
		t.Fatalf("expected rollback hint")
	}
}

func TestMutationSafetyGrantLive(t *testing.T) {
	got := mutationSafetyForPath(http.MethodPost, "/api/v1/players/grant-live")
	if got.Risk != "high" || !got.RequiresReason || !got.RequiresPreview {
		t.Fatalf("unexpected grant-live classification: %#v", got)
	}
	if !strings.Contains(got.RecommendedPath, "plain template-plus-amount") {
		t.Fatalf("expected claim reward guidance, got %q", got.RecommendedPath)
	}
}

func TestMutationSafetyTeleport(t *testing.T) {
	got := mutationSafetyForPath(http.MethodPost, "/api/v1/players/teleport")
	if got.Risk != "high" || got.RollbackHint == "" || len(got.OperatorWarnings) == 0 {
		t.Fatalf("unexpected teleport classification: %#v", got)
	}
}

func TestMutationSafetyStorage(t *testing.T) {
	got := mutationSafetyForPath(http.MethodPost, "/api/v1/storage/123/give-item")
	if got.Risk != "high" || len(got.OperatorWarnings) == 0 {
		t.Fatalf("unexpected storage classification: %#v", got)
	}
}

func TestMutationSafetyHighRiskCoverage(t *testing.T) {
	tests := []struct {
		method string
		path   string
	}{
		{http.MethodPost, "/api/v1/reconnect"},
		{http.MethodPost, "/api/v1/battlegroup/exec"},
		{http.MethodPost, "/api/v1/database/sql"},
		{http.MethodPost, "/api/v1/logs/stream-ticket"},
		{http.MethodPost, "/api/v1/notify"},
		{http.MethodPost, "/api/v1/players/item/stack-size"},
		{http.MethodPost, "/api/v1/players/repair-item"},
		{http.MethodPost, "/api/v1/players/teleport"},
		{http.MethodPost, "/api/v1/players/journey/complete"},
		{http.MethodPost, "/api/v1/players/set-spec-xp"},
		{http.MethodPost, "/api/v1/storage/123/give-item"},
	}
	for _, tc := range tests {
		got := mutationSafetyForPath(tc.method, tc.path)
		if got.Risk != "high" || !got.RequiresReason || !got.RequiresPreview || got.Destructive {
			t.Fatalf("%s %s: expected high risk requiring reason/preview, got %#v", tc.method, tc.path, got)
		}
	}
}

func TestMutationSafetyDestructiveCoverage(t *testing.T) {
	tests := []struct {
		method string
		path   string
	}{
		{http.MethodDelete, "/api/v1/auth/discord/player-links/123"},
		{http.MethodDelete, "/api/v1/players/item/99"},
		{http.MethodPost, "/api/v1/players/journey/reset"},
		{http.MethodPost, "/api/v1/players/journey/wipe"},
		{http.MethodPost, "/api/v1/players/delete-tutorials"},
		{http.MethodPost, "/api/v1/players/wipe-codex"},
		{http.MethodPost, "/api/v1/blueprints/import"},
	}
	for _, tc := range tests {
		got := mutationSafetyForPath(tc.method, tc.path)
		if got.Risk != "destructive" || !got.Destructive || !got.RequiresReason || !got.RequiresPreview {
			t.Fatalf("%s %s: expected destructive risk requiring reason/preview, got %#v", tc.method, tc.path, got)
		}
	}
}

func TestMutationSafetyReasonEnforcementFlag(t *testing.T) {
	t.Setenv("ADMIN_REQUIRE_REASON", "true")
	got := mutationSafetyForPath(http.MethodPost, "/api/v1/players/give-item")
	if !got.ReasonEnforcementEnabled {
		t.Fatalf("expected enforcement enabled in classification: %#v", got)
	}

	t.Setenv("ADMIN_REQUIRE_REASON", "false")
	got = mutationSafetyForPath(http.MethodPost, "/api/v1/players/give-item")
	if got.ReasonEnforcementEnabled {
		t.Fatalf("expected enforcement disabled in classification: %#v", got)
	}
}

func TestMutationSafetyClassifyHandler(t *testing.T) {
	t.Setenv("ADMIN_REQUIRE_REASON", "true")
	req := httptest.NewRequest(http.MethodGet, "/api/v1/mutation-safety/classify?method=DELETE&path=/api/v1/players/item/99", nil)
	res := httptest.NewRecorder()

	handleMutationSafetyClassify(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", res.Code, res.Body.String())
	}
	var got mutationSafetyClass
	if err := json.NewDecoder(res.Body).Decode(&got); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if got.Risk != "destructive" || !got.Destructive || !got.RequiresReason {
		t.Fatalf("unexpected handler response: %#v", got)
	}
	if !got.ReasonEnforcementEnabled {
		t.Fatalf("expected handler response to expose reason enforcement state: %#v", got)
	}
}

func TestMutationSafetyClassifyHandlerRequiresPath(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/mutation-safety/classify", nil)
	res := httptest.NewRecorder()

	handleMutationSafetyClassify(res, req)

	if res.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", res.Code)
	}
}

func TestMutationSafetyMiddlewareDisabledAllowsMissingReason(t *testing.T) {
	t.Setenv("ADMIN_REQUIRE_REASON", "false")
	called := false
	handler := mutationSafetyMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusNoContent)
	}))

	req := httptest.NewRequest(http.MethodPost, "/api/v1/players/give-item", strings.NewReader(`{"player_id":1}`))
	res := httptest.NewRecorder()
	handler.ServeHTTP(res, req)

	if !called {
		t.Fatalf("expected wrapped handler to be called")
	}
	if res.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", res.Code)
	}
}

func TestMutationSafetyMiddlewareRequiresHeaderReasonWhenEnabled(t *testing.T) {
	t.Setenv("ADMIN_REQUIRE_REASON", "true")
	handler := mutationSafetyMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))

	req := httptest.NewRequest(http.MethodPost, "/api/v1/players/give-item", strings.NewReader(`{"player_id":1}`))
	res := httptest.NewRecorder()
	handler.ServeHTTP(res, req)

	if res.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 without reason, got %d", res.Code)
	}

	req = httptest.NewRequest(http.MethodPost, "/api/v1/players/give-item", strings.NewReader(`{"player_id":1}`))
	req.Header.Set("X-Admin-Reason", "support correction")
	res = httptest.NewRecorder()
	handler.ServeHTTP(res, req)

	if res.Code != http.StatusNoContent {
		t.Fatalf("expected 204 with header reason, got %d", res.Code)
	}
}

func TestMutationSafetyMiddlewareAcceptsBodyReasonWhenEnabled(t *testing.T) {
	t.Setenv("ADMIN_REQUIRE_REASON", "true")
	var bodyAfterMiddleware map[string]any
	handler := mutationSafetyMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&bodyAfterMiddleware); err != nil {
			t.Fatalf("decode body after middleware: %v", err)
		}
		w.WriteHeader(http.StatusNoContent)
	}))

	req := httptest.NewRequest(http.MethodPost, "/api/v1/players/give-item", strings.NewReader(`{"player_id":1,"reason":"support correction"}`))
	res := httptest.NewRecorder()
	handler.ServeHTTP(res, req)

	if res.Code != http.StatusNoContent {
		t.Fatalf("expected 204 with body reason, got %d", res.Code)
	}
	if bodyAfterMiddleware["player_id"].(float64) != 1 {
		t.Fatalf("request body was not restored after middleware: %#v", bodyAfterMiddleware)
	}
}

func TestMutationSafetyMiddlewareRejectsOversizedInspectableBodyWhenReasonRequired(t *testing.T) {
	t.Setenv("ADMIN_REQUIRE_REASON", "true")
	handler := mutationSafetyMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))

	body := strings.NewReader(strings.Repeat("a", int(maxAuditInspectableBodyBytes)+1))
	req := httptest.NewRequest(http.MethodPost, "/api/v1/players/give-item", body)
	req.ContentLength = maxAuditInspectableBodyBytes + 1
	res := httptest.NewRecorder()
	handler.ServeHTTP(res, req)

	if res.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for oversized body when reason inspection is required, got %d", res.Code)
	}
}
