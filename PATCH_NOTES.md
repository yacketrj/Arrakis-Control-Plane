# Arrakis Control Panel Release Notes

## Current update: Executable rename to Arrakis Control Panel

### Why this update was made

The project label has been renamed to Arrakis Control Panel. The compiled executable should match that product name instead of continuing to build as the upstream-compatible `dune-admin` binary.

### What changed

- Updated Bash update defaults so the default backend binary name is `arrakis-control-panel`.
- Updated `update.sh` so canonical builds emit:
  - `arrakis-control-panel` on Linux/macOS
  - `arrakis-control-panel.exe` on Windows
- Updated `update.ps1` so PowerShell builds emit `arrakis-control-panel.exe`.
- Updated `scripts/linux/build-linux.sh` so Linux helper builds emit `dist/linux/arrakis-control-panel`.
- Updated `README.md` build-output documentation to match the new executable name.

### Compatibility note

The Linux systemd installer may still retain legacy upstream-compatible service path/unit defaults until that service migration is completed and validated. This slice changes compiled build outputs, not installed service semantics.

### Security and operator impact

- No route behavior changed.
- No mutation behavior changed.
- No new endpoint was added.
- Player 360 remains read-only.
- Existing operators should confirm any local scripts that launch `dune-admin` directly and update them to `arrakis-control-panel` after pulling this change.

### Validation

Validated from the canonical local update path:

```bash
./update.sh
```

### Remaining rename work

- Complete systemd installer/path migration if desired as a separate validated slice.
- Continue repo-wide verification for stale `DA Manager`, `Arrakis Control Plane`, and outdated workflow labels.
- Complete the full documentation review in `docs/documentation-review-plan.md`.

---

## Previous update: README correction and documentation review plan

### Validation

Validated from the canonical local update path:

```bash
./update.sh
```
