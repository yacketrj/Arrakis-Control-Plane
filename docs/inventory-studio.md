# Inventory Studio v2

## Purpose

Inventory Studio v2 is the protected operator surface for inspecting and safely adjusting player inventory state.

The feature started as a read-only inspection tool and now includes narrow confirmed item workflows. Every write path must preserve the same safety pattern:

1. Operator selects a specific player and target.
2. Operator reviews current inventory state.
3. A local before-action inventory snapshot is exported before mutation.
4. Shared mutation confirmation is displayed.
5. An admin reason is required.
6. The mutation request includes the captured reason.
7. Inventory is reloaded after success.
8. Audit visibility remains available through the admin audit log.

## Current UI location

Inventory Studio is available from the main app navigation as:

```text
Inventory Studio
```

Implementation file:

```text
web/src/tabs/InventoryStudioTab.tsx
```

## Current capabilities

### Player selection

- Loads players from the existing protected player list endpoint.
- Supports search by:
  - player name
  - actor ID
  - account ID
  - controller ID
  - class
  - map

### Inventory inspection

- Loads inventory rows for the selected player.
- Supports inventory filtering by:
  - item template
  - item display name
  - item row ID
  - quality
- Displays selected-item details:
  - item ID
  - template ID
  - stack size
  - quality
  - durability
  - max durability
  - raw JSON inspection

### Snapshot export

Operators can export the currently loaded inventory as a local JSON snapshot.

Snapshot exports include:

- export timestamp
- selected player metadata
- item count
- full inventory row list

### Snapshot comparison

Operators can load a previously exported Inventory Studio snapshot and compare it against the currently loaded inventory.

The comparison detects:

- added item rows
- removed item rows
- changed item rows

Changed field detection currently covers:

- template ID
- name
- stack size
- quality
- durability
- max durability

Comparison is local in the browser. It does not upload the prior snapshot to the backend.

### Item catalog browser

Inventory Studio loads item templates from the existing protected templates endpoint.

The catalog supports:

- refresh
- search by template ID
- search by display name
- selected-template detail display

### Confirmed item add

Inventory Studio supports adding a selected catalog template to the selected player's inventory.

Inputs:

- selected player
- selected catalog template
- quantity
- quality

Safety behavior:

- client-side quantity and quality clamping
- local before-action inventory snapshot export
- shared mutation confirmation
- required admin reason capture
- reason passed into `api.players.giveItem`
- inventory reload after success

### Confirmed item repair

Inventory Studio supports repairing the selected inventory item.

Safety behavior:

- local before-action inventory snapshot export
- shared mutation confirmation
- required admin reason capture
- reason passed into `api.players.repairItem`
- inventory reload after success

### Confirmed item removal

Inventory Studio supports removing the selected inventory item.

Safety behavior:

- local before-action inventory snapshot export
- shared mutation confirmation
- required admin reason capture
- reason passed into `api.players.deleteItem`
- inventory reload after success

## Current non-goals

The current implementation does not yet include:

- stack-size editing for an existing row
- quality editing for an existing row
- template replacement for an existing row
- augment/stat editing
- server-side before-change snapshot persistence
- named operator identity beyond the current admin-token model
- batch inventory operations
- rollback automation

## Safety rules

- Do not add inventory writes without shared mutation confirmation.
- Do not send inventory writes without an admin reason.
- Do not bypass the before-action snapshot behavior for item add, repair, removal, or future edit workflows.
- Do not add Player 360 quick actions that mutate inventory until they reuse this same safety pattern.
- Prefer one narrow confirmed workflow at a time over broad inventory editing surfaces.
- Keep catalog browsing and snapshot comparison usable without requiring mutation confirmation.

## Recommended next work

1. Add automatic post-action comparison by retaining the before-action snapshot in UI state and comparing it against the reloaded inventory.
2. Add a dedicated Inventory Studio action history panel for the current browser session.
3. Add stack-size edit only after before/after preview is displayed.
4. Add quality edit only after before/after preview is displayed.
5. Add server-side before-change snapshot persistence for inventory mutations.
6. Add typed backend mutation wrappers for inventory operations.
7. Add audit export and filtering improvements.

## Validation

Use the local Windows validation path:

```powershell
.\update.ps1
```

GitHub Actions also runs Linux and Windows validation on push and pull request.