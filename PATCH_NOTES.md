# Dune Admin Release Notes

## Current update: Discord self-session route remediation

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
- Updated `docs/appsec-endpoint-audit.md` so `ASEA-002` is remediated pending validation.

### Security and operator impact

- This is a narrow auth-boundary change for self-session behavior.
- Normal Discord sessions can now inspect their own auth context and clear their own session cookie/session record.
- Normal Discord sessions still cannot access admin review, player, database, infrastructure, or mutation routes.
- No Player 360 mutation, inventory mutation, guild mutation, or direct game-state mutation was added.

### Validation

Required from the canonical local update path:

```bash
./update.sh
```

This should run the new AppSec Discord self-session route tests.

---

## Previous update: AppSec auth boundary regression tests

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
- Updated `docs/appsec-endpoint-audit.md` so `ASEA-001` is marked as validated partial remediation.

### Security and operator impact

- Test/documentation change only. No route behavior, auth behavior, endpoint implementation, UI behavior, or data mutation behavior changed.
- This reduces regression risk around public route allowlisting, protected route enforcement, self-service boundaries, and log-stream ticket handling.
- `ASEA-001` is validated as partial remediation; generated full-route auth-boundary coverage remains a future hardening step.

### Validation

Verified from the canonical local update path:

```bash
./update.sh
```

This validated the new Go auth-boundary regression tests.

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
