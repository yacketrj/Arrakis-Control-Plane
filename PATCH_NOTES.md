# Dune Admin Release Notes

## Current update: Confirmed player move workflow

### Why this update was made

Player movement is a high-impact support action that can create location/state drift if used while the player is online. This update isolates the move workflow in a dedicated confirmed modal so the action uses shared mutation confirmation and required admin reason capture.

### What changed

- Added `web/src/tabs/PlayerTeleportModal.tsx` as a dedicated confirmed player move modal.
- Added destination loading from the existing partitions endpoint.
- Added a Players-table `Move` launcher through `PlayersTabWith360Launcher.tsx`.
- Added shared mutation confirmation and required admin reason capture before move requests are sent.
- Added player target metadata, current map, online state, destination, and drift warning details to the confirmation flow.
- Passed captured reasons into `api.players.teleport`.
- Disabled move submission while the player is online.

### Security and operator impact

- Active player move actions now use the same shared confirmation foundation as Give Item, Inventory repair/delete, resource/spec actions, and journey node actions.
- Admin reason capture is required before move requests are sent.
- Player 360 remains read-only.
- No Player 360 quick actions were added.
- Remaining action migrations are journey wipe, tutorial deletion, codex wipe, and kick.

### Validation

Validation required in the Windows development environment:

```powershell
.\update.ps1
```

---

## Previous update: Journey node shared mutation confirmation migration

Journey node complete/reset actions were migrated to the confirmed Player Actions modal with required admin reason capture.

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
