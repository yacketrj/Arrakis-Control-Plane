# Changelog

All notable changes to this project will be documented in this file.

This file is intentionally compact. Detailed change records are stored as small per-slice files under `docs/changelog/unreleased/` and archive indexes under `docs/changelog/archive/`.

See `docs/changelog/README.md` for the changelog and ledger policy.

## [Unreleased]

### Current focus

- README prerequisites and required tooling guidance.
- Final `v0.1.0` release-readiness gate disposition.
- Active Arrakis Control Panel product-identity cleanup.

### Added

- README guidance for prerequisite local applications, command-line tools, server-side access, optional integrations, optional release evidence tools, and quick prerequisite checks.
- `docs/final-v0.1.0-gate-status.md` for final release gate disposition.
- Release-deviation entries for deferred broad Go refactor review and deferred broad documentation review.

### Changed

- Recorded clean build validation for the README prerequisites update.
- Added `docs/final-v0.1.0-gate-status.md` to the README release-reference list.
- Recorded clean local validation for the final gate-disposition update.
- Kept the current README setup path pointed at `SECURITY.md` for strict token guidance.
- Earlier active-identity, validation, audit-helper, allowed-origin, and update-script cleanup remains recorded in detailed change records and Git history.

### Validation

- Operator-reported clean build validation for the README prerequisites update:

```bash
./update.sh
```

- PowerShell validation path remains available on Windows:

```powershell
.\update.ps1 -SkipAutoPush
```

### Known issues

- Broad Go refactor review is deferred to `v0.1.1` or the next hardening slice.
- Broad documentation review is deferred to `v0.1.1` or the next documentation-hardening slice.
- Post-release verification remains pending after tag/artifact install or launch.

## Detailed change records

- `README.md`
- `PATCH_NOTES.md`
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
