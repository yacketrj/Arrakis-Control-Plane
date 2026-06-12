package main

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
)

const maxJSONBodyBytes int64 = 1 << 20 // 1 MiB
const adminTokenRawBytes = 32
const adminTokenEncodedLength = 43

var k8sNamePattern = regexp.MustCompile(`^[a-z0-9]([-a-z0-9]*[a-z0-9])?$`)
var listenPortPattern = regexp.MustCompile(`^[0-9]{1,5}$`)
var sqlDangerPattern = regexp.MustCompile(`(?i)\b(insert|update|delete|drop|alter|truncate|create|grant|revoke|copy|call|do|execute|merge|vacuum|analyze)\b`)
var adminTokenPattern = regexp.MustCompile(`^[A-Za-z0-9_-]{43}$`)
var adminToken string
var allowedOrigins = "http://127.0.0.1:5173,http://localhost:5173,http://127.0.0.1:4173,http://localhost:4173"

func init() {
	if v := os.Getenv("ADMIN_TOKEN"); v != "" {
		adminToken = v
	}
	if v := os.Getenv("ALLOWED_ORIGINS"); v != "" {
		allowedOrigins = v
	}
}

func generateAdminToken() string {
	b := make([]byte, adminTokenRawBytes)
	if _, err := rand.Read(b); err != nil {
		panic(fmt.Errorf("generate admin token: %w", err))
	}
	return base64.RawURLEncoding.EncodeToString(b)
}

func validateStrictAdminToken(value string) error {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return fmt.Errorf("ADMIN_TOKEN is required")
	}
	if trimmed != value || containsUnsafeControl(value) {
		return fmt.Errorf("ADMIN_TOKEN contains unsupported whitespace or control characters")
	}
	if !adminTokenPattern.MatchString(value) {
		return fmt.Errorf("ADMIN_TOKEN must be exactly %d base64url characters generated from %d random bytes", adminTokenEncodedLength, adminTokenRawBytes)
	}
	rejected := map[string]bool{
		"changeme":          true,
		"change-me":        true,
		"password":          true,
		"admin":             true,
		"admin-token":       true,
		"replace-me":        true,
		"replace_with_token": true,
		"<your_admin_token>": true,
	}
	if rejected[strings.ToLower(value)] {
		return fmt.Errorf("ADMIN_TOKEN uses a forbidden placeholder value")
	}
	return nil
}

func isWebSocketLogStreamRequest(r *http.Request) bool {
	return r != nil && strings.EqualFold(r.Header.Get("Upgrade"), "websocket") && r.URL != nil && r.URL.Path == "/api/v1/logs/stream"
}

func isSelfServicePath(path string) bool {
	return strings.HasPrefix(path, "/api/v1/self/")
}

func isDiscordSelfSessionRoute(method, path string) bool {
	switch {
	case method == http.MethodGet && path == "/api/v1/auth/discord/me":
		return true
	case method == http.MethodPost && path == "/api/v1/auth/discord/logout":
		return true
	default:
		return false
	}
}

func isSelfServiceRoute(method, path string) bool {
	return isSelfServicePath(path) || isDiscordSelfSessionRoute(method, path)
}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions || isPublicPath(r.URL.Path) {
			next.ServeHTTP(w, r)
			return
		}
		if isWebSocketLogStreamRequest(r) {
			if !validateAndConsumeLogStreamTicket(r) {
				jsonErr(w, fmt.Errorf("unauthorized"), http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
			return
		}
		if adminToken == "" {
			adminToken = os.Getenv("ADMIN_TOKEN")
		}
		if err := validateStrictAdminToken(adminToken); err != nil {
			log.Printf("security: rejecting request because backend admin token is invalid: %v", err)
			jsonErr(w, fmt.Errorf("backend admin token is not securely configured"), http.StatusServiceUnavailable)
			return
		}
		provided := bearerToken(r.Header.Get("Authorization"))
		if provided == "" {
			provided = r.Header.Get("X-Admin-Token")
		}
		if validateStrictAdminToken(provided) == nil && subtle.ConstantTimeCompare([]byte(provided), []byte(adminToken)) == 1 {
			next.ServeHTTP(w, r)
			return
		}
		if discordAuthEnabled() {
			if discordSessionIsAdmin(r) {
				next.ServeHTTP(w, r)
				return
			}
			if isSelfServiceRoute(r.Method, r.URL.Path) && discordSessionIsRegistered(r) {
				next.ServeHTTP(w, r)
				return
			}
		}
		jsonErr(w, fmt.Errorf("unauthorized"), http.StatusUnauthorized)
	})
}

func isPublicPath(path string) bool {
	switch path {
	case "/api/v1/public/status", "/api/v1/auth/discord/login", "/api/v1/auth/discord/callback":
		return true
	default:
		return false
	}
}

func bearerToken(header string) string {
	parts := strings.Fields(header)
	if len(parts) == 2 && strings.EqualFold(parts[0], "Bearer") {
		return parts[1]
	}
	return ""
}

func parseAllowedOrigins(raw string) map[string]bool {
	allowed := map[string]bool{}
	for _, origin := range strings.Split(raw, ",") {
		origin = strings.TrimSpace(origin)
		if isAllowedOriginValue(origin) {
			allowed[origin] = true
		}
	}
	return allowed
}

func isAllowedOriginValue(origin string) bool {
	if origin == "" || origin == "*" || strings.EqualFold(origin, "null") || containsUnsafeControl(origin) {
		return false
	}
	parsed, err := url.Parse(origin)
	if err != nil {
		return false
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return false
	}
	if parsed.Host == "" || parsed.User != nil || parsed.RawQuery != "" || parsed.Fragment != "" {
		return false
	}
	if allowedOriginHostHasWildcard(parsed) {
		return false
	}
	path := strings.Trim(parsed.EscapedPath(), "/")
	return path == ""
}

func allowedOriginHostHasWildcard(parsed *url.URL) bool {
	if parsed == nil {
		return false
	}
	return strings.Contains(parsed.Host, "*") || strings.Contains(parsed.Hostname(), "*")
}

func originAllowed(origin string) bool {
	if origin == "" {
		return true
	}
	if !isAllowedOriginValue(origin) {
		return false
	}
	if v := os.Getenv("ALLOWED_ORIGINS"); v != "" {
		allowedOrigins = v
	}
	return parseAllowedOrigins(allowedOrigins)[origin]
}

func normalizeListenAddr(addr string) string {
	addr = strings.TrimSpace(addr)
	if addr == "" {
		return "127.0.0.1:8080"
	}
	if listenPortPattern.MatchString(addr) {
		return "127.0.0.1:" + addr
	}
	if strings.HasPrefix(addr, ":") {
		return "127.0.0.1" + addr
	}
	host, port, err := net.SplitHostPort(addr)
	if err != nil || port == "" {
		return addr
	}
	if host == "" {
		return "127.0.0.1:" + port
	}
	return addr
}

func isValidK8sName(v string) bool { return k8sNamePattern.MatchString(v) }

func isReadOnlySQL(sql string) bool {
	trimmed := strings.TrimSpace(sql)
	if trimmed == "" {
		return false
	}
	lower := strings.ToLower(trimmed)
	if strings.HasPrefix(lower, "with") || strings.HasPrefix(lower, "select") || strings.HasPrefix(lower, "show") || strings.HasPrefix(lower, "explain") {
		return !sqlDangerPattern.MatchString(lower)
	}
	return false
}
