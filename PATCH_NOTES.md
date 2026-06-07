# Arrakis Control Panel Release Notes

## Current update: README correction and documentation review plan

### Why this update was made

The canonical update path emitted a non-blocking ledger-size warning for `docs/admin-feature-design-and-priorities.md`. That warning raised a valid concern: large planning documents can become stale, inaccurate, or out of sync without failing validation.

The README was also stale. It did not describe the current `./update.sh` workflow accurately and still contained old product/workflow assumptions.

### What changed

- Updated `README.md` for Arrakis Control Panel naming and current workflow.
- Documented `./update.sh` as the canonical validation/build path.
- Corrected strict admin-token generation guidance to use a 43-character base64url token generated from 32 random bytes.
- Clarified that manual Go/npm/Linux helper commands are fallback/debug/platform-specific paths, not the canonical release workflow.
- Added release evidence locations and release workflow guidance to the README.
- Added upstream attribution in the README.
- Added `docs/documentation-review-plan.md`.

### Documentation review scope

The new documentation review plan covers:

- accuracy against current code, scripts, routes, config, and release workflow
- authenticity of ownership, attribution, release maturity, and compliance claims
- comprehensiveness of operator, security, release, and troubleshooting workflows
- consistency of product naming, version labels, commands, paths, and release scope
- maintainability of large mutable planning/audit documents

### Initial concern

`docs/admin-feature-design-and-priorities.md` is currently large enough to trigger:

```text
Large Markdown notice: docs/admin-feature-design-and-priorities.md has 405 lines. Consider splitting if it is a mutable ledger.
```

The warning is non-blocking, but this file now requires review before final `v0.1.0` acceptance.

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

### Remaining documentation work

- Complete the full documentation review in `docs/documentation-review-plan.md`.
- Review and either split, archive, or correct `docs/admin-feature-design-and-priorities.md`.
- Continue repo-wide verification for stale `DA Manager`, `Arrakis Control Plane`, and outdated workflow labels.

---

## Previous update: Product rename, release-plan goals, and deviation tracking

### Validation

Validation pending from the canonical local update path:

```bash
./update.sh
```
