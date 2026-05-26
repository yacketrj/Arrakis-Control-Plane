# Dune Admin Release Notes

## Current update: Inventory Studio v2 confirmed catalog item add

### Why this update was made

Inventory Studio v2 now supports a controlled add workflow from the validated item catalog. This update adds confirmed catalog-item add with quantity and quality inputs, a before-action inventory snapshot, shared mutation confirmation, and required admin reason capture.

### What changed

- Updated `web/src/tabs/InventoryStudioTab.tsx` with a confirmed catalog item add action.
- Added quantity and quality inputs for the selected catalog template.
- Added client-side clamping for quantity and quality values.
- Added automatic before-action snapshot export before the add request is sent.
- Added shared mutation confirmation through `useMutationConfirmation`.
- Added required admin reason capture.
- Added selected player, online state, template, catalog name, quantity, and quality details to the confirmation flow.
- Passed the captured reason into `api.players.giveItem`.
- Reloaded the selected player inventory after successful add.

### Security and operator impact

- Inventory Studio v2 now supports confirmed add, repair, and removal workflows.
- Add uses the direct inventory write path and remains scoped to one selected player and one selected catalog template.
- A local before-action snapshot is exported before the mutation request is sent.
- Player 360 remains read-only.

### Validation

Validation required in the Windows development environment:

```powershell
.\update.ps1
```

GitHub Actions also runs Linux and Windows validation on push.

---

## Previous update: Inventory Studio v2 confirmed item removal

Inventory Studio v2 added confirmed selected-item removal with a before-action inventory snapshot.

---

## Previous update: Inventory Studio v2 confirmed item repair

Inventory Studio v2 added confirmed selected-item repair with a before-action inventory snapshot.

---

## Previous update: Inventory Studio v2 item catalog browser

Inventory Studio v2 added a read-only item catalog browser.

---

## Previous update: Inventory Studio v2 snapshot comparison

Inventory Studio v2 added local comparison against a previously exported inventory snapshot while remaining read-only.

---

## Previous update: Inventory Studio v2 read-only foundation

Inventory Studio v2 was added as a read-only player inventory inspection and snapshot page.
