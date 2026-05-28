# DA Manager Release Checklist

## Purpose

This checklist defines the minimum release gate for DA Manager. It follows NIST SSDF principles and ITIL-style release/change management practices.

A release is not considered production-ready until validation evidence, risk notes, rollback instructions, and documentation updates are complete.

## Release metadata

| Field | Value |
|---|---|
| Release version | TBD |
| Release date | TBD |
| Release owner | TBD |
| Security reviewer | TBD |
| Change type | Standard / Normal / Emergency |
| Risk rating | Low / Medium / High / Critical |
| Target environment | Development / Test / Production |
| Release branch/tag | TBD |
| Build artifact | TBD |

## Change summary

Document the business and operational reason for the release.

```text
Summary:

Scope:

Out of scope:
```

## Risk and impact assessment

| Question | Response |
|---|---|
| Does this affect authentication or authorization? | TBD |
| Does this affect SSH, tunneling, or database connectivity? | TBD |
| Does this affect player mutation workflows? | TBD |
| Does this affect audit logging? | TBD |
| Does this affect frontend token handling? | TBD |
| Does this affect deployment/network exposure? | TBD |
| Does this introduce new dependencies? | TBD |
| Does this require secret rotation? | TBD |

## Required validation gates

### Build and test

| Gate | Command | Result | Evidence location |
|---|---|---|---|
| Go tests | `go test -v ./...` | TBD | TBD |
| Go build | `go build -trimpath ./...` or update script | TBD | TBD |
| Frontend install | `cd web && npm install` | TBD | TBD |
| Frontend typecheck | `cd web && npm run typecheck` | TBD | TBD |
| Frontend lint | `cd web && npm run lint` | TBD | TBD |
| Frontend build | `cd web && npm run build` | TBD | TBD |

### Security validation

| Gate | Command | Result | Evidence location |
|---|---|---|---|
| Go vulnerability scan | `govulncheck ./...` | TBD | TBD |
| Go static security scan | `gosec ./...` | TBD | TBD |
| Secret scan | `gitleaks detect --source .` | TBD | TBD |
| Filesystem/container scan | `trivy fs --severity HIGH,CRITICAL .` | TBD | TBD |
| npm audit | `cd web && npm audit --audit-level=high` | TBD | TBD |
| Dependency review | GitHub dependency review or equivalent | TBD | TBD |
| SBOM generation | `syft . -o spdx-json` or equivalent | TBD | TBD |

### Manual security checks

| Check | Expected result | Result |
|---|---|---|
| Browser Network tab contains no `ws_token` | Pass | TBD |
| WebSocket URL contains only one-time stream ticket | Pass | TBD |
| Reusing a consumed WebSocket ticket fails | Pass | TBD |
| Expired WebSocket ticket fails | Pass | TBD |
| `LISTEN_ADDR=0.0.0.0:8080` without override fails startup | Pass | TBD |
| Non-loopback bind with documented override starts only when intentionally configured | Pass | TBD |
| Unknown or changed SSH host key is rejected | Pass | TBD |
| Non-Ed25519 SSH client key is rejected | Pass | TBD |
| Weak or malformed `ADMIN_TOKEN` is rejected | Pass | TBD |
| Frontend rejects weak or malformed Browser Access Key | Pass | TBD |
| High/destructive mutation without reason is rejected when enforcement is enabled | Pass | TBD |

## Documentation updates

| Document | Required update | Status |
|---|---|---|
| `README.md` | Update if startup, setup, or usage changed. | TBD |
| `docs/USER_GUIDE.md` | Update operator workflow changes. | TBD |
| `docs/ADMIN_GUIDE.md` | Update deployment/security changes. | TBD |
| `docs/NIST_SSDF_ALIGNMENT.md` | Update control status. | TBD |
| `SECURITY_REMEDIATION_TODO.md` | Mark completed/remediated items. | TBD |
| `docs/COMPLIANCE_READINESS.md` | Update readiness posture if material. | TBD |
| `CHANGELOG.md` | Add technical change entries. | TBD |
| `RELEASE_NOTES.md` | Add operator-facing release notes. | TBD |
| `docs/RISK_REGISTER.md` | Add/update risks and residual risk. | TBD |
| `docs/CONTROL_MATRIX.md` | Update control implementation/evidence. | TBD |

## Rollback plan

```text
Rollback trigger:

Rollback steps:

Expected rollback duration:

Data compatibility concerns:

Operator communications:
```

## Secret rotation checklist

Complete when the release changes token handling, transport exposure, SSH configuration, logging, or credential controls.

| Secret | Rotate? | Evidence |
|---|---|---|
| `ADMIN_TOKEN` | TBD | TBD |
| DA Manager SSH key | TBD | TBD |
| DB password | TBD | TBD |
| GitHub tokens/secrets | TBD | TBD |
| Reverse proxy credentials/certificates | TBD | TBD |

## Known issues and risk acceptance

| Issue | Severity | Risk decision | Owner | Review date |
|---|---|---|---|---|
| TBD | TBD | Mitigate / Accept / Transfer / Avoid | TBD | TBD |

## Approval

| Role | Name | Date | Approval |
|---|---|---|---|
| Release owner | TBD | TBD | Pending |
| Security reviewer | TBD | TBD | Pending |
| Operations owner | TBD | TBD | Pending |

## Post-release verification

| Check | Result | Notes |
|---|---|---|
| Backend starts cleanly | TBD | TBD |
| Frontend loads cleanly | TBD | TBD |
| Connectivity diagnostics pass | TBD | TBD |
| Audit events are written | TBD | TBD |
| Logs stream using one-time tickets | TBD | TBD |
| No new high/critical findings | TBD | TBD |

## Release decision

```text
Approved / Rejected / Deferred

Decision rationale:

Follow-up actions:
```
