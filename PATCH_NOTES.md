# Arrakis Control Panel Release Notes

## Current update: Final v0.1.0 gate disposition

### Why this update was made

The final release-readiness checklist still had three open gates after update-script modularization was closed:

- Go code-quality/refactor review or explicit deferral
- full documentation review beyond primary release/security docs, or explicit deferral
- post-release verification after tag/artifact install or launch

### What changed

- Added `docs/final-v0.1.0-gate-status.md`.
- Updated `docs/documentation-review-status.md` so it no longer lists update-script modularization as open.
- Added explicit release-deviation entries for:
  - broad Go code-quality/refactor review deferral
  - broad documentation review deferral beyond primary release/security docs
- Kept post-release verification pending because it requires actual tag/artifact install or launch evidence.

### Security and operator impact

- No route behavior changed.
- No mutation behavior changed.
- No endpoint was added.
- No Live Admin / RMQ execution was added.
- No Player 360 mutation behavior was added.
- No Welcome Kits behavior was added.
- Release status is clearer for final `v0.1.0` decision-making.

### Validation

This was a documentation-only gate-disposition update prepared through the GitHub connector. Local validation still must be run before final tagging:

```bash
./update.sh
```

On Windows, also run:

```powershell
.\update.ps1 -SkipAutoPush
```

### Remaining final-`v0.1.0` gates

- Local validation for this documentation-only gate-disposition update.
- Post-release verification after tag/artifact install or launch.

---

## Previous update: Update-script modularization closure

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

### Remaining final-`v0.1.0` gates at that time

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
