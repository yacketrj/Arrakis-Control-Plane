# Arrakis Control Panel Release Notes

## Current update: Linux systemd service migration

### Why this update was made

The compiled backend executable has already been renamed to `arrakis-control-panel`. The Linux systemd installer still used legacy upstream-compatible `dune-admin` defaults for service name, install path, service user, installed binary path, and unit `ExecStart`. This slice aligns the Linux service installer with the Arrakis Control Panel product name.

### What changed

- Updated `scripts/linux/install-systemd.sh` defaults:
  - service name: `arrakis-control-panel`
  - install directory: `/opt/arrakis-control-panel`
  - service user/group: `arrakis-control-panel`
  - source binary: `dist/linux/arrakis-control-panel`
  - installed binary: `/opt/arrakis-control-panel/arrakis-control-panel`
  - unit description: `Arrakis Control Panel backend`
  - `ExecStart`: `/opt/arrakis-control-panel/arrakis-control-panel`
- Updated `README.md` systemd install commands and default service/path documentation.

### Compatibility note

Existing systems that already installed the legacy `dune-admin` service should migrate intentionally. Stop the legacy service, copy/verify `.env`, then enable/start the new `arrakis-control-panel` service. Do not run both services against the same host/configuration unless intentionally testing.

### Security and operator impact

- No route behavior changed.
- No mutation behavior changed.
- No new endpoint was added.
- Player 360 remains read-only.
- Linux service installs now use product-aligned names by default.
- Existing operators must verify migration steps before replacing a working legacy service.

### Validation

Validation pending from the canonical local update path:

```bash
./update.sh
```

### Remaining rename work

- Continue repo-wide verification for stale `DA Manager`, `Arrakis Control Plane`, and outdated workflow labels.
- Complete the full documentation review in `docs/documentation-review-plan.md`.

---

## Previous update: Executable rename to Arrakis Control Panel

### Validation

Validated from the canonical local update path:

```bash
./update.sh
```
