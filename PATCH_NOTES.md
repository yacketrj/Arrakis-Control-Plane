# Dune Admin Release Notes

## Current update: Inventory Studio v2 read-only foundation

### Why this update was made

Inventory Studio v2 is the next P1 operator-support feature after Player 360 and the shared confirmation work. This update starts Inventory Studio as a read-only inspection and snapshot page.

### What changed

- Added `web/src/tabs/InventoryStudioTab.tsx`.
- Added an `Inventory Studio` tab to app navigation.
- Added player search and selection.
- Added inventory loading for the selected player.
- Added inventory filtering by item template, item name, item ID, and quality.
- Added selected-item read-only detail view.
- Added raw selected-item JSON inspection.
- Added inventory snapshot export as JSON.

### Security and operator impact

- Inventory Studio v2 starts with visibility and snapshot export only.
- No item edit controls were added in this slice.
- Player 360 remains read-only.

### Validation

Validation required in the Windows development environment:

```powershell
.\update.ps1
```

GitHub Actions also runs Linux and Windows validation on push.

---

## Previous update: GitHub CI validation workflows

GitHub-hosted Linux and Windows validation workflows were added for push, pull request, and manual dispatch.

---

## Previous update: Blueprint import shared mutation confirmation migration

Blueprint import was migrated to shared mutation confirmation with required admin reason capture.

---

## Previous update: Battlegroup Exec shared mutation confirmation migration

Battlegroup Exec server-control actions were migrated to shared mutation confirmation with required admin reason capture.

---

## Previous update: Database SQL shared mutation confirmation migration

Database SQL was migrated to shared mutation confirmation with required admin reason capture.

---

## Previous update: Storage shared mutation confirmation migration

Storage container add/remove item operations were migrated to shared mutation confirmation with required admin reason capture.
