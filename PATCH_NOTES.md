# Dune Admin Release Notes

## Current update: Infrastructure and log endpoint security hardening

### Why this update was made

The AppSec endpoint audit item `ASEA-005` requires review of infrastructure command and log endpoints for command allowlisting, runtime target validation, log-stream ticket replay/TTL behavior, and output redaction. Initial review found that namespace validation and output redaction should be applied more consistently across Battlegroup and log paths.

### What changed

- Hardened `handlers_battlegroup.go`:
  - validates Kubernetes namespace before status, health, and pod command construction
  - normalizes Battlegroup command input by trimming and lowercasing it
  - rejects Battlegroup command control characters
  - enforces the static Battlegroup command allowlist
  - redacts status, health, exec, and pod-list output before returning it
- Hardened `handlers_logs.go`:
  - validates runtime namespace before log target discovery
  - redacts Docker display names
  - redacts log stream error and line output
  - redacts returned cheat-log fields before returning rows
- Added `infrastructure_security_test.go` covering:
  - Battlegroup command normalization and strict allowlist behavior
  - command control-character/metacharacter rejection
  - Kubernetes namespace validation
  - Docker runtime namespace bypass behavior
  - split-and-redact line handling
  - Docker/Kubernetes log target rejection for unsafe targets
  - log-stream ticket single-use behavior
  - log-stream ticket wrong-target behavior
  - expired-ticket rejection
  - invalid-target ticket issuance rejection
  - cheat-log field redaction
- Added `docs/infrastructure-log-endpoint-security.md` to capture the `ASEA-005` review state, guardrails, tests, and remaining work.
- Updated `docs/appsec-endpoint-audit.md` so `ASEA-005` is partially remediated pending validation.

### Security and operator impact

- Infrastructure and log endpoints remain admin-only.
- Battlegroup exec remains restricted to the static allowlist: `start`, `stop`, `restart`, `update`, `backup`, and `restore`.
- Kubernetes namespace command interpolation is now guarded by shared validation.
- Log-stream tickets remain one-time, scoped, and 60-second TTL limited.
- Wrong-target ticket use consumes and rejects the ticket.
- Returned remote output is redacted before reaching the browser.
- No new infrastructure command, direct game-state mutation, Player 360 mutation, inventory mutation, or self-service log access was added.
- `ASEA-005` remains partially remediated pending validation and further handler-level SSH/database-stub tests, command timeout review, WebSocket origin review, live runtime/manual validation, and real-output redaction review.

### Validation

Required from the canonical local update path:

```bash
./update.sh
```

This should run the new infrastructure/log security tests.

---

## Previous update: Database endpoint security hardening

### Why this update was made

The AppSec endpoint audit item `ASEA-004` requires review of database endpoints for SQL injection, read-only SQL enforcement, timeout/result-limit behavior, and redaction. Initial review found the main SQL paths already used parameterization or safe table quoting, but handler-level query parameter validation and response redaction needed stronger durable coverage.

### What changed

- Hardened `handlers_database.go` database endpoint inputs:
  - trims database query parameters
  - rejects overlong database query parameters above 128 characters
  - rejects unsafe control characters in database query parameters
  - requires numeric `oid` for function inspection
  - trims submitted manual SQL before validation
- Redacts sampled/search database rows before returning them.
- Redacts manual SQL output before returning it.
- Preserves existing sample/manual SQL row limits:
  - sample endpoint clamps to 200 rows
  - manual SQL output remains capped at 200 rows in `cmdRunSQL`
- Added `handlers_database_test.go` covering:
  - overlong database query parameters
  - unsafe control-character parameters
  - trimmed parameter behavior
  - database row redaction helper behavior
  - non-numeric function OID rejection
  - unsafe SQL rejection before database use
  - trimmed unsafe SQL rejection
  - blank SQL rejection
  - overlong search-term rejection before database use
  - redacted SQL response payload shape
- Added `docs/database-endpoint-security.md` to capture the `ASEA-004` review state, guardrails, tests, and remaining work.
- Updated `docs/appsec-endpoint-audit.md` so `ASEA-004` is validated as partial remediation.

### Security and operator impact

- Database routes remain admin-only.
- Manual SQL remains restricted to single-statement read-only SQL by `isReadOnlySQL`.
- Returned sampled/search rows and manual SQL text now pass through `RedactSensitiveText` before reaching the browser.
- No database mutation capability, new admin route, Player 360 mutation, inventory mutation, or self-service database access was added.
- `ASEA-004` is validated as partial remediation; SQL timeout/manual abuse-case review remains open.

### Validation

Verified from the canonical local update path:

```bash
./update.sh
```

This validated the new database handler security tests and database endpoint hardening changes.

---

## Previous update: Mutation safety classification coverage

### Why this update was made

The AppSec endpoint audit item `ASEA-003` requires verification of high-risk mutation endpoints for admin reason handling, audit visibility, mutation-safety classification, request-size limits, and pre/post-change safety behavior. Initial review found several high-risk mutation routes were classified only as medium risk.

A follow-up validation run also found an invalid JSON payload in `audit_log_test.go` that prevented the audit metadata reason from being parsed. The update script output was also improved so test status lines are easier to read during local validation.

### What changed

- Tightened mutation-safety risk classification in `audit_log.go`.
- Marked these mutation paths high risk:
  - `POST /api/v1/reconnect`
  - `POST /api/v1/database/sql`
  - `POST /api/v1/logs/stream-ticket`
  - `POST /api/v1/notify`
  - `POST /api/v1/players/item/stack-size`
- Preserved destructive classification for reset, wipe, delete, and blueprint import paths.
- Added regression coverage in `mutation_safety_test.go` for high-risk mutation paths.
- Added regression coverage for destructive mutation paths.
- Added an oversized-body reason-enforcement test for high-risk mutations when `ADMIN_REQUIRE_REASON=true`.
- Fixed `audit_log_test.go` so the audit metadata payload uses valid JSON with an escaped newline in the `reason` field.
- Updated `update.sh` to colorize validation output:
  - `=== RUN` in cyan
  - `PASS` and `--- PASS:` in green
  - `FAIL` and `--- FAIL:` in red
  - `Update failed.` in red
- Updated `docs/appsec-endpoint-audit.md` so `ASEA-003` is validated as partial remediation.

### Security and operator impact

- High-risk routes now correctly require reason and preview metadata in mutation-safety classification.
- When reason enforcement is enabled, these routes participate in `X-Admin-Reason` / body `reason` checks like other high-risk mutations.
- Local validation output now makes failures more visually obvious; `FAIL` is red.
- No new mutation route, UI workflow, game-state operation, or Player 360 mutation was added.
- Full endpoint-by-endpoint audit-event assertion coverage is still required before `ASEA-003` can be closed.

### Validation

Verified from the canonical local update path:

```bash
./update.sh
```

This validated the updated Go mutation-safety and audit-log tests, including the corrected audit metadata JSON payload and the colored update output path.
