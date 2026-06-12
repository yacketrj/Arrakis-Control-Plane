# Arrakis Control Panel Release Notes

## Current update: Active identity migration and origin validation cleanup

### Why this update was made

Several active build, runtime, test, and setup references still used the legacy `dune-admin` name. Upstream attribution to Icehunter's original `dune-admin` project should remain, but active project identity should use Arrakis Control Panel naming.

Validation also exposed a real allowed-origin validation gap: wildcard hosts such as `http://*` were accepted by the origin validator and should be rejected.

### Naming decision

- Product/UI name: `Arrakis Control Panel`
- Binary/service/deploy name: `arrakis-control-panel`
- Go module path: `arrakis-control-plane`

This preserves the project/repository concept of a control plane while keeping the operator-facing application name aligned with the current README and systemd service naming.

### What changed

- Renamed the Go module from `dune-admin` to `arrakis-control-plane`.
- Renamed Makefile build outputs:
  - `arrakis-control-panel`
  - `arrakis-control-panel-linux`
- Renamed the Cloudflare Pages project target in the Makefile to `arrakis-control-panel`.
- Updated diagnostic export service identity to `arrakis-control-panel`.
- Updated auth test allowed-origin fixture to `https://arrakis-control-panel.layout.tools`.
- Updated package-lock root name to `arrakis-control-panel`.
- Updated setup banner, build/run hint, and generated `.env` comment to Arrakis Control Panel naming.
- Updated the security remediation tracker from legacy `dune-admin-fork` / `dune-admin.exe` references to Arrakis Control Panel / `arrakis-control-panel.exe`.
- Corrected the allowed-origin test helper reference from the stale `validateAllowedOriginValue` name to the active `isAllowedOriginValue` helper.
- Hardened allowed-origin validation to reject wildcard host components through `allowedOriginHostHasWildcard`.

### Intentionally retained legacy references

The remaining `dune-admin` references should be limited to upstream attribution or historical release notes, especially references to Icehunter's original `dune-admin` project by Ryan Wilson.

### Security and operator impact

- Active runtime, build, setup, package, and diagnostic identity now use Arrakis Control Panel naming.
- Backend Go validation excludes third-party frontend dependency trees.
- Allowed origins reject wildcard host values such as `http://*`.
- No route behavior changed.
- No mutation behavior changed.
- No new endpoint was added.

### Validation

Validated from the canonical local update path:

```bash
./update.sh
```

PowerShell validation remains pending unless separately run:

```powershell
.\update.ps1 -SkipAutoPush
```

### Remaining cleanup work

- Re-run grep after pulling this change:

```bash
grep -R "dune-admin" . --exclude-dir=.git --exclude-dir=dist --exclude-dir=web/node_modules
```

- Review any remaining hits and classify them as either upstream attribution/history or active identity that still needs migration.
- `SECURITY.md` still needs a small active-identity cleanup; the connector blocked a full-file replacement, so that file may require a smaller local edit.

---

## Previous update: Go test package filter fix

### Validation

Validated from the canonical local update path:

```bash
./update.sh
```

PowerShell validation remains pending unless separately run:

```powershell
.\update.ps1 -SkipAutoPush
```

---

## Previous update: Audit metadata helper refactor

### Validation

Validated from the canonical local update path:

```bash
./update.sh
```
