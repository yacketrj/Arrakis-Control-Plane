# Dune Admin Release Notes

## Current update: Inventory repair/delete shared mutation confirmation migration

### Why this update was made

Inventory repair and delete actions are high-impact player inventory mutations. This update routes the active Players inventory workflow through an extracted inventory modal that uses the shared mutation-confirmation hook and required admin reason capture.

### What changed

- Added `web/src/tabs/InventoryModal.tsx` as an extracted inventory modal.
- Wired inventory repair/delete through `useMutationConfirmation`.
- Added required admin reason capture for repair and delete operations.
- Passed captured reasons into `api.players.repairItem` and `api.players.deleteItem`.
- Added player/item target metadata to confirmation dialogs.
- Preserved online-player desync warnings in confirmation details.
- Updated `web/src/tabs/PlayersTabWith360Launcher.tsx` so active Players-table Inventory clicks open the extracted confirmed modal.
- Avoided direct large-file replacement of `web/src/tabs/PlayersTab.tsx` in this slice.

### Security and operator impact

- Inventory repair/delete now use the same shared confirmation foundation as Give Item.
- Admin reason capture is required before inventory mutation requests are sent.
- Player 360 remains read-only.
- No Player 360 quick actions were added.
- The old inline `InventoryModal` remains in `PlayersTab.tsx` as cleanup debt, but active app routing now uses the extracted confirmed modal through the wrapper.
- Remaining PlayersTab action migrations are resource grants, XP/spec changes, journey reset/wipe, admin wipes, kick, and teleport.

### Validation

Validated in the Windows development environment with:

```powershell
.\update.ps1
```

---

## Previous update: Give Item shared mutation confirmation migration

The Give Item modal was migrated from local browser confirm/prompt handling to the shared frontend mutation-confirmation hook with required admin reason capture.

---

## Previous update: Shared frontend mutation confirmation foundation

The reusable frontend mutation-confirmation hook was added so future high-risk UI actions can display mutation safety metadata and capture an admin reason before sending the mutation request.

---

## Previous update: Player 360 validated read-only profile

Player 360 v1 was validated as a protected read-only profile with a standalone tab and Players-table launcher.

---

## Previous update: Player 360 launcher from Players table

Added the read-only Players-table `360` launcher and Player 360 auto-load behavior.

---

## Previous update: Player 360 read-only frontend tab

Added the standalone read-only Player 360 frontend tab and navigation entry.

---

## Previous update: Player 360 backend profile foundation

Added the protected read-only Player 360 backend profile endpoint, route registration, helper tests, and backend documentation.

---

## Previous update: Player 360 roadmap and design foundation

The roadmap and design documents were updated so Player 360 is the next P1 read-only implementation slice after the P0 safety foundation.

---

## Previous update: Admin audit and mutation-safety documentation sync

Documentation and tracking were synced with the landed audit foundation and in-progress mutation-safety foundation.

---

## Previous update: SSH tunnel management foundation

Managed SSH tunnel behavior was added for protected game-management database access.

---

## Previous update: Battlegroup Health Diagnostics

Read-only Battlegroup Health Diagnostics and support-bundle export were added for operator troubleshooting.
