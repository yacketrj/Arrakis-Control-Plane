# Arrakis Control Panel Release Notes

## Current update: PowerShell common helper modularization and colored status output

### Why this update was made

PowerShell update support remains part of the pre-`v0.1.0` refactor/modularization gate. `update.ps1` was still mostly monolithic, even though the Bash update path had already been partially modularized under `scripts/update/`.

This slice starts PowerShell modularization by extracting common helper functions while preserving `update.ps1` as the entry point and keeping the validation/build order unchanged. It also aligns PowerShell update output with the Bash status style.

### What changed

- Added `scripts/update/powershell-common.ps1`.
- Moved common helper functions out of `update.ps1`:
  - `Write-Section`
  - `Invoke-Step`
  - `Invoke-Native`
  - `Update-ProcessPath`
  - `Install-PrerequisiteForCommand`
  - `Assert-CommandAvailable`
  - `Resolve-OutputDirectory`
- Updated `update.ps1` to dot-source `scripts/update/powershell-common.ps1`.
- Added shared `Write-StepStatus` helper.
- Added PowerShell colored status output:
  - `RUN` = cyan
  - `PASS` = green
  - `FAIL` = red
  - `WARN` = yellow
- Preserved existing Git, Go, npm, build, auto-commit, and auto-push flow.

### Security and operator impact

- No route behavior changed.
- No mutation behavior changed.
- No new endpoint was added.
- Bash `./update.sh` remains the canonical validated release workflow.
- PowerShell update support is more modular and easier to read, but this slice still requires validation on Windows PowerShell.

### Validation

Validation pending.

Recommended validation:

```powershell
.\update.ps1 -SkipAutoPush
```

Canonical Bash validation should also still pass:

```bash
./update.sh
```

### Remaining refactor work

- Continue PowerShell modularization for Git and npm helper groups.
- Continue Go review for handler boundaries, audit/mutation-safety helper boundaries, and typed execution surfaces.

---

## Previous update: Route registration grouping

### Validation

Validated from the canonical local update path:

```bash
./update.sh
```
