# Dune Admin Release Notes

## Current update: Give Item shared mutation confirmation migration

### Why this update was made

The Give Item modal is a high-impact player mutation surface. It previously used local browser confirmation and prompt handling. This update moves Give Item confirmation and admin reason capture onto the shared frontend mutation-confirmation hook.

### What changed

- Updated `web/src/tabs/GiveItemModalAugmented.tsx` to use `useMutationConfirmation`.
- Removed the modal-local mutation safety classification, `window.confirm`, and `window.prompt` flow from Give Item.
- Added shared confirmation metadata for both delivery modes:
  - Inventory Write
  - Live Claim Rewards
- Added item-row count, delivery-mode context, and player target identifiers to the shared confirmation dialog.
- Kept admin reason capture mandatory for Give Item operations.
- Passed the captured reason into `api.players.giveItems` and `api.players.grantLive`.

### Security and operator impact

- Give Item now uses the same frontend confirmation foundation planned for other high-risk mutations.
- Player 360 remains read-only.
- No new Player 360 quick actions were added.
- Existing PlayersTab action migrations remain pending for inventory repair/delete, resource grants, XP/spec changes, journey reset/wipe, admin wipes, kick, and teleport.

### Validation

Validation required in the Windows development environment:

```powershell
.\update.ps1
```

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
