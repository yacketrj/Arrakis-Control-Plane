# Dune Admin Release Notes

## Current update: Blueprint import shared mutation confirmation migration

### Why this update was made

Blueprint import changes player construction data and should not rely on one-off browser prompts. This update moves blueprint import onto the shared mutation-confirmation hook with required admin reason capture.

### What changed

- Updated `web/src/tabs/BlueprintsTab.tsx` to use `useMutationConfirmation` for Import Blueprint.
- Removed the local `window.prompt` and `window.confirm` sequence from blueprint import.
- Preserved blueprint list and export behavior as read-only flows.
- Added shared mutation confirmation and required admin reason capture before `api.blueprints.import` is called.
- Added player ID, file name, file size, audit-log guidance, and target verification details to the confirmation flow.
- Passed captured reasons into `api.blueprints.import`.

### Security and operator impact

- Blueprint import now uses the same shared confirmation foundation as Players, Storage, Database SQL, and Battlegroup Exec mutations.
- Admin reason capture is required before blueprint import requests are sent.
- Player 360 remains read-only.
- The explicit current mutation-safety review targets are now covered. Future Inventory Studio v2 workflows and Player 360 quick actions must still be added only as confirmed workflows.

### Validation

Validation required in the Windows development environment:

```powershell
.\update.ps1
```

---

## Previous update: Battlegroup Exec shared mutation confirmation migration

Battlegroup Exec server-control actions were migrated to shared mutation confirmation with required admin reason capture.

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
