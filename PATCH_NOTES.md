# Arrakis Control Panel Release Notes

## Current update: Active identity migration

### Why this update was made

Several active build, runtime, test, and setup references still used the legacy `dune-admin` name. Upstream attribution to Icehunter's original `dune-admin` project should remain, but active project identity should use Arrakis Control Panel naming.

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

### Intentionally retained legacy references

The remaining `dune-admin` references should be limited to upstream attribution or historical release notes, especially references to Icehunter's original `dune-admin` project by Ryan Wilson.

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

Validation pending.

Recommended validation:

```bash
./update.sh
```

PowerShell validation should also be rerun:

```powershell
.\update.ps1 -SkipAutoPush
```
