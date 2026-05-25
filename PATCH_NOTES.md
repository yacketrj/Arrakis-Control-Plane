# Dune Admin Release Notes

## Current update: Battlegroup Exec shared mutation confirmation migration

### Why this update was made

Battlegroup Exec controls server lifecycle actions such as start, stop, restart, update, backup, and restore. These are privileged operations that can disrupt connected players or alter active server state. This update moves Server Control actions onto the shared mutation-confirmation hook with required admin reason capture.

### What changed

- Updated `web/src/tabs/BattlegroupTab.tsx` to use `useMutationConfirmation` for Server Control actions.
- Removed the local one-off confirmation modal for battlegroup commands.
- Preserved read-only status, pod, health diagnostics, and support-bundle export flows.
- Added shared mutation confirmation and required admin reason capture for:
  - Start
  - Stop
  - Restart
  - Update
  - Backup
  - Restore
- Added command, namespace, disruption, backup/restore, and audit-log guidance to confirmation details.
- Passed captured reasons into `api.battlegroup.exec`.
- Preserved the existing command-running/output modal.

### Security and operator impact

- Battlegroup Exec now uses the same shared confirmation foundation as Players, Storage, and Database SQL mutations.
- Admin reason capture is required before server-control requests are sent.
- Player 360 remains read-only.
- Remaining shared-confirmation review targets include Blueprint import and future Inventory Studio/Player 360 quick-action workflows.

### Validation

Validation required in the Windows development environment:

```powershell
.\update.ps1
```

---

## Previous update: Database SQL shared mutation confirmation migration

Database SQL was migrated to shared mutation confirmation with required admin reason capture.

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
