# Arrakis Control Panel Release Notes

## Current update: PowerShell Git helper modularization

### Why this update was made

PowerShell update support remains part of the pre-`v0.1.0` refactor/modularization gate. The common helper module has already been extracted, but Git helper functions still needed to be split out of `update.ps1`.

This slice moves Git status, pull, auto-commit, and auto-push helpers into a dedicated PowerShell module while preserving `update.ps1` as the entry point and keeping the validation/build order unchanged.

### What changed

- Added `scripts/update/powershell-git.ps1`.
- Moved Git helper functions out of `update.ps1`:
  - `Get-GitStatusLines`
  - `Write-GitStatusPreview`
  - `Invoke-GitPullIfSafe`
  - `Invoke-AutoCommitIfNeeded`
  - `Invoke-AutoPushIfNeeded`
- Updated `update.ps1` to dot-source `scripts/update/powershell-git.ps1`.
- Preserved existing Git pull, dirty-worktree handling, auto-commit, and auto-push behavior.

### Security and operator impact

- No route behavior changed.
- No mutation behavior changed.
- No new endpoint was added.
- Bash `./update.sh` remains the canonical validated release workflow.
- PowerShell update support is more modular and easier to review.

### Validation

Validated from both update paths:

```powershell
.\update.ps1 -SkipAutoPush
```

```bash
./update.sh
```

### Remaining refactor work

- Continue PowerShell modularization for npm helper group.
- Continue Go review for handler boundaries, audit/mutation-safety helper boundaries, and typed execution surfaces.

---

## Previous update: Update toolchain checks, PowerShell color output, and PATH fix

### Validation

Validated from both update paths:

```powershell
.\update.ps1 -SkipAutoPush
```

```bash
./update.sh
```
