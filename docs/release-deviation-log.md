# Release Deviation Log

## Purpose

This log records deviations from the planned Arrakis Control Panel release train, release scope, validation gates, naming, and security expectations.

A deviation is not automatically a failure. It is a documented decision that changes expected release scope, timing, validation, naming, or risk posture.

## Required fields

Each entry must include:

- date
- release or planned release
- deviation type
- decision
- rationale
- risk impact
- mitigation
- owner
- follow-up target

## Entries

### 2026-06-15 — Broad Go code-quality/refactor review deferred for final 0.1.0

| Field | Value |
|---|---|
| Release or planned release | planned `v0.1.0` |
| Deviation type | Code-quality/refactor review deferral |
| Decision | Defer broad Go code-quality/refactor review beyond current validated secure-baseline slices. |
| Rationale | Final `v0.1.0` should prioritize release stability and avoid late broad structural churn unless a critical issue is found. Existing AppSec, audit, auth-boundary, and blocked-mutation work remains the controlling safety posture. |
| Risk impact | Medium. Deferring a broad review may leave maintainability issues for a later hardening slice, but it avoids introducing unvalidated code churn before final release. |
| Mitigation | Require clean local `./update.sh` validation before tagging. Do not add new endpoints, Live Admin/RMQ execution, Player 360 mutations, or Welcome Kits in final `v0.1.0`. Schedule focused Go code-quality/refactor review for `v0.1.1` or the next hardening slice. |
| Owner | Ron Yacketta |
| Follow-up target | `v0.1.1` or next hardening slice. |

### 2026-06-15 — Broad documentation review deferred beyond primary release/security docs

| Field | Value |
|---|---|
| Release or planned release | planned `v0.1.0` |
| Deviation type | Documentation review deferral |
| Decision | Defer full repository-wide documentation review beyond the primary release, security, roadmap, changelog, patch-note, and deviation-log documents. |
| Rationale | Primary release/security documentation has been reviewed and corrected for the active Arrakis Control Panel identity. A full doc sweep remains useful, but it does not materially change runtime safety for final `v0.1.0`. |
| Risk impact | Low to medium. Some stale or less-used documentation may remain, but primary operator and release-control references are aligned. |
| Mitigation | Keep primary release documents consistent before tagging. Preserve legitimate upstream `dune-admin` attribution. Avoid stale DA Manager wording in active operator guidance. Track full documentation review for `v0.1.1` or the next documentation-hardening slice. |
| Owner | Ron Yacketta |
| Follow-up target | `v0.1.1` or next documentation-hardening slice. |

### 2026-06-06 — Product rename before final 0.1.0

| Field | Value |
|---|---|
| Release or planned release | `v0.1.0-rc.1` / planned `v0.1.0` |
| Deviation type | Product identity / release-label sync |
| Decision | Rename the application from DA Manager to Arrakis Control Panel. |
| Rationale | The project now has a public product identity aligned with the repository name and Dune-inspired operating model. |
| Risk impact | Low for runtime behavior; medium for documentation consistency because stale labels can confuse release notes, support, and operator guidance. |
| Mitigation | Update release policy, release checklist, changelog, patch notes, and remaining documentation/code labels. Treat stale DA Manager labels as a pre-`v0.1.0` cleanup item. |
| Owner | Ron Yacketta |
| Follow-up target | Complete repo-wide label sync before final `v0.1.0`. |

### 2026-06-06 — Security scan evidence deferred for release candidate

| Field | Value |
|---|---|
| Release or planned release | `v0.1.0-rc.1` |
| Deviation type | Validation evidence deferral |
| Decision | Allow `v0.1.0-rc.1` with `govulncheck`, `gosec`, `gitleaks`, `trivy`, and SBOM evidence deferred. |
| Rationale | The release is a pre-1.0 release candidate with clean `./update.sh`, passing npm audit, and documented AppSec slices. |
| Risk impact | Medium. External scan gaps reduce confidence in supply-chain and static-analysis posture. |
| Mitigation | Record deferral in the release checklist and require scan evidence before final `v0.1.0` if feasible. |
| Owner | Ron Yacketta |
| Follow-up target | Before final `v0.1.0`. |

### 2026-06-06 — Update-script modularization accepted as incomplete for RC

| Field | Value |
|---|---|
| Release or planned release | `v0.1.0-rc.1` |
| Deviation type | Tooling/refactor completion deferral |
| Decision | Accept partial update-script modularization for RC scope. |
| Rationale | The canonical update path passed clean, but Bash/PowerShell modularization is not fully completed. |
| Risk impact | Medium. Update-script drift or partial refactor could affect repeatable builds if not completed. |
| Mitigation | Keep this documented as a known issue and continue validation through `./update.sh`; complete Bash and PowerShell modularization before broader release confidence is claimed. |
| Owner | Ron Yacketta |
| Follow-up target | Before final `v0.1.0` or `v0.1.1`, depending on RC findings. |
