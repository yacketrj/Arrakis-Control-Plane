package main

import "testing"

func TestRedactSensitiveTextAssignments(t *testing.T) {
	input := "ADMIN_TOKEN=abcdefghijklmnopqrstuvwxyz0123456789ABCDEFG DB_PASS=supersecret funcom_token:abc123456789"
	got := RedactSensitiveText(input)
	for _, leaked := range []string{"abcdefghijklmnopqrstuvwxyz0123456789ABCDEFG", "supersecret", "abc123456789"} {
		if contains := containsSubstring(got, leaked); contains {
			t.Fatalf("redacted output leaked %q: %s", leaked, got)
		}
	}
	if !containsRedaction(got) {
		t.Fatalf("expected redaction marker in %q", got)
	}
}

func TestRedactSensitiveTextJSONAndBearer(t *testing.T) {
	input := `{"admin_token":"abcdefghijklmnopqrstuvwxyz0123456789ABCDEFG","password":"hunter2","ok":"visible"} Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.payload.sig`
	got := RedactSensitiveText(input)
	for _, leaked := range []string{"abcdefghijklmnopqrstuvwxyz0123456789ABCDEFG", "hunter2", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.payload.sig"} {
		if containsSubstring(got, leaked) {
			t.Fatalf("redacted output leaked %q: %s", leaked, got)
		}
	}
	if !containsSubstring(got, `"ok":"visible"`) {
		t.Fatalf("expected non-sensitive field to remain visible: %s", got)
	}
}

func TestRedactPostgresURLPassword(t *testing.T) {
	input := "postgres://dune:supersecret@127.0.0.1:15432/dune"
	got := RedactSensitiveText(input)
	if containsSubstring(got, "supersecret") {
		t.Fatalf("redacted output leaked postgres password: %s", got)
	}
	if !containsSubstring(got, "postgres://dune:[REDACTED]@127.0.0.1:15432/dune") {
		t.Fatalf("unexpected postgres URL redaction: %s", got)
	}
}

func containsSubstring(value, needle string) bool {
	return len(needle) == 0 || (len(value) >= len(needle) && indexSubstring(value, needle) >= 0)
}

func indexSubstring(value, needle string) int {
	for i := 0; i+len(needle) <= len(value); i++ {
		if value[i:i+len(needle)] == needle {
			return i
		}
	}
	return -1
}
