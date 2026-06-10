# Changelog

All notable changes to this project will be documented in this file.

This file is intentionally compact. Detailed change records are stored as small per-slice files under `docs/changelog/unreleased/` and archive indexes under `docs/changelog/archive/`.

See `docs/changelog/README.md` for the changelog and ledger policy.

## [Unreleased]

### Added

- Added a direct `docs/roadmap.md` pointer to `README.md`.
- Added `scripts/update/powershell-npm.ps1` for PowerShell npm/web update helper functions.
- Added `scripts/update/powershell-git.ps1` for PowerShell Git update helper functions.
- Added Bash and PowerShell frontend package toolchain checks for local `tsc`, `eslint`, and `vite` package binaries.
- Added colored PowerShell update status output with `RUN`, `PASS`, `FAIL`, and `WARN` states.
- Added `scripts/update/powershell-common.ps1` for shared PowerShell update helper functions.
- Added shared Go application identity constants in `app_identity.go`.
- Added final-`v0.1.0` release gates for update-script modularization and Go code-quality/refactor review in `docs/release-versioning.md`.
- Added deep documentation review plan in `docs/documentation-review-plan.md`.
- Added release deviation log in `docs/release-deviation-log.md`.
- Added release-train goals, label-sync rules, deviation policy, and industry-standard gap assessment to `docs/release-versioning.md`.
- Added upstream attribution requirement for Icehunter's `dune-admin` project by Ryan Wilson.
- Added `VERSION` with initial release candidate version `0.1.0-rc.1`.
- Added first release checklist instance in `docs/releases/v0.1.0-rc.1.md`.
- Added route-specific audit target assertions in `audit_log_target_test.go`.
- Added per-slice changelog record for route-specific audit targets in `docs/changelog/unreleased/2026-06-route-specific-audit-targets.md`.
- Added blocked high-risk/destructive mutation audit coverage in `audit_log_negative_test.go`.
- Added per-slice changelog record for blocked mutation audit coverage in `docs/changelog/unreleased/2026-06-blocked-mutation-audit-coverage.md`.
- Added changelog and ledger policy in `docs/changelog/README.md`.
- Added June 2026 archive index in `docs/changelog/archive/2026-06.md` with a pointer to the last full pre-compaction changelog commit.
- Added per-slice changelog record for high-risk mutation audit-event coverage in `docs/changelog/unreleased/2026-06-high-risk-mutation-audit-coverage.md`.
- Added `scripts/check-ledger-size.sh` to detect oversized mutable Markdown ledgers and warn on large Markdown files.

### Changed

- Canonical roadmap and feature priorities are now discoverable from the README through `docs/roadmap.md`.
- Confirmed the old ambiguous roadmap path `docs/admin-feature-design-and-priorities.md` is no longer present.
- Updated `update.ps1` to dot-source npm/web helpers from `scripts/update/powershell-npm.ps1`.
- Updated `update.ps1` to dot-source Git helpers from `scripts/update/powershell-git.ps1`.
- Updated missing frontend package binaries to trigger npm install/repair before typecheck/lint/build instead of failing later with missing `tsc`, `eslint`, or `vite` commands.
- Updated PowerShell PATH refresh to de-duplicate PATH entries before assigning `$env:Path`, preventing runaway environment block growth.
- Updated `update.ps1` to dot-source common PowerShell update helpers from `scripts/update/powershell-common.ps1`.
- Grouped Go route registration by API domain while preserving the existing `registerRoutes` entry point and route mappings.
- Updated Go startup logging, public status identity, and setup repair guidance to use current Arrakis Control Panel identity values.
- Removed developer-specific local checkout path from `README.md` and replaced it with repository-root-relative guidance.
- Refreshed `docs/admin-feature-design-and-priorities.md` for Arrakis Control Panel, the current release train, and validated product/service naming.
- Migrated Linux systemd installer defaults from legacy `dune-admin` names to `arrakis-control-panel` service, user/group, install path, binary path, unit description, and `ExecStart` values.
- Updated `README.md` systemd install commands and service defaults for Arrakis Control Panel.
- Renamed compiled backend executable output to `arrakis-control-panel` / `arrakis-control-panel.exe` in the Bash update workflow, PowerShell update workflow, and Linux build helper.
- Updated `README.md` build-output documentation for the Arrakis Control Panel executable name.
- Updated `README.md` to document Arrakis Control Panel naming, upstream attribution, strict token generation, canonical `./update.sh` validation/build workflow, release evidence locations, and release tagging workflow.
- Renamed the application/product label from DA Manager to Arrakis Control Panel in active release-governance documentation.
- Updated `docs/releases/v0.1.0-rc.1.md` with Arrakis Control Panel naming, upstream attribution, and current industry-standard release gaps.
- Established `v0.1.0-rc.1` as the first controlled release-candidate target.
- Updated `docs/releases/v0.1.0-rc.1.md` with clean build validation and approval to tag.
- Expanded admin audit target metadata extraction for player identity, item template, quantity, quality, vehicle, guild, rank, command, and command-path fields.
- Replaced oversized root `CHANGELOG.md` with a compact index and current summary.
- Moved durable detail out of the root changelog pattern and into per-slice/archive changelog records.

### Security

- Required update-script modularization and Go code-quality/refactor review before final `v0.1.0`, unless explicitly deferred in the release deviation log.
- Clarified in the roadmap that full Discord server management, Live Admin RMQ, Welcome Kits, and arbitrary raw command publishing remain out of scope for `v0.1.0`.
- Added a pre-`v0.1.0` documentation review gate for accuracy, authenticity, comprehensiveness, naming consistency, release evidence, and stale large mutable documents.
- Kept Live Admin / RMQ / Discord full server management out of the initial release-candidate scope.
- Added route-specific audit target assertions for high-risk mutation attempts.
- Added negative-path audit assertions for high-risk/destructive mutations blocked by missing admin reason or oversized reason-inspection body.
- Reduced future edit risk for security/audit release records by preventing `CHANGELOG.md` from continuing to grow as a giant mutable ledger.
- Added guardrails for other mutable ledger files including `PATCH_NOTES.md`, `docs/appsec-endpoint-audit.md`, and per-slice changelog records.

### Validation

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
- Validated route registration grouping from the canonical local update path:
  - `./update.sh`
- Validated refactor gate and Go app identity cleanup from the canonical local update path:
  - `./update.sh`
- Validated roadmap documentation refresh and README local-path cleanup from the canonical local update path:
  - `./update.sh`
- Validated Linux systemd service migration from the canonical local update path:
  - `./update.sh`
- Validated executable rename from the canonical local update path:
  - `./update.sh`
- Validated README correction and documentation-review-plan updates from the canonical local update path:
  - `./update.sh`
- Validated Arrakis Control Panel rename and release-plan documentation from the canonical local update path:
  - `./update.sh`
- Validated release-candidate setup from the canonical local update path:
  - `./update.sh`
- Validated route-specific audit target assertions from the canonical local update path:
  - `./update.sh`
  - Non-blocking build-performance warning observed: `[PLUGIN_TIMINGS] Your build spent significant time in plugin @tailwindcss/vite:generate:build`.
- Validated blocked mutation audit coverage from the canonical local update path:
  - `./update.sh`
- Ledger-specific validation can be run directly with:
  - `bash scripts/check-ledger-size.sh`

### Known issues

- Full documentation review is required before final `v0.1.0` or must be explicitly deferred in `docs/release-deviation-log.md`.
- Update-script modularization remains required before final `v0.1.0` or must be explicitly deferred in `docs/release-deviation-log.md`.
- Go code-quality/refactor review remains required before final `v0.1.0` or must be explicitly deferred in `docs/release-deviation-log.md`.
- Continue reviewing other long-lived docs for stale implemented-vs-planned claims.
- Repo-wide stale label verification for `DA Manager` and `Arrakis Control Plane` should continue before final `v0.1.0`.
- Existing Linux installs using legacy `dune-admin` service/path defaults require intentional migration before switching to the new `arrakis-control-panel` service defaults.
- `v0.1.0-rc.1` is approved to tag; post-release verification checks remain pending after tag/artifact install or launch.

## Detailed change records

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

File:

```text
CHANGELOG.md
```
