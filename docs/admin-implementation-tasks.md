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
| P0 | Database routine discovery backend | Done | `PATCH_NOTES.md`, `CHANGELOG.md` pending sync after workflow cleanup | Go Quality |
| P0 | Database routine inspection UI | Done | `PATCH_NOTES.md`, `CHANGELOG.md` pending sync after workflow cleanup | Frontend Quality |
| P0 | Public-safe vs protected admin portal design | In Progress | `docs/portal-separation-design.md` | Go + frontend validation |
| P0 | Admin Action Audit Log | Done | `docs/admin-audit-log.md` | Go tests |
| P0 | Mutation Safety Framework | In Progress | `docs/mutation-safety-framework.md` | Go + frontend tests |
| P1 | Player 360 Profile | In Progress | `docs/player-360-profile.md` | Go tests pending local validation; frontend next |
| P1 | Inventory Studio v2 | Planned | `docs/inventory-studio.md` required | Go + frontend tests |
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

### 1. Validate Player 360 backend foundation

- Confirm `player_profile.go`, `routes.go`, and `server.go` compile together.
- Run Go tests locally or through CI.
- Fix any formatting, compile, or test failures before frontend work.

### 2. Build Player 360 frontend surface

- Add frontend API client support for `GET /api/v1/players/{id}/profile`.
- Add Player 360 detail page or panel from the existing Players view.
- Fold Currency and Online Status into Player Info.
- Keep the first frontend slice read-only.

### 3. Continue quality-gate cleanup

- Confirm remediation workflow applies Go formatting.
- Confirm Go Quality passes formatting, module verification, vet, and tests.
- Confirm Frontend Quality passes install, audit, typecheck, lint, and build.

## Documentation requirement per feature

Each feature document must include:

1. purpose
2. user/operator problem
3. supported workflows
4. backend routes or commands
5. frontend UI behavior
6. security controls
7. audit requirements
8. validation steps
9. known limitations
10. follow-up tasks

## Validation command set

```bash
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

The DB routine discovery feature is intended to answer whether the database has a safer function or routine for gameplay-like item delivery. Discovery and inspection are read-only. Any future routine execution must remain allowlisted, typed, previewed, audited, and protected by the mutation safety framework.

The DA Manager workstream also tracks runtime/operator usability requirements. SSH tunnel status should remain visible at startup so operators can confirm managed database forwarding before using protected admin workflows.
