# Dune Admin Release Notes

## Current update: Live Claim Rewards delivery mode

### Why this update was made

Direct inventory database writes support the richest Give Item workflow, including stack placement, grades, and augmented item stats. Online players may not see those direct inventory changes until they relog because the game client/server inventory state is already loaded.

This update exposes the existing Claim Rewards delivery path directly in the Give Item modal as a separate operator-selectable delivery mode.

### What changed

- Added a delivery mode selector to `GiveItemModalAugmented.tsx`.
- Added `Inventory Write` mode for direct database item creation with item grade, stack size, and augment support.
- Added `Live Claim Rewards` mode for online-friendly plain item delivery through the existing live reward API path.
- Added delivery-mode-aware payload preview.
- Blocked Live Claim Rewards mode when selected rows contain item grade greater than zero or augments, because the live reward path only supports plain template-and-amount entries.
- Disabled grade and augment controls while Live Claim Rewards mode is selected.

### Operator guidance

Use `Live Claim Rewards` when the player is online and you need them to receive plain items without logging out. The player should see a Claim Rewards prompt.

Use `Inventory Write` when you need exact inventory placement, item grade, stack control, or augmented item stats. For online players, this may still require logout/login before the game refreshes client-visible inventory state.

### Validation

Expected validation:

```bash
go mod tidy
go test ./...
cd web
npm install
npm audit --audit-level=high
npm run typecheck
npm run lint
npm run build
```

---

## Previous update: Go Quality run 26326549538 remediation

### Why this update was made

GitHub Actions run `26326549538` failed in the expanded template merge test because the test expected three merged templates, while the hybrid merge correctly returns four: two live database templates, one curated name-only template, and one curated item-rule template.

### Remediation

- Updated the test to expect four merged templates.
- Added explicit verification that the curated `item_rule_template` entry is present in the serialized template response.
- Preserved the lower-case friendly-name lookup assertion for `DB_TEMPLATE`.

### Validation

The next Go Quality run should proceed past `TestMergeItemTemplatesAndHandleGetTemplates`.

---

## Previous update: Give Item helper cleanup and expanded test coverage

### Why this update was made

The previous polish pass added a shared frontend Give Item payload helper, but the active augmented Give Item modal still duplicated clamping, preset, roll parsing, and payload mapping logic inline. This update finishes that cleanup and expands backend test coverage around augmented item and template-list behavior.

### Security and operator impact

- Updated `GiveItemModalAugmented.tsx` to use `web/src/tabs/giveItemPayload.ts` directly.
- Removed duplicated frontend payload helper logic from the modal so the UI and future tests exercise one payload mapping path.
- Added backend tests for augment quality-to-grade alias handling.
- Added backend tests for augment roll defaulting, repeated roll generation, and explicit roll precedence.
- Added backend tests for database-plus-JSON template merging and `/api/v1/players/templates` response serialization.
- Preserved the augmented Give Item workflow with item grade, augment grade, roll values, explicit roll arrays, roll count, and generated payload preview.

### Validation

Expected validation:

```bash
go mod tidy
go test ./...
cd web
npm install
npm audit --audit-level=high
npm run typecheck
npm run lint
npm run build
```

### Remaining polish work

- Watch the next Go Quality and Frontend Quality workflow runs and fix any compile, typecheck, lint, or build failures.
- Continue splitting `PlayersTab.tsx` into smaller player table, inventory, action, and modal components.
- Regenerate and commit `web/package-lock.json` from a clean local `npm install` once the frontend manifest is stable.
- Expand the augment preset catalog as more verified in-game augmented item examples are captured.
