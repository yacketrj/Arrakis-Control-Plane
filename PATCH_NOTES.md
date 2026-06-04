# Dune Admin Release Notes

## Current update: Database endpoint security hardening

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

### Security and operator impact

- Database routes remain admin-only.
- Manual SQL remains restricted to single-statement read-only SQL by `isReadOnlySQL`.
- Returned sampled/search rows and manual SQL text now pass through `RedactSensitiveText` before reaching the browser.
- No database mutation capability, new admin route, Player 360 mutation, inventory mutation, or self-service database access was added.
- `ASEA-004` remains partially remediated pending validation and further timeout/manual abuse-case review.

### Validation

Required from the canonical local update path:

```bash
./update.sh
```

This should run the new database handler security tests.

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

---

## Previous update: Discord self-session route remediation

### Why this update was made

The AppSec endpoint audit identified `ASEA-002`: the Discord `me` and `logout` handlers supported session-cookie behavior, but middleware allowed registered non-admin Discord sessions only through `/api/v1/self/*`. That meant normal Discord users could use self-service player-card routes but could not reliably inspect or clear their own auth session through the intended Discord session endpoints.

### What changed

- Added `isDiscordSelfSessionRoute` and `isSelfServiceRoute` middleware helpers.
- Allowed registered non-admin Discord sessions to reach only:
  - `GET /api/v1/auth/discord/me`
  - `POST /api/v1/auth/discord/logout`
  - `/api/v1/self/*`
- Kept `GET /api/v1/auth/discord/users`, Discord player-link admin endpoints, player routes, database routes, infrastructure routes, and all admin mutation routes admin-only.
- Added AppSec regression tests confirming normal Discord sessions can reach `me`, `logout`, and self-service routes but cannot reach representative admin routes.
- Updated `docs/discord-auth.md` with the explicit normal-session route boundary.
- Updated `docs/appsec-endpoint-audit.md` so `ASEA-002` is validated as remediated.

### Security and operator impact

- This is a narrow auth-boundary change for self-session behavior.
- Normal Discord sessions can now inspect their own auth context and clear their own session cookie/session record.
- Normal Discord sessions still cannot access admin review, player, database, infrastructure, or mutation routes.
- No Player 360 mutation, inventory mutation, guild mutation, or direct game-state mutation was added.

### Validation

Verified from the canonical local update path:

```bash
./update.sh
```

The run was clean. It emitted this non-blocking build-performance warning:

```text
[PLUGIN_TIMINGS] Your build spent significant time in plugin `@tailwindcss/vite:generate:build`. See https://rolldown.rs/options/checks#plugintimings for more details.
```

This warning did not fail the validation gate.
