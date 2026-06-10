# Arrakis Control Panel Release Notes

## Current update: PowerShell npm/web helper modularization

### Why this update was made

PowerShell update support remains part of the pre-`v0.1.0` refactor/modularization gate. Common helpers and Git helpers have already been extracted. The npm/web repair and frontend toolchain helper functions still needed to be split out of `update.ps1`.

This slice moves npm install/repair, Node process summary, frontend package binary checks, and npm lock guidance into a dedicated PowerShell module while preserving `update.ps1` as the entry point and keeping the validation/build order unchanged.

### What changed

- Added `scripts/update/powershell-npm.ps1`.
- Moved npm/web helper functions out of `update.ps1`:
  - `Get-NodeProcessSummary`
  - `Remove-NodeModulesForRepair`
  - `Invoke-NpmInstallWithRepair`
  - `Test-WebPackageBinary`
  - `Assert-WebPackageToolchain`
  - `Show-NpmLockHelp`
- Updated `update.ps1` to dot-source `scripts/update/powershell-npm.ps1`.
- Preserved existing npm install, npm repair, frontend package toolchain, audit, typecheck, lint, and build behavior.

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

- Continue Go review for handler boundaries, audit/mutation-safety helper boundaries, and typed execution surfaces.
- Consider splitting remaining PowerShell main flow only if it improves readability without obscuring execution order.

---

## Previous update: PowerShell Git helper modularization

### Validation

Validated from both update paths:

```powershell
.\update.ps1 -SkipAutoPush
```

```bash
./update.sh
```
