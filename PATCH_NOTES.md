# Dune Admin Release Notes

## Current update: High-risk mutation audit-event coverage

### Why this update was made

The AppSec endpoint audit identified endpoint-by-endpoint audit-event assertion coverage as an open follow-up under `ASEA-003`. Representative audit tests existed, but the high-risk/destructive mutation route set needed durable coverage proving that audit events carry the expected accountability fields.

### What changed

- Expanded `audit_log_test.go` with `TestAuditMiddlewareHighRiskAndDestructiveRouteCoverage`.
- The new table-driven test enumerates high-risk and destructive mutation routes and verifies each route emits exactly one audit event.
- The test asserts:
  - request method
  - request path
  - mutation-safety action
  - mutation-safety risk
  - destructive flag
  - requires-reason flag
  - requires-preview flag
  - HTTP status
  - result outcome
  - sanitized reason
  - common target metadata
  - request ID
- Added `docs/high-risk-mutation-audit-coverage.md` to document the coverage model, covered routes, security impact, and remaining work.
- Updated `docs/appsec-endpoint-audit.md` so `ASEA-003` audit-event assertion coverage is validated as partial remediation.

### Security and operator impact

- No route behavior changed.
- No mutation behavior changed.
- No new endpoint was added.
- Player 360 remains read-only.
- This adds regression evidence that high-risk/destructive mutation routes produce expected audit accountability fields.
- `ASEA-003` is validated as partial remediation; route-specific target assertions and pre/post-change review verification remain open.

### Validation

Verified from the canonical local update path:

```bash
./update.sh
```

This validated the expanded audit-log coverage.

---

## Previous update: Generated route auth-boundary coverage

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
