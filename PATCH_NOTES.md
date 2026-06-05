# Dune Admin Release Notes

## Current update: Generated route auth-boundary coverage

### Why this update was made

The AppSec endpoint audit identified generated full-route auth-boundary coverage as an open follow-up under `ASEA-001`. Representative auth-boundary tests existed, but new routes could still be added to `routes.go` without an explicit AppSec exposure expectation.

### What changed

- Added `appsec_route_inventory_test.go`.
- The new test parses `routes.go` using Go's AST parser and extracts every registered `mux.HandleFunc("METHOD /path", handler)` route.
- Every registered route must now have an explicit auth-boundary expectation in `appsecExpectedRouteAuth`.
- The test fails if:
  - a registered route has no auth-boundary expectation
  - an expectation references a route no longer registered in `routes.go`
  - a public route does not bypass auth as expected
  - a self-service route is not denied without auth, allowed with admin token, and allowed with a normal registered Discord session
  - an admin route is not denied without auth, allowed with admin token, and denied to a normal registered Discord session
  - the WebSocket log-stream upgrade path does not require a one-time ticket before admin-token fallback
- Fixed the route-inventory test helper name from `containsString` to `appsecContainsString` to avoid colliding with the existing production helper in `db_functions.go`.
- Added `docs/generated-route-auth-boundary-coverage.md` to document the route-inventory/auth-boundary coverage model.
- Updated `docs/appsec-endpoint-audit.md` so `ASEA-001` generated full-route coverage is validated as partial remediation.

### Security and operator impact

- No route behavior changed.
- No new route was added.
- No mutation path was added.
- Player 360 remains read-only.
- Future route additions now fail local validation until their auth-boundary expectation is explicitly reviewed and added.

### Validation

Verified from the canonical local update path:

```bash
./update.sh
```

This validated the generated route inventory/auth-boundary coverage test after the test-helper rename fix.

---

## Previous update: Browser token and CORS security hardening

### Why this update was made

The AppSec endpoint audit item `ASEA-006` requires review of browser token handling and CORS behavior. The current Browser Access Key flow is already session-scoped and strict-token validated, but CORS allowed-origin parsing needed stronger protection against unsafe misconfiguration values.

### What changed

- Hardened `auth.go` CORS origin parsing:
  - rejects wildcard `*`
  - rejects `null`
  - rejects control characters
  - rejects non-HTTP(S) schemes such as `file://` and `javascript:`
  - rejects origins with userinfo
  - rejects origins with path, query, or fragment components
  - continues exact-match origin allowlisting only
- Added CORS/origin tests in `auth_test.go` covering:
  - unsafe origin value rejection
  - safe origin value acceptance
  - mixed safe/unsafe allowed-origin parsing
  - disallowed-origin preflight behavior
  - allowed-origin preflight reflection and `Vary: Origin`
- Added `docs/browser-token-cors-security.md` to capture the `ASEA-006` review state, current guardrails, tests, and remaining work.
- Updated `docs/appsec-endpoint-audit.md` so `ASEA-006` is validated as partial remediation.

### Security and operator impact

- CORS remains exact-match only.
- Disallowed origins are not reflected.
- Unsafe `ALLOWED_ORIGINS` values are ignored instead of accepted.
- Browser Access Key storage remains session-scoped in `sessionStorage`, with legacy `localStorage` cleanup still in place.
- No new admin route, auth mode, mutation path, Player 360 mutation, inventory mutation, or self-service admin access was added.
- `ASEA-006` is validated as partial remediation; future memory-only or HttpOnly secure session-cookie authentication design remains open.

### Validation

Verified from the canonical local update path:

```bash
./update.sh
```

This validated the updated CORS/origin tests and browser-token/CORS hardening changes.

---

## Previous update: Infrastructure and log endpoint security hardening

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
- Updated `docs/appsec-endpoint-audit.md` so `ASEA-005` is validated as partial remediation.

### Security and operator impact

- Infrastructure and log endpoints remain admin-only.
- Battlegroup exec remains restricted to the static allowlist: `start`, `stop`, `restart`, `update`, `backup`, and `restore`.
- Kubernetes namespace command interpolation is now guarded by shared validation.
- Log-stream tickets remain one-time, scoped, and 60-second TTL limited.
- Wrong-target ticket use consumes and rejects the ticket.
- Returned remote output is redacted before reaching the browser.
- No new infrastructure command, direct game-state mutation, Player 360 mutation, inventory mutation, or self-service log access was added.
- `ASEA-005` is validated as partial remediation; further handler-level SSH/database-stub tests, command timeout review, WebSocket origin review, live runtime/manual validation, and real-output redaction review remain open.

### Validation

Verified from the canonical local update path:

```bash
./update.sh
```

This validated the new infrastructure/log security tests and infrastructure/log endpoint hardening changes.
