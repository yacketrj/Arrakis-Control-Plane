# Arrakis Control Panel Release Notes

## Current update: Update toolchain checks, PowerShell color output, and PATH fix

### Why this update was made

PowerShell update support remains part of the pre-`v0.1.0` refactor/modularization gate. `update.ps1` was still mostly monolithic, even though the Bash update path had already been partially modularized under `scripts/update/`.

This slice starts PowerShell modularization by extracting common helper functions while preserving `update.ps1` as the entry point and keeping the validation/build order unchanged. It also aligns PowerShell update output with the Bash status style.

A validation failure also showed that frontend package tools such as `tsc` can be missing even when Node/npm are installed. The update scripts now check the local frontend package toolchain before typecheck/lint/build. Another PowerShell failure showed PATH growth could trigger `Environment variable name or value is too long`; PowerShell PATH refresh now de-duplicates entries.

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
- Added Bash frontend package toolchain checks for:
  - `tsc`
  - `eslint`
  - `vite`
- Added PowerShell frontend package toolchain checks for:
  - `tsc`
  - `eslint`
  - `vite`
- Missing frontend package binaries now trigger npm install/repair before typecheck/lint/build.
- Updated PowerShell PATH refresh to de-duplicate PATH entries before assigning `$env:Path`.
- Preserved existing Git, Go, npm, build, auto-commit, and auto-push flow.

### Security and operator impact

- No route behavior changed.
- No mutation behavior changed.
- No new endpoint was added.
- Bash `./update.sh` remains the canonical validated release workflow.
- PowerShell update support is more modular and easier to read, but this slice still requires validation on Windows PowerShell.
- Missing local frontend tools should now be repaired before failing with `tsc is not recognized`.

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
