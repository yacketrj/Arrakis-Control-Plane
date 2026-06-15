# Arrakis Control Panel Release Notes

## Current update: PowerShell backend helper extraction

### Why this update was made

The final release gate still includes update-script modularization. Bash update logic is already split across helper modules, but PowerShell backend test/build/copy logic was still embedded in `update.ps1`.

### What changed

- Added `scripts/update/powershell-backend.ps1`.
- Moved PowerShell backend Go package discovery and test execution into the backend helper module.
- Moved PowerShell backend build execution into the backend helper module.
- Moved PowerShell backend binary and asset copy logic into the backend helper module.
- Updated `update.ps1` to source the backend helper module and call the extracted functions.

### Security and operator impact

- No route behavior changed.
- No mutation behavior changed.
- No new endpoint was added.
- PowerShell update behavior should remain equivalent while reducing `update.ps1` monolith size.

### Validation

Validated from both update paths:

```bash
./update.sh
```

```powershell
.\update.ps1 -SkipAutoPush
```

### Remaining final-`v0.1.0` gates

- Remaining PowerShell web-validation modularization or explicit deferral.
- Remaining Go code-quality/refactor review or explicit deferral.
- Full documentation review beyond primary release/security docs, or explicit deferral.
- Post-release verification after tag/artifact install or launch.

---

## Previous update: README validation wording correction

### Validation

Validated from both update paths:

```bash
./update.sh
```

```powershell
.\update.ps1 -SkipAutoPush
```
