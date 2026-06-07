# Arrakis Control Panel Feature Design and Priorities

## Purpose

This document defines the current feature roadmap, priority order, safety model, and implementation dependencies for Arrakis Control Panel.

It replaces older Dune Admin-era planning language. Historical upstream references to `dune-admin` remain valid only when referring to Icehunter's upstream project, legacy executable/service compatibility, or migration notes.

## Current roadmap position

Arrakis Control Panel is in the `v0.1.x` secure-baseline release track.

Current status:

- `v0.1.0-rc.1` exists as the first controlled release candidate.
- The canonical validation/build workflow is `./update.sh`.
- Product naming has been migrated to Arrakis Control Panel across active release governance.
- Build outputs now use `arrakis-control-panel` / `arrakis-control-panel.exe`.
- Linux systemd defaults now use `arrakis-control-panel.service`, `/opt/arrakis-control-panel`, and the `arrakis-control-panel` service user/group.
- Player 360 remains read-only.
- AppSec hardening includes audit coverage, blocked-mutation coverage, and route-specific audit target metadata coverage.
- Full Discord server management, Live Admin RMQ, Welcome Kits, and arbitrary raw command publishing remain out of scope for `v0.1.0`.

## Release-train priorities

| Target | Theme | Goal | Status |
|---|---|---|---|
| `v0.1.0` | Secure baseline | Accept the current hardened admin foundation after RC verification. | Current track |
| `v0.2.0` | Discord Admin Foundation | Add Discord RBAC, command registry, previews, reason capture, and audit without raw command execution. | Planned |
| `v0.3.0` | Server Lifecycle Management | Expose allowlisted Docker/K8s/local lifecycle operations through UI/Discord. | Planned |
| `v0.4.0` | Live Admin / RMQ Foundation | Add controlled RMQ envelope/command execution for online-safe operations only. | Planned |
| `v0.5.0` | Welcome Kits / Player Requests | Add audited kits, eligibility, ledger, and Discord approval workflows. | Planned |
| `v0.6.0` | Guild Operations | Add audited guild create/delete, membership, and rank workflows. | Planned |

## Guiding principles

1. **Safety before power**: high-impact mutations require audit logging, previews, validation, and rollback guidance.
2. **Read-only before write**: new domains start with visibility and diagnostics before mutating actions.
3. **Operator clarity**: workflows must state whether changes are instant, delayed, claim-based, or require relog.
4. **Live-player awareness**: workflows must distinguish online-safe actions from direct database writes.
5. **Least surprise**: button labels must describe the actual mechanism, not just the desired outcome.
6. **Typed operations only**: prefer typed handlers and allowlisted command adapters over raw SQL, raw shell, or raw RMQ access.
7. **Continuous validation**: every durable change must update `PATCH_NOTES.md`, `CHANGELOG.md`, relevant docs, tests, and validation notes.

## Priority model

| Priority | Meaning | Rule |
|---|---|---|
| P0 | Foundation / blocker | Required before expanding high-risk operations. |
| P1 | High-value core feature | Daily operator support surface with strong safety controls. |
| P2 | Advanced operations | Useful and important, but depends on P0/P1 controls. |
| P3 | Specialized / future | Defer until core workflows are stable. |
| P5 | Documentation / enablement backlog | Important support docs, training, or operator runbooks. |

## Current baseline scope: `v0.1.0`

`v0.1.0` should finalize the secure baseline without adding major new feature scope.

Included or expected in the baseline:

- localhost-default backend posture
- strict admin-token behavior
- canonical `./update.sh` validation/build workflow
- release checklist and release-deviation tracking
- compact changelog and per-slice records
- upstream attribution preservation
- read-only Player 360 foundation
- Discord auth/session/player-link groundwork
- inventory request/order coordination
- AppSec audit and mutation-safety foundations
- route auth-boundary coverage
- high-risk/destructive mutation audit-event coverage
- blocked-mutation audit coverage
- route-specific audit target metadata coverage

Out of scope for `v0.1.0`:

- Discord-driven full server management
- Live Admin RMQ command execution
- Welcome Kits
- raw command publishing UI
- Player 360 mutation/self-service expansion
- formal SOC 2, ISO, FedRAMP, or NIST compliance claims

## High-risk feature tracks

### Discord Admin Foundation — target `v0.2.0`

Goal: allow Discord to request approved admin actions safely.

Required controls:

- Discord RBAC
- command registry
- command allowlist
- reason capture
- preview/dry-run support
- audit events
- clear operator feedback
- no raw shell, raw SQL, or raw RMQ command exposure

### Server Lifecycle Management — target `v0.3.0`

Goal: expose safe server lifecycle operations through typed workflows.

Candidate workflows:

- status
- start/stop/restart
- update
- backup
- restore guidance
- health checks
- Docker/K8s/local provider detection
- approval gates for destructive operations

Required controls:

- command allowlists
- provider-specific validation
- reason/preview/audit
- backup-before-dangerous-action where applicable
- rollback/runbook links

### Live Admin / RMQ Foundation — target `v0.4.0`

Goal: add controlled live-server actions without exposing arbitrary raw command publishing.

Required modules remain planned until implemented:

- `rmq_commands.go`
- `rmq_envelope.go`
- `rmq_player_identity.go`
- `handlers_live_admin.go`
- frontend API/client support for live admin actions

Core rule:

- Use RMQ only for online-safe, allowlisted operations.
- Use DB fallback only when required and guarded.
- Audit every attempt.
- Do not expose arbitrary raw RMQ publishing in the UI or Discord.

### Welcome Kits / Player Requests — target `v0.5.0`

Goal: provide structured, auditable starter/support kits and player request workflows.

Required controls:

- kit CRUD
- activation/deactivation
- eligibility rules
- one-time claim controls
- ledger
- Discord approval path
- audit and rollback notes

### Guild Operations — target `v0.6.0`

Goal: support guild administration without bypassing safety controls.

Requested capabilities:

- create/delete guild
- add/remove player from guild
- promote/demote or change rank
- Player tab guild workflows
- Discord command support for approved guild operations

Required controls:

- reason/preview/audit
- route-specific target metadata
- destructive guardrails
- backup requirements for DB fallback mutations where applicable

## Item delivery paths

There are multiple item paths, and they are not equivalent.

### Gameplay inventory path

This is what happens when the game itself gives an item through normal gameplay systems such as crafting, looting, harvesting, vendor purchase, or quest reward handling.

Arrakis Control Panel does not currently have a verified general-purpose gameplay RPC/API for arbitrary admin item injection through that same server code path.

### Direct Inventory Write

This is the database inventory mutation path.

Characteristics:

- writes directly into inventory-related database state
- supports higher-fidelity item creation and editing when implemented safely
- may require online players to logout/login before client-visible state refreshes
- must be guarded by preview, reason, audit, and rollback guidance

### Claim Rewards Queue

This is the live reward/claim path used by backend live grant behavior.

Characteristics:

- queues a plain template-and-amount reward for the player to claim
- useful for online-friendly delivery when supported
- not the same as crafting, looting, or exact inventory insertion
- should not be used for graded items, augmented items, exact slot placement, durability editing, or custom stats

### Operator rule of thumb

Use Direct Inventory Write for full-fidelity item creation and edits. Use Claim Rewards Queue only when the player is online and the grant is a plain item template plus amount.

## Documentation and validation requirements

Every implementation slice must:

- update `PATCH_NOTES.md`
- update `CHANGELOG.md`
- update relevant docs
- update or add tests when code changes
- run `./update.sh`
- record validation status
- add a deviation entry when scope, risk acceptance, naming, or release plan changes

## Known documentation follow-up

This document was refreshed because it triggered a large Markdown warning and contained stale roadmap/product naming language.

Remaining documentation work is tracked in:

```text
docs/documentation-review-plan.md
```
