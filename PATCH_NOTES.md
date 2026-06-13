# Arrakis Control Panel Release Notes

## Current update: README validation wording correction

### Why this update was made

The documentation review found one stale README statement: it still said the current validated release workflow was only `./update.sh`, even after PowerShell validation was completed for the current hardening slice.

### What changed

- Updated README PowerShell wording to state that `update.ps1` is validated for the current hardening slice.
- Added the validated PowerShell command:

```powershell
.\update.ps1 -SkipAutoPush
```

- Removed the duplicate inline admin-token generation example from README.
- Pointed admin-token setup guidance to `SECURITY.md` so strict token guidance has a single source of truth.

### Security and operator impact

- No route behavior changed.
- No mutation behavior changed.
- No new endpoint was added.
- README now more accurately reflects the current validation status and security documentation boundaries.

### Validation

Validated from the local update path requested for this slice.

```bash
./update.sh
```

```powershell
.\update.ps1 -SkipAutoPush
```

### Remaining final-`v0.1.0` gates

- Full documentation review beyond the primary release/security docs already checked, or explicit deferral.
- Remaining update-script modularization or explicit deferral.
- Remaining Go code-quality/refactor review or explicit deferral.
- Post-release verification after tag/artifact install or launch.

---

## Previous update: Documentation review status record

### Validation

Validated from both update paths:

```bash
./update.sh
```

```powershell
.\update.ps1 -SkipAutoPush
```
