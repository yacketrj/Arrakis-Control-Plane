# Dune Admin Release Notes

## Current update: Inventory Studio v2 confirmed item removal

### Why this update was made

Inventory Studio v2 now supports the second narrow confirmed edit workflow. This update adds confirmed selected-item removal with a before-action inventory snapshot, shared mutation confirmation, and required admin reason capture.

### What changed

- Updated `web/src/tabs/InventoryStudioTab.tsx` with a confirmed selected-item removal action.
- Added automatic before-action snapshot export before the delete request is sent.
- Added shared mutation confirmation through `useMutationConfirmation`.
- Added required admin reason capture.
- Added selected player, online state, item ID, template, stack size, and quality details to the confirmation flow.
- Passed the captured reason into `api.players.deleteItem`.
- Reloaded the selected player inventory after successful removal.

### Security and operator impact

- Inventory Studio v2 now supports confirmed repair and confirmed removal for the selected item.
- Removal is destructive and remains scoped to one selected inventory row.
- A local before-action snapshot is exported before the mutation request is sent.
- Player 360 remains read-only.

### Validation

Validation required in the Windows development environment:

```powershell
.\update.ps1
```

GitHub Actions also runs Linux and Windows validation on push.

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

---

## Previous update: GitHub CI validation workflows

GitHub-hosted Linux and Windows validation workflows were added for push, pull request, and manual dispatch.
