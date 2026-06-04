package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestDBQueryParamValidation(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/database/describe?table="+strings.Repeat("a", maxDBQueryParamLength+1), nil)
	if _, err := requiredDBQueryParam(req, "table"); err == nil {
		t.Fatalf("expected overlong table parameter to be rejected")
	}

	req = httptest.NewRequest(http.MethodGet, "/api/v1/database/search?term=player%00name", nil)
	if _, err := requiredDBQueryParam(req, "term"); err == nil {
		t.Fatalf("expected control-character query parameter to be rejected")
	}

	req = httptest.NewRequest(http.MethodGet, "/api/v1/database/search?term=%20player%20", nil)
	got, err := requiredDBQueryParam(req, "term")
	if err != nil {
		t.Fatalf("expected trimmed query parameter to be accepted: %v", err)
	}
	if got != "player" {
		t.Fatalf("expected trimmed value, got %q", got)
	}
}

func TestRedactDBStringRows(t *testing.T) {
	rows := [][]string{{"player", "ADMIN_TOKEN=abcdefghijklmnopqrstuvwxyz0123456789ABCDEFG"}, {"password", "DB_PASS=supersecret"}}
	got := redactDBStringRows(rows)
	flat := strings.Join([]string{got[0][1], got[1][1]}, " ")
	if strings.Contains(flat, "abcdefghijklmnopqrstuvwxyz0123456789ABCDEFG") || strings.Contains(flat, "supersecret") {
		t.Fatalf("expected sensitive values to be redacted, got %#v", got)
	}
	if !strings.Contains(flat, "[REDACTED]") {
		t.Fatalf("expected redaction marker, got %#v", got)
	}
	if rows[0][1] == got[0][1] {
		t.Fatalf("expected returned rows to be redacted copy")
	}
}

func TestHandleDBFunctionInspectRejectsNonNumericOID(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/database/functions/inspect?oid=1;drop", nil)
	res := httptest.NewRecorder()

	handleDBFunctionInspect(res, req)

	if res.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for non-numeric oid, got %d", res.Code)
	}
}

func TestHandleDBSQLRejectsUnsafeSQLBeforeDatabaseUse(t *testing.T) {
	oldDB := globalDB
	globalDB = nil
	t.Cleanup(func() { globalDB = oldDB })

	cases := []string{
		`{"sql":"delete from dune.actors"}`,
		`{"sql":"select * from dune.actors; drop table dune.items"}`,
		`{"sql":"update dune.actors set id = id"}`,
	}
	for _, body := range cases {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/database/sql", strings.NewReader(body))
		res := httptest.NewRecorder()
		handleDBSQL(res, req)
		if res.Code != http.StatusBadRequest {
			t.Fatalf("%s: expected 400 for unsafe SQL, got %d", body, res.Code)
		}
	}
}

func TestHandleDBSQLTrimsSQLBeforeReadOnlyCheck(t *testing.T) {
	oldDB := globalDB
	globalDB = nil
	t.Cleanup(func() { globalDB = oldDB })

	req := httptest.NewRequest(http.MethodPost, "/api/v1/database/sql", strings.NewReader(`{"sql":"   delete from dune.actors   "}`))
	res := httptest.NewRecorder()

	handleDBSQL(res, req)

	if res.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for trimmed unsafe SQL, got %d", res.Code)
	}
}

func TestHandleDBSQLRequiresSQL(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/v1/database/sql", strings.NewReader(`{"sql":"   "}`))
	res := httptest.NewRecorder()

	handleDBSQL(res, req)

	if res.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for blank SQL, got %d", res.Code)
	}
}

func TestHandleDBSearchRejectsBadTermBeforeDatabaseUse(t *testing.T) {
	oldDB := globalDB
	globalDB = nil
	t.Cleanup(func() { globalDB = oldDB })

	req := httptest.NewRequest(http.MethodGet, "/api/v1/database/search?term="+strings.Repeat("a", maxDBQueryParamLength+1), nil)
	res := httptest.NewRecorder()

	handleDBSearch(res, req)

	if res.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for overlong search term, got %d", res.Code)
	}
}

func TestDBSQLResponseRedactionShape(t *testing.T) {
	result := RedactSensitiveText("DB_PASS=supersecret")
	payload := map[string]string{"result": result}
	data, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("marshal payload: %v", err)
	}
	if strings.Contains(string(data), "supersecret") || !strings.Contains(string(data), "[REDACTED]") {
		t.Fatalf("expected redacted SQL result payload, got %s", string(data))
	}
}
