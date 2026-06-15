# PowerShell web helper extraction

## Summary

Extracted PowerShell web validation and build logic from `update.ps1` into `scripts/update/powershell-web.ps1`.

## Changes

- Added `scripts/update/powershell-web.ps1`.
- Moved web folder detection and package-file checks into a helper function.
- Moved Node/npm prerequisite checks and version probes into the helper module.
- Moved npm install/repair, toolchain validation, audit, typecheck, lint, and build execution into the helper module.
- Updated `update.ps1` to source and call `Invoke-WebValidationAndBuild`.

## Impact

- No route behavior changed.
- No mutation behavior changed.
- No new endpoint was added.
- PowerShell update behavior should remain equivalent while reducing `update.ps1` monolith size.

## Validation

Pending local validation:

```bash
./update.sh
```

```powershell
.\update.ps1 -SkipAutoPush
```
