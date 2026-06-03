package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

type routeExpectation struct {
	method string
	path   string
}

func appsecNoopHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
}

func TestAppSecPublicPathAllowlist(t *testing.T) {
	allowed := []string{
		"/api/v1/public/status",
		"/api/v1/auth/discord/login",
		"/api/v1/auth/discord/callback",
	}
	for _, path := range allowed {
		if !isPublicPath(path) {
			t.Fatalf("expected %s to be public", path)
		}
	}

	denied := []string{
		"/api/v1/status",
		"/api/v1/auth/discord/me",
		"/api/v1/auth/discord/logout",
		"/api/v1/auth/discord/users",
		"/api/v1/auth/discord/player-links",
		"/api/v1/self/player-card",
		"/api/v1/logs/stream",
		"/api/v1/public/status/extra",
	}
	for _, path := range denied {
		if isPublicPath(path) {
			t.Fatalf("expected %s to require non-public auth handling", path)
		}
	}
}

func TestAppSecSelfServicePathClassification(t *testing.T) {
	allowed := []string{
		"/api/v1/self/player-link",
		"/api/v1/self/player-card",
	}
	for _, path := range allowed {
		if !isSelfServicePath(path) {
			t.Fatalf("expected %s to be classified as self-service", path)
		}
	}

	denied := []string{
		"/api/v1/self",
		"/api/v1/auth/discord/me",
		"/api/v1/auth/discord/logout",
		"/api/v1/auth/discord/player-links",
		"/api/v1/players/123/profile",
	}
	for _, path := range denied {
		if isSelfServicePath(path) {
			t.Fatalf("expected %s not to be classified as self-service", path)
		}
	}
}

func TestAppSecPublicRoutesBypassInvalidBackendAdminToken(t *testing.T) {
	old := adminToken
	adminToken = "invalid-test-token"
	t.Cleanup(func() { adminToken = old })

	h := authMiddleware(appsecNoopHandler())
	for _, tc := range []routeExpectation{
		{http.MethodGet, "/api/v1/public/status"},
		{http.MethodGet, "/api/v1/auth/discord/login"},
		{http.MethodGet, "/api/v1/auth/discord/callback"},
	} {
		resp := httptest.NewRecorder()
		h.ServeHTTP(resp, httptest.NewRequest(tc.method, tc.path, nil))
		if resp.Code != http.StatusNoContent {
			t.Fatalf("%s %s: expected public route to bypass backend token validation, got %d", tc.method, tc.path, resp.Code)
		}
	}
}

func TestAppSecAdminRoutesRequireConfiguredAdminTokenAndPresentedToken(t *testing.T) {
	old := adminToken
	adminToken = testStrictAdminToken
	t.Cleanup(func() { adminToken = old })

	h := authMiddleware(appsecNoopHandler())
	routes := []routeExpectation{
		{http.MethodGet, "/api/v1/status"},
		{http.MethodPost, "/api/v1/reconnect"},
		{http.MethodGet, "/api/v1/connectivity/diagnostics"},
		{http.MethodGet, "/api/v1/audit/events"},
		{http.MethodGet, "/api/v1/battlegroup/status"},
		{http.MethodPost, "/api/v1/battlegroup/exec"},
		{http.MethodGet, "/api/v1/players"},
		{http.MethodGet, "/api/v1/players/123/inventory"},
		{http.MethodPost, "/api/v1/players/give-item"},
		{http.MethodDelete, "/api/v1/players/item/123"},
		{http.MethodPost, "/api/v1/players/item/stack-size"},
		{http.MethodGet, "/api/v1/inventory/requests"},
		{http.MethodPatch, "/api/v1/inventory/requests/request-1"},
		{http.MethodGet, "/api/v1/database/tables"},
		{http.MethodPost, "/api/v1/database/sql"},
		{http.MethodGet, "/api/v1/logs/pods"},
		{http.MethodPost, "/api/v1/logs/stream-ticket"},
		{http.MethodPost, "/api/v1/notify"},
		{http.MethodGet, "/api/v1/storage"},
		{http.MethodPost, "/api/v1/storage/123/give-item"},
		{http.MethodGet, "/api/v1/blueprints"},
		{http.MethodPost, "/api/v1/blueprints/import"},
	}

	for _, tc := range routes {
		missing := httptest.NewRecorder()
		h.ServeHTTP(missing, httptest.NewRequest(tc.method, tc.path, nil))
		if missing.Code != http.StatusUnauthorized {
			t.Fatalf("%s %s: missing token got %d", tc.method, tc.path, missing.Code)
		}

		allowed := httptest.NewRecorder()
		req := httptest.NewRequest(tc.method, tc.path, nil)
		req.Header.Set("X-Admin-Token", testStrictAdminToken)
		h.ServeHTTP(allowed, req)
		if allowed.Code != http.StatusNoContent {
			t.Fatalf("%s %s: valid token got %d", tc.method, tc.path, allowed.Code)
		}
	}
}

func TestAppSecSelfServiceRoutesDenyWithoutDiscordSessionOrAdminToken(t *testing.T) {
	old := adminToken
	adminToken = testStrictAdminToken
	t.Cleanup(func() { adminToken = old })

	h := authMiddleware(appsecNoopHandler())
	for _, path := range []string{"/api/v1/self/player-link", "/api/v1/self/player-card"} {
		missing := httptest.NewRecorder()
		h.ServeHTTP(missing, httptest.NewRequest(http.MethodGet, path, nil))
		if missing.Code != http.StatusUnauthorized {
			t.Fatalf("%s: missing session/token got %d", path, missing.Code)
		}

		allowed := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, path, nil)
		req.Header.Set("Authorization", "Bearer "+testStrictAdminToken)
		h.ServeHTTP(allowed, req)
		if allowed.Code != http.StatusNoContent {
			t.Fatalf("%s: admin token got %d", path, allowed.Code)
		}
	}
}

func TestAppSecWebSocketLogStreamRequiresOneTimeTicket(t *testing.T) {
	old := adminToken
	adminToken = testStrictAdminToken
	t.Cleanup(func() { adminToken = old })

	h := authMiddleware(appsecNoopHandler())
	req := httptest.NewRequest(http.MethodGet, "/api/v1/logs/stream?ns=default&pod=pod-0", nil)
	req.Header.Set("Upgrade", "websocket")
	req.Header.Set("Connection", "Upgrade")
	req.Header.Set("X-Admin-Token", testStrictAdminToken)

	resp := httptest.NewRecorder()
	h.ServeHTTP(resp, req)
	if resp.Code != http.StatusUnauthorized {
		t.Fatalf("expected WebSocket log stream without ticket to be denied before admin-token fallback, got %d", resp.Code)
	}
}
