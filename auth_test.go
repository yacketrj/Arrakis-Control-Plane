package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

const testStrictAdminToken = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOP_"

func TestBearerToken(t *testing.T) {
	if got := bearerToken("Bearer secret"); got != "secret" {
		t.Fatalf("expected bearer token, got %q", got)
	}
	if got := bearerToken("Basic secret"); got != "" {
		t.Fatalf("expected empty token for non-bearer auth, got %q", got)
	}
	if got := bearerToken("Bearer"); got != "" {
		t.Fatalf("expected empty token for malformed bearer auth, got %q", got)
	}
}

func TestAuthMiddleware(t *testing.T) {
	old := adminToken
	adminToken = testStrictAdminToken
	t.Cleanup(func() { adminToken = old })

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
	h := authMiddleware(next)

	missing := httptest.NewRecorder()
	h.ServeHTTP(missing, httptest.NewRequest(http.MethodGet, "/api/v1/status", nil))
	if missing.Code != http.StatusUnauthorized {
		t.Fatalf("missing token: got %d", missing.Code)
	}

	bad := httptest.NewRecorder()
	badReq := httptest.NewRequest(http.MethodGet, "/api/v1/status", nil)
	badReq.Header.Set("Authorization", "Bearer wrong")
	h.ServeHTTP(bad, badReq)
	if bad.Code != http.StatusUnauthorized {
		t.Fatalf("bad token: got %d", bad.Code)
	}

	good := httptest.NewRecorder()
	goodReq := httptest.NewRequest(http.MethodGet, "/api/v1/status", nil)
	goodReq.Header.Set("Authorization", "Bearer "+testStrictAdminToken)
	h.ServeHTTP(good, goodReq)
	if good.Code != http.StatusNoContent {
		t.Fatalf("good token: got %d", good.Code)
	}
}

func TestAuthMiddlewareRejectsInvalidBackendToken(t *testing.T) {
	old := adminToken
	adminToken = "test-token"
	t.Cleanup(func() { adminToken = old })

	h := authMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))

	resp := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/status", nil)
	req.Header.Set("Authorization", "Bearer test-token")
	h.ServeHTTP(resp, req)

	if resp.Code != http.StatusServiceUnavailable {
		t.Fatalf("invalid backend token: got %d", resp.Code)
	}
}

func TestValidateStrictAdminToken(t *testing.T) {
	valid := []string{
		testStrictAdminToken,
		generateAdminToken(),
	}
	for _, token := range valid {
		if err := validateStrictAdminToken(token); err != nil {
			t.Fatalf("expected valid strict token %q: %v", token, err)
		}
	}

	invalid := []string{
		"",
		"test-token",
		"abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOP",
		"abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNO+",
		"abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNO/",
		" abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNO_",
		"abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNO_\n",
	}
	for _, token := range invalid {
		if err := validateStrictAdminToken(token); err == nil {
			t.Fatalf("expected invalid strict token %q", token)
		}
	}
}

func TestOriginAllowed(t *testing.T) {
	old := allowedOrigins
	allowedOrigins = "http://localhost:5173,https://arrakis-control-panel.layout.tools"
	t.Cleanup(func() { allowedOrigins = old })

	if !originAllowed("http://localhost:5173") {
		t.Fatal("expected configured localhost origin to be allowed")
	}
	if originAllowed("https://evil.example") {
		t.Fatal("expected unconfigured origin to be denied")
	}
	if !originAllowed("") {
		t.Fatal("expected non-browser requests without Origin to be allowed")
	}
}

func TestAllowedOriginValueRejectsUnsafeOrigins(t *testing.T) {
	invalid := []string{
		"*",
		"http://*",
		"https://*",
		"http://example.com/*",
		"javascript:alert(1)",
	}
	for _, origin := range invalid {
		if err := validateAllowedOriginValue(origin); err == nil {
			t.Fatalf("expected unsafe origin %q to be rejected", origin)
		}
	}
}
