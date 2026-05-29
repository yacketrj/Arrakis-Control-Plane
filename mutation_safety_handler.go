package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func handleMutationSafetyClassify(w http.ResponseWriter, r *http.Request) {
	method := strings.ToUpper(strings.TrimSpace(r.URL.Query().Get("method")))
	path := strings.TrimSpace(r.URL.Query().Get("path"))
	if method == "" {
		method = http.MethodPost
	}
	if path == "" {
		jsonErr(w, fmt.Errorf("path is required"), http.StatusBadRequest)
		return
	}
	jsonOK(w, mutationSafetyForPath(method, path))
}

func mutationSafetyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !adminReasonEnforcementEnabled() || !isAuditableRequest(r) {
			next.ServeHTTP(w, r)
			return
		}
		safety := mutationSafetyForPath(r.Method, r.URL.Path)
		if !safety.RequiresReason {
			next.ServeHTTP(w, r)
			return
		}
		reason, err := mutationReasonFromRequest(r)
		if err != nil {
			jsonErr(w, err, http.StatusBadRequest)
			return
		}
		if strings.TrimSpace(reason) == "" {
			jsonErr(w, fmt.Errorf("admin reason required for %s mutation", safety.Risk), http.StatusBadRequest)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func adminReasonEnforcementEnabled() bool {
	switch strings.ToLower(strings.TrimSpace(os.Getenv("ADMIN_REQUIRE_REASON"))) {
	case "0", "false", "no", "off", "disabled", "disable":
		return false
	default:
		return true
	}
}

func mutationReasonFromRequest(r *http.Request) (string, error) {
	if reason := strings.TrimSpace(r.Header.Get("X-Admin-Reason")); reason != "" {
		return reason, nil
	}
	if r.Body == nil || r.ContentLength == 0 {
		return "", nil
	}
	if r.ContentLength > maxAuditInspectableBodyBytes {
		return "", fmt.Errorf("request body too large to inspect for admin reason")
	}
	body, err := io.ReadAll(io.LimitReader(r.Body, maxAuditInspectableBodyBytes+1))
	if err != nil {
		r.Body = io.NopCloser(bytes.NewReader(body))
		return "", err
	}
	r.Body = io.NopCloser(bytes.NewReader(body))
	if int64(len(body)) > maxAuditInspectableBodyBytes || len(bytes.TrimSpace(body)) == 0 {
		return "", nil
	}
	var payload map[string]any
	if err := json.Unmarshal(body, &payload); err != nil {
		return "", nil
	}
	if raw, ok := payload["reason"].(string); ok {
		return strings.TrimSpace(raw), nil
	}
	return "", nil
}
