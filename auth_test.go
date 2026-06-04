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
	allowedOrigins = "http://localhost:5173,https://dune-admin.layout.tools"
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
		"null",
		"file://localhost",
		"javascript:alert(1)",
		"https://user:pass@example.com",
		"https://example.com/path",
		"https://example.com?x=1",
		"https://example.com#frag",
		"https://example.com\n",
	}
	for _, origin := range invalid {
		if isAllowedOriginValue(origin) {
			t.Fatalf("expected origin to be rejected: %q", origin)
		}
	}

	valid := []string{"http://localhost:5173", "https://example.com", "https://example.com:8443"}
	for _, origin := range valid {
		if !isAllowedOriginValue(origin) {
			t.Fatalf("expected origin to be accepted: %q", origin)
		}
	}
}

func TestParseAllowedOriginsIgnoresUnsafeValues(t *testing.T) {
	allowed := parseAllowedOrigins("*,null,http://localhost:5173,https://example.com/path,https://admin.example")
	if allowed["*"] || allowed["null"] || allowed["https://example.com/path"] {
		t.Fatalf("unsafe origins should not be parsed as allowed: %#v", allowed)
	}
	if !allowed["http://localhost:5173"] || !allowed["https://admin.example"] {
		t.Fatalf("expected safe origins to be allowed: %#v", allowed)
	}
}

func TestCorsMiddlewareDoesNotReflectDisallowedOrigin(t *testing.T) {
	old := allowedOrigins
	allowedOrigins = "http://localhost:5173"
	t.Cleanup(func() { allowedOrigins = old })

	h := corsMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodOptions, "/api/v1/status", nil)
	req.Header.Set("Origin", "https://evil.example")
	res := httptest.NewRecorder()
	h.ServeHTTP(res, req)

	if res.Code != http.StatusNoContent {
		t.Fatalf("expected preflight no-content, got %d", res.Code)
	}
	if got := res.Header().Get("Access-Control-Allow-Origin"); got != "" {
		t.Fatalf("expected no reflected disallowed origin, got %q", got)
	}
}

func TestCorsMiddlewareReflectsAllowedOrigin(t *testing.T) {
	old := allowedOrigins
	allowedOrigins = "http://localhost:5173"
	t.Cleanup(func() { allowedOrigins = old })

	h := corsMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodOptions, "/api/v1/status", nil)
	req.Header.Set("Origin", "http://localhost:5173")
	res := httptest.NewRecorder()
	h.ServeHTTP(res, req)

	if res.Code != http.StatusNoContent {
		t.Fatalf("expected preflight no-content, got %d", res.Code)
	}
	if got := res.Header().Get("Access-Control-Allow-Origin"); got != "http://localhost:5173" {
		t.Fatalf("expected reflected allowed origin, got %q", got)
	}
	if got := res.Header().Get("Vary"); got != "Origin" {
		t.Fatalf("expected Vary Origin, got %q", got)
	}
}

func TestK8sNameValidation(t *testing.T) {
	valid := []string{"bgd-0", "funcom-operators", "pod123"}
	invalid := []string{"", "PodUpper", "pod;rm-rf", "../pod", "pod name", "-bad", "bad-"}
	for _, v := range valid {
		if !isValidK8sName(v) {
			t.Fatalf("expected %q to be valid", v)
		}
	}
	for _, v := range invalid {
		if isValidK8sName(v) {
			t.Fatalf("expected %q to be invalid", v)
		}
	}
}

func TestReadOnlySQL(t *testing.T) {
	allowed := []string{
		"select * from dune.actors",
		"WITH x AS (SELECT 1) SELECT * FROM x",
		"show search_path",
		"explain select * from dune.actors",
	}
	denied := []string{
		"delete from dune.actors",
		"update dune.actors set id = id",
		"drop table dune.actors",
		"select * from dune.actors; drop table dune.items",
		"insert into dune.items values (1)",
	}
	for _, q := range allowed {
		if !isReadOnlySQL(q) {
			t.Fatalf("expected query to be allowed: %q", q)
		}
	}
	for _, q := range denied {
		if isReadOnlySQL(q) {
			t.Fatalf("expected query to be denied: %q", q)
		}
	}
}
