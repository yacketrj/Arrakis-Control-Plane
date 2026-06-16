# Changelog

All notable changes to this project will be documented in this file.

This file is intentionally compact. Detailed change records are stored as small per-slice files under `docs/changelog/unreleased/` and archive indexes under `docs/changelog/archive/`.

See `docs/changelog/README.md` for the changelog and ledger policy.

## [Unreleased]

### Current focus

- Final `v0.1.0` gate disposition and release-readiness cleanup.
- Active project identity migration to Arrakis Control Panel / `arrakis-control-panel` / `arrakis-control-plane`.
- Go validation package filtering to avoid frontend dependency trees.
- Audit helper modularization for risk and request metadata boundaries.
- Allowed-origin validation hardening.
- Documentation review gate status tracking.

### Added

- `docs/final-v0.1.0-gate-status.md` for the final release-readiness gate disposition.
- Explicit release-deviation entries for broad Go code-quality/refactor review deferral and broad documentation review deferral beyond primary release/security docs.
- `docs/documentation-review-status.md` for the current documentation-review findings, validation evidence, and remaining gate status.
- `audit_metadata.go` for audit request metadata extraction and sanitization helpers.
- `audit_risk.go` for audit action and mutation-risk classification helpers.
- `docs/roadmap.md` as the canonical roadmap pointer from README.
- PowerShell update helper modules under `scripts/update/`.
- Frontend package toolchain checks for Bash and PowerShell update paths.

### Changed

- Updated `docs/documentation-review-status.md` to reflect that update-script modularization is closed for final `v0.1.0` readiness.
- Updated `PATCH_NOTES.md` to record the final gate-disposition update and remaining validation/runtime evidence needs.
- Updated README to state that `update.ps1` is validated for the current hardening slice.
- Consolidated README admin-token setup guidance by pointing strict token instructions to `SECURITY.md`.
- Linked `docs/documentation-review-status.md` from `docs/documentation-review-plan.md`.
- Restored `audit_metadata.go` after validation showed `audit_log.go` expected extracted metadata helpers that were absent from `main`.
- Updated `SECURITY.md` active product wording, allowed-origin example, and browser token-storage description for current Arrakis Control Panel identity.
- Renamed active Go module identity from `dune-admin` to `arrakis-control-plane`.
- Renamed active build outputs, deploy target, setup output, package metadata, diagnostics, and remediation tracker references to Arrakis Control Panel naming.
- Hardened allowed-origin validation to reject wildcard host components such as `http://*`.
- Updated Bash and PowerShell Go test validation to enumerate backend packages through `go list ./...` and exclude frontend dependency/build paths.
- Moved audit metadata extraction out of `audit_log.go` into `audit_metadata.go`.
- Moved audit action/risk classification out of `audit_log.go` into `audit_risk.go`.

### Security

- Final `v0.1.0` will not add Live Admin / RMQ execution, Player 360 mutations, Welcome Kits, or raw command publishing.
- Broad Go code-quality/refactor review is explicitly deferred to avoid late unvalidated structural churn before final `v0.1.0`.
- Broad documentation review beyond primary release/security docs is explicitly deferred; primary operator and release-control references remain the trusted final-release set.
- Post-release verification remains pending until tag/artifact install or launch evidence exists.
- Security guidance now uses active Arrakis Control Panel identity and current browser token-storage wording.
- README now delegates strict token guidance to `SECURITY.md` to avoid conflicting token-generation instructions.
- Backend Go test validation excludes third-party frontend dependency trees such as `web/node_modules`.
- Allowed-origin validation rejects wildcard host values.
- High/destructive mutation audit and blocked-mutation audit coverage remain tracked.

### Validation

- Documentation-only gate-disposition update was prepared through the GitHub connector; local validation remains required before final tagging:
  - `./update.sh`
  - `.\update.ps1 -SkipAutoPush` on Windows
- Validation pending for README validation wording correction:
  - `./update.sh`
  - `.\update.ps1 -SkipAutoPush`
- Validated documentation-review status tracking from both update paths:
  - `./update.sh`
  - `.\update.ps1 -SkipAutoPush`
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

- Local validation is required before final `v0.1.0` tagging because the latest gate-disposition update was made through the GitHub connector.
- Broad Go code-quality/refactor review is deferred to `v0.1.1` or the next hardening slice.
- Broad documentation review beyond primary release/security docs is deferred to `v0.1.1` or the next documentation-hardening slice.
- `v0.1.0-rc.1` is approved to tag; post-release verification checks remain pending after tag/artifact install or launch.

## Detailed change records

- `docs/final-v0.1.0-gate-status.md`
- `docs/documentation-review-status.md`
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
