# Arrakis Control Panel Release Notes

## Current update: Refactor gate and application identity cleanup

### Why this update was made

Update-script modularization is now explicitly required before final `v0.1.0` unless it is intentionally deferred in the release deviation log. The Go backend also needs a code-quality and modularization review before final `v0.1.0` so the project does not carry stale product strings, scattered constants, or avoidable structural debt into the first accepted release.

During the first Go review pass, stale `dune-admin` identity strings were found in startup logging, public status, and setup repair guidance.

### What changed

- Added `app_identity.go` with shared application identity constants:
  - `appDisplayName`
  - `appServiceName`
  - `appWindowsExecutable`
- Updated `server.go` startup logging to use the shared service identity.
- Updated `/api/v1/public/status` to return Arrakis Control Panel identity values.
- Updated `main.go` setup repair guidance to use the shared Windows executable name.
- Updated `docs/release-versioning.md` to make these final-`v0.1.0` gates explicit:
  - update-script modularization decision
  - PowerShell modularization or documented deferral
  - Go refactoring/code-quality review
  - handler/route/audit/mutation-safety review
  - typed/allowlisted execution surface review

### Security and operator impact

- No route behavior changed.
- No mutation behavior changed.
- No new endpoint was added.
- Player 360 remains read-only.
- Public status now reports the current service/product identity instead of the legacy upstream-compatible name.
- This reduces stale-label risk and centralizes product identity for future refactors.

### Validation

Validated from the canonical local update path:

```bash
./update.sh
```

### Remaining refactor work

- Continue update-script modularization review.
- Review PowerShell script modularization or document a final `v0.1.0` deferral.
- Continue Go review for route grouping, handler boundaries, audit/mutation-safety helper boundaries, and typed execution surfaces.

---

## Previous update: Roadmap documentation refresh and README path cleanup

### Validation

Validated from the canonical local update path:

```bash
./update.sh
```
