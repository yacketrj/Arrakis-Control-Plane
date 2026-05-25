package main

func buildStatusPayload() map[string]any {
	return map[string]any{
		"ssh_connected":         globalSSH != nil,
		"db_connected":          globalDB != nil,
		"pod_ns":                globalPodNS,
		"ssh_host":              sshHost,
		"runtime":               normalizeRuntime(serverRuntime),
		"tunnel_mode":           normalizedTunnelMode(),
		"tunnels":               tunnelStatus(),
		"admin_reason_required": adminReasonEnforcementEnabled(),
	}
}