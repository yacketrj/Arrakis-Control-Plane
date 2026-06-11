# Arrakis Control Panel Release Notes

## Current update: Go test package filter fix

### Why this update was made

Validation failed because `go test ./...` traversed frontend dependency directories and discovered a third-party Go package under `web/node_modules/flatted/golang/pkg/flatted`. Frontend dependency trees should not be included in backend Go validation.

### What changed

- Updated Bash backend validation to enumerate Go packages through `go list ./...` and filter out frontend dependency/build paths:
  - `/web/node_modules/`
  - `/web/dist/`
- Updated PowerShell backend validation to apply the same package filtering before running `go test`.
- Preserved Go backend build behavior.
- Preserved frontend npm audit, typecheck, lint, and build behavior.

### Security and operator impact

- No route behavior changed.
- No mutation behavior changed.
- No new endpoint was added.
- Backend Go validation now excludes third-party frontend dependency trees.
- This prevents unrelated vendored/frontend dependency packages from failing backend test gates.

### Validation

Validation pending.

Recommended validation:

```bash
./update.sh
```

PowerShell validation should also be rerun:

```powershell
.\update.ps1 -SkipAutoPush
```

### Previous validation issue

The failure looked like:

```text
FAIL    dune-admin [build failed]
?       dune-admin/web/node_modules/flatted/golang/pkg/flatted [no test files]
FAIL
```

### Remaining refactor work

- Re-run validation for the audit metadata helper refactor after this test package filter fix.
- Continue Go review for mutation-safety helper boundaries.
- Continue Go review for handler boundaries and typed/allowlisted execution surfaces.

---

## Previous update: Audit metadata helper refactor

### Validation

Validation pending.

Recommended validation:

```bash
./update.sh
```
