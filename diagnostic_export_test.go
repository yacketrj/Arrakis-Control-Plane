package main

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestSafeDiagnosticJSONRedactsSecretsAndPII(t *testing.T) {
	input := map[string]any{
		"ssh_host":      "203.0.113.44:22",
		"player_name":   "Sihaya",
		"account_id":    "12345",
		"admin_token":   "abcdefghijklmnopqrstuvwxyz0123456789ABCDEFG",
		"database_url":  "postgres://dune:supersecret@10.0.0.5:15432/dune",
		"support_email": "operator@example.com",
	}

	raw := safeDiagnosticJSON(input)
	if !json.Valid(raw) {
		t.Fatalf("expected valid JSON, got %s", string(raw))
	}
	out := string(raw)
	for _, leaked := range []string{
		"203.0.113.44",
		"10.0.0.5",
		"Sihaya",
		"12345",
		"abcdefghijklmnopqrstuvwxyz0123456789ABCDEFG",
		"supersecret",
		"operator@example.com",
	} {
		if strings.Contains(out, leaked) {
			t.Fatalf("diagnostic export leaked %q: %s", leaked, out)
		}
	}
	if !strings.Contains(out, "[REDACTED") {
		t.Fatalf("expected redaction markers in diagnostic export: %s", out)
	}
}

func TestBuildDiagnosticExportPayloadRedactionPolicy(t *testing.T) {
	listenAddr = "203.0.113.44:8080"
	startupConnectErr = "connect failed for player_name=Sihaya ADMIN_TOKEN=abcdefghijklmnopqrstuvwxyz0123456789ABCDEFG remote_addr=203.0.113.44:54321"
	payload := buildDiagnosticExportPayload()

	if !payload.Redaction.SecretsRedacted || !payload.Redaction.PIIRedacted {
		t.Fatalf("expected diagnostic export to advertise secret and PII redaction: %#v", payload.Redaction)
	}
	combined := payload.ListenAddr + " " + payload.StartupConnectStatus + " " + string(payload.Connectivity)
	for _, leaked := range []string{"203.0.113.44", "Sihaya", "abcdefghijklmnopqrstuvwxyz0123456789ABCDEFG"} {
		if strings.Contains(combined, leaked) {
			t.Fatalf("diagnostic export payload leaked %q: %#v", leaked, payload)
		}
	}
}
