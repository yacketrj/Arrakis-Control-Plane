# Dune Admin Release Notes

## Current update: Inventory Studio action history

### Why this update was made

Inventory Studio already had the post-action diff panel implemented and documented, but it only preserved the latest completed action diff. Operators need short-term browser-session context when doing several confirmed inventory actions in sequence.

### What changed

- Added browser-session action history state to `web/src/tabs/InventoryStudioTab.tsx`.
- Completed add, repair, and removal workflows now append their post-action diff records to action history.
- The history keeps the latest 10 completed action diffs for the selected player session.
- Added an **Action History** panel with action name, target, timestamp, before/after counts, diff count, and a short changed-row preview.
- Added action-history JSON export for local review.
- Added a clear action that resets the browser-session action history and the latest post-action diff.
- Selecting a different player clears action history to avoid mixing records across players.
- Updated `docs/inventory-studio.md` with action history behavior and safety notes.
- Updated `docs/admin-implementation-tasks.md` so the existing post-action diff panel is marked done and stack-size edit is the next planned Inventory Studio workflow after validation.

### Security and operator impact

- This is browser-local review state only. It does not create a new backend mutation, server-side persistence path, or rollback automation.
- The existing Inventory Studio add/repair/remove workflows still require shared mutation confirmation, before-action snapshot export, and admin reason capture.
- Action history does not replace the server-side audit log; it supplements operator review during the current browser session.
- Player 360 remains read-only. No Player 360 quick action or self-service mutation was added.

### Validation

Required from the canonical local update path:

```bash
./update.sh
```

Manual browser validation should confirm add, repair, and removal actions append history entries, history resets when a new player is selected, export downloads JSON, and clear removes local history.

---

## Previous update: Discord self-service frontend tabs

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

Verified from the canonical local update path after the Discord self-service frontend tab changes:

```bash
./update.sh
```

This covered backend tests/build plus frontend install/audit/typecheck/lint/build.

Manual browser validation has also been verified for **Discord Links** list/create/edit/delete behavior and **My Player Card** loading through Discord session cookies without a Browser Access Key.

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
