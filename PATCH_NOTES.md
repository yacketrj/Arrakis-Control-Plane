# Dune Admin Release Notes

## Current update: GitHub CI validation workflows

### Why this update was made

Local `update.ps1` validation remains important because it tests the operator's Windows development environment, but it cannot be run from the assistant runtime. This update adds GitHub-hosted validation workflows so pushes and pull requests receive independent Linux and Windows compile/test signals.

### What changed

- Added `.github/workflows/ci-linux.yml`.
- Added `.github/workflows/ci-windows.yml`.
- Both workflows run on push, pull request, and manual dispatch.
- Both workflows validate:
  - Go formatting
  - Go module tidiness
  - Go module verification
  - `go vet ./...`
  - `go test -v ./...`
  - backend build
  - frontend dependency install with `npm ci`
  - frontend audit with `npm audit --audit-level=high`
  - frontend typecheck
  - frontend lint
  - frontend build
- Linux builds `dist/linux/dune-admin`.
- Windows builds `dist/windows/dune-admin.exe`.

### Security and operator impact

- GitHub CI now provides a VM-like validation signal for repository pushes.
- Local `update.ps1` remains the final environment-specific Windows validation path.
- CI does not replace local testing for machine-specific PATH, locked files, local credentials, or operator runtime configuration.

### Validation

Validation runs automatically on push and pull request. Manual validation is available from the GitHub Actions tab with:

```text
CI Linux Validation
CI Windows Validation
```

Local validation remains:

```powershell
.\update.ps1
```

---

## Previous update: Blueprint import shared mutation confirmation migration

Blueprint import was migrated to shared mutation confirmation with required admin reason capture.

---

## Previous update: Battlegroup Exec shared mutation confirmation migration

Battlegroup Exec server-control actions were migrated to shared mutation confirmation with required admin reason capture.

---

## Previous update: Database SQL shared mutation confirmation migration

Database SQL was migrated to shared mutation confirmation with required admin reason capture.

---

## Previous update: Storage shared mutation confirmation migration

Storage container add/remove item operations were migrated to shared mutation confirmation with required admin reason capture.

---

## Previous update: Confirmed player admin actions workflow

Player admin actions were migrated to a dedicated confirmed modal with required admin reason capture.

---

## Previous update: Confirmed player move workflow

Player move actions were migrated to a dedicated confirmed modal with required admin reason capture and online-state safeguards.

---

## Previous update: Journey node shared mutation confirmation migration

Journey node complete/reset actions were migrated to the confirmed Player Actions modal with required admin reason capture.

---

## Previous update: Player resource/action shared mutation confirmation migration

Resource, XP, specialization, and faction reputation actions were migrated to the confirmed Player Actions modal with required admin reason capture.
