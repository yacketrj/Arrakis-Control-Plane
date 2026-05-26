# Dune Admin Release Notes

## Current update: Inventory Studio v2 confirmed item repair

### Why this update was made

Inventory Studio v2 now has enough read-only context for its first narrow edit workflow. This update adds confirmed selected-item repair with a before-action inventory snapshot, shared mutation confirmation, and required admin reason capture.

### What changed

- Updated `web/src/tabs/InventoryStudioTab.tsx` with a confirmed selected-item repair action.
- Added automatic before-action snapshot export before the repair request is sent.
- Added shared mutation confirmation through `useMutationConfirmation`.
- Added required admin reason capture.
- Added selected player, online state, item ID, template, and durability details to the confirmation flow.
- Passed the captured reason into `api.players.repairItem`.
- Reloaded the selected player inventory after successful repair.

### Security and operator impact

- Inventory Studio v2 now has its first confirmed edit workflow.
- Repair remains narrow: it only targets the selected inventory item.
- A local before-action snapshot is exported before the mutation request is sent.
- Player 360 remains read-only.

### Validation

Validation required in the Windows development environment:

```powershell
.\update.ps1
```

GitHub Actions also runs Linux and Windows validation on push.

---

## Previous update: Inventory Studio v2 item catalog browser

Inventory Studio v2 added a read-only item catalog browser.

---

## Previous update: Inventory Studio v2 snapshot comparison

Inventory Studio v2 added local comparison against a previously exported inventory snapshot while remaining read-only.

---

## Previous update: Inventory Studio v2 read-only foundation

Inventory Studio v2 was added as a read-only player inventory inspection and snapshot page.

---

## Previous update: GitHub CI validation workflows

GitHub-hosted Linux and Windows validation workflows were added for push, pull request, and manual dispatch.

---

## Previous update: Blueprint import shared mutation confirmation migration

Blueprint import was migrated to shared mutation confirmation with required admin reason capture.
