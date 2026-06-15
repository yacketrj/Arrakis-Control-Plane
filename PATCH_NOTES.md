# Arrakis Control Panel Release Notes

## Current update: PowerShell web helper extraction

### Why this update was made

The final release gate still includes update-script modularization. PowerShell backend logic has been extracted, but web dependency validation, audit, typecheck, lint, and build execution were still embedded in `update.ps1`.

### What changed

- Added `scripts/update/powershell-web.ps1`.
- Moved PowerShell web folder detection and package-file checks into the web helper module.
- Moved Node/npm prerequisite checks and version probes into the web helper module.
- Moved npm install/repair, toolchain validation, audit, typecheck, lint, and build execution into the web helper module.
- Updated `update.ps1` to source the web helper module and call `Invoke-WebValidationAndBuild`.

### Security and operator impact

- No route behavior changed.
- No mutation behavior changed.
- No new endpoint was added.
- PowerShell update behavior should remain equivalent while reducing `update.ps1` monolith size.

### Validation

Validation pending.

Recommended validation:

```bash
./update.sh
```

```powershell
.\update.ps1 -SkipAutoPush
```

### Remaining final-`v0.1.0` gates

- Confirm whether update-script modularization is now sufficient for final `v0.1.0` or document any remaining deferral.
- Remaining Go code-quality/refactor review or explicit deferral.
- Full documentation review beyond primary release/security docs, or explicit deferral.
- Post-release verification after tag/artifact install or launch.

---

## Previous update: PowerShell backend helper extraction

### Validation

Validated from both update paths:

```bash
./update.sh
```

```powershell
.\update.ps1 -SkipAutoPush
```
