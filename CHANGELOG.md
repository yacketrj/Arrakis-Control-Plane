# Changelog

All notable changes to this project will be documented in this file.

This file is intentionally compact. Detailed change records are stored as small per-slice files under `docs/changelog/unreleased/` and archive indexes under `docs/changelog/archive/`.

See `docs/changelog/README.md` for the changelog and ledger policy.

## [Unreleased]

### Current focus

- Runtime-aware status handling.
- Setup wizard `.env` generation before SSH validation.
- Ledger-size compliance for current operator-facing notes.
- Final `v0.1.0` release-readiness gate disposition.

### Added

- `docs/changelog/unreleased/2026-06-16-runtime-aware-status.md` for runtime-aware status detail.
- `docs/changelog/unreleased/2026-06-16-setup-env-before-ssh.md` for setup `.env` generation fix detail.
- `docs/changelog/unreleased/2026-06-16-readme-prerequisites-clean-build.md` for README prerequisite and clean-build detail.

### Changed

- Runtime status now renders the correct target type for the detected runtime.
- Database discovery for container deployments now prefers the published service mapping.
- Setup writes a preliminary `.env` before remote SSH validation.
- `PATCH_NOTES.md` remains compact and focused on the current operator-facing update.

### Validation

Run local validation:

```bash
./update.sh
```

On Windows:

```powershell
.\update.ps1 -SkipAutoPush
```

### Known issues

- Container runtime server-control scripts remain disabled until that command path is safe.
- Broad Go refactor review is deferred to `v0.1.1` or the next hardening slice.
- Broad documentation review is deferred to `v0.1.1` or the next documentation-hardening slice.
- Post-release verification remains pending after tag/artifact install or launch.

## Detailed change records

- `docker_discovery.go`
- `web/src/tabs/BattlegroupTab.tsx`
- `setup.go`
- `README.md`
- `PATCH_NOTES.md`
- `docs/changelog/unreleased/2026-06-16-runtime-aware-status.md`
- `docs/changelog/unreleased/2026-06-16-setup-env-before-ssh.md`
- `docs/changelog/unreleased/2026-06-16-readme-prerequisites-clean-build.md`
- `docs/final-v0.1.0-gate-status.md`
- `docs/documentation-review-status.md`
- `docs/releases/v0.1.0-rc.1.md`
- `docs/roadmap.md`
- `docs/changelog/README.md`
- `docs/changelog/archive/`
- `docs/release-deviation-log.md`

## Last full pre-compaction changelog

The last full inline changelog before compaction is preserved in Git history at commit:

```text
05c6cc17d133a0815af0fb1be0fc5cc1e8d53d40
```
