package main

func buildStatusPayload() map[string]any {
	return map[string]any{
		"ssh_connected": globalSSH != nil,
		"db_connected":  globalDB != nil,
		"pod_ns":        globalPodNS,
		"ssh_host":      sshHost,
		"tunnel_mode":   normalizedTunnelMode(),
		"tunnels":       tunnelStatus(),
	}
}
