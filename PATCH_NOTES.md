# Dune Admin Release Notes

## Current update: Changelog and ledger compaction

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

### Other large log/audit files

The same rule applies to:

- `docs/appsec-endpoint-audit.md`
- future AppSec finding files
- release checklist records
- risk-register updates
- validation evidence logs
- any Markdown file acting as a mutable ledger

The preferred remediation is not to create another giant archive file. Instead:

- keep a compact index
- preserve full historical states through immutable commit references
- add small per-topic/per-slice records
- enforce line-count guardrails before validation/builds

### Security and operator impact

- No route behavior changed.
- No mutation behavior changed.
- No new endpoint was added.
- Player 360 remains read-only.
- This reduces future edit risk for audit/security release records.

### Validation

Validation pending from the canonical local update path:

```bash
./update.sh
```

Ledger-specific validation can be run directly with:

```bash
bash scripts/check-ledger-size.sh
```

### Known limitation

`update.sh` is also large enough that full-file connector reads can truncate. I added the ledger-size script, but did not wire it into `update.sh` in this slice to avoid a risky full-file replacement. The safe next step is to patch `update.sh` locally or through a patch-capable edit path to run:

```bash
step "Ledger size check" run bash scripts/check-ledger-size.sh
```

after `invoke_git_pull_if_safe` and before Go tests.
