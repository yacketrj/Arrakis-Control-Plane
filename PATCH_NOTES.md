# Dune Admin Release Notes

## Current update: Backlog planning additions

### Why this update was made

The implementation tracker needed durable entries for newly requested roadmap work so future coding passes do not lose the intended scope.

### What changed

- Added a P5 documentation task for a detailed Discord bot setup and usage guide.
- Added a P2 Guild Management feature covering create/delete guild, add/remove player membership, and player guild-rank changes.
- Added a P2 Player tab guild workflow feature covering add/remove selected player from a guild and promote/change guild rank from the Player tab.
- Updated `docs/admin-implementation-tasks.md` with guild-management backlog scope, Player tab guild workflow scope, and Discord bot documentation guide scope.

### Security and operator impact

- Planning-only change. No code, route, schema write, bot runtime, or UI mutation path was added.
- Future guild-management work must perform schema discovery before writes are implemented.
- Future guild mutations must use shared mutation confirmation, admin reason capture, before-change snapshot/review, post-action refresh/diff where practical, and audit visibility.
- Future Discord bot documentation must include setup, permissions, secret handling, command behavior, troubleshooting, and security boundaries.

### Validation

Documentation/planning review only. No build validation is required for this planning-only update.

---

## Previous update: Inventory Studio action history

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

Verified from the canonical local update path after the Inventory Studio action-history changes:

```bash
./update.sh
```

Manual browser validation has also been verified for add/repair/removal history append behavior, reset on player change, JSON export, and clear-history behavior.

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
