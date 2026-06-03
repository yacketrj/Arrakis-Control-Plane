# Dune Admin Release Notes

## Current update: Inventory Studio stack-size edit workflow

### Why this update was made

Inventory Studio needed the next narrow confirmed edit workflow after action-history validation. Stack-size editing is the smallest direct item-row edit and can reuse the existing safety pattern: before-action snapshot, shared confirmation, admin reason, reload, post-action diff, and browser-session action history.

### What changed

- Added `inventory_stack_size.go` with a protected backend command and handler for item stack-size updates.
- Registered `POST /api/v1/players/item/stack-size` in `routes.go`.
- Added `web/src/api/inventoryStudioMutations.ts` with a frontend helper for stack-size edits.
- Updated `web/src/tabs/InventoryStudioTab.tsx` with a confirmed selected-item stack-size edit control.
- Stack-size edits clamp and validate values to `1..9999`.
- Stack-size edits now export a before-action inventory snapshot, require shared mutation confirmation and admin reason capture, reload inventory after success, display post-action diff, and append to action history.
- Updated `docs/inventory-studio.md` with the stack-size endpoint, payload, UI behavior, and safety notes.
- Updated `docs/admin-implementation-tasks.md` so stack-size edit is In Progress and quality edit is the next planned Inventory Studio workflow after validation.

### Security and operator impact

- This adds one direct item-row mutation path for `dune.items.stack_size` only.
- The backend validates item ID and stack size and rejects missing, zero, negative, or overly large stack sizes.
- The frontend sends the admin reason through `X-Admin-Reason` and preserves the same operator confirmation flow as add, repair, and remove.
- No quality edit, template replacement, augment/stat edit, Player 360 mutation, self-service mutation, or rollback automation was added.
- Player 360 remains read-only.

### Validation

Required from the canonical local update path:

```bash
./update.sh
```

Manual browser validation should confirm selected-item stack-size edit, before-action snapshot export, required reason capture, post-action diff, action-history append, unchanged-value guard, and inventory reload behavior.

---

## Previous update: Backlog planning additions

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
