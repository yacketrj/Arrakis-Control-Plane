# Dune Admin Release Notes

## Current update: Inventory Studio v2 snapshot comparison

### Why this update was made

Inventory Studio v2 needs safe before/after visibility before item editing workflows are added. This update adds local comparison against a previously exported inventory snapshot while keeping the feature read-only.

### What changed

- Updated `web/src/tabs/InventoryStudioTab.tsx` with snapshot comparison support.
- Added local JSON snapshot loading from a previously exported Inventory Studio snapshot.
- Added inventory diff detection for added, removed, and changed item rows.
- Added diff details for template, name, stack, quality, durability, and max durability changes.
- Added current-vs-snapshot summary counts.
- Kept comparison entirely local in the browser.

### Security and operator impact

- Inventory Studio v2 remains read-only.
- No item edit controls were added in this slice.
- Snapshot compare gives operators a safer review path before future confirmed edit workflows are added.
- Player 360 remains read-only.

### Validation

Validation required in the Windows development environment:

```powershell
.\update.ps1
```

GitHub Actions also runs Linux and Windows validation on push.

---

## Previous update: Inventory Studio v2 read-only foundation

Inventory Studio v2 was added as a read-only player inventory inspection and snapshot page.

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
