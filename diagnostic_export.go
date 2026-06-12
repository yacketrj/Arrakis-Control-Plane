package main

import (
	"encoding/json"
	"net/http"
	goruntime "runtime"
	"time"
)

type diagnosticExportPayload struct {
	GeneratedAt          string          `json:"generated_at"`
	Service              string          `json:"service"`
	Version              string          `json:"version"`
	GoOS                 string          `json:"go_os"`
	GoArch               string          `json:"go_arch"`
	Runtime              string          `json:"runtime"`
	ListenAddr           string          `json:"listen_addr"`
	ReasonEnforced       bool            `json:"reason_enforced"`
	Redaction            redactionPolicy `json:"redaction"`
	Connectivity         json.RawMessage `json:"connectivity,omitempty"`
	StartupConnectStatus string          `json:"startup_connect_status,omitempty"`
}

type redactionPolicy struct {
	SecretsRedacted bool     `json:"secrets_redacted"`
	PIIRedacted     bool     `json:"pii_redacted"`
	Excluded        []string `json:"excluded"`
}

func handleDiagnosticExport(w http.ResponseWriter, r *http.Request) {
	payload := buildDiagnosticExportPayload()
	w.Header().Set("Cache-Control", "no-store")
	jsonOK(w, payload)
}

func buildDiagnosticExportPayload() diagnosticExportPayload {
	payload := diagnosticExportPayload{
		GeneratedAt:    time.Now().UTC().Format(time.RFC3339Nano),
		Service:        "arrakis-control-panel",
		Version:        RedactPIIText(version),
		GoOS:           goruntime.GOOS,
		GoArch:         goruntime.GOARCH,
		Runtime:        RedactPIIText(normalizeRuntime(serverRuntime)),
		ListenAddr:     RedactPIIText(normalizeListenAddr(listenAddr)),
		ReasonEnforced: adminReasonEnforcementEnabled(),
		Redaction: redactionPolicy{
			SecretsRedacted: true,
			PIIRedacted:     true,
			Excluded: []string{
				"raw player names",
				"raw character names",
				"raw account identifiers",
				"raw FLS identifiers",
				"raw IP addresses",
				"secrets and tokens",
			},
		},
	}

	if startupConnectErr != "" {
		payload.StartupConnectStatus = RedactPIIText(startupConnectErr)
	}
	payload.Connectivity = safeDiagnosticJSON(runConnectivityDiagnostics())
	return payload
}

func safeDiagnosticJSON(value any) json.RawMessage {
	data, err := json.Marshal(value)
	if err != nil {
		fallback, _ := json.Marshal(map[string]string{"error": "diagnostic serialization failed"})
		return fallback
	}
	redacted := []byte(RedactPIIText(string(data)))
	if !json.Valid(redacted) {
		fallback, _ := json.Marshal(map[string]string{"error": "diagnostic redaction produced invalid JSON"})
		return fallback
	}
	return json.RawMessage(redacted)
}
