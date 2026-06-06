# Changelog

All notable changes to this project will be documented in this file.

This file is intentionally compact. Detailed change records are stored as small per-slice files under `docs/changelog/unreleased/` and archive indexes under `docs/changelog/archive/`.

See `docs/changelog/README.md` for the changelog and ledger policy.

## [Unreleased]

### Added

- Added route-specific audit target assertions in `audit_log_target_test.go`.
- Added per-slice changelog record for route-specific audit targets in `docs/changelog/unreleased/2026-06-route-specific-audit-targets.md`.
- Added blocked high-risk/destructive mutation audit coverage in `audit_log_negative_test.go`.
- Added per-slice changelog record for blocked mutation audit coverage in `docs/changelog/unreleased/2026-06-blocked-mutation-audit-coverage.md`.
- Added changelog and ledger policy in `docs/changelog/README.md`.
- Added June 2026 archive index in `docs/changelog/archive/2026-06.md` with a pointer to the last full pre-compaction changelog commit.
- Added per-slice changelog record for high-risk mutation audit-event coverage in `docs/changelog/unreleased/2026-06-high-risk-mutation-audit-coverage.md`.
- Added `scripts/check-ledger-size.sh` to detect oversized mutable Markdown ledgers and warn on large Markdown files.

### Changed

- Expanded admin audit target metadata extraction for player identity, item template, quantity, quality, vehicle, guild, rank, command, and command-path fields.
- Replaced oversized root `CHANGELOG.md` with a compact index and current summary.
- Moved durable detail out of the root changelog pattern and into per-slice/archive changelog records.

### Security

- Added route-specific audit target assertions for high-risk mutation attempts.
- Added negative-path audit assertions for high-risk/destructive mutations blocked by missing admin reason or oversized reason-inspection body.
- Reduced future edit risk for security/audit release records by preventing `CHANGELOG.md` from continuing to grow as a giant mutable ledger.
- Added guardrails for other mutable ledger files including `PATCH_NOTES.md`, `docs/appsec-endpoint-audit.md`, and per-slice changelog records.

### Validation

- Validated route-specific audit target assertions from the canonical local update path:
  - `./update.sh`
  - Non-blocking build-performance warning observed: `[PLUGIN_TIMINGS] Your build spent significant time in plugin @tailwindcss/vite:generate:build`.
- Validated blocked mutation audit coverage from the canonical local update path:
  - `./update.sh`
- Ledger-specific validation can be run directly with:
  - `bash scripts/check-ledger-size.sh`

### Known issues

- Update-script modularization has started and needs continued validation through the canonical update path before further refactor expansion.

## Detailed change records

- `docs/changelog/README.md`
- `docs/changelog/unreleased/`
- `docs/changelog/archive/`

## Last full pre-compaction changelog

The last full inline changelog before compaction is preserved in Git history at commit:

```text
05c6cc17d133a0815af0fb1be0fc5cc1e8d53d40
```

File:

```text
CHANGELOG.md
```
