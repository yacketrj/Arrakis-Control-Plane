# Dune Admin Release Notes

## Current update: Discord player link foundation

### Why this update was made

Player Cards and future player-safe self-service need a durable identity-to-player mapping before any self-service behavior can be considered safe. This slice adds the mapping foundation and read-only self-service endpoints while preserving the rule that Player 360 and player mutations remain locked down.

### What changed

- Added `discord_player_links.go` with a local JSON-backed Discord-to-player link store.
- Added admin-managed link endpoints:
  - `GET /api/v1/auth/discord/player-links`
  - `POST /api/v1/auth/discord/player-links`
  - `DELETE /api/v1/auth/discord/player-links/{discord_id}`
- Added read-only self-service endpoints for linked Discord sessions:
  - `GET /api/v1/self/player-link`
  - `GET /api/v1/self/player-card`
- Updated auth middleware so normal Discord sessions can access only `/api/v1/self/*`; admin-token and Discord-admin access remain required elsewhere.
- Added `discord_player_links_test.go` coverage for link payload validation, store helper behavior, route handlers, current session link lookup, and self-service auth gating.
- Added `docs/discord-player-links.md` with model, storage, endpoint, auth-boundary, validation, and safety notes.
- Fixed Discord-player link text validation so raw control characters are rejected before trimming.

### Security and operator impact

- Normal Discord sessions are scoped to `/api/v1/self/*` only.
- Admin link management requires admin token or Discord admin session.
- Self-service player card output is read-only and derived from the existing Player 360 profile builder.
- This feature does not write player inventory, guild storage, claim rewards, currency, XP, Player 360, or any game-state table.
- Player 360 remains read-only. Any future self-service mutation must enforce the mapped player ID, mutation classification, auditability, and explicit safety workflows.

### Validation

Verified from the local checkout after the Discord-player link validation fix:

```bash
go test ./...
go build ./...
```

Manual release validation has also been verified for admin link CRUD, normal Discord self-service access, normal Discord denial from admin paths, unlinked Discord safe failures, and read-only self player-card behavior.

---

## Previous update: Farming Requests UI

### Why this update was made

The inventory request/order backend needs a protected operator UI so request and farming-order coordination can be managed without using raw API calls. This UI is intentionally coordination-only and does not deliver items or mutate game-state tables.

### What changed

- Added `web/src/api/inventoryRequests.ts` with typed frontend helpers for inventory request/order endpoints.
- Added `web/src/tabs/FarmingRequestsTab.tsx` with a protected Farming Requests operator tab.
- Wired the Farming Requests tab into `web/src/App.tsx` navigation.
- Updated `docs/inventory-requests-orders.md` with frontend UI behavior and validation notes.
- The tab supports request/order filtering, personal/guild request creation, open-request selection, farming-order creation, and order fill/cancel status updates.

### Security and operator impact

- The Farming Requests tab uses the coordination-only request/order backend and does not write player inventory, guild storage, claim rewards, currency, XP, Player 360, or game-state tables.
- The tab uses a separate frontend API module instead of extending the high-risk player/admin API surface.
- Player 360 remains read-only. Self-service player-card actions remain blocked until Discord identity-to-player mapping exists and is explicitly enforced.
- The Farming Requests UI validation has been verified from the local checkout.

### Validation

Verified from the local checkout after the Farming Requests UI changes:

```bash
./update.sh
```

Equivalent frontend validation gates covered by this verification:

```bash
cd web
npm run typecheck
npm run lint
npm run build
```

Manual browser validation remains recommended before release: confirm the Farming Requests tab can list, create, group, fill, and cancel request/order records without touching player inventory or Player 360 data.

---

## Previous update: Discord bot command adapter skeleton

### Why this update was made

The farming request/order backend needs a Discord-facing adapter layer before an actual bot runtime is introduced. This slice adds command-shape normalization and payload conversion without adding a networked bot process, Discord gateway client, or new runtime dependency.

### What changed

- Added `discord_bot_adapter.go` with a non-network command adapter for Discord-style farming request/order commands.
- Added adapter support for personal item requests, guild item requests, farm-order creation, filled-order updates, and cancelled-order updates.
- Reused existing inventory request/order validation and normalization instead of creating a second validation path.
- Added `discord_bot_adapter_test.go` coverage for personal requests, guild requests, farm orders, fill updates, cancel updates, and unsupported command rejection.

### Security and operator impact

- This is an adapter skeleton only. It does not connect to Discord, register slash commands, open a gateway connection, or execute bot actions on its own.
- The adapter does not mutate player inventory, guild storage, claim rewards, currency, XP, Player 360, or any game-state table.
- Player 360 remains read-only. Self-service player-card actions remain blocked until Discord identity-to-player mapping exists and is explicitly enforced.
- The adapter deliberately maps bot-style inputs into the existing request/order coordination model so future bot runtime work can stay thin and testable.

### Validation

User-provided local validation completed cleanly after the scoped frontend lint fix:

```bash
go test -v ./...
go build ./...
npm install
npm audit --audit-level=high
npm run typecheck
npm run lint
npm run build
```

The local `update.sh` run failed only at the Git auto-commit step because Git author identity was not configured in the local checkout.
