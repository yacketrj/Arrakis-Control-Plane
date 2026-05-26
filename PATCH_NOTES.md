# Dune Admin Release Notes

## Current update: Inventory Studio v2 item catalog browser

### Why this update was made

Inventory Studio v2 needs a validated item catalog browsing surface before future inventory edit workflows are introduced. This update adds a read-only catalog browser inside Inventory Studio.

### What changed

- Updated `web/src/tabs/InventoryStudioTab.tsx` with item catalog browsing.
- Loaded item templates from the existing player template endpoint.
- Added catalog refresh.
- Added catalog search by template ID and display name.
- Added selected-template detail display.
- Kept the catalog browser read-only.

### Security and operator impact

- Inventory Studio v2 remains read-only.
- No item edit controls were added in this slice.
- Catalog browsing prepares the UI for future confirmed add/edit workflows without introducing new mutation paths.
- Player 360 remains read-only.

### Validation

Validation required in the Windows development environment:

```powershell
.\update.ps1
```

GitHub Actions also runs Linux and Windows validation on push.

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

---

## Previous update: Battlegroup Exec shared mutation confirmation migration

Battlegroup Exec server-control actions were migrated to shared mutation confirmation with required admin reason capture.

---

## Previous update: Database SQL shared mutation confirmation migration

Database SQL was migrated to shared mutation confirmation with required admin reason capture.
