package main

import "testing"

func TestNormalizedTunnelMode(t *testing.T) {
	tests := []struct {
		name string
		mode string
		want string
	}{
		{name: "empty defaults auto", mode: "", want: "auto"},
		{name: "auto", mode: "auto", want: "auto"},
		{name: "managed alias", mode: "managed", want: "auto"},
		{name: "on alias", mode: "on", want: "auto"},
		{name: "existing", mode: "existing", want: "existing"},
		{name: "external alias", mode: "external", want: "existing"},
		{name: "local alias", mode: "local", want: "existing"},
		{name: "off", mode: "off", want: "off"},
		{name: "direct alias", mode: "direct", want: "off"},
		{name: "disabled alias", mode: "disabled", want: "off"},
		{name: "unknown defaults auto", mode: "unexpected", want: "auto"},
		{name: "case and trim", mode: " Existing ", want: "existing"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			old := sshTunnelMode
			t.Cleanup(func() { sshTunnelMode = old })
			sshTunnelMode = tt.mode
			if got := normalizedTunnelMode(); got != tt.want {
				t.Fatalf("normalizedTunnelMode() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestTunnelStatusStartsEmpty(t *testing.T) {
	old := managedTunnels
	t.Cleanup(func() { managedTunnels = old })
	managedTunnels = nil
	if got := tunnelStatus(); len(got) != 0 {
		t.Fatalf("expected no active tunnels, got %#v", got)
	}
}
