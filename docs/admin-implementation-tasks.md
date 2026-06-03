# Dune Admin Implementation Task Tracker

## Purpose

This tracker turns the feature roadmap into implementation tasks with status, documentation requirements, validation requirements, and release-note requirements.

Every feature task must update:

- `PATCH_NOTES.md`
- `CHANGELOG.md`
- the relevant design or operating document under `docs/`
- tests or validation notes

## Status legend

| Status | Meaning |
|---|---|
| Done | Code or documentation has landed on `main` |
| In Progress | Work has started but needs validation or follow-up commits |
| Next | Highest priority work not yet started |
| Planned | Accepted roadmap item awaiting earlier dependencies |
| Blocked | Needs research, schema discovery, or a green quality gate first |

## Current task board

| Priority | Task | Status | Documentation | Validation |
|---|---|---|---|---|
| P0 | Feature design and priorities | Done | `docs/admin-feature-design-and-priorities.md` | Review only |
| P0 | Item delivery architecture clarification | Done | `docs/admin-feature-design-and-priorities.md` | Review only |
| P0 | Public-safe vs protected admin portal design | In Progress | `docs/portal-separation-design.md` | Go + frontend validation |
| P0 | Admin Action Audit Log | Done | `docs/admin-audit-log.md` | Go tests |
| P0 | Mutation Safety Framework | In Progress | `docs/mutation-safety-framework.md` | Go + frontend tests |
| P0 | Comprehensive AppSec endpoint audit | In Progress | `docs/appsec-endpoint-audit.md` | Endpoint inventory + SAST/DAST/manual abuse-case review |
| P0 | Active Players-table mutation confirmation migration | Done | `docs/mutation-safety-framework.md` | Windows `update.ps1` validation |
| P0 | PlayersTab inline modal cleanup | Done | `docs/mutation-safety-framework.md` | Windows `update.ps1` validation |
| P0 | Storage mutation confirmation migration | Done | `docs/mutation-safety-framework.md` | Windows `update.ps1` validation |
| P0 | Database SQL confirmation migration | Done | `docs/mutation-safety-framework.md` | Windows `update.ps1` validation |
| P0 | Battlegroup Exec confirmation migration | Done | `docs/mutation-safety-framework.md` | Windows `update.ps1` validation |
| P0 | Blueprint import confirmation migration | Done | `docs/mutation-safety-framework.md` | Windows `update.ps1` validation |
| P1 | Player 360 Profile | Done | `docs/player-360-profile.md` | Go + frontend validation clean |
| P1 | Inventory Studio v2 read-only foundation | Done | `docs/inventory-studio.md` | Windows `update.ps1` validation |
| P1 | Inventory Studio v2 snapshot compare | Done | `docs/inventory-studio.md` | Windows `update.ps1` validation |
| P1 | Inventory Studio v2 item catalog browser | Done | `docs/inventory-studio.md` | Windows `update.ps1` validation |
| P1 | Inventory Studio v2 confirmed add/repair/remove workflows | Done | `docs/inventory-studio.md` | Windows `update.ps1` validation |
| P1 | Inventory Studio v2 post-action diff panel | Done | `docs/inventory-studio.md` | Windows `update.ps1` validation |
| P1 | Inventory Studio v2 browser-session action history | Done | `docs/inventory-studio.md` | `./update.sh` validation |
| P1 | Inventory Studio v2 stack-size edit workflow | In Progress | `docs/inventory-studio.md` | `./update.sh` required |
| P1 | Inventory Studio v2 quality edit workflow | Next | `docs/inventory-studio.md` | Go + frontend tests |
| P1 | Battlegroup Status v2 | Planned | `docs/battlegroup-status-v2.md` required | Go + frontend tests |
| P1 | Broadcast Center | Planned | `docs/broadcast-center.md` required | Go + frontend tests |
| P1 | Safe Offline Teleport / Rescue | Planned | `docs/safe-teleport-rescue.md` required | Go + frontend tests |
| P1 | Server Health Command Center | Planned | `docs/server-health-command-center.md` required | Go + frontend tests |
| P2 | RMQ Operations Dashboard | Planned | `docs/rmq-operations-dashboard.md` required | Go + frontend tests |
| P2 | Journey / Progression Manager | Planned | `docs/journey-progression-manager.md` required | Go + frontend tests |
| P2 | Guild / Faction Admin | Planned | `docs/guild-faction-admin.md` required | Go + frontend tests |
| P2 | Guild Management: create/delete guild, membership, ranks | Planned | `docs/guild-management.md` required | Go + frontend tests |
| P2 | Player tab guild workflows: add/remove/promote | Planned | `docs/guild-management.md` required | Go + frontend tests |
| P2 | Augment Preset Manager | Planned | `docs/augment-preset-manager.md` required | Go + frontend tests |
| P2 | Maintenance Mode Assistant | Planned | `docs/maintenance-mode-assistant.md` required | Go + frontend tests |
| P3 | Settings Diff / Config Manager | Planned | `docs/settings-diff-config-manager.md` required | Research + tests |
| P3 | Chat Moderation / Intercept Viewer | Planned | `docs/chat-moderation-intercept-viewer.md` required | Research + tests |
| P3 | Safe Stored Procedure Runner | Planned | `docs/safe-stored-procedure-runner.md` required | Requires typed wrappers and snapshots |
| P5 | Discord bot setup and usage guide | Planned | `docs/discord-bot-setup-and-usage.md` required | Documentation review |

## Current implementation focus

### 1. Inventory Studio v2

Inventory Studio v2 now includes:

- Player search and selection.
- Inventory loading for the selected player.
- Inventory filtering.
- Selected-item details and raw JSON inspection.
- Snapshot export.
- Local comparison against a previously exported inventory snapshot.
- Item catalog browsing.
- Confirmed catalog item add.
- Confirmed selected-item repair.
- Confirmed selected-item removal.
- Confirmed selected-item stack-size editing.
- Before-action snapshot export for confirmed workflows.
- Shared mutation confirmation and required admin reason capture for confirmed workflows.
- Automatic post-action diff panel after confirmed add, stack-size edit, repair, or removal.
- Browser-session action history for recent confirmed action diffs.

Next Inventory Studio work should validate the stack-size edit workflow, then add a quality edit workflow only after before/after preview and confirmed mutation behavior are preserved.

### 2. Harden Mutation Safety Framework v1

The explicit current mutation-confirmation migration set is complete for:

- Active Players-table mutation workflows.
- Storage add/remove item operations.
- Database SQL execution.
- Battlegroup Exec server-control operations.
- Blueprint import.

Next framework hardening tasks:

- Typed backend mutation wrappers.
- Server-side before-change snapshot helpers.
- Audit export and filtering.
- UI visibility for reason-enforcement state.

### 3. Comprehensive AppSec endpoint audit

P0 AppSec audit work now has an initial document at `docs/appsec-endpoint-audit.md`.

The audit should continue through:

- Endpoint-by-endpoint handler review.
- Automated public/admin/self-service route boundary tests.
- SAST, DAST, dependency, and manual abuse-case validation.
- Remediation issue tracking for all findings in the audit document.
- Closure evidence for each finding.

### 4. Preserve Player 360 validated state

- Keep the standalone `Player 360` tab available.
- Keep the Players-table `360` launcher read-only.
- Do not add Player 360 quick actions until they are implemented as confirmed workflows with reason capture, target metadata, and audit visibility.

### 5. Guild management backlog

P2 guild-management work now includes dedicated guild-management features.

Guild administration should cover:

- Create guild.
- Delete guild.
- Add player to a guild.
- Remove player from a guild.
- Change a player's guild rank.
- Schema discovery for guild, member, and rank tables/functions before writes are implemented.
- Confirmation, admin reason capture, before-change snapshot, post-action diff/reload, and audit visibility for every guild mutation.
- Clear distinction between guild management and faction reputation management.

Player tab guild workflows should cover:

- Add the selected player to a guild from the Player tab.
- Remove the selected player from their guild from the Player tab.
- Promote or otherwise change the selected player's rank within their guild from the Player tab.
- Show current guild membership and rank context before any mutation.
- Reuse the same confirmed mutation, admin reason, audit, and post-action refresh pattern as other player-management workflows.

### 6. Documentation backlog

P5 documentation work now includes a detailed Discord bot setup and usage guide.

That guide should cover:

- Discord application and bot creation.
- Required OAuth2 scopes and bot permissions.
- Bot invite workflow.
- Required environment variables and secret-handling expectations.
- Local development setup.
- Production deployment expectations.
- Slash-command registration or command adapter behavior.
- Role and permission model.
- Player-link/self-service prerequisites.
- Farming request/order command examples.
- Troubleshooting and validation steps.
- Security boundaries and explicit non-goals.

## Validation command set

```bash
gofmt -w *.go
go mod tidy
go test ./...
cd web
npm install
npm audit --audit-level=high
npm run typecheck
npm run lint
npm run build
```

Canonical local validation path:

```bash
./update.sh
```

On Windows PowerShell, use:

```powershell
.\update.ps1
```

## Notes

The DA Manager workstream also tracks runtime/operator usability requirements. SSH tunnel status should remain visible at startup so operators can confirm managed database forwarding before using protected admin workflows.
