# README validation wording correction

## Summary

Updated README documentation so the PowerShell update path is not understated after successful local validation.

## Changes

- README now states that `update.ps1` is validated for the current hardening slice.
- README shows the validated PowerShell command:

```powershell
.\update.ps1 -SkipAutoPush
```

- README now points setup guidance to `SECURITY.md` instead of duplicating strict credential setup details.

## Validation

Validated locally:

```bash
./update.sh
```

```powershell
.\update.ps1 -SkipAutoPush
```

## Impact

- Documentation-only change.
- No route behavior changed.
- No mutation behavior changed.
- No new endpoint was added.
