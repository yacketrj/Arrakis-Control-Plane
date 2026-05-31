package main

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strings"
	"sync"
	"time"
)

const discordAPIBase = "https://discord.com/api/v10"
const discordSessionCookieName = "dune_admin_session"
const discordStateCookieName = "dune_admin_oauth_state"
const discordSessionTTL = 12 * time.Hour
const discordOAuthStateTTL = 10 * time.Minute

type appRole string

const (
	appRoleNone   appRole = "none"
	appRoleNormal appRole = "normal"
	appRoleAdmin  appRole = "admin"
)

type discordUser struct {
	ID            string `json:"id"`
	Username      string `json:"username"`
	GlobalName    string `json:"global_name"`
	Discriminator string `json:"discriminator"`
	Avatar        string `json:"avatar"`
}

type discordMember struct {
	User  *discordUser `json:"user,omitempty"`
	Nick  string       `json:"nick,omitempty"`
	Roles []string     `json:"roles"`
}

type discordTokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
}

type registeredDiscordUser struct {
	DiscordID      string   `json:"discord_id"`
	Username       string   `json:"username,omitempty"`
	GlobalName     string   `json:"global_name,omitempty"`
	Role           appRole  `json:"role"`
	DiscordRoleIDs []string `json:"discord_role_ids,omitempty"`
	RegisteredAt   string   `json:"registered_at"`
	LastLoginAt    string   `json:"last_login_at"`
}

type discordUserStore struct {
	Users []registeredDiscordUser `json:"users"`
}

type discordSession struct {
	ID        string
	DiscordID string
	Role      appRole
	ExpiresAt time.Time
}

type discordAuthContext struct {
	DiscordID string  `json:"discord_id,omitempty"`
	Username  string  `json:"username,omitempty"`
	Role      appRole `json:"role"`
	AuthType  string  `json:"auth_type"`
}

var (
	discordSessionsMu sync.Mutex
	discordSessions   = map[string]discordSession{}
)

func discordAuthEnabled() bool {
	return truthyEnv("DISCORD_AUTH_ENABLED")
}

func truthyEnv(key string) bool {
	switch strings.ToLower(strings.TrimSpace(os.Getenv(key))) {
	case "1", "true", "yes", "on", "enabled", "enable":
		return true
	default:
		return false
	}
}

func discordConfigValue(key string) string {
	return strings.TrimSpace(os.Getenv(key))
}

func discordConfigured() error {
	if !discordAuthEnabled() {
		return fmt.Errorf("Discord auth is disabled")
	}
	for _, key := range []string{"DISCORD_CLIENT_ID", "DISCORD_CLIENT_SECRET", "DISCORD_REDIRECT_URI", "DISCORD_GUILD_ID"} {
		if discordConfigValue(key) == "" || containsUnsafeControl(discordConfigValue(key)) {
			return fmt.Errorf("%s is required for Discord auth", key)
		}
	}
	return nil
}

func discordUserStorePath() string {
	if path := strings.TrimSpace(os.Getenv("DISCORD_USER_STORE")); path != "" {
		return path
	}
	return "discord-users.json"
}

func randomURLToken(byteLen int) (string, error) {
	b := make([]byte, byteLen)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func discordScopes() string {
	// guilds.members.read returns the caller's member object for the configured guild,
	// including Discord role IDs. identify provides the stable Discord user ID.
	return "identify guilds.members.read"
}

func handleDiscordLogin(w http.ResponseWriter, r *http.Request) {
	if err := discordConfigured(); err != nil {
		jsonErr(w, err, http.StatusServiceUnavailable)
		return
	}
	state, err := randomURLToken(32)
	if err != nil {
		jsonErr(w, fmt.Errorf("create Discord OAuth state: %w", err), http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     discordStateCookieName,
		Value:    state,
		Path:     "/",
		MaxAge:   int(discordOAuthStateTTL.Seconds()),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   discordCookieSecure(),
	})
	authURL, err := discordAuthorizeURL(state)
	if err != nil {
		jsonErr(w, err, http.StatusInternalServerError)
		return
	}
	if r.URL.Query().Get("json") == "1" {
		jsonOK(w, map[string]string{"url": authURL})
		return
	}
	http.Redirect(w, r, authURL, http.StatusFound)
}

func discordAuthorizeURL(state string) (string, error) {
	redirectURI := discordConfigValue("DISCORD_REDIRECT_URI")
	if _, err := url.ParseRequestURI(redirectURI); err != nil {
		return "", fmt.Errorf("DISCORD_REDIRECT_URI is invalid")
	}
	values := url.Values{}
	values.Set("response_type", "code")
	values.Set("client_id", discordConfigValue("DISCORD_CLIENT_ID"))
	values.Set("redirect_uri", redirectURI)
	values.Set("scope", discordScopes())
	values.Set("state", state)
	values.Set("prompt", "consent")
	return "https://discord.com/oauth2/authorize?" + values.Encode(), nil
}

func handleDiscordCallback(w http.ResponseWriter, r *http.Request) {
	if err := discordConfigured(); err != nil {
		jsonErr(w, err, http.StatusServiceUnavailable)
		return
	}
	code := strings.TrimSpace(r.URL.Query().Get("code"))
	state := strings.TrimSpace(r.URL.Query().Get("state"))
	if code == "" || state == "" {
		jsonErr(w, fmt.Errorf("Discord callback requires code and state"), http.StatusBadRequest)
		return
	}
	stateCookie, err := r.Cookie(discordStateCookieName)
	if err != nil || stateCookie.Value == "" || subtleCompareString(stateCookie.Value, state) != 1 {
		jsonErr(w, fmt.Errorf("invalid Discord OAuth state"), http.StatusBadRequest)
		return
	}
	http.SetCookie(w, &http.Cookie{Name: discordStateCookieName, Value: "", Path: "/", MaxAge: -1, HttpOnly: true, SameSite: http.SameSiteLaxMode, Secure: discordCookieSecure()})

	token, err := exchangeDiscordCode(code)
	if err != nil {
		jsonErr(w, err, http.StatusBadGateway)
		return
	}
	user, err := fetchDiscordUser(token.AccessToken)
	if err != nil {
		jsonErr(w, err, http.StatusBadGateway)
		return
	}
	member, err := fetchDiscordCurrentGuildMember(token.AccessToken)
	if err != nil {
		jsonErr(w, err, http.StatusForbidden)
		return
	}
	role := mapDiscordRoles(member.Roles)
	if role == appRoleNone {
		jsonErr(w, fmt.Errorf("Discord user is not authorized for this application"), http.StatusForbidden)
		return
	}
	if err := upsertRegisteredDiscordUser(user, member.Roles, role); err != nil {
		jsonErr(w, err, http.StatusInternalServerError)
		return
	}
	sessionID, err := createDiscordSession(user.ID, role)
	if err != nil {
		jsonErr(w, err, http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     discordSessionCookieName,
		Value:    sessionID,
		Path:     "/",
		MaxAge:   int(discordSessionTTL.Seconds()),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   discordCookieSecure(),
	})

	redirect := strings.TrimSpace(os.Getenv("DISCORD_POST_LOGIN_REDIRECT"))
	if redirect == "" {
		redirect = "/"
	}
	http.Redirect(w, r, redirect, http.StatusFound)
}

func subtleCompareString(a, b string) int {
	if len(a) != len(b) {
		return 0
	}
	ah := sha256.Sum256([]byte(a))
	bh := sha256.Sum256([]byte(b))
	if ah == bh {
		return 1
	}
	return 0
}

func exchangeDiscordCode(code string) (discordTokenResponse, error) {
	values := url.Values{}
	values.Set("grant_type", "authorization_code")
	values.Set("code", code)
	values.Set("redirect_uri", discordConfigValue("DISCORD_REDIRECT_URI"))
	request, err := http.NewRequest(http.MethodPost, discordAPIBase+"/oauth2/token", strings.NewReader(values.Encode()))
	if err != nil {
		return discordTokenResponse{}, err
	}
	request.SetBasicAuth(discordConfigValue("DISCORD_CLIENT_ID"), discordConfigValue("DISCORD_CLIENT_SECRET"))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	var token discordTokenResponse
	if err := doDiscordJSON(request, &token); err != nil {
		return discordTokenResponse{}, fmt.Errorf("Discord token exchange failed: %w", err)
	}
	if token.AccessToken == "" {
		return discordTokenResponse{}, fmt.Errorf("Discord token exchange returned no access token")
	}
	return token, nil
}

func fetchDiscordUser(accessToken string) (discordUser, error) {
	request, err := http.NewRequest(http.MethodGet, discordAPIBase+"/users/@me", nil)
	if err != nil {
		return discordUser{}, err
	}
	request.Header.Set("Authorization", "Bearer "+accessToken)
	var user discordUser
	if err := doDiscordJSON(request, &user); err != nil {
		return discordUser{}, fmt.Errorf("Discord user lookup failed: %w", err)
	}
	if user.ID == "" {
		return discordUser{}, fmt.Errorf("Discord user lookup returned no id")
	}
	return user, nil
}

func fetchDiscordCurrentGuildMember(accessToken string) (discordMember, error) {
	guildID := discordConfigValue("DISCORD_GUILD_ID")
	request, err := http.NewRequest(http.MethodGet, discordAPIBase+"/users/@me/guilds/"+url.PathEscape(guildID)+"/member", nil)
	if err != nil {
		return discordMember{}, err
	}
	request.Header.Set("Authorization", "Bearer "+accessToken)
	var member discordMember
	if err := doDiscordJSON(request, &member); err != nil {
		return discordMember{}, fmt.Errorf("Discord guild membership lookup failed: %w", err)
	}
	return member, nil
}

func doDiscordJSON(request *http.Request, target any) error {
	client := &http.Client{Timeout: 10 * time.Second}
	res, err := client.Do(request)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	body, _ := io.ReadAll(io.LimitReader(res.Body, 1<<20))
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return fmt.Errorf("Discord API returned %d: %s", res.StatusCode, RedactPIIText(string(body)))
	}
	if target == nil {
		return nil
	}
	if err := json.Unmarshal(body, target); err != nil {
		return fmt.Errorf("decode Discord API response: %w", err)
	}
	return nil
}

func mapDiscordRoles(roleIDs []string) appRole {
	roleSet := map[string]bool{}
	for _, id := range roleIDs {
		id = strings.TrimSpace(id)
		if id != "" {
			roleSet[id] = true
		}
	}
	for _, id := range splitCSVEnv("DISCORD_ADMIN_ROLE_IDS") {
		if roleSet[id] {
			return appRoleAdmin
		}
	}
	normalRoles := splitCSVEnv("DISCORD_NORMAL_ROLE_IDS")
	if len(normalRoles) == 0 {
		return appRoleNormal
	}
	for _, id := range normalRoles {
		if roleSet[id] {
			return appRoleNormal
		}
	}
	return appRoleNone
}

func splitCSVEnv(key string) []string {
	var out []string
	for _, part := range strings.Split(os.Getenv(key), ",") {
		part = strings.TrimSpace(part)
		if part != "" && !containsUnsafeControl(part) {
			out = append(out, part)
		}
	}
	return out
}

func upsertRegisteredDiscordUser(user discordUser, roleIDs []string, role appRole) error {
	path := discordUserStorePath()
	store := discordUserStore{}
	if data, err := os.ReadFile(path); err == nil && len(strings.TrimSpace(string(data))) > 0 {
		_ = json.Unmarshal(data, &store)
	}
	now := time.Now().UTC().Format(time.RFC3339Nano)
	sort.Strings(roleIDs)
	updated := false
	for i := range store.Users {
		if store.Users[i].DiscordID == user.ID {
			store.Users[i].Username = user.Username
			store.Users[i].GlobalName = user.GlobalName
			store.Users[i].Role = role
			store.Users[i].DiscordRoleIDs = roleIDs
			store.Users[i].LastLoginAt = now
			updated = true
			break
		}
	}
	if !updated {
		store.Users = append(store.Users, registeredDiscordUser{
			DiscordID:      user.ID,
			Username:       user.Username,
			GlobalName:     user.GlobalName,
			Role:           role,
			DiscordRoleIDs: roleIDs,
			RegisteredAt:   now,
			LastLoginAt:    now,
		})
	}
	data, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o600)
}

func createDiscordSession(discordID string, role appRole) (string, error) {
	sessionID, err := randomURLToken(32)
	if err != nil {
		return "", err
	}
	discordSessionsMu.Lock()
	defer discordSessionsMu.Unlock()
	discordSessions[sessionID] = discordSession{ID: sessionID, DiscordID: discordID, Role: role, ExpiresAt: time.Now().Add(discordSessionTTL)}
	return sessionID, nil
}

func discordCookieSecure() bool {
	return truthyEnv("SESSION_COOKIE_SECURE") || strings.HasPrefix(strings.ToLower(os.Getenv("DISCORD_REDIRECT_URI")), "https://")
}

func discordSessionFromRequest(r *http.Request) (discordSession, bool) {
	cookie, err := r.Cookie(discordSessionCookieName)
	if err != nil || cookie.Value == "" {
		return discordSession{}, false
	}
	discordSessionsMu.Lock()
	defer discordSessionsMu.Unlock()
	session, ok := discordSessions[cookie.Value]
	if !ok || time.Now().After(session.ExpiresAt) {
		delete(discordSessions, cookie.Value)
		return discordSession{}, false
	}
	return session, true
}

func handleDiscordMe(w http.ResponseWriter, r *http.Request) {
	if provided := bearerToken(r.Header.Get("Authorization")); provided != "" || r.Header.Get("X-Admin-Token") != "" {
		jsonOK(w, discordAuthContext{Role: appRoleAdmin, AuthType: "admin-token"})
		return
	}
	session, ok := discordSessionFromRequest(r)
	if !ok {
		jsonErr(w, fmt.Errorf("unauthorized"), http.StatusUnauthorized)
		return
	}
	user := findRegisteredDiscordUser(session.DiscordID)
	jsonOK(w, discordAuthContext{DiscordID: session.DiscordID, Username: user.Username, Role: session.Role, AuthType: "discord"})
}

func handleDiscordLogout(w http.ResponseWriter, r *http.Request) {
	if cookie, err := r.Cookie(discordSessionCookieName); err == nil {
		discordSessionsMu.Lock()
		delete(discordSessions, cookie.Value)
		discordSessionsMu.Unlock()
	}
	http.SetCookie(w, &http.Cookie{Name: discordSessionCookieName, Value: "", Path: "/", MaxAge: -1, HttpOnly: true, SameSite: http.SameSiteLaxMode, Secure: discordCookieSecure()})
	jsonOK(w, map[string]string{"ok": "logged out"})
}

func findRegisteredDiscordUser(discordID string) registeredDiscordUser {
	store := discordUserStore{}
	if data, err := os.ReadFile(discordUserStorePath()); err == nil {
		_ = json.Unmarshal(data, &store)
	}
	for _, user := range store.Users {
		if user.DiscordID == discordID {
			return user
		}
	}
	return registeredDiscordUser{DiscordID: discordID, Role: appRoleNone}
}

func handleDiscordUsers(w http.ResponseWriter, r *http.Request) {
	store := discordUserStore{}
	if data, err := os.ReadFile(discordUserStorePath()); err == nil && len(strings.TrimSpace(string(data))) > 0 {
		if err := json.Unmarshal(data, &store); err != nil {
			jsonErr(w, err, http.StatusInternalServerError)
			return
		}
	}
	jsonOK(w, store.Users)
}

func discordSessionIsAdmin(r *http.Request) bool {
	session, ok := discordSessionFromRequest(r)
	return ok && session.Role == appRoleAdmin
}

func discordSessionIsRegistered(r *http.Request) bool {
	_, ok := discordSessionFromRequest(r)
	return ok
}

func discordSessionRole(r *http.Request) appRole {
	session, ok := discordSessionFromRequest(r)
	if !ok {
		return appRoleNone
	}
	return session.Role
}

func discordSessionHash(r *http.Request) string {
	cookie, err := r.Cookie(discordSessionCookieName)
	if err != nil || cookie.Value == "" {
		return ""
	}
	sum := sha256.Sum256([]byte(cookie.Value))
	return hex.EncodeToString(sum[:])[:16]
}
