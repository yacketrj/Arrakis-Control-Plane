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
| P1 | Player 360 Profile | In Progress | `docs/player-360-profile.md` | Go + frontend validation pending |
| P1 | Inventory Studio v2 | Planned | `docs/inventory-studio.md` required | Go + frontend tests |
| P1 | Battlegroup Status v2 | Planned | `docs/battlegroup-status-v2.md` required | Go + frontend tests |
| P1 | Broadcast Center | Planned | `docs/broadcast-center.md` required | Go + frontend tests |
| P1 | Safe Offline Teleport / Rescue | Planned | `docs/safe-teleport-rescue.md` required | Go + frontend tests |
| P1 | Server Health Command Center | Planned | `docs/server-health-command-center.md` required | Go + frontend tests |
| P2 | RMQ Operations Dashboard | Planned | `docs/rmq-operations-dashboard.md` required | Go + frontend tests |
| P2 | Journey / Progression Manager | Planned | `docs/journey-progression-manager.md` required | Go + frontend tests |
| P2 | Guild / Faction Admin | Planned | `docs/guild-faction-admin.md` required | Go + frontend tests |
| P2 | Augment Preset Manager | Planned | `docs/augment-preset-manager.md` required | Go + frontend tests |
| P2 | Maintenance Mode Assistant | Planned | `docs/maintenance-mode-assistant.md` required | Go + frontend tests |
| P3 | Settings Diff / Config Manager | Planned | `docs/settings-diff-config-manager.md` required | Research + tests |
| P3 | Chat Moderation / Intercept Viewer | Planned | `docs/chat-moderation-intercept-viewer.md` required | Research + tests |
| P3 | Safe Stored Procedure Runner | Blocked | `docs/safe-stored-procedure-runner.md` required | Requires audit + mutation safety |

## Current implementation focus

### 1. Validate Player 360 read-only foundation

- Confirm backend compile and Go tests for `player_profile.go`, `routes.go`, and `server.go`.
- Confirm frontend typecheck, lint, and build for `web/src/api/playerProfile.ts`, `web/src/tabs/Player360Tab.tsx`, and `web/src/App.tsx`.
- Fix any validation failures before adding more Player 360 UI affordances.

### 2. Stabilize Player 360 operator workflow

- Keep the standalone `Player 360` tab read-only.
- Keep existing Players, Inventory, Give Item, and Actions workflows unchanged until validation is green.
- Consider adding a `Player 360` launch button to the existing Players table only after standalone tab validation succeeds.

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

## Notes

The DA Manager workstream also tracks runtime/operator usability requirements. SSH tunnel status should remain visible at startup so operators can confirm managed database forwarding before using protected admin workflows.
