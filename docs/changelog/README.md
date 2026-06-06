# Changelog and Ledger Policy

## Purpose

This directory prevents large mutable release/audit ledgers from becoming unsafe to edit.

`CHANGELOG.md`, `PATCH_NOTES.md`, and active audit trackers should stay compact. Detailed records belong in small dated files under `docs/changelog/unreleased/` or immutable archive/index files under `docs/changelog/archive/`.

## File roles

| File or directory | Role | Edit pattern |
|---|---|---|
| `CHANGELOG.md` | Compact release index and current summary | Small edits only |
| `PATCH_NOTES.md` | Current operator-facing update only | Replaced per active slice |
| `docs/changelog/unreleased/*.md` | Detailed per-slice records | One new file per meaningful slice |
| `docs/changelog/archive/*.md` | Immutable monthly archive/index records | Rare append or no edits |
| `docs/appsec-endpoint-audit.md` | Current AppSec audit index/status | Keep concise; split detailed findings out when it grows |

## Rules

- Do not paste long implementation details into `CHANGELOG.md`.
- Do not let `PATCH_NOTES.md` become historical archive content.
- Prefer one small detailed file per validated work slice.
- Archive historical detail by month or by source commit instead of repeatedly editing a giant Markdown file.
- Keep active audit tracker files as indexes. Move long finding details to dedicated files such as `docs/appsec/findings/ASEA-003.md` when needed.

## Size guardrails

The canonical update path runs:

```bash
scripts/check-ledger-size.sh
```

The script enforces line-count limits for mutable ledger files and small per-slice changelog records. It also reports large Markdown files for operator visibility.

## Current archive

The last full inline changelog before compaction is preserved in Git history at commit:

```text
05c6cc17d133a0815af0fb1be0fc5cc1e8d53d40
```

See `docs/changelog/archive/2026-06.md` for the archive index.
