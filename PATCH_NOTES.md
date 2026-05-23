# Dune Admin Release Notes

## Current update: Give Item helper cleanup and expanded test coverage

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

---

## Previous update: Frontend Quality workflow stabilization

### Why this update was made

GitHub Actions run `26325532482` failed before frontend dependency installation because `actions/setup-node` referenced a frontend lockfile that is not currently committed.

### Security and operator impact

- Removed the frontend workflow cache dependency on the missing lockfile.
- Kept frontend validation active: install, high-severity npm audit, TypeScript typecheck, lint, and build.
- Added `FORCE_JAVASCRIPT_ACTIONS_TO_NODE24=true` to the frontend workflow.

### Validation

```bash
cd web
npm install
npm audit --audit-level=high
npm run typecheck
npm run lint
npm run build
```
