# Arrakis Control Panel Release Notes

## Current update: Roadmap discoverability update

### Why this update was made

The feature roadmap exists at `docs/roadmap.md`, but it was not obvious from the README. The previous roadmap filename was also too indirect, which made the feature list hard to locate during release planning.

### What changed

- Confirmed `docs/roadmap.md` is the canonical roadmap and feature-priority file.
- Confirmed the old ambiguous path `docs/admin-feature-design-and-priorities.md` is no longer present.
- Added a direct roadmap pointer to `README.md` under the current release section:

```text
docs/roadmap.md
```

### Security and operator impact

- No route behavior changed.
- No mutation behavior changed.
- No new endpoint was added.
- No update-script behavior changed.
- Documentation is easier to navigate for release planning and feature tracking.

### Validation

Validated from the canonical local update path:

```bash
./update.sh
```

### Remaining refactor work

- Continue Go review for handler boundaries, audit/mutation-safety helper boundaries, and typed execution surfaces.
- Consider splitting remaining PowerShell main flow only if it improves readability without obscuring execution order.

---

## Previous update: PowerShell npm/web helper modularization

### Validation

Validated from both update paths:

```powershell
.\update.ps1 -SkipAutoPush
```

```bash
./update.sh
```
