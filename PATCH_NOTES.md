# Dune Admin Release Notes

## Current update: Initial release candidate setup

### Why this update was made

The project is in a reasonable development cycle to prepare its first controlled release, but it should be treated as a pre-1.0 release candidate rather than a stable `1.0.0` release. The current security posture has meaningful hardening and validation evidence, while some release evidence and security scans still need to be recorded before a final `0.1.0` release.

### Version decision

- Initial release candidate: `0.1.0-rc.1`
- Git tag target after validation: `v0.1.0-rc.1`
- First accepted release target after release-candidate validation: `v0.1.0`

### What changed

- Added `VERSION` with `0.1.0-rc.1`.
- Added `docs/release-versioning.md` with the release numbering, tag, release-candidate, and pre-1.0 policy.
- Added `docs/releases/v0.1.0-rc.1.md` as the first release checklist instance.
- The release checklist records:
  - release metadata
  - scope and out-of-scope items
  - risk and impact assessment
  - build/test gates
  - security validation gates
  - manual security checks
  - rollback plan
  - secret rotation considerations
  - known risks
  - approval and post-release verification fields

### Security and operator impact

- No route behavior changed.
- No mutation behavior changed.
- No new endpoint was added.
- Player 360 remains read-only.
- Live Admin / RMQ / Discord full server management remains out of scope for this release candidate.
- This creates release discipline before expanding high-risk Discord/server-management features.

### Validation

Validation pending from the canonical local update path:

```bash
./update.sh
```

The release candidate should not be tagged until validation evidence is recorded in:

```text
docs/releases/v0.1.0-rc.1.md
```

---

## Previous update: Route-specific audit target assertions

### Why this update was made

The AppSec hardening track is continuing before new Live Admin / RMQ / Welcome Kit features. The previous slices proved that high-risk/destructive mutations emit audit events and that blocked mutations are still auditable. This slice tightened target accountability by ensuring audit events capture route-specific identifiers needed to investigate player, item, server-command, vehicle, and guild mutations.

### Validation

Validated from the canonical local update path:

```bash
./update.sh
```

Non-blocking build-performance warning observed:

```text
[PLUGIN_TIMINGS] Your build spent significant time in plugin `@tailwindcss/vite:generate:build`.
```
