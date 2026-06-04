package main

import (
	"strings"
	"testing"
	"time"
)

func TestNormalizeBattlegroupCommandStrictAllowlist(t *testing.T) {
	cmd, err := normalizeBattlegroupCommand(" Restart ")
	if err != nil {
		t.Fatalf("expected normalized command to be accepted: %v", err)
	}
	if cmd != "restart" {
		t.Fatalf("expected restart, got %q", cmd)
	}

	denied := []string{"", "restart; rm -rf /", "restart\nstop", "status"}
	for _, raw := range denied {
		if got, err := normalizeBattlegroupCommand(raw); err == nil {
			t.Fatalf("expected %q to be rejected, got %q", raw, got)
		}
	}
}

func TestValidateRuntimeCommandNamespace(t *testing.T) {
	oldRuntime := serverRuntime
	oldNS := globalPodNS
	t.Cleanup(func() {
		serverRuntime = oldRuntime
		globalPodNS = oldNS
	})

	serverRuntime = runtimeModeKubernetes
	globalPodNS = "funcom-operators"
	if err := validateRuntimeCommandNamespace(); err != nil {
		t.Fatalf("expected valid kubernetes namespace: %v", err)
	}

	globalPodNS = "bad;namespace"
	if err := validateRuntimeCommandNamespace(); err == nil {
		t.Fatalf("expected invalid kubernetes namespace to be rejected")
	}

	serverRuntime = runtimeModeDocker
	if err := validateRuntimeCommandNamespace(); err != nil {
		t.Fatalf("docker runtime should not require kubernetes namespace validation: %v", err)
	}
}

func TestSplitAndRedactLines(t *testing.T) {
	lines := splitAndRedactLines("ok\nDB_PASS=supersecret\n\nADMIN_TOKEN=abcdefghijklmnopqrstuvwxyz0123456789ABCDEFG")
	joined := strings.Join(lines, " ")
	if strings.Contains(joined, "supersecret") || strings.Contains(joined, "abcdefghijklmnopqrstuvwxyz0123456789ABCDEFG") {
		t.Fatalf("expected sensitive values to be redacted, got %#v", lines)
	}
	if !strings.Contains(joined, "[REDACTED]") {
		t.Fatalf("expected redaction marker, got %#v", lines)
	}
	if len(lines) != 3 {
		t.Fatalf("expected blank lines to be dropped, got %#v", lines)
	}
}

func TestRuntimeLogTargetValidation(t *testing.T) {
	oldRuntime := serverRuntime
	oldNS := globalPodNS
	t.Cleanup(func() {
		serverRuntime = oldRuntime
		globalPodNS = oldNS
	})

	serverRuntime = runtimeModeDocker
	if !isValidRuntimeLogTarget("docker", "abcdef123456") {
		t.Fatalf("expected docker container id to be valid")
	}
	if isValidRuntimeLogTarget("docker", "bad;target") {
		t.Fatalf("expected shell metacharacter docker target to be rejected")
	}

	serverRuntime = runtimeModeKubernetes
	globalPodNS = "dune"
	if isValidRuntimeLogTarget("dune", "bad;pod") {
		t.Fatalf("expected shell metacharacter kubernetes target to be rejected")
	}
}

func TestLogStreamTicketSingleUseWrongTargetAndExpiry(t *testing.T) {
	oldRuntime := serverRuntime
	t.Cleanup(func() { serverRuntime = oldRuntime })
	serverRuntime = runtimeModeDocker

	logStreamTickets.Lock()
	logStreamTickets.values = map[string]logStreamTicket{}
	logStreamTickets.Unlock()

	issued, err := issueLogStreamTicket("docker", "abcdef123456")
	if err != nil {
		t.Fatalf("issue ticket: %v", err)
	}
	if !consumeLogStreamTicket(issued.Ticket, "docker", "abcdef123456") {
		t.Fatalf("expected first consume to succeed")
	}
	if consumeLogStreamTicket(issued.Ticket, "docker", "abcdef123456") {
		t.Fatalf("expected replayed ticket to fail")
	}

	wrongTarget, err := issueLogStreamTicket("docker", "abcdef123456")
	if err != nil {
		t.Fatalf("issue second ticket: %v", err)
	}
	if consumeLogStreamTicket(wrongTarget.Ticket, "docker", "abcdef654321") {
		t.Fatalf("expected wrong target consume to fail")
	}
	if consumeLogStreamTicket(wrongTarget.Ticket, "docker", "abcdef123456") {
		t.Fatalf("expected ticket to be consumed after wrong-target attempt")
	}

	expiredTicket, err := generateLogStreamTicketValue()
	if err != nil {
		t.Fatalf("generate expired ticket value: %v", err)
	}
	logStreamTickets.Lock()
	logStreamTickets.values[expiredTicket] = logStreamTicket{Namespace: "docker", Pod: "abcdef123456", ExpiresAt: time.Now().UTC().Add(-time.Second)}
	logStreamTickets.Unlock()
	if consumeLogStreamTicket(expiredTicket, "docker", "abcdef123456") {
		t.Fatalf("expected expired ticket to fail")
	}
}

func TestIssueLogStreamTicketRejectsInvalidTarget(t *testing.T) {
	oldRuntime := serverRuntime
	t.Cleanup(func() { serverRuntime = oldRuntime })
	serverRuntime = runtimeModeDocker

	if _, err := issueLogStreamTicket("docker", "bad;target"); err == nil {
		t.Fatalf("expected invalid docker target to be rejected")
	}
	if _, err := issueLogStreamTicket("", "abcdef123456"); err == nil {
		t.Fatalf("expected missing namespace to be rejected")
	}
}

func TestRedactCheatEntries(t *testing.T) {
	rows := []cheatEntry{{
		FLSID:         "ADMIN_TOKEN=abcdefghijklmnopqrstuvwxyz0123456789ABCDEFG",
		CheatType:     "speed",
		EventTime:     "2026-06-04 00:00:00",
		CharacterName: "DB_PASS=supersecret",
	}}
	got := redactCheatEntries(rows)
	combined := got[0].FLSID + " " + got[0].CharacterName
	if strings.Contains(combined, "abcdefghijklmnopqrstuvwxyz0123456789ABCDEFG") || strings.Contains(combined, "supersecret") {
		t.Fatalf("expected cheat log fields to be redacted, got %#v", got[0])
	}
	if !strings.Contains(combined, "[REDACTED]") {
		t.Fatalf("expected redaction marker, got %#v", got[0])
	}
	if rows[0].FLSID == got[0].FLSID {
		t.Fatalf("expected helper to return redacted copy")
	}
}
