# Arrakis Control Panel Release Notes

## Current update: Roadmap documentation refresh

### Why this update was made

The documentation review plan identified `docs/admin-feature-design-and-priorities.md` as a large mutable planning document that could become stale. The file still used older Dune Admin-era product naming, described Player 360 as the next implementation slice, and did not reflect the current Arrakis Control Panel release train.

### What changed

- Refreshed `docs/admin-feature-design-and-priorities.md` for Arrakis Control Panel.
- Replaced stale Dune Admin-era planning language with the current `v0.1.x` secure-baseline release track.
- Added the current release-train priorities:
  - `v0.1.0` secure baseline
  - `v0.2.0` Discord Admin Foundation
  - `v0.3.0` Server Lifecycle Management
  - `v0.4.0` Live Admin / RMQ Foundation
  - `v0.5.0` Welcome Kits / Player Requests
  - `v0.6.0` Guild Operations
- Preserved safety principles, item delivery path distinctions, and validation requirements.
- Clarified that full Discord server management, Live Admin RMQ, Welcome Kits, and arbitrary raw command publishing remain out of scope for `v0.1.0`.

### Security and operator impact

- No route behavior changed.
- No mutation behavior changed.
- No new endpoint was added.
- Player 360 remains read-only.
- This improves roadmap accuracy and reduces stale-document risk before final `v0.1.0`.

### Validation

Validation pending from the canonical local update path:

```bash
./update.sh
```

### Remaining documentation work

- Complete the full documentation review in `docs/documentation-review-plan.md`.
- Continue repo-wide verification for stale `DA Manager`, `Arrakis Control Plane`, and outdated workflow labels.
- Review other long-lived docs for stale implemented-vs-planned claims.

---

## Previous update: Linux systemd service migration

### Validation

Validated from the canonical local update path:

```bash
./update.sh
```
