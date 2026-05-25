# Dune Admin Release Notes

## Current update: Confirmed player admin actions workflow

### Why this update was made

The remaining player admin actions modify or disrupt player state and should not run through unstructured prompts. This update isolates them in a dedicated confirmed admin-actions modal with shared mutation confirmation and required admin reason capture.

### What changed

- Added `web/src/tabs/PlayerAdminActionsModal.tsx` as a dedicated confirmed player admin-actions modal.
- Added a Players-table `Admin` launcher through `PlayersTabWith360Launcher.tsx`.
- Added shared mutation confirmation and required admin reason capture for:
  - Clear all journey progress
  - Remove tutorial records
  - Clear codex discoveries
  - Disconnect player session
- Added player target metadata, current map, online state, operator warnings, and per-action support details to the confirmation flow.
- Passed captured reasons into the matching `api.players.*` mutation calls.

### Security and operator impact

- The active Players-table mutation surface now routes Give Item, Inventory repair/delete, resource/spec actions, journey node actions, player move, and admin actions through shared confirmation and reason capture.
- Player 360 remains read-only.
- No Player 360 quick actions were added.
- `PlayersTab.tsx` still contains legacy inline modal code as cleanup debt, but the active wrapper path now routes these workflows through confirmed extracted modals.

### Validation

Validation required in the Windows development environment:

```powershell
.\update.ps1
```

---

## Previous update: Confirmed player move workflow

Player move actions were migrated to a dedicated confirmed modal with required admin reason capture and online-state safeguards.

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
