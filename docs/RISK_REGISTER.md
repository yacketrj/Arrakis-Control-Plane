# DA Manager Risk Register

## Purpose

This register tracks security, operational, release, and compliance risks for DA Manager. It supports NIST SSDF vulnerability response, ITIL-style change management, and future SOC 2 / ISO/IEC 27001 evidence mapping.

## Risk rating model

| Rating | Criteria |
|---|---|
| Critical | Likely compromise of admin control plane, secrets, database, or live service integrity. |
| High | Material exploitability or operational outage risk requiring release-blocking remediation. |
| Medium | Meaningful security or operational risk with compensating controls or limited blast radius. |
| Low | Hygiene, documentation, or process issue with limited immediate impact. |

## Risk treatment options

| Treatment | Meaning |
|---|---|
| Mitigate | Implement controls to reduce likelihood or impact. |
| Accept | Formally accept residual risk with owner and review date. |
| Transfer | Move risk to another party or service model. |
| Avoid | Remove the activity or feature causing the risk. |

## Active risks

| ID | Risk | Rating | Affected area | Treatment | Owner | Status | Review date | Notes |
|---|---|---|---|---|---|---|---|---|
| R-001 | SSH host identity is not correctly validated, enabling man-in-the-middle risk against the privileged SSH control path. | High | SSH / DB tunnel / runtime commands | Mitigate | Security reviewer | In progress | TBD | Ed25519 and known_hosts hardening is in progress; requires validation evidence. |
| R-002 | Static admin token could be exposed through browser storage, logs, screenshots, or operator handling. | High | Auth / frontend / operations | Mitigate | Release owner | In progress | TBD | Moved from localStorage to sessionStorage as interim control; future memory-only or HttpOnly session model remains planned. |
| R-003 | Admin token in WebSocket URLs could leak through logs, browser history, proxies, or diagnostics. | High | Logs / WebSocket / auth | Mitigate | Release owner | Implemented - needs evidence | TBD | One-time WebSocket tickets implemented; requires validation evidence. |
| R-004 | Backend could be exposed beyond loopback without TLS/VPN/reverse proxy controls. | High | Deployment / transport | Mitigate | Operations owner | Implemented - needs evidence | TBD | Non-loopback binding now fails closed unless explicit secure exposure setting is configured. |
| R-005 | High/destructive mutations may lack sufficient reason/audit context. | Medium | Player/server mutation workflows | Mitigate | Release owner | Planned | TBD | Default reason enforcement should be enabled. |
| R-006 | Generated binaries in source control could create provenance and trust ambiguity. | Medium | Release / supply chain | Mitigate | Release owner | Planned | TBD | Remove committed binaries or formalize signed/attested artifact process. |
| R-007 | CI/CD supply-chain risk from unpinned third-party GitHub Actions. | Medium | CI/CD | Mitigate | Release owner | Planned | TBD | Pin Actions to full commit SHAs. |
| R-008 | Insufficient SBOM and artifact attestation evidence. | Medium | Supply chain / release | Mitigate | Release owner | Planned | TBD | Generate SBOM and provenance artifacts during release. |
| R-009 | Secrets could be committed or logged during active troubleshooting. | High | Secrets / logs / Git | Mitigate | Security reviewer | Planned | TBD | Require Gitleaks and redaction tests before release. |
| R-010 | Backup/restore process is not proven. | Medium | Operations / availability | Mitigate | Operations owner | Planned | TBD | Requires documented restore test and evidence. |
| R-011 | Home-hosted deployment may lack volumetric DDoS protection. | Medium | Availability / hosting | Accept or transfer | Operations owner | Open | TBD | Use ISP mitigation, VPN, or datacenter/colo if public exposure risk increases. |
| R-012 | Documentation may drift from implemented behavior during rapid hardening. | Medium | Governance / operations | Mitigate | Release owner | In progress | TBD | User guide, admin guide, release checklist, control matrix, changelog, release notes must remain current. |
| R-013 | No formal incident response procedure evidence exists yet. | Medium | Incident response | Mitigate | Security reviewer | Planned | TBD | Admin guide contains baseline process; dedicated IR playbook may be needed. |
| R-014 | SOC 2 / ISO readiness could be overstated without governance and evidence maturity. | Medium | Compliance | Mitigate | Security reviewer | Implemented - monitor | TBD | Compliance readiness document states not audit-ready; SSDF is active baseline. |

## Residual risk acceptance log

| ID | Accepted risk | Rating | Rationale | Owner | Review date |
|---|---|---|---|---|---|
| TBD | TBD | TBD | TBD | TBD | TBD |

## Risk review cadence

Risk register review is required:

- Before every release.
- After any security incident.
- After external audit or security review intake.
- After changes to auth, SSH, transport exposure, logging, CI/CD, or deployment topology.
- At least quarterly for ongoing operation.

## Release gate linkage

Any open **Critical** or **High** risk requires one of the following before release:

1. Mitigation completed and validated.
2. Formal risk acceptance with owner and review date.
3. Release deferred.

Risk acceptance must not be implicit.
