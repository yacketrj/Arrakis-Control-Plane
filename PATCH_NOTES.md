# Arrakis Control Panel Release Notes

## Current update: Update-script modularization closure

### Why this update was made

The final release checklist still needed a clear decision on whether update-script modularization was sufficient for final `v0.1.0` readiness.

### What changed

- Added `docs/update-script-modularization-status.md`.
- Recorded Bash helper-module coverage for `update.sh`.
- Recorded PowerShell helper-module coverage for `update.ps1`.
- Recorded the recently validated PowerShell backend and web helper extraction work.
- Closed the update-script modularization gate as sufficient for final `v0.1.0` readiness.

### Security and operator impact

- No route behavior changed.
- No mutation behavior changed.
- No new endpoint was added.
- Release tooling status is clearer for final readiness tracking.

### Validation

Validated from both update paths:

```bash
./update.sh
```

```powershell
.\update.ps1 -SkipAutoPush
```

### Remaining final-`v0.1.0` gates

- Remaining Go code-quality/refactor review or explicit deferral.
- Full documentation review beyond primary release/security docs, or explicit deferral.
- Post-release verification after tag/artifact install or launch.

---

## Previous update: PowerShell web helper extraction

### Validation

Validated from both update paths:

```bash
./update.sh
```

```powershell
.\update.ps1 -SkipAutoPush
```
