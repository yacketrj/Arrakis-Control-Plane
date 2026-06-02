# Inventory Studio v2

## Purpose

Inventory Studio v2 is the protected operator surface for inspecting and safely adjusting player inventory state.

Every write path must preserve the same safety pattern:

1. Operator selects a specific player and target.
2. Operator reviews current inventory state.
3. A local before-action inventory snapshot is exported before mutation.
4. Shared mutation confirmation is displayed.
5. An admin reason is required.
6. The mutation request includes the captured reason.
7. Inventory is reloaded after success.
8. A post-action diff is displayed in the UI.
9. The browser-session action history records the completed action diff.
10. Audit visibility remains available through the admin audit log.

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
- Supports search by player name, actor ID, account ID, controller ID, class, and map.

### Inventory inspection

- Loads inventory rows for the selected player.
- Supports inventory filtering by item template, item display name, item row ID, and quality.
- Displays selected-item details, including raw JSON inspection.

### Snapshot export

Operators can export the currently loaded inventory as a local JSON snapshot.

Snapshot exports include:

- export timestamp
- selected player metadata
- item count
- full inventory row list

### Snapshot comparison

Operators can load a previously exported Inventory Studio snapshot and compare it against the currently loaded inventory.

The comparison detects added, removed, and changed item rows. Changed field detection covers template ID, name, stack size, quality, durability, and max durability.

Comparison is local in the browser. It does not upload the prior snapshot to the backend.

### Item catalog browser

Inventory Studio loads item templates from the existing protected templates endpoint.

The catalog supports refresh, search by template ID or display name, and selected-template detail display.

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
- post-action diff display
- action history capture

### Confirmed item repair

Inventory Studio supports repairing the selected inventory item.

Safety behavior:

- local before-action inventory snapshot export
- shared mutation confirmation
- required admin reason capture
- reason passed into `api.players.repairItem`
- inventory reload after success
- post-action diff display
- action history capture

### Confirmed item removal

Inventory Studio supports removing the selected inventory item.

Safety behavior:

- local before-action inventory snapshot export
- shared mutation confirmation
- required admin reason capture
- reason passed into `api.players.deleteItem`
- inventory reload after success
- post-action diff display
- action history capture

### Post-action diff panel

After a confirmed add, repair, or removal completes, Inventory Studio compares the before-action in-memory inventory list against the reloaded inventory.

The panel displays:

- action name
- target
- before item count
- after item count
- diff count
- checked timestamp
- added, removed, and changed item rows

This is browser-local review state. It is not persisted server-side.

### Action history panel

Inventory Studio keeps the most recent completed action diffs in browser memory for the selected player session.

The panel displays:

- latest action first
- action name and target
- checked timestamp
- before and after row counts
- diff count
- short preview of changed item rows

Operators can clear the browser-session action history or export it as a local JSON file. The exported history contains the selected player metadata and the recorded action diff entries.

Selecting a different player resets the local action history so the panel does not mix records across players.

This is browser-local review state. It does not replace the server-side audit log.

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
- Do not bypass post-action diff review for confirmed Inventory Studio workflows.
- Do not bypass action history capture for confirmed Inventory Studio workflows.
- Do not add Player 360 quick actions that mutate inventory until they reuse this same safety pattern.
- Prefer one narrow confirmed workflow at a time over broad inventory editing surfaces.
- Keep catalog browsing and snapshot comparison usable without requiring mutation confirmation.

## Recommended next work

1. Add stack-size edit only after before/after preview is displayed.
2. Add quality edit only after before/after preview is displayed.
3. Add server-side before-change snapshot persistence for inventory mutations.
4. Add typed backend mutation wrappers for inventory operations.
5. Add audit export and filtering improvements.

## Validation

Use the canonical local validation path:

```bash
./update.sh
```

GitHub Actions also runs Linux and Windows validation on push and pull request.
