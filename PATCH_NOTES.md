# Dune Admin Release Notes

## Current update: AppSec auth boundary regression tests

### Why this update was made

The AppSec endpoint audit identified `ASEA-001`: route and auth-boundary expectations were documented, but representative automated regression tests were missing for public, self-service, admin-only, and WebSocket-ticket boundaries.

### What changed

- Added `appsec_auth_boundary_test.go`.
- Added tests for the public path allowlist.
- Added tests for self-service path classification.
- Added tests confirming public routes bypass backend admin-token validation.
- Added tests confirming representative admin routes reject missing tokens and accept a valid strict admin token.
- Added tests confirming self-service routes are denied without a Discord session or admin token.
- Added tests confirming WebSocket log-stream upgrades require a one-time ticket before admin-token fallback.
- Updated `docs/appsec-endpoint-audit.md` so `ASEA-001` is marked partially remediated pending local validation.

### Security and operator impact

- Test/documentation change only. No route behavior, auth behavior, endpoint implementation, UI behavior, or data mutation behavior changed.
- This reduces regression risk around public route allowlisting, protected route enforcement, self-service boundaries, and log-stream ticket handling.
- `ASEA-001` remains partially remediated until the new test coverage passes the canonical validation path.

### Validation

Required from the canonical local update path:

```bash
./update.sh
```

This should run the new Go auth-boundary regression tests.

---

## Previous update: Farming Requests lint warning fix

### Why this update was made

Frontend lint reported a `react-hooks/exhaustive-deps` warning in `web/src/tabs/FarmingRequestsTab.tsx` because the `useEffect` that reloads requests and orders referenced `load` without including it as a dependency.

### What changed

- Updated `FarmingRequestsTab.tsx` to import and use `useCallback`.
- Wrapped `load` in `useCallback` with `scopeFilter`, `requestStatus`, and `orderStatus` as explicit dependencies.
- Updated the reload effect to depend on `load` directly.

### Security and operator impact

- UI behavior should remain unchanged.
- This is a frontend lint/quality fix only.
- No route, backend behavior, auth behavior, inventory mutation, request/order model, or Player 360 behavior changed.

### Validation

Verified from the canonical local update path:

```bash
./update.sh
```

This cleared the reported `react-hooks/exhaustive-deps` warning.

---

## Previous update: Initial AppSec endpoint audit pass

### Why this update was made

The newly added P0 AppSec audit task needed a concrete starting point: a route inventory, auth-boundary summary, first findings, and a remediation checklist for public and protected backend endpoints.

### What changed

- Added `docs/appsec-endpoint-audit.md`.
- Documented the current global middleware/auth boundary from `auth.go` and `server.go`.
- Inventoried endpoints from `routes.go` across public, Discord/self-service, core status/diagnostics/audit, Battlegroup, player read, player mutation, inventory request/order, database, log, notification, storage, and blueprint groups.
- Added initial findings `ASEA-001` through `ASEA-006` covering endpoint auth-boundary regression tests, Discord session UX review, mutation reason/audit verification, database endpoint injection review, infrastructure/log endpoint review, and browser-token/CORS follow-up.
- Added a manual abuse-case checklist for future endpoint-by-endpoint validation.
- Marked the AppSec audit task as In Progress in `docs/admin-implementation-tasks.md`.

### Security and operator impact

- Documentation/audit pass only. No route, auth behavior, endpoint implementation, validation gate, or UI behavior changed.
- The audit document is intentionally not marked complete; handler-by-handler review, SAST, DAST, dependency review, and manual abuse-case validation remain open.
- Current Inventory Studio stack-size validation remains pending and unchanged.

### Validation

Documentation/audit review only. No build validation is required for this documentation-only update.
