# Dune Admin Release Notes

## Current update: Inventory Studio v2 post-action diff panel

### Why this update was made

Inventory Studio v2 now has confirmed add, repair, and removal workflows. Operators need an immediate view of what changed after a confirmed inventory action completes. This update adds a post-action diff panel that compares the in-memory before-action inventory state against the reloaded inventory after mutation.

### What changed

- Updated `web/src/tabs/InventoryStudioTab.tsx` with a post-action diff panel.
- Retained the in-memory before-action inventory list for add, repair, and removal operations.
- Compared the before-action inventory against the reloaded inventory after successful mutation.
- Added post-action summary fields:
  - action
  - target
  - before item count
  - after item count
  - diff count
  - checked timestamp
- Reused the existing diff rendering for added, removed, and changed item rows.
- Refactored the Inventory Studio component into smaller local helper components to reduce future edit risk.

### Security and operator impact

- Confirmed Inventory Studio actions now provide immediate before/after review in the UI.
- Local before-action snapshot export remains in place before each mutation request is sent.
- Shared mutation confirmation and required admin reason capture remain in place.
- Player 360 remains read-only.

### Validation

Validation required in the Windows development environment:

```powershell
.\update.ps1
```

GitHub Actions also runs Linux and Windows validation on push.

---

## Previous update: Inventory Studio v2 confirmed catalog item add

Inventory Studio v2 added confirmed catalog-item add with quantity and quality inputs.

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
