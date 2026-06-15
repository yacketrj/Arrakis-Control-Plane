# PowerShell backend helper extraction

## Summary

Extracted PowerShell backend update logic from `update.ps1` into `scripts/update/powershell-backend.ps1`.

## Changes

- Added `scripts/update/powershell-backend.ps1`.
- Moved backend Go package discovery and test execution into helper functions.
- Moved backend build execution into a helper function.
- Moved backend binary and asset copy behavior into a helper function.
- Updated `update.ps1` to source and call the backend helper module.

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
