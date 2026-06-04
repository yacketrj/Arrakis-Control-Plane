# Dune Admin Release Notes

## Current update: Mutation safety classification coverage

### Why this update was made

The AppSec endpoint audit item `ASEA-003` requires verification of high-risk mutation endpoints for admin reason handling, audit visibility, mutation-safety classification, request-size limits, and pre/post-change safety behavior. Initial review found several high-risk mutation routes were classified only as medium risk.

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
- Updated `docs/appsec-endpoint-audit.md` so `ASEA-003` is partially remediated pending local validation.

### Security and operator impact

- High-risk routes now correctly require reason and preview metadata in mutation-safety classification.
- When reason enforcement is enabled, these routes participate in `X-Admin-Reason` / body `reason` checks like other high-risk mutations.
- No new mutation route, UI workflow, game-state operation, or Player 360 mutation was added.
- Full endpoint-by-endpoint audit-event assertion coverage is still required before `ASEA-003` can be closed.

### Validation

Required from the canonical local update path:

```bash
./update.sh
```

This should run the updated Go mutation-safety tests.

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
