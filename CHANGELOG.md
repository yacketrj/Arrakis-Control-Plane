# Changelog

All notable changes to this project will be documented in this file.

This file is intentionally compact. Detailed change records are stored as small per-slice files under `docs/changelog/unreleased/` and archive indexes under `docs/changelog/archive/`.

See `docs/changelog/README.md` for the changelog and ledger policy.

## [Unreleased]

### Current focus

- Runtime evidence audit for Docker and Kubernetes detection.
- Frontend typecheck fix for runtime-aware Battlegroup UI.
- Runtime-aware Docker and Kubernetes status handling.
- Setup wizard `.env` generation before SSH validation.
- Ledger-size compliance for current operator-facing notes.
- Final `v0.1.0` release-readiness gate disposition.

### Added

- `docs/changelog/unreleased/2026-06-16-runtime-evidence-audit.md` for runtime evidence audit detail.
- `docs/changelog/unreleased/2026-06-16-heroui-button-title-typecheck.md` for the frontend typecheck fix.
- `docs/changelog/unreleased/2026-06-16-runtime-aware-status.md` for runtime-aware status detail.
- `docs/changelog/unreleased/2026-06-16-setup-env-before-ssh.md` for setup `.env` generation fix detail.
- `docs/changelog/unreleased/2026-06-16-readme-prerequisites-clean-build.md` for README prerequisite and clean-build detail.

### Changed

- Runtime auto-detection now checks active Dune workload evidence before selecting Docker or Kubernetes.
- Status payload generation refreshes runtime from active workload evidence when SSH is connected.
- Removed an unsupported UI prop from the Battlegroup server-control button.
- Docker database discovery now prefers the published Docker port mapping.
- Battlegroup status UI now reads backend runtime and renders Docker as containers or Kubernetes as pods.
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

- Docker runtime keeps battlegroup script controls disabled until that command path is Docker-safe.
- Broad Go refactor review is deferred to `v0.1.1` or the next hardening slice.
- Broad documentation review is deferred to `v0.1.1` or the next documentation-hardening slice.
- Post-release verification remains pending after tag/artifact install or launch.

## Detailed change records

- `runtime_discovery.go`
- `status_payload.go`
- `docker_discovery.go`
- `web/src/tabs/BattlegroupTab.tsx`
- `setup.go`
- `README.md`
- `PATCH_NOTES.md`
- `docs/changelog/unreleased/2026-06-16-runtime-evidence-audit.md`
- `docs/changelog/unreleased/2026-06-16-heroui-button-title-typecheck.md`
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
