# Arrakis Control Panel Release Notes

## Current update: README prerequisites and required tooling

### Why this update was made

The README needed a clear, operator-facing prerequisites section that lists the local applications, command-line tools, server-side access, and optional release/security tools required to build, validate, and operate Arrakis Control Panel.

### What changed

- Added `Final release-gate disposition` to the README release-reference list.
- Added a new `Prerequisites and required tooling` section to `README.md`.
- Documented required local applications:
  - Git
  - Go 1.26.3 or compatible newer Go toolchain
  - Node.js with npm
  - Bash
  - PowerShell 7+
  - OpenSSH client tools
  - modern browser
- Documented required server-side access:
  - Dune: Awakening self-hosted server host
  - SSH user and Ed25519 private key
  - PostgreSQL-compatible game database access
  - admin token
  - allowed frontend origins
- Documented optional integrations and release tools:
  - Discord OAuth and bot credentials
  - GitHub CLI
  - optional release/security evidence tools such as `govulncheck`, `gosec`, `gitleaks`, `trivy`, and `syft`
  - reverse proxy, VPN, and TLS tooling
- Added quick prerequisite check commands for Git, Go, Node, npm, OpenSSH, and PowerShell.

### Security and operator impact

- No route behavior changed.
- No mutation behavior changed.
- No endpoint was added.
- No Live Admin / RMQ execution was added.
- No Player 360 mutation behavior was added.
- No Welcome Kits behavior was added.
- Operator setup expectations are clearer before running validation or connecting to a server.

### Validation

This was a documentation-only README update prepared through the GitHub connector. Local validation should be run before tagging or releasing:

```bash
./update.sh
```

On Windows, also run:

```powershell
.\update.ps1 -SkipAutoPush
```

### Remaining final-`v0.1.0` gates

- Post-release verification after tag/artifact install or launch.

---

## Previous update: Clean local validation recorded for final v0.1.0 gate disposition

### Why this update was made

The release owner reported that local validation completed cleanly after the documentation-only final gate-disposition update.

### What changed

- Updated `docs/final-v0.1.0-gate-status.md` to close the local-validation gate for the gate-disposition update.
- Recorded the validation as operator-reported clean local validation on 2026-06-15.
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

Operator-reported clean local validation:

```bash
./update.sh
```

PowerShell validation path remains available on Windows:

```powershell
.\update.ps1 -SkipAutoPush
```

### Remaining final-`v0.1.0` gates

- Post-release verification after tag/artifact install or launch.

---

## Previous update: Final v0.1.0 gate disposition

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

### Remaining final-`v0.1.0` gates at that time

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
