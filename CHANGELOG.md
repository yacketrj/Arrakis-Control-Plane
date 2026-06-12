# Changelog

All notable changes to this project will be documented in this file.

This file is intentionally compact. Detailed change records are stored as small per-slice files under `docs/changelog/unreleased/` and archive indexes under `docs/changelog/archive/`.

See `docs/changelog/README.md` for the changelog and ledger policy.

## [Unreleased]

### Current focus

- Active project identity migration to Arrakis Control Panel / `arrakis-control-panel` / `arrakis-control-plane`.
- Go validation package filtering to avoid frontend dependency trees.
- Audit helper modularization for risk and request metadata boundaries.
- Allowed-origin validation hardening.

### Added

- `audit_metadata.go` for audit request metadata extraction and sanitization helpers.
- `audit_risk.go` for audit action and mutation-risk classification helpers.
- `docs/roadmap.md` as the canonical roadmap pointer from README.
- PowerShell update helper modules under `scripts/update/`.
- Frontend package toolchain checks for Bash and PowerShell update paths.

### Changed

- Restored `audit_metadata.go` after validation showed `audit_log.go` expected extracted metadata helpers that were absent from `main`.
- Updated `SECURITY.md` active product wording, allowed-origin example, and browser token-storage description for current Arrakis Control Panel identity.
- Renamed active Go module identity from `dune-admin` to `arrakis-control-plane`.
- Renamed active build outputs, deploy target, setup output, package metadata, diagnostics, and remediation tracker references to Arrakis Control Panel naming.
- Hardened allowed-origin validation to reject wildcard host components such as `http://*`.
- Updated Bash and PowerShell Go test validation to enumerate backend packages through `go list ./...` and exclude frontend dependency/build paths.
- Moved audit metadata extraction out of `audit_log.go` into `audit_metadata.go`.
- Moved audit action/risk classification out of `audit_log.go` into `audit_risk.go`.

### Security

- Security guidance now uses active Arrakis Control Panel identity and current browser token-storage wording.
- Backend Go test validation excludes third-party frontend dependency trees such as `web/node_modules`.
- Allowed-origin validation rejects wildcard host values.
- High/destructive mutation audit and blocked-mutation audit coverage remain tracked.
- Final `v0.1.0` gates still include documentation review, update-script modularization or explicit deferral, and Go code-quality/refactor review or explicit deferral.

### Validation

- Validated `SECURITY.md` active identity cleanup and restored `audit_metadata.go` from both update paths:
  - `./update.sh`
  - `.\update.ps1 -SkipAutoPush`
- Validated active identity migration, Go test package filter fix, audit metadata helper refactor, and allowed-origin wildcard hardening from both update paths:
  - `./update.sh`
  - `.\update.ps1 -SkipAutoPush`
- Validated audit risk helper refactor from the canonical local update path:
  - `./update.sh`
- Validated roadmap discoverability update from the canonical local update path:
  - `./update.sh`
- Validated PowerShell npm/web helper modularization from both update paths:
  - `.\update.ps1 -SkipAutoPush`
  - `./update.sh`
- Validated PowerShell Git helper modularization from both update paths:
  - `.\update.ps1 -SkipAutoPush`
  - `./update.sh`
- Validated update toolchain checks, PowerShell PATH de-duplication, and colored status output from both update paths:
  - `.\update.ps1 -SkipAutoPush`
  - `./update.sh`
- Earlier validated release-candidate setup and route/audit safety work remains recorded in per-slice changelog records and release evidence.

### Known issues

- Full documentation review is required before final `v0.1.0` or must be explicitly deferred in `docs/release-deviation-log.md`.
- Update-script modularization remains required before final `v0.1.0` or must be explicitly deferred in `docs/release-deviation-log.md`.
- Go code-quality/refactor review remains required before final `v0.1.0` or must be explicitly deferred in `docs/release-deviation-log.md`.
- `v0.1.0-rc.1` is approved to tag; post-release verification checks remain pending after tag/artifact install or launch.

## Detailed change records

- `docs/releases/v0.1.0-rc.1.md`
- `docs/roadmap.md`
- `docs/documentation-review-plan.md`
- `docs/changelog/README.md`
- `docs/changelog/unreleased/`
- `docs/changelog/archive/`
- `docs/release-deviation-log.md`

## Last full pre-compaction changelog

The last full inline changelog before compaction is preserved in Git history at commit:

```text
05c6cc17d133a0815af0fb1be0fc5cc1e8d53d40
```
