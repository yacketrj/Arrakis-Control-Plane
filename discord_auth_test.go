package main

import (
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
	"time"
)

func resetDiscordSessionsForTest(t *testing.T) {
	t.Helper()
	discordSessionsMu.Lock()
	oldSessions := discordSessions
	discordSessions = map[string]discordSession{}
	discordSessionsMu.Unlock()

	t.Cleanup(func() {
		discordSessionsMu.Lock()
		discordSessions = oldSessions
		discordSessionsMu.Unlock()
	})
}

func TestRegisterRoutesIncludesDiscordAuthEndpoints(t *testing.T) {
	resetDiscordSessionsForTest(t)
	t.Setenv("DISCORD_AUTH_ENABLED", "0")
	t.Setenv("DISCORD_USER_STORE", filepath.Join(t.TempDir(), "discord-users.json"))

	mux := http.NewServeMux()
	registerRoutes(mux)

	tests := []struct {
		method string
		path   string
		want   int
	}{
		{method: http.MethodGet, path: "/api/v1/auth/discord/login", want: http.StatusServiceUnavailable},
		{method: http.MethodGet, path: "/api/v1/auth/discord/callback", want: http.StatusServiceUnavailable},
		{method: http.MethodGet, path: "/api/v1/auth/discord/me", want: http.StatusUnauthorized},
		{method: http.MethodPost, path: "/api/v1/auth/discord/logout", want: http.StatusOK},
		{method: http.MethodGet, path: "/api/v1/auth/discord/users", want: http.StatusOK},
	}

	for _, tc := range tests {
		recorder := httptest.NewRecorder()
		req := httptest.NewRequest(tc.method, tc.path, nil)
		mux.ServeHTTP(recorder, req)
		if recorder.Code != tc.want {
			t.Fatalf("%s %s: got %d, want %d", tc.method, tc.path, recorder.Code, tc.want)
		}
	}
}

func TestMapDiscordRoles(t *testing.T) {
	t.Setenv("DISCORD_ADMIN_ROLE_IDS", "admin-role, owner-role")
	t.Setenv("DISCORD_NORMAL_ROLE_IDS", "member-role, helper-role")

	if got := mapDiscordRoles([]string{"member-role", "admin-role"}); got != appRoleAdmin {
		t.Fatalf("admin role should win over normal role, got %q", got)
	}
	if got := mapDiscordRoles([]string{"helper-role"}); got != appRoleNormal {
		t.Fatalf("normal role not mapped, got %q", got)
	}
	if got := mapDiscordRoles([]string{"unknown-role"}); got != appRoleNone {
		t.Fatalf("unknown role should map to none, got %q", got)
	}
}

func TestMapDiscordRolesDefaultsToNormalWhenNoNormalRoleListConfigured(t *testing.T) {
	t.Setenv("DISCORD_ADMIN_ROLE_IDS", "admin-role")
	t.Setenv("DISCORD_NORMAL_ROLE_IDS", "")

	if got := mapDiscordRoles([]string{"some-guild-role"}); got != appRoleNormal {
		t.Fatalf("expected default normal role when no normal role allowlist is configured, got %q", got)
	}
}

func TestDiscordSessionFromRequestAndRoleHelpers(t *testing.T) {
	resetDiscordSessionsForTest(t)

	sessionID := "session-token"
	discordSessionsMu.Lock()
	discordSessions[sessionID] = discordSession{ID: sessionID, DiscordID: "discord-123", Role: appRoleAdmin, ExpiresAt: time.Now().Add(time.Hour)}
	discordSessionsMu.Unlock()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/status", nil)
	req.AddCookie(&http.Cookie{Name: discordSessionCookieName, Value: sessionID})

	session, ok := discordSessionFromRequest(req)
	if !ok {
		t.Fatal("expected valid Discord session")
	}
	if session.DiscordID != "discord-123" || session.Role != appRoleAdmin {
		t.Fatalf("unexpected session: %#v", session)
	}
	if !discordSessionIsAdmin(req) {
		t.Fatal("expected admin session helper to return true")
	}
	if !discordSessionIsRegistered(req) {
		t.Fatal("expected registered session helper to return true")
	}
	if got := discordSessionRole(req); got != appRoleAdmin {
		t.Fatalf("expected admin role helper, got %q", got)
	}
	if got := discordSessionHash(req); got == "" || len(got) != 16 {
		t.Fatalf("expected stable 16-character session hash, got %q", got)
	}
}

func TestDiscordSessionFromRequestExpiresAndEvicts(t *testing.T) {
	resetDiscordSessionsForTest(t)

	sessionID := "expired-session-token"
	discordSessionsMu.Lock()
	discordSessions[sessionID] = discordSession{ID: sessionID, DiscordID: "discord-123", Role: appRoleAdmin, ExpiresAt: time.Now().Add(-time.Minute)}
	discordSessionsMu.Unlock()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/status", nil)
	req.AddCookie(&http.Cookie{Name: discordSessionCookieName, Value: sessionID})

	if _, ok := discordSessionFromRequest(req); ok {
		t.Fatal("expected expired Discord session to be rejected")
	}

	discordSessionsMu.Lock()
	_, exists := discordSessions[sessionID]
	discordSessionsMu.Unlock()
	if exists {
		t.Fatal("expected expired Discord session to be evicted")
	}
}

func TestHandleDiscordLogoutClearsSession(t *testing.T) {
	resetDiscordSessionsForTest(t)

	sessionID := "logout-session-token"
	discordSessionsMu.Lock()
	discordSessions[sessionID] = discordSession{ID: sessionID, DiscordID: "discord-123", Role: appRoleNormal, ExpiresAt: time.Now().Add(time.Hour)}
	discordSessionsMu.Unlock()

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/discord/logout", nil)
	req.AddCookie(&http.Cookie{Name: discordSessionCookieName, Value: sessionID})
	recorder := httptest.NewRecorder()

	handleDiscordLogout(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf("logout got %d", recorder.Code)
	}
	discordSessionsMu.Lock()
	_, exists := discordSessions[sessionID]
	discordSessionsMu.Unlock()
	if exists {
		t.Fatal("expected logout to remove Discord session")
	}
}
