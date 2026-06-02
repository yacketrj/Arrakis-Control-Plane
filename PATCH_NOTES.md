# Dune Admin Release Notes

## Current update: Discord self-service frontend tabs

### Why this update was made

The validated Discord-to-player mapping foundation now needs frontend surfaces so admins can manage links and linked Discord users can view their own read-only player card. This slice exposes that functionality without adding any player mutation or Player 360 write path.

### What changed

- Added `web/src/api/discordSelfService.ts` with typed helpers for Discord player links and `/api/v1/self/*` endpoints.
- Added `web/src/tabs/DiscordPlayerLinksTab.tsx` for admin link management.
- Added `web/src/tabs/SelfPlayerCardTab.tsx` for read-only linked-player self-service.
- Wired **Discord Links** and **My Player Card** into `web/src/App.tsx` navigation.
- Updated `docs/discord-player-links.md` with frontend tab behavior, cookie-aware self-service calls, and validation notes.
- Adjusted tab gating so **My Player Card** can load through Discord session cookies even when no Browser Access Key is configured.

### Security and operator impact

- **Discord Links** remains an admin surface and still requires Browser Access Key or Discord admin authorization through the backend.
- **My Player Card** calls only `/api/v1/self/player-link` and `/api/v1/self/player-card` with browser session cookies.
- No player inventory, guild storage, claim rewards, currency, XP, Player 360, or game-state mutation path was added.
- Player 360 remains read-only. Future self-service mutation remains blocked until explicit mapped-player enforcement, auditability, and mutation-safety workflows are implemented.

### Validation

Required from the canonical local update path:

```bash
./update.sh
```

This must cover backend tests/build plus frontend install/audit/typecheck/lint/build. Manual browser validation should confirm admin link management and read-only self player-card behavior through Discord session cookies.

---

## Previous update: Discord player link foundation

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
