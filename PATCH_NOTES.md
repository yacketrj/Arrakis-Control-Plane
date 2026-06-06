# Dune Admin Release Notes

## Current update: Blocked mutation audit coverage

### Why this update was made

The AppSec hardening track is continuing before new Live Admin / RMQ / Welcome Kit features. The previous audit slice verified successful high-risk/destructive mutation audit events. This slice adds negative-path coverage so blocked high-risk/destructive mutations are also auditable when admin-reason enforcement rejects them before the downstream handler runs.

### What changed

- Added `audit_log_negative_test.go`.
- Added table-driven coverage for high-risk and destructive mutation routes blocked by missing admin reason.
- Verified blocked mutations:
  - do not reach the downstream handler
  - return `400 Bad Request`
  - still emit exactly one audit event
  - record the expected method and path
  - record mutation-safety action, risk, destructive flag, reason flag, and preview flag
  - record failure result and `400` status
  - preserve request ID
  - preserve common target metadata such as `player_id`, `account_id`, and `actor_id`
- Added oversized-body negative-path coverage for reason inspection.
- Added `docs/changelog/unreleased/2026-06-blocked-mutation-audit-coverage.md` as the durable per-slice record.

### Security and operator impact

- No route behavior changed.
- No mutation behavior changed.
- No new endpoint was added.
- Player 360 remains read-only.
- This adds regression evidence that blocked high-risk/destructive mutations remain visible in the admin audit trail.

### Validation

Validation pending from the canonical local update path:

```bash
./update.sh
```

### Remaining AppSec work

- route-specific target assertions beyond shared metadata
- pre/post-change review verification where practical
- SAST/DAST/dependency evidence
- manual abuse-case validation

---

## Previous update: Changelog and ledger compaction

### Why this update was made

`CHANGELOG.md` had grown into a large mutable release ledger. This created recurring connector/edit risk: each update required replacing a large file, and truncated tool output made it unsafe to guarantee that unrelated changelog entries would not be dropped.

The same pattern could affect other mutable log/audit files if they are allowed to grow indefinitely.

### What changed

- Added `docs/changelog/README.md` with the new changelog and ledger policy.
- Added `docs/changelog/archive/2026-06.md` as an archive index for June 2026 work.
- Added `docs/changelog/unreleased/2026-06-high-risk-mutation-audit-coverage.md` as the first detailed per-slice changelog record.
- Added `scripts/check-ledger-size.sh` to detect oversized mutable Markdown ledgers.
- Replaced the oversized root `CHANGELOG.md` with a compact index and current summary.
- Preserved the last full pre-compaction changelog in Git history at commit:
  - `05c6cc17d133a0815af0fb1be0fc5cc1e8d53d40`

### Policy going forward

- `CHANGELOG.md` stays compact and index-like.
- `PATCH_NOTES.md` remains current-update only.
- Detailed work-slice records go under `docs/changelog/unreleased/`.
- Monthly or release archives go under `docs/changelog/archive/`.
- Large audit trackers should become indexes; detailed findings should move to dedicated smaller files.

### Security and operator impact

- No route behavior changed.
- No mutation behavior changed.
- No new endpoint was added.
- Player 360 remains read-only.
- This reduces future edit risk for audit/security release records.
