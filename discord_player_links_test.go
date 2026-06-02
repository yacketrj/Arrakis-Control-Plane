package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
	"time"
)

func useTempDiscordPlayerLinkStore(t *testing.T) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "discord-player-links.json")
	t.Setenv("DISCORD_PLAYER_LINK_STORE", path)
	return path
}

func TestNormalizeDiscordPlayerLinkPayload(t *testing.T) {
	link, err := normalizeDiscordPlayerLinkPayload(discordPlayerLinkPayload{
		DiscordID:  "discord-1",
		PlayerID:   123,
		PlayerName: "Sihaya",
		Notes:      "verified by admin",
	}, nil, "admin-discord", "discord", time.Date(2026, 6, 1, 12, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatalf("normalize link: %v", err)
	}
	if link.DiscordID != "discord-1" || link.PlayerID != 123 || link.PlayerName != "Sihaya" {
		t.Fatalf("unexpected link: %#v", link)
	}
	if link.LinkedByDiscordID != "admin-discord" || link.LinkedByAuthType != "discord" {
		t.Fatalf("unexpected attribution: %#v", link)
	}
	if link.CreatedAt == "" || link.UpdatedAt == "" {
		t.Fatalf("expected timestamps: %#v", link)
	}
}

func TestNormalizeDiscordPlayerLinkPayloadRejectsInvalidValues(t *testing.T) {
	tests := []discordPlayerLinkPayload{
		{DiscordID: "", PlayerID: 123},
		{DiscordID: "discord-1", PlayerID: 0},
		{DiscordID: "discord-1\n", PlayerID: 123},
	}
	for _, payload := range tests {
		if _, err := normalizeDiscordPlayerLinkPayload(payload, nil, "", "admin-token", time.Now()); err == nil {
			t.Fatalf("expected invalid payload to fail: %#v", payload)
		}
	}
}

func TestDiscordPlayerLinkStoreHelpers(t *testing.T) {
	store := discordPlayerLinkStore{}
	first := discordPlayerLink{DiscordID: "discord-1", PlayerID: 100, CreatedAt: "old", UpdatedAt: "old"}
	upsertDiscordPlayerLink(&store, first)
	if len(store.Links) != 1 {
		t.Fatalf("expected one link: %#v", store)
	}
	updated := discordPlayerLink{DiscordID: "discord-1", PlayerID: 200, CreatedAt: "old", UpdatedAt: "new"}
	upsertDiscordPlayerLink(&store, updated)
	link, found := findDiscordPlayerLink(store, "discord-1")
	if !found || link.PlayerID != 200 || len(store.Links) != 1 {
		t.Fatalf("expected updated link: found=%v link=%#v store=%#v", found, link, store)
	}
	if !deleteDiscordPlayerLink(&store, "discord-1") || len(store.Links) != 0 {
		t.Fatalf("expected link deletion: %#v", store)
	}
	if deleteDiscordPlayerLink(&store, "discord-1") {
		t.Fatal("expected second delete to return false")
	}
}

func TestDiscordPlayerLinkHandlersCreateListDelete(t *testing.T) {
	useTempDiscordPlayerLinkStore(t)
	mux := http.NewServeMux()
	registerRoutes(mux)

	createRecorder := httptest.NewRecorder()
	mux.ServeHTTP(createRecorder, jsonBodyRequest(http.MethodPost, "/api/v1/auth/discord/player-links", discordPlayerLinkPayload{
		DiscordID:  "discord-1",
		PlayerID:   123,
		PlayerName: "Sihaya",
		Notes:      "verified",
	}))
	if createRecorder.Code != http.StatusOK {
		t.Fatalf("create link got %d: %s", createRecorder.Code, createRecorder.Body.String())
	}
	created := decodeJSONBody[discordPlayerLink](t, createRecorder)
	if created.DiscordID != "discord-1" || created.PlayerID != 123 {
		t.Fatalf("unexpected created link: %#v", created)
	}

	listRecorder := httptest.NewRecorder()
	mux.ServeHTTP(listRecorder, httptest.NewRequest(http.MethodGet, "/api/v1/auth/discord/player-links", nil))
	if listRecorder.Code != http.StatusOK {
		t.Fatalf("list links got %d: %s", listRecorder.Code, listRecorder.Body.String())
	}
	links := decodeJSONBody[[]discordPlayerLink](t, listRecorder)
	if len(links) != 1 || links[0].DiscordID != "discord-1" {
		t.Fatalf("unexpected links: %#v", links)
	}

	deleteRecorder := httptest.NewRecorder()
	mux.ServeHTTP(deleteRecorder, httptest.NewRequest(http.MethodDelete, "/api/v1/auth/discord/player-links/discord-1", nil))
	if deleteRecorder.Code != http.StatusOK {
		t.Fatalf("delete link got %d: %s", deleteRecorder.Code, deleteRecorder.Body.String())
	}

	store, err := loadDiscordPlayerLinkStore()
	if err != nil {
		t.Fatalf("load link store: %v", err)
	}
	if len(store.Links) != 0 {
		t.Fatalf("expected empty link store: %#v", store)
	}
}

func TestCurrentDiscordPlayerLink(t *testing.T) {
	useTempDiscordPlayerLinkStore(t)
	resetDiscordSessionsForTest(t)

	store := discordPlayerLinkStore{Links: []discordPlayerLink{{DiscordID: "discord-1", PlayerID: 123}}}
	if err := saveDiscordPlayerLinkStore(store); err != nil {
		t.Fatalf("save link store: %v", err)
	}
	sessionID := "normal-session"
	discordSessionsMu.Lock()
	discordSessions[sessionID] = discordSession{ID: sessionID, DiscordID: "discord-1", Role: appRoleNormal, ExpiresAt: time.Now().Add(time.Hour)}
	discordSessionsMu.Unlock()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/self/player-link", nil)
	req.AddCookie(&http.Cookie{Name: discordSessionCookieName, Value: sessionID})
	link, found, err := currentDiscordPlayerLink(req)
	if err != nil {
		t.Fatalf("current link: %v", err)
	}
	if !found || link.PlayerID != 123 {
		t.Fatalf("unexpected current link: found=%v link=%#v", found, link)
	}
}

func TestAuthMiddlewareAllowsNormalDiscordOnlyForSelfEndpoints(t *testing.T) {
	resetDiscordSessionsForTest(t)
	t.Setenv("DISCORD_AUTH_ENABLED", "1")
	old := adminToken
	adminToken = testStrictAdminToken
	t.Cleanup(func() { adminToken = old })

	sessionID := "normal-session"
	discordSessionsMu.Lock()
	discordSessions[sessionID] = discordSession{ID: sessionID, DiscordID: "discord-1", Role: appRoleNormal, ExpiresAt: time.Now().Add(time.Hour)}
	discordSessionsMu.Unlock()

	h := authMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))

	selfReq := httptest.NewRequest(http.MethodGet, "/api/v1/self/player-card", nil)
	selfReq.AddCookie(&http.Cookie{Name: discordSessionCookieName, Value: sessionID})
	selfRecorder := httptest.NewRecorder()
	h.ServeHTTP(selfRecorder, selfReq)
	if selfRecorder.Code != http.StatusNoContent {
		t.Fatalf("normal Discord self request got %d", selfRecorder.Code)
	}

	adminReq := httptest.NewRequest(http.MethodGet, "/api/v1/status", nil)
	adminReq.AddCookie(&http.Cookie{Name: discordSessionCookieName, Value: sessionID})
	adminRecorder := httptest.NewRecorder()
	h.ServeHTTP(adminRecorder, adminReq)
	if adminRecorder.Code != http.StatusUnauthorized {
		t.Fatalf("normal Discord admin request got %d", adminRecorder.Code)
	}
}

func jsonBodyRequest(method string, path string, value any) *http.Request {
	body, _ := json.Marshal(value)
	req := httptest.NewRequest(method, path, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	return req
}

func decodeJSONBody[T any](t *testing.T, recorder *httptest.ResponseRecorder) T {
	t.Helper()
	var value T
	if err := json.Unmarshal(recorder.Body.Bytes(), &value); err != nil {
		t.Fatalf("decode response: %v; body=%s", err, recorder.Body.String())
	}
	return value
}
