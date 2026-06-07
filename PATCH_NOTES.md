# Dune Admin Release Notes

## Current update: Initial release candidate setup

### Why this update was made

The project is in a reasonable development cycle to prepare its first controlled release, but it should be treated as a pre-1.0 release candidate rather than a stable `1.0.0` release. The current security posture has meaningful hardening and validation evidence, while some release evidence and security scans still need to be recorded before a final `0.1.0` release.

### Version decision

- Initial release candidate: `0.1.0-rc.1`
- Git tag target: `v0.1.0-rc.1`
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
- Updated the release checklist with clean `./update.sh` validation and approval to tag `v0.1.0-rc.1`.

### Security and operator impact

- No route behavior changed.
- No mutation behavior changed.
- No new endpoint was added.
- Player 360 remains read-only.
- Live Admin / RMQ / Discord full server management remains out of scope for this release candidate.
- This creates release discipline before expanding high-risk Discord/server-management features.

### Validation

Validated from the canonical local update path:

```bash
./update.sh
```

The release checklist now marks `v0.1.0-rc.1` as approved to tag:

```text
docs/releases/v0.1.0-rc.1.md
```

### Next release action

Create and push the annotated tag from the validated local checkout:

```bash
git tag -a v0.1.0-rc.1 -m "DA Manager v0.1.0-rc.1"
git push origin v0.1.0-rc.1
```

---

## Previous update: Route-specific audit target assertions

### Validation

Validated from the canonical local update path:

```bash
./update.sh
```

Non-blocking build-performance warning observed:

```text
[PLUGIN_TIMINGS] Your build spent significant time in plugin `@tailwindcss/vite:generate:build`.
```
