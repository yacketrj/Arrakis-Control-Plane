# Dune Admin Release Notes

## Current update: Storage shared mutation confirmation migration

### Why this update was made

Storage container item add/remove operations mutate persisted server state and require a server zone restart before visibility is consistent for other players. This update moves the active Storage tab mutation flows onto the shared mutation-confirmation hook with required admin reason capture.

### What changed

- Updated `web/src/tabs/StorageTab.tsx` to use `useMutationConfirmation`.
- Added shared mutation confirmation and required admin reason capture for:
  - Add item to storage container
  - Remove item from storage container
- Added storage container target metadata, item identifiers, container type, map, stack size, quality, and restart visibility warning details to the confirmation flow.
- Passed captured reasons into `api.storage.giveItem` and `api.players.deleteItem`.
- Preserved existing storage warning that item add/remove operations require a server zone restart before other players see the change.

### Security and operator impact

- Storage add/remove mutations now use the same shared confirmation foundation as active Players-table mutations.
- Admin reason capture is required before storage mutation requests are sent.
- Player 360 remains read-only.
- Remaining shared-confirmation review targets include Database SQL, Battlegroup Exec, Blueprint import, and future Inventory Studio/Player 360 quick-action workflows.

### Validation

Validation required in the Windows development environment:

```powershell
.\update.ps1
```

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

---

## Previous update: Player 360 validated read-only profile

Player 360 v1 was validated as a protected read-only profile with a standalone tab and Players-table launcher.
