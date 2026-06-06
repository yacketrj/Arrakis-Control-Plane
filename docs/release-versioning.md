# Release Versioning Policy

## Purpose

DA Manager uses Semantic Versioning-style release numbers with pre-1.0 stability rules.

The project is currently in pre-1.0 hardening. Releases are intended to produce repeatable artifacts, validation evidence, rollback instructions, and release notes before expanding new high-risk features such as Discord-driven server management or Live Admin RMQ controls.

## Version format

Use:

```text
MAJOR.MINOR.PATCH[-PRERELEASE][+BUILD]
```

Examples:

```text
0.1.0-rc.1
0.1.0
0.1.1
0.2.0
1.0.0
```

## Pre-1.0 rules

Until `1.0.0`, the public/operator contract is still stabilizing.

| Change type | Version action | Example |
|---|---|---|
| First release candidate | `0.1.0-rc.1` | First complete validation candidate |
| Release-candidate fix | Increment rc number | `0.1.0-rc.2` |
| First accepted release | Drop rc suffix | `0.1.0` |
| Patch/security fix after release | Increment patch | `0.1.1` |
| New validated feature set | Increment minor | `0.2.0` |
| Stable operator contract | Promote to 1.0 | `1.0.0` |

## Tag format

Git tags must use a leading `v`:

```text
v0.1.0-rc.1
v0.1.0
v0.1.1
```

## Branching model

For the current project size:

- `main` remains the integration branch.
- Create a release candidate tag from `main` only after `./update.sh` passes.
- Do not add new feature scope after the release candidate tag unless the release is intentionally deferred.
- Only security, build, validation, release-note, and documentation fixes should land between `v0.1.0-rc.N` and `v0.1.0`.

## Release gates

A release candidate requires:

- `VERSION` updated.
- `CHANGELOG.md` current and compact.
- `PATCH_NOTES.md` current for the active release slice.
- `docs/releases/<version>.md` created from the release checklist template.
- `./update.sh` passes.
- Security validation status recorded.
- Known issues and residual risks documented.
- Rollback plan documented.

## First release recommendation

Use:

```text
0.1.0-rc.1
```

Rationale:

- The app has meaningful hardening and validation evidence.
- The release checklist exists.
- Player 360 remains read-only.
- Several AppSec slices are validated.
- Some governance and security gates still need formal evidence before a stable `0.1.0` release.

## Do not use 1.0 yet

Do not promote to `1.0.0` until:

- update scripts are stable and fully validated on intended platforms
- release artifacts are generated outside the repo or are signed/attested
- SAST/DAST/dependency evidence is consistently recorded
- backup/restore evidence is recorded
- admin auth/session model is stable
- Discord-driven server management, if added, has complete RBAC/audit/allowlist controls
