# Config environment helper extraction

## Summary

Extracted environment helper logic from `main.go` into `config_env.go` as part of the Go code-quality/refactor review gate.

## Changes

- Added `config_env.go`.
- Moved dotenv loading and environment default helpers out of `main.go`.
- Reduced `main.go` imports to remove helper-only dependencies.

## Impact

- No route behavior changed.
- No mutation behavior changed.
- No new endpoint was added.
- Environment defaulting behavior should remain equivalent.

## Validation

Validated locally:

```bash
./update.sh
```

```powershell
.\update.ps1 -SkipAutoPush
```
