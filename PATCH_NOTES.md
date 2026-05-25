# Dune Admin Release Notes

## Current update: Player resource/action shared mutation confirmation migration

### Why this update was made

The active Players Actions workflow still allowed high-impact player resource, XP, specialization, and faction reputation mutations to flow through local action handlers. This update adds a confirmed Player Actions modal for the non-destructive/high-impact player action slice and routes the active Actions button through that modal.

### What changed

- Added `web/src/tabs/PlayerActionsModalConfirmed.tsx` for confirmed player resource and specialization actions.
- Wired the active Players-table Actions button through `PlayersTabWith360Launcher.tsx`.
- Added shared mutation confirmation and required admin reason capture for:
  - Give Currency
  - Give Scrip
  - Award Intel
  - Award Character XP
  - Give Faction Reputation
  - Set Specialization XP
- Added player target metadata and online-state warnings to confirmation details.
- Passed captured reasons into the matching `api.players.*` mutation calls.
- Avoided direct large-file replacement of `web/src/tabs/PlayersTab.tsx` in this slice.

### Security and operator impact

- Active resource, XP, specialization, and faction reputation actions now use the same shared confirmation foundation as Give Item and Inventory repair/delete.
- Admin reason capture is required before these mutation requests are sent.
- Player 360 remains read-only.
- No Player 360 quick actions were added.
- Remaining action migrations are journey complete/reset/wipe, tutorial deletion, codex wipe, kick, and teleport.

### Validation

Validation required in the Windows development environment:

```powershell
.\update.ps1
```

---

## Previous update: Inventory repair/delete shared mutation confirmation migration

Inventory repair/delete was migrated to the extracted confirmed inventory modal with required admin reason capture.

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
