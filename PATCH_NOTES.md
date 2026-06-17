# Arrakis Control Panel Release Notes

## Current update: Ledger cleanup after README prerequisites validation

### Why this update was made

`./update.sh` failed because `PATCH_NOTES.md` had 250 lines. The ledger check allows a maximum of 220 lines and requires this file to stay focused on the current operator-facing update.

### What changed

- Replaced the long historical patch notes with this compact current update.
- Moved durable detail to:

```text
docs/changelog/unreleased/2026-06-16-readme-prerequisites-clean-build.md
```

- Kept the README prerequisites and required tooling update intact.
- Kept the clean build validation record in the detailed unreleased changelog file.

### Impact

This is a documentation-only ledger cleanup.

### Validation

Run the canonical validation path again:

```bash
./update.sh
```

On Windows, also run:

```powershell
.\update.ps1 -SkipAutoPush
```

### Remaining final-`v0.1.0` gate

- Post-release verification after tag/artifact install or launch.
