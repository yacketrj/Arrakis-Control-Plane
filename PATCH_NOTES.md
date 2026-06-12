# Arrakis Control Panel Release Notes

## Current update: Security guidance cleanup and audit metadata restore

### Why this update was made

`SECURITY.md` still contained active legacy identity examples after the main Arrakis Control Panel rename work. Those references were no longer only upstream attribution/history; they were current operator-facing security guidance.

Validation also exposed that `audit_log.go` expected extracted audit metadata helpers while `audit_metadata.go` was missing from `main`. The missing helper file caused undefined-symbol failures in Go tests.

### What changed

- Updated `SECURITY.md` active identity examples:
  - Product wording now uses `Arrakis Control Panel`.
  - Allowed-origin example now uses `https://arrakis-control-panel.layout.tools`.
  - Browser token-storage wording now reflects the current interim `sessionStorage` state.
- Restored `audit_metadata.go` with the audit metadata helper definitions required by `audit_log.go`, `handlers_battlegroup.go`, and audit tests.

### Security and operator impact

- Security guidance now matches the active project identity.
- Security guidance now reflects the current browser token-storage posture.
- Audit metadata extraction helpers are present again, restoring Go package compilation.
- No route behavior changed.
- No mutation behavior changed.
- No new endpoint was added.

### Validation

Validated from both update paths:

```bash
./update.sh
```

```powershell
.\update.ps1 -SkipAutoPush
```

### Remaining cleanup work

- Re-run grep after pulling this change:

```bash
grep -R "dune-admin" . --exclude-dir=.git --exclude-dir=dist --exclude-dir=web/node_modules
```

- Remaining `dune-admin` references should now be upstream attribution or historical release notes only.

---

## Previous update: Active identity migration and origin validation cleanup

### Validation

Validated from both update paths:

```bash
./update.sh
```

```powershell
.\update.ps1 -SkipAutoPush
```
