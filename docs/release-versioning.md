# Release Versioning Policy

## Purpose

Arrakis Control Panel uses Semantic Versioning-style release numbers with pre-1.0 stability rules.

The project is currently in pre-1.0 hardening. Releases are intended to produce repeatable artifacts, validation evidence, rollback instructions, release notes, and explicit scope control before expanding new high-risk features such as Discord-driven server management or Live Admin RMQ controls.

## Upstream attribution requirement

Arrakis Control Panel is a fork of Icehunter's `dune-admin` project by Ryan Wilson:

```text
https://github.com/Icehunter/dune-admin
```

Every release must preserve clear upstream attribution in release notes or release evidence. The attribution should state that Arrakis Control Panel builds on Icehunter's original `dune-admin` work.

Future RMQ/live-admin work should also preserve the upstream acknowledgement that the original `dune-admin` README gives to `@adainrivers` and the `dune-dedicated-server-manager` project for RabbitMQ server-command research.

## Release goals

### Primary goals before `v1.0.0`

- Keep the backend localhost-bound by default unless explicitly configured otherwise.
- Keep Player 360 read-only until identity mapping, self-service controls, and mutation safety are fully validated.
- Require reason/preview/audit coverage for high-risk and destructive mutations.
- Maintain generated route auth-boundary coverage.
- Maintain compact release records and per-slice durable changelog records.
- Preserve upstream attribution.
- Treat Discord-driven server management and RMQ/live-admin controls as high-risk features requiring dedicated release trains.

### Release-train themes

| Target | Theme | Goal |
|---|---|---|
| `v0.1.0` | Secure baseline | Accept the current hardened admin foundation after RC verification. |
| `v0.2.0` | Discord Admin Foundation | Add safe Discord RBAC, command registry, previews, and audit without raw command execution. |
| `v0.3.0` | Server Lifecycle Management | Expose safe allowlisted Docker/K8s/local lifecycle operations through UI/Discord. |
| `v0.4.0` | Live Admin / RMQ Foundation | Add controlled RMQ envelope/command execution for online-safe operations only. |
| `v0.5.0` | Welcome Kits / Player Requests | Add audited kits, eligibility, ledger, and Discord approval workflows. |
| `v0.6.0` | Guild Operations | Add audited guild create/delete, membership, and rank workflows. |

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

## Label and scope sync rules

For every release train, keep these labels synchronized:

- `VERSION`
- Git tag name
- GitHub Release title
- `docs/releases/<version>.md`
- `CHANGELOG.md`
- `PATCH_NOTES.md`
- active per-slice changelog records under `docs/changelog/unreleased/`

A release is out of sync if any of these refer to different versions, different scope, or different validation status.

## Deviation policy

Any deviation from the planned release train must be cataloged before continuing work.

A deviation includes:

- adding a feature that is outside the current release theme
- accepting a skipped security gate
- changing release scope after an RC tag
- moving a high-risk feature earlier than planned
- changing version numbering or tag policy
- releasing without planned artifacts or evidence
- changing the backend exposure model
- adding Discord, RMQ, database, or infrastructure mutation capability without a dedicated security review
- product naming changes that affect documentation, code labels, release evidence, or GitHub release text

Deviation entries must include:

- date
- release or planned release
- deviation type
- decision
- rationale
- risk impact
- mitigation
- owner
- follow-up target

Use `docs/release-deviation-log.md` for the durable deviation log.

## Industry-standard alignment

The current process is aligned with common release-management practices in these areas:

- semantic versioning-style labels
- immutable Git tags for release points
- release-candidate flow before final release
- release checklist with risk, validation, rollback, and approval sections
- compact changelog plus durable per-slice records
- explicit known-risk acceptance for RC scope
- upstream attribution preservation

Current gaps before a stronger final/stable release:

- signed release artifacts are not yet required
- SBOM generation is deferred
- SAST, DAST, secret scanning, and vulnerability scan evidence is not consistently attached
- artifact checksums are not yet generated and published
- update-script modularization is not fully completed across Bash and PowerShell
- post-release verification evidence is still manual
- GitHub Release artifact attachment is not automated

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
- Upstream attribution included or referenced.
- Deviation log updated when scope or process differs from the plan.

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
