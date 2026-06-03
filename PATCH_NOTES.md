# Dune Admin Release Notes

## Current update: Initial AppSec endpoint audit pass

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

---

## Previous update: Inventory Studio stack-size edit workflow

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
