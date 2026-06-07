# Documentation Review Plan

## Purpose

This plan defines a deep documentation review for Arrakis Control Panel before final `v0.1.0` acceptance.

The trigger for this review was a ledger-size warning for `docs/admin-feature-design-and-priorities.md`. The warning did not fail validation, but it indicates that large planning documents may become stale, duplicative, or inconsistent with the actual code and release workflow.

## Review objectives

The documentation review must verify:

- accuracy: documentation matches the current code, scripts, routes, configuration, and release workflow
- authenticity: documentation preserves upstream attribution and does not overstate ownership, security posture, release maturity, or compliance readiness
- comprehensiveness: critical operator, security, release, and troubleshooting workflows are documented without relying on tribal knowledge
- consistency: product naming, version labels, paths, commands, and release scope are aligned across documents
- maintainability: large mutable documents are split, indexed, or converted into durable per-slice records when appropriate

## Required naming standard

Current product name:

```text
Arrakis Control Panel
```

Deprecated labels that must be reviewed and replaced where they refer to the current product:

```text
DA Manager
Arrakis Control Plane
Dune Admin Fork
```

Compatibility labels such as `dune-admin`, `dune-admin.exe`, systemd unit names, and upstream project references may remain when they refer to executable names, existing service paths, upstream attribution, or backward-compatible implementation details.

## Required upstream attribution check

Every release and primary operator-facing document should preserve clear attribution:

- Arrakis Control Panel is a fork of Icehunter's `dune-admin` project by Ryan Wilson.
- Future RMQ/live-admin documentation must preserve the upstream acknowledgement to `@adainrivers` and the `dune-dedicated-server-manager` project for RabbitMQ server-command research.

## Scope

Review all Markdown documentation and release records, including but not limited to:

- `README.md`
- `PATCH_NOTES.md`
- `CHANGELOG.md`
- `docs/RELEASE_CHECKLIST.md`
- `docs/release-versioning.md`
- `docs/release-deviation-log.md`
- `docs/releases/*.md`
- `docs/changelog/**/*.md`
- `docs/appsec-endpoint-audit.md`
- `docs/admin-feature-design-and-priorities.md`
- `docs/admin-implementation-tasks.md`
- `docs/discord-auth.md`
- `docs/discord-player-links.md`
- `docs/inventory-requests-orders.md`
- `docs/inventory-studio.md`
- `docs/player-360-profile.md`
- `docs/browser-token-cors-security.md`
- `docs/database-endpoint-security.md`
- `docs/infrastructure-log-endpoint-security.md`
- `docs/linux.md`
- setup/deployment docs

## Review checklist

### Product and release labels

- [ ] Product name is `Arrakis Control Panel` where referring to the current application.
- [ ] Deprecated labels are removed or explicitly marked as historical/internal.
- [ ] `VERSION`, tags, release checklist, changelog, patch notes, and GitHub Release title agree.
- [ ] `v0.1.0-rc.1` evidence is not confused with final `v0.1.0` evidence.

### Workflow accuracy

- [ ] `./update.sh` is documented as the canonical validation/build path.
- [ ] Manual `go`, `npm`, and Linux helper commands are clearly marked as debugging or platform-specific alternatives.
- [ ] PowerShell support is described accurately and not overstated until fully validated.
- [ ] Build output paths and executable names match current behavior.
- [ ] Token-generation instructions produce a strict-valid 43-character base64url token.

### Security accuracy

- [ ] Backend localhost-default binding is documented correctly.
- [ ] Remote exposure requires explicit reverse proxy/TLS/identity controls.
- [ ] Player 360 remains read-only.
- [ ] Browser token/session limitations are not understated.
- [ ] High-risk/destructive mutation reason, preview, and audit behavior is documented accurately.
- [ ] No document suggests arbitrary raw command publishing through UI or Discord.
- [ ] Discord/server-management/RMQ features are identified as future high-risk work unless implemented and validated.

### AppSec and compliance claims

- [ ] Documentation does not claim SOC 2, ISO, FedRAMP, or NIST compliance as complete.
- [ ] Control mappings are described as readiness/alignment evidence, not certification.
- [ ] Deferred SAST, DAST, secret scan, vulnerability scan, and SBOM evidence is identified clearly.
- [ ] Release-candidate risk acceptance is distinct from final-release acceptance.

### Feature accuracy

- [ ] Current implemented features are distinguished from roadmap/planned features.
- [ ] Discord auth/player-link/self-service boundaries are accurate.
- [ ] Inventory requests/orders are described as coordination-only unless code proves mutation behavior.
- [ ] Inventory Studio mutation workflows match implemented safety controls.
- [ ] Guild, Welcome Kit, RMQ, and full Discord server-management features are not documented as complete until implemented.

### Large-document maintainability

- [ ] Any mutable Markdown file over the ledger warning threshold is reviewed for stale or duplicate content.
- [ ] Large planning documents are split into index plus smaller per-topic records when they are active mutable ledgers.
- [ ] Long historical content is archived by commit reference or compact monthly archive, not repeatedly edited as one large file.

## Initial concern: `docs/admin-feature-design-and-priorities.md`

This file triggered a large Markdown notice at 405 lines.

Required review actions:

1. Determine whether it is current strategy, historical planning, or mixed content.
2. Verify it matches the current release train in `docs/release-versioning.md`.
3. Move obsolete or superseded planning into archive notes.
4. Split active roadmap items into smaller topic files if it remains mutable.
5. Ensure planned features are not presented as implemented features.

## Review deliverables

The review should produce:

- a summary of reviewed files
- a list of stale or inaccurate statements found
- a list of documents requiring immediate correction
- a list of documents to split or archive
- release-impact classification for each finding
- validation notes after corrections

## Release gate

Before final `v0.1.0`, this review should be completed or explicitly deferred in `docs/release-deviation-log.md`.

For `v0.1.0-rc.1`, the review is an accepted follow-up item.

## Validation

After documentation corrections, run:

```bash
./update.sh
```

If the change is documentation-only and the full validation path was already run for the same release candidate, at minimum run:

```bash
bash scripts/check-ledger-size.sh
```
