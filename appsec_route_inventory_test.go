package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"net/http"
	"net/http/httptest"
	"sort"
	"strconv"
	"strings"
	"testing"
	"time"
)

type appsecRouteAuthClass string

const (
	appsecRoutePublic          appsecRouteAuthClass = "public"
	appsecRouteSelfService     appsecRouteAuthClass = "self-service"
	appsecRouteAdmin           appsecRouteAuthClass = "admin"
	appsecRouteWebSocketTicket appsecRouteAuthClass = "websocket-ticket"
)

var appsecExpectedRouteAuth = map[string]appsecRouteAuthClass{
	"GET /api/v1/public/status":                       appsecRoutePublic,
	"GET /api/v1/auth/discord/login":                  appsecRoutePublic,
	"GET /api/v1/auth/discord/callback":               appsecRoutePublic,
	"GET /api/v1/auth/discord/me":                     appsecRouteSelfService,
	"POST /api/v1/auth/discord/logout":                appsecRouteSelfService,
	"GET /api/v1/auth/discord/users":                  appsecRouteAdmin,
	"GET /api/v1/auth/discord/player-links":           appsecRouteAdmin,
	"POST /api/v1/auth/discord/player-links":          appsecRouteAdmin,
	"DELETE /api/v1/auth/discord/player-links/{discord_id}": appsecRouteAdmin,
	"GET /api/v1/self/player-link":                    appsecRouteSelfService,
	"GET /api/v1/self/player-card":                    appsecRouteSelfService,
	"GET /api/v1/status":                              appsecRouteAdmin,
	"POST /api/v1/reconnect":                          appsecRouteAdmin,
	"GET /api/v1/connectivity/diagnostics":            appsecRouteAdmin,
	"GET /api/v1/diagnostics/export":                  appsecRouteAdmin,
	"GET /api/v1/audit/events":                        appsecRouteAdmin,
	"GET /api/v1/mutation-safety/classify":            appsecRouteAdmin,
	"GET /api/v1/battlegroup/status":                  appsecRouteAdmin,
	"GET /api/v1/battlegroup/health":                  appsecRouteAdmin,
	"POST /api/v1/battlegroup/exec":                   appsecRouteAdmin,
	"GET /api/v1/battlegroup/pods":                    appsecRouteAdmin,
	"GET /api/v1/players":                             appsecRouteAdmin,
	"GET /api/v1/players/online":                      appsecRouteAdmin,
	"GET /api/v1/players/currency":                    appsecRouteAdmin,
	"GET /api/v1/players/factions":                    appsecRouteAdmin,
	"GET /api/v1/players/specs":                       appsecRouteAdmin,
	"GET /api/v1/players/templates":                   appsecRouteAdmin,
	"POST /api/v1/players/templates/refresh":          appsecRouteAdmin,
	"GET /api/v1/players/{id}/profile":                appsecRouteAdmin,
	"GET /api/v1/players/{id}/inventory":              appsecRouteAdmin,
	"GET /api/v1/players/{id}/journey":                appsecRouteAdmin,
	"POST /api/v1/players/give-item":                  appsecRouteAdmin,
	"POST /api/v1/players/give-currency":              appsecRouteAdmin,
	"POST /api/v1/players/grant-live":                 appsecRouteAdmin,
	"POST /api/v1/players/give-faction-rep":           appsecRouteAdmin,
	"POST /api/v1/players/give-scrip":                 appsecRouteAdmin,
	"POST /api/v1/players/award-xp":                   appsecRouteAdmin,
	"POST /api/v1/players/award-char-xp":              appsecRouteAdmin,
	"POST /api/v1/players/award-intel":                appsecRouteAdmin,
	"POST /api/v1/players/kick":                       appsecRouteAdmin,
	"DELETE /api/v1/players/item/{id}":                appsecRouteAdmin,
	"POST /api/v1/players/item/stack-size":            appsecRouteAdmin,
	"POST /api/v1/players/reset-spec":                 appsecRouteAdmin,
	"POST /api/v1/players/set-faction-tier":           appsecRouteAdmin,
	"POST /api/v1/players/journey/complete":           appsecRouteAdmin,
	"POST /api/v1/players/journey/reset":              appsecRouteAdmin,
	"POST /api/v1/players/journey/wipe":               appsecRouteAdmin,
	"POST /api/v1/players/delete-tutorials":           appsecRouteAdmin,
	"POST /api/v1/players/wipe-codex":                 appsecRouteAdmin,
	"GET /api/v1/players/{id}/char-xp":                appsecRouteAdmin,
	"GET /api/v1/players/{id}/specs":                  appsecRouteAdmin,
	"POST /api/v1/players/set-spec-xp":                appsecRouteAdmin,
	"GET /api/v1/players/{id}/vehicles":               appsecRouteAdmin,
	"POST /api/v1/players/repair-item":                appsecRouteAdmin,
	"GET /api/v1/players/partitions":                  appsecRouteAdmin,
	"POST /api/v1/players/teleport":                   appsecRouteAdmin,
	"GET /api/v1/players/{id}/events":                 appsecRouteAdmin,
	"GET /api/v1/players/{id}/dungeons":               appsecRouteAdmin,
	"GET /api/v1/inventory/requests":                  appsecRouteAdmin,
	"POST /api/v1/inventory/requests":                 appsecRouteAdmin,
	"PATCH /api/v1/inventory/requests/{id}":           appsecRouteAdmin,
	"GET /api/v1/inventory/orders":                    appsecRouteAdmin,
	"POST /api/v1/inventory/orders":                   appsecRouteAdmin,
	"PATCH /api/v1/inventory/orders/{id}":             appsecRouteAdmin,
	"GET /api/v1/database/tables":                     appsecRouteAdmin,
	"GET /api/v1/database/describe":                   appsecRouteAdmin,
	"GET /api/v1/database/sample":                     appsecRouteAdmin,
	"GET /api/v1/database/search":                     appsecRouteAdmin,
	"GET /api/v1/database/functions":                  appsecRouteAdmin,
	"GET /api/v1/database/functions/inspect":          appsecRouteAdmin,
	"POST /api/v1/database/sql":                       appsecRouteAdmin,
	"GET /api/v1/logs/pods":                           appsecRouteAdmin,
	"POST /api/v1/logs/stream-ticket":                 appsecRouteAdmin,
	"GET /api/v1/logs/stream":                         appsecRouteWebSocketTicket,
	"GET /api/v1/logs/cheats":                         appsecRouteAdmin,
	"POST /api/v1/notify":                             appsecRouteAdmin,
	"GET /api/v1/storage":                             appsecRouteAdmin,
	"GET /api/v1/storage/{id}/items":                  appsecRouteAdmin,
	"POST /api/v1/storage/{id}/give-item":             appsecRouteAdmin,
	"GET /api/v1/blueprints":                          appsecRouteAdmin,
	"GET /api/v1/blueprints/{id}/export":              appsecRouteAdmin,
	"POST /api/v1/blueprints/import":                  appsecRouteAdmin,
}

func TestAppSecRegisteredRoutesHaveAuthBoundaryExpectations(t *testing.T) {
	registered := appsecRegisteredRoutePatterns(t)
	for _, route := range registered {
		if _, ok := appsecExpectedRouteAuth[route]; !ok {
			t.Fatalf("registered route missing AppSec auth-boundary expectation: %s", route)
		}
	}
	for route := range appsecExpectedRouteAuth {
		if !containsString(registered, route) {
			t.Fatalf("AppSec auth-boundary expectation references unregistered route: %s", route)
		}
	}
}

func TestAppSecRegisteredRouteAuthBoundaryEnforcement(t *testing.T) {
	resetDiscordSessionsForTest(t)
	t.Setenv("DISCORD_AUTH_ENABLED", "1")
	old := adminToken
	adminToken = testStrictAdminToken
	t.Cleanup(func() { adminToken = old })

	sessionID := "generated-route-boundary-normal-session"
	discordSessionsMu.Lock()
	discordSessions[sessionID] = discordSession{ID: sessionID, DiscordID: "discord-normal", Role: appRoleNormal, ExpiresAt: time.Now().Add(time.Hour)}
	discordSessionsMu.Unlock()

	h := authMiddleware(appsecNoopHandler())
	routes := appsecRegisteredRoutePatterns(t)
	for _, pattern := range routes {
		method, concretePath := appsecConcreteRoute(t, pattern)
		class := appsecExpectedRouteAuth[pattern]

		switch class {
		case appsecRoutePublic:
			assertRouteStatus(t, h, method, concretePath, nil, http.StatusNoContent, pattern+" public")
			assertPublicPathClassification(t, pattern, concretePath)
		case appsecRouteSelfService:
			assertRouteStatus(t, h, method, concretePath, nil, http.StatusUnauthorized, pattern+" missing self-service auth")
			assertRouteStatus(t, h, method, concretePath, map[string]string{"X-Admin-Token": testStrictAdminToken}, http.StatusNoContent, pattern+" admin self-service auth")
			assertRouteWithDiscordSessionStatus(t, h, method, concretePath, sessionID, http.StatusNoContent, pattern+" normal Discord self-service auth")
		case appsecRouteAdmin:
			assertRouteStatus(t, h, method, concretePath, nil, http.StatusUnauthorized, pattern+" missing admin auth")
			assertRouteStatus(t, h, method, concretePath, map[string]string{"X-Admin-Token": testStrictAdminToken}, http.StatusNoContent, pattern+" admin token auth")
			assertRouteWithDiscordSessionStatus(t, h, method, concretePath, sessionID, http.StatusUnauthorized, pattern+" normal Discord blocked from admin route")
		case appsecRouteWebSocketTicket:
			assertRouteStatus(t, h, method, concretePath, nil, http.StatusUnauthorized, pattern+" missing normal auth")
			assertRouteStatus(t, h, method, concretePath, map[string]string{"X-Admin-Token": testStrictAdminToken}, http.StatusNoContent, pattern+" normal admin auth")
			assertWebSocketRouteStatus(t, h, concretePath, map[string]string{"X-Admin-Token": testStrictAdminToken}, http.StatusUnauthorized, pattern+" websocket requires ticket before admin fallback")
		default:
			t.Fatalf("%s: unsupported route auth class %q", pattern, class)
		}
	}
}

func appsecRegisteredRoutePatterns(t *testing.T) []string {
	t.Helper()
	fileSet := token.NewFileSet()
	parsed, err := parser.ParseFile(fileSet, "routes.go", nil, 0)
	if err != nil {
		t.Fatalf("parse routes.go: %v", err)
	}
	var routes []string
	ast.Inspect(parsed, func(node ast.Node) bool {
		call, ok := node.(*ast.CallExpr)
		if !ok || len(call.Args) == 0 {
			return true
		}
		selector, ok := call.Fun.(*ast.SelectorExpr)
		if !ok || selector.Sel.Name != "HandleFunc" {
			return true
		}
		literal, ok := call.Args[0].(*ast.BasicLit)
		if !ok || literal.Kind != token.STRING {
			return true
		}
		route, err := strconv.Unquote(literal.Value)
		if err != nil {
			t.Fatalf("unquote route literal %s: %v", literal.Value, err)
		}
		routes = append(routes, route)
		return true
	})
	if len(routes) == 0 {
		t.Fatal("expected at least one route in routes.go")
	}
	sort.Strings(routes)
	return routes
}

func appsecConcreteRoute(t *testing.T, pattern string) (string, string) {
	t.Helper()
	parts := strings.SplitN(pattern, " ", 2)
	if len(parts) != 2 {
		t.Fatalf("route pattern must be METHOD /path: %q", pattern)
	}
	path := parts[1]
	path = strings.ReplaceAll(path, "{id}", "123")
	path = strings.ReplaceAll(path, "{discord_id}", "discord-123")
	if strings.Contains(path, "{") || strings.Contains(path, "}") {
		t.Fatalf("route pattern has unhandled path placeholder: %q", pattern)
	}
	return parts[0], path
}

func assertRouteStatus(t *testing.T, h http.Handler, method, path string, headers map[string]string, want int, label string) {
	t.Helper()
	req := httptest.NewRequest(method, path, nil)
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	resp := httptest.NewRecorder()
	h.ServeHTTP(resp, req)
	if resp.Code != want {
		t.Fatalf("%s: got status %d, want %d", label, resp.Code, want)
	}
}

func assertRouteWithDiscordSessionStatus(t *testing.T, h http.Handler, method, path, sessionID string, want int, label string) {
	t.Helper()
	req := httptest.NewRequest(method, path, nil)
	req.AddCookie(&http.Cookie{Name: discordSessionCookieName, Value: sessionID})
	resp := httptest.NewRecorder()
	h.ServeHTTP(resp, req)
	if resp.Code != want {
		t.Fatalf("%s: got status %d, want %d", label, resp.Code, want)
	}
}

func assertWebSocketRouteStatus(t *testing.T, h http.Handler, path string, headers map[string]string, want int, label string) {
	t.Helper()
	req := httptest.NewRequest(http.MethodGet, path, nil)
	req.Header.Set("Upgrade", "websocket")
	req.Header.Set("Connection", "Upgrade")
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	resp := httptest.NewRecorder()
	h.ServeHTTP(resp, req)
	if resp.Code != want {
		t.Fatalf("%s: got status %d, want %d", label, resp.Code, want)
	}
}

func assertPublicPathClassification(t *testing.T, pattern, concretePath string) {
	t.Helper()
	_, path := appsecConcreteRoute(t, pattern)
	if path != concretePath {
		t.Fatalf("internal route conversion mismatch for %s: %s != %s", pattern, path, concretePath)
	}
	if !isPublicPath(concretePath) {
		t.Fatalf("%s: expected concrete path %s to be public", pattern, concretePath)
	}
}

func containsString(values []string, needle string) bool {
	for _, value := range values {
		if value == needle {
			return true
		}
	}
	return false
}
