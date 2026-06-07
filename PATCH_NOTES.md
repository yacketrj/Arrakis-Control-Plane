# Arrakis Control Panel Release Notes

## Current update: Product rename, release-plan goals, and deviation tracking

### Why this update was made

The application has been renamed to **Arrakis Control Panel**. It is not named DA Manager. Release documentation, planning language, and operator-facing labels must use the new product name going forward.

This update also documents the release plan, release goals, release-label sync rules, industry-standard alignment, current gaps, and deviation logging expectations.

### Product identity

- Product name: `Arrakis Control Panel`
- Repository: `yacketrj/Arrakis-Control-Plane`
- Prior internal/project label: `DA Manager` — deprecated and should be removed from current documentation and code labels.

### Upstream attribution

Arrakis Control Panel is a fork of Icehunter's `dune-admin` project by Ryan Wilson:

```text
https://github.com/Icehunter/dune-admin
```

Every release should preserve clear upstream attribution and state that Arrakis Control Panel builds on Icehunter's original `dune-admin` work.

Future RMQ/live-admin work should also preserve the upstream acknowledgement that the original `dune-admin` README gives to `@adainrivers` and the `dune-dedicated-server-manager` project for RabbitMQ server-command research.

### What changed

- Updated `docs/release-versioning.md` to use Arrakis Control Panel as the canonical product name.
- Added release-train goals for `v0.1.0` through `v0.6.0`.
- Added release-label sync rules for `VERSION`, Git tags, GitHub Release titles, release checklist files, `CHANGELOG.md`, `PATCH_NOTES.md`, and per-slice changelog records.
- Added release deviation policy and criteria.
- Added `docs/release-deviation-log.md`.
- Updated `docs/releases/v0.1.0-rc.1.md` to use Arrakis Control Panel, include upstream attribution, and document current industry-standard release gaps.

### Industry-standard alignment

The current process is aligned with common release-management expectations in these areas:

- semantic-versioning-style labels
- immutable Git tags for release points
- release-candidate flow before final release
- release checklist with validation, rollback, approval, and risk sections
- compact changelog plus durable per-slice records
- explicit known-risk acceptance for RC scope
- upstream attribution preservation

### Current gaps

The current release process is still lacking:

- signed release artifacts
- artifact checksums
- SBOM generation
- consistently attached SAST, DAST, secret-scan, and vulnerability-scan evidence
- automated GitHub Release artifact attachment
- fully completed Bash and PowerShell update-script modularization
- automated post-release verification evidence

### Deviation tracking

Any deviation from the release plan must be cataloged in:

```text
docs/release-deviation-log.md
```

The initial deviation log records:

- product rename from DA Manager to Arrakis Control Panel before final `v0.1.0`
- security scan evidence deferred for `v0.1.0-rc.1`
- update-script modularization accepted as incomplete for RC scope

### Security and operator impact

- No route behavior changed.
- No mutation behavior changed.
- No new endpoint was added.
- Player 360 remains read-only.
- Live Admin / RMQ / Discord full server management remains out of scope for this release candidate.

### Validation

Validation pending from the canonical local update path:

```bash
./update.sh
```

### Remaining rename work

- Continue repo-wide verification for stale `DA Manager` and `Arrakis Control Plane` labels.
- Update remaining documentation/code labels to `Arrakis Control Panel`.
- Keep the release deviation log synchronized if release scope changes.

---

## Previous update: Initial release candidate setup

### Validation

Validated from the canonical local update path:

```bash
./update.sh
```
