# Dune Admin Release Notes

## Current update: Journey node shared mutation confirmation migration

### Why this update was made

Journey node complete/reset actions are high-impact player progression mutations. This update extends the confirmed Player Actions modal with a Journey section so node-level journey mutations use shared confirmation and required admin reason capture.

### What changed

- Extended `web/src/tabs/PlayerActionsModalConfirmed.tsx` with a Journey section.
- Added journey node loading and filtering to the confirmed modal.
- Added shared mutation confirmation and required admin reason capture for:
  - Journey node complete
  - Journey node reset
- Added player target metadata, online-state context, and journey node identifiers to confirmation details.
- Passed captured reasons into `api.players.journeyComplete` and `api.players.journeyReset`.
- Kept bulk journey wipe and admin actions out of this slice so destructive operations can be migrated separately.

### Security and operator impact

- Active journey node complete/reset actions now use the same shared confirmation foundation as Give Item, Inventory repair/delete, and resource/spec actions.
- Admin reason capture is required before node-level journey mutation requests are sent.
- Player 360 remains read-only.
- No Player 360 quick actions were added.
- Remaining action migrations are journey wipe, tutorial deletion, codex wipe, kick, and teleport.

### Validation

Validation required in the Windows development environment:

```powershell
.\update.ps1
```

---

## Previous update: Player resource/action shared mutation confirmation migration

Resource, XP, specialization, and faction reputation actions were migrated to the confirmed Player Actions modal with required admin reason capture.

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
