# Dune Admin Feature Design and Priorities

## Purpose

This document defines the planned Dune Admin feature roadmap, priority order, safety model, and implementation dependencies for the expanded operator toolkit.

The goal is to turn Dune Admin into a polished server-operations console for player support, live operations, diagnostics, inventory administration, moderation support, and safe high-risk state changes.

## Important clarification: item delivery paths

There are multiple item paths, and they are not equivalent.

### Gameplay inventory path

This is what happens when the game itself gives an item through normal gameplay systems such as crafting, looting, harvesting, vendor purchase, quest reward handling, or other authoritative server-side gameplay code.

Characteristics:

- The game server owns the inventory mutation.
- The online player's live inventory state can replicate through normal game systems.
- Item rules, stack behavior, durability, grade, stats, and side effects are handled by game code.

This is the best instant path, but Dune Admin does not currently have a verified general-purpose gameplay RPC/API for arbitrary admin item injection through that same server code path.

### Direct Inventory Write

This is Dune Admin's database inventory mutation path.

Characteristics:

- Writes directly into `dune.items` and related inventory state.
- Supports item grade, stack size, and augmented item stats.
- Good for offline support and precise inventory editing.
- Online players may need logout/login before their client-visible inventory state refreshes.

### Claim Rewards Queue

This is the existing live reward/claim path used by the backend live grant function.

Characteristics:

- Queues a plain template-and-amount reward for the player to claim.
- It is useful for online-friendly delivery because the player can see Claim Rewards without a relog.
- It is not the same as crafting, looting, finding an item, or direct inventory insertion.
- It should be treated as an inbox/reward-claim workflow, not a full inventory editor.
- It currently supports plain item template plus amount only.
- It should not be used for graded items, augmented items, exact slot placement, durability editing, or custom stats.

### Operator rule of thumb

Use Direct Inventory Write for full-fidelity item creation and edits. Use Claim Rewards Queue only when a player is online and the grant is a plain item template plus amount.

## Guiding principles

1. **Safety before power**: high-impact mutations must be gated by audit logging, clear previews, validation, and rollback guidance.
2. **Read-only before write**: new domains should start with visibility and diagnostics before destructive or mutating actions.
3. **Operator clarity**: workflows must explain whether a change is instant, delayed, requires relog, or creates a claimable reward.
4. **Live-player awareness**: workflows must distinguish online-safe paths from direct database writes that may require relog.
5. **Least surprise**: button labels must describe the actual mechanism, not just the desired outcome.
6. **Composable backend commands**: feature work should prefer typed commands and handlers over raw SQL in the UI.
7. **Continuous validation**: every feature must keep `PATCH_NOTES.md`, `CHANGELOG.md`, Go tests, frontend typecheck, lint, build, and security scans current.

## Priority model

| Priority | Meaning | Rule |
|---|---|---|
| P0 | Foundation / blocker | Required before expanding high-risk operations |
| P1 | High-value core feature | Directly improves daily operations and support outcomes |
| P2 | Advanced operations | Useful and important, but depends on P0/P1 controls |
| P3 | Specialized / future | Valuable after core workflows are stable |

## Prioritized roadmap

| Rank | Priority | Feature | Primary value | Dependency |
|---:|---|---|---|---|
| 1 | P0 | Admin Action Audit Log | Makes high-risk actions accountable | None |
| 2 | P0 | Mutation Safety Framework | Shared preview, reason, audit, validation, rollback hooks | Audit Log |
| 3 | P1 | Player 360 Profile | Single support view for player identity, state, inventory, history | Existing player APIs |
| 4 | P1 | Inventory Studio v2 | Safer inventory snapshots, diffs, item edits, augment inspection | Audit + Give Item helpers |
| 5 | P1 | Broadcast Center | Template-based maintenance and live-ops communication | Notification/RMQ path |
| 6 | P1 | Safe Offline Teleport / Rescue | Stuck-player rescue with guardrails and rollback hints | Audit + partition APIs |
| 7 | P1 | Server Health Command Center | Unified operational status and diagnostic bundle | Battlegroup + DB status |
| 8 | P2 | RMQ Operations Dashboard | Queue/exchange/consumer diagnostics | Read-only RMQ helpers |
| 9 | P2 | Journey / Progression Manager | Safer quest/progression support | Audit + journey commands |
| 10 | P2 | Guild / Faction Admin | Guild/faction support tooling | Audit + DB procedures |
| 11 | P2 | Augment Preset Manager | Better augmented item creation and reverse-engineering | Inventory Studio v2 |
| 12 | P2 | Maintenance Mode Assistant | Guided downtime, broadcasts, health snapshots | Broadcast + Health Center |
| 13 | P3 | Settings Diff / Config Manager | Read-only config baseline and future safe writes | Health + RMQ format research |
| 14 | P3 | Chat Moderation / Intercept Viewer | Moderation evidence and chat diagnostics | RMQ dashboard foundation |
| 15 | P3 | Safe Stored Procedure Runner | Typed allowlisted admin procedure execution | Audit + safety framework |

## P0 foundation features

### Admin Action Audit Log

Record every high-impact action with:

- timestamp
- operator identity or auth mode
- action type
- target player/account/controller/item/guild where available
- sanitized payload summary
- success/failure result
- error message when applicable
- rollback hint when available

Do not log secrets, admin tokens, database passwords, SSH keys, or raw credential-bearing environment values.

### Mutation Safety Framework

Provide one shared backend/frontend pattern for:

- required reason for high-risk actions
- preview before mutation
- server-side validation
- audit event creation
- rollback hint generation
- consistent operator errors

## P1 core operator features

### Player 360 Profile

A single player detail page should aggregate:

- identity and IDs
- online state
- current map/partition
- inventory summary
- vehicles
- currencies
- faction status
- specializations
- journey summary
- recent events
- dungeon history
- quick actions

### Inventory Studio v2

Inventory Studio should include:

- inventory snapshot before mutation
- search/filter by item template, name, grade, stack, durability, augment presence
- item add/delete/repair
- partial stack removal
- clone item stats from existing item
- augment JSON viewer and export
- before/after diff
- clear delivery path labels: Direct Inventory Write vs Claim Rewards Queue

### Broadcast Center

Broadcast Center should include:

- reusable templates
- maintenance countdowns
- preview
- routing key selection
- send history
- audit records

### Safe Offline Teleport / Rescue Tool

This tool should include:

- current location display
- safe destination presets
- online/offline warnings
- prior-location rollback hint
- audit record

### Server Health Command Center

The health center should include:

- battlegroup overview
- map/server cards
- player counts by map/partition
- database/SSH/pod/log health
- stale heartbeat warnings
- diagnostic bundle export

## P2 advanced operations

### RMQ Operations Dashboard

Read-only first:

- overview
- nodes
- queues
- exchanges
- bindings
- consumers
- connections
- channels
- alarms
- safe queue peek with requeue

### Journey / Progression Manager

- journey node search
- complete/reset selected nodes
- reset workflows with snapshot
- codex/tutorial support
- destructive warnings

### Guild / Faction Admin

- view guild for player
- view guild roster
- promote/demote members
- disband guild with destructive guardrails
- change faction or faction tier through typed helpers

### Augment Preset Manager

- UI-managed augment preset catalog
- roll-count metadata
- default roll arrays
- import/export presets
- reverse-engineer augments from existing item stats
- clone augments from inventory item into Give Item workflow

### Maintenance Mode Assistant

- guided checklist
- broadcast countdown
- online player view
- queue and health checks
- diagnostic snapshot before/after

## P3 specialized/future work

### Settings Diff / Config Manager

Start read-only:

- current gameplay settings
- saved baseline
- diff view
- export/import desired baseline

Future write support requires fully validated message formats and audit controls.

### Chat Moderation / Intercept Viewer

Start read-only:

- intercepted chat stream when available
- filters by player/channel/map
- evidence export
- moderation notes

### Safe Stored Procedure Runner

Only allow typed, allowlisted procedures. No arbitrary procedure execution from UI.

Minimum controls:

- allowlist
- typed parameters
- risk level
- preview
- required reason
- audit event
- rollback guidance

## Implementation phases

### Phase 0: Stabilize existing workflows

- Keep Actions green.
- Resolve frontend typecheck/lint/build issues.
- Commit frontend lockfile once dependency state is stable.
- Keep `PATCH_NOTES.md` and `CHANGELOG.md` updated for every change.

### Phase 1: Safety foundation

- Admin Action Audit Log.
- Mutation Safety Framework.
- Backfill audit calls into existing mutation handlers.

### Phase 2: Operator support surface

- Player 360 Profile.
- Inventory Studio v2.
- Safe Offline Teleport / Rescue.
- Broadcast Center.

### Phase 3: Server operations

- Server Health Command Center.
- RMQ Operations Dashboard read-only.
- Maintenance Mode Assistant.

### Phase 4: Progression and organization tooling

- Journey / Progression Manager.
- Guild / Faction Admin.
- Augment Preset Manager.

### Phase 5: Specialized tools

- Settings Diff / Config Manager.
- Chat Moderation / Intercept Viewer.
- Safe Stored Procedure Runner.

## Next implementation slice

The next implementation slice should be foundational:

1. Add backend audit event model and append-only JSONL audit sink.
2. Add `recordAdminAction` helper with redaction.
3. Wrap existing mutation handlers.
4. Add read-only Audit tab.
5. Add tests for redaction and success/failure audit entries.

This unlocks safer implementation of every later feature.
