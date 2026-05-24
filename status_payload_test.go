package main

import "testing"

func TestBuildStatusPayloadIncludesTunnelFields(t *testing.T) {
	oldMode := sshTunnelMode
	oldHost := sshHost
	oldNS := globalPodNS
	oldSSH := globalSSH
	oldDB := globalDB
	oldTunnels := managedTunnels
	t.Cleanup(func() {
		sshTunnelMode = oldMode
		sshHost = oldHost
		globalPodNS = oldNS
		globalSSH = oldSSH
		globalDB = oldDB
		managedTunnels = oldTunnels
	})

	sshTunnelMode = "existing"
	sshHost = "example.invalid:22"
	globalPodNS = "dune-test"
	globalSSH = nil
	globalDB = nil
	managedTunnels = nil

	payload := buildStatusPayload()
	if payload["ssh_connected"] != false {
		t.Fatalf("ssh_connected = %#v, want false", payload["ssh_connected"])
	}
	if payload["db_connected"] != false {
		t.Fatalf("db_connected = %#v, want false", payload["db_connected"])
	}
	if payload["pod_ns"] != "dune-test" {
		t.Fatalf("pod_ns = %#v, want dune-test", payload["pod_ns"])
	}
	if payload["ssh_host"] != "example.invalid:22" {
		t.Fatalf("ssh_host = %#v, want example.invalid:22", payload["ssh_host"])
	}
	if payload["tunnel_mode"] != "existing" {
		t.Fatalf("tunnel_mode = %#v, want existing", payload["tunnel_mode"])
	}
	tunnels, ok := payload["tunnels"].([]map[string]string)
	if !ok {
		t.Fatalf("tunnels payload has type %T, want []map[string]string", payload["tunnels"])
	}
	if len(tunnels) != 0 {
		t.Fatalf("expected no tunnels, got %#v", tunnels)
	}
}
