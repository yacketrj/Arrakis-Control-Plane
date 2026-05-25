# Dune Admin Release Notes

## Current update: Database SQL shared mutation confirmation migration

### Why this update was made

The Database tab SQL runner is a privileged database operation surface. Even read-looking SQL is submitted through a protected mutation endpoint, so this update requires shared confirmation and an admin reason before any SQL is sent.

### What changed

- Updated `web/src/tabs/DatabaseTab.tsx` to use `useMutationConfirmation` for Run SQL.
- Kept read-only database views frictionless:
  - List Tables
  - Describe
  - Sample
  - Search Columns
- Added shared mutation confirmation and required admin reason capture before `api.database.sql` is called.
- Added SQL preview, database target metadata, common mutating-keyword detection, and blast-radius verification guidance to the confirmation details.
- Passed captured reasons into `api.database.sql`.
- Updated the SQL UI helper text so operators know confirmation and a reason are required.

### Security and operator impact

- Database SQL now uses the same shared confirmation foundation as active Players-table and Storage mutations.
- Admin reason capture is required before SQL requests are sent.
- Player 360 remains read-only.
- Remaining shared-confirmation review targets include Battlegroup Exec, Blueprint import, and future Inventory Studio/Player 360 quick-action workflows.

### Validation

Validation required in the Windows development environment:

```powershell
.\update.ps1
```

---

## Previous update: Storage shared mutation confirmation migration

Storage container add/remove item operations were migrated to shared mutation confirmation with required admin reason capture.

---

## Previous update: Confirmed player admin actions workflow

Player admin actions were migrated to a dedicated confirmed modal with required admin reason capture.

---

## Previous update: Confirmed player move workflow

Player move actions were migrated to a dedicated confirmed modal with required admin reason capture and online-state safeguards.

---

## Previous update: Journey node shared mutation confirmation migration

Journey node complete/reset actions were migrated to the confirmed Player Actions modal with required admin reason capture.

---

## Previous update: Player resource/action shared mutation confirmation migration

Resource, XP, specialization, and faction reputation actions were migrated to the confirmed Player Actions modal with required admin reason capture.

---

## Previous update: Inventory repair/delete shared mutation confirmation migration

Inventory repair/delete was migrated to the extracted confirmed inventory modal with required admin reason capture.

---

## Previous update: Give Item shared mutation confirmation migration

The Give Item modal was migrated from local browser confirm/prompt handling to the shared frontend mutation-confirmation hook with required admin reason capture.

---

## Previous update: Shared frontend mutation confirmation foundation

The reusable frontend mutation-confirmation hook was added so future high-risk UI actions can display mutation safety metadata and capture an admin reason before sending the mutation request.
