# NIST SSDF Alignment Plan

## Purpose

This document establishes NIST SP 800-218, Secure Software Development Framework (SSDF) Version 1.1, as the primary secure-development baseline for DA Manager.

DA Manager is a high-privilege administrative control plane for a self-hosted Dune: Awakening environment. The project should therefore be managed as a corporate development effort with repeatable controls, release evidence, secure defaults, and documented risk decisions.

## Positioning

SSDF is the best near-term alignment target because it is product- and SDLC-focused. SOC 2 and ISO/IEC 27001 can be mapped later, but both require broader organizational governance and operating evidence. SSDF gives this project an actionable secure-development structure now.

Primary references:

- NIST SP 800-218: https://doi.org/10.6028/NIST.SP.800-218
- NIST SSDF project page: https://csrc.nist.gov/projects/ssdf

## SSDF practice groups

| SSDF group | DA Manager interpretation |
|---|---|
| PO — Prepare the Organization | Define security requirements, roles, release process, development standards, and evidence expectations. |
| PS — Protect the Software | Protect source, secrets, builds, CI/CD, artifacts, and release provenance. |
| PW — Produce Well-Secured Software | Build secure-by-default features, review code, test security controls, and prevent regressions. |
| RV — Respond to Vulnerabilities | Track findings, prioritize remediation, validate fixes, rotate secrets, and publish release notes. |

## Current control posture

| Area | Current posture | SSDF group |
|---|---|---|
| Security remediation tracker | Present and maintained. | PO, RV |
| Audit intake | External security audit findings converted into prioritized remediation. | RV |
| Secure defaults | Loopback default, fail-closed non-loopback exposure, strict token format, one-time WebSocket tickets. | PW |
| SSH trust | Known-hosts validation and Ed25519-only direction established. | PW |
| Token handling | Browser storage moved from localStorage to strict sessionStorage as an interim step. | PW, RV |
| Release validation | Update scripts run Go tests, frontend checks, audit, build, and Git safety behavior. | PW, PS |
| Documentation | User/admin/compliance/security docs are being built and maintained. | PO |
| Supply chain | Scanning exists; action pinning, SBOM, attestations, and binary provenance remain open. | PS |
| Vulnerability response | Findings are tracked, but formal SLA and evidence record are still needed. | RV |

## DA Manager SSDF control objectives

### PO — Prepare the Organization

| Objective | Required implementation | Evidence |
|---|---|---|
| Define secure-development policy | Adopt SSDF as project baseline. Maintain security posture rule: strict secure defaults unless usability impact is material. | This document; security remediation tracker. |
| Define roles and ownership | Identify code owner, release owner, security reviewer, operator/admin role. | Admin guide; release checklist. |
| Define release process | Use ITIL-style change/release workflow: change summary, risk, validation, rollback, known issues, approval. | Release notes; changelog; release checklist. |
| Define security requirements | Maintain requirements for SSH, token handling, transport exposure, mutation accountability, logging, and supply chain. | Security docs; control matrix. |
| Define evidence expectations | Every release must capture validation commands, scan output, known risks, and rollback instructions. | Release checklist and release notes. |

### PS — Protect the Software

| Objective | Required implementation | Evidence |
|---|---|---|
| Protect source code | Use GitHub branch protections, reviewed changes, and no generated binaries in source. | Repository settings evidence; changelog. |
| Protect secrets | Never commit `.env`, SSH keys, DB passwords, admin tokens, logs, or crash dumps. | `.gitignore`; gitleaks output. |
| Protect build pipeline | Pin GitHub Actions to full commit SHAs. | Workflow review evidence. |
| Protect artifacts | Generate checksums, SBOM, and attestations for release artifacts. | Release artifacts and provenance records. |
| Protect dependencies | Run `govulncheck`, `npm audit`, Trivy, Gitleaks, CodeQL, and dependency review. | CI evidence; release checklist. |

### PW — Produce Well-Secured Software

| Objective | Required implementation | Evidence |
|---|---|---|
| Secure authentication | Strict admin token format now; future memory-only or HttpOnly session-cookie auth. | Auth tests; browser storage verification. |
| Secure SSH | Ed25519 private keys, Ed25519 host keys, mandatory known_hosts validation, fail closed on mismatch. | Unit/integration tests; gosec G106 enabled. |
| Secure transport | Loopback-only default; non-loopback requires explicit reverse-proxy/TLS acknowledgment. | Startup validation test. |
| Secure WebSockets | Static admin token must never appear in WebSocket URLs; one-time scoped tickets only. | Browser Network tab check; backend test. |
| Secure mutation workflows | High/destructive mutations require reasons and audit records by default. | Mutation tests; audit log review. |
| Secure logging | Redact credentials, tokens, account IDs where appropriate, and sensitive paths before logging. | Log review evidence; tests. |
| Secure SQL access | Read-only SQL by default, explicit denylist for dangerous operations, schema identifier sanitization. | SQL tests. |

### RV — Respond to Vulnerabilities

| Objective | Required implementation | Evidence |
|---|---|---|
| Intake findings | Capture external audit results and operator-reported defects into a remediation tracker. | Security remediation tracker. |
| Prioritize risk | Rank by severity, exploitability, impact, and release blocking status. | Remediation tracker. |
| Validate remediation | Every security fix needs test/build evidence and manual verification where applicable. | Release checklist. |
| Rotate secrets | Rotate admin token, SSH keys, and DB password after token/transport or log exposure fixes. | Admin guide and release notes. |
| Prevent recurrence | Add regression tests, CI gates, and documentation updates after each fix. | Test results; CI workflow. |
| Communicate release risk | Release notes must include security changes, known issues, and operator action items. | Release notes. |

## SSDF-aligned release gate

A DA Manager release should not be treated as production-ready until the following are true:

```text
[ ] Go tests pass
[ ] Frontend typecheck passes
[ ] Frontend lint passes
[ ] Frontend build passes
[ ] npm audit --audit-level=high passes
[ ] govulncheck ./... passes
[ ] gosec ./... passes without excluding SSH host-key validation findings
[ ] gitleaks detect --source . passes
[ ] trivy fs --severity HIGH,CRITICAL . passes or exceptions are documented
[ ] Static admin token is absent from WebSocket URLs
[ ] Non-loopback LISTEN_ADDR fails closed unless explicit secure exposure setting is present
[ ] SSH host-key mismatch is rejected
[ ] High/destructive mutation without reason is rejected when enforcement is enabled
[ ] Release notes updated
[ ] CHANGELOG.md updated
[ ] Known risks and rollback plan documented
```

## Current priority backlog under SSDF

| Priority | SSDF group | Work item | Status |
|---|---|---|---|
| P0 | PW | Validate Ed25519-only SSH and mandatory known_hosts implementation. | In progress |
| P0 | PW | Validate strict backend admin-token enforcement. | In progress |
| P0 | PW | Validate one-time WebSocket ticket flow. | Implemented, needs local validation |
| P0 | PW | Validate fail-closed backend exposure behavior. | Implemented, needs local validation |
| P0 | RV | Rotate admin token and SSH/DB secrets after token-path changes. | Pending operator execution |
| P1 | PW/RV | Default mutation reason enforcement to enabled. | Pending |
| P1 | PS | Pin GitHub Actions to full commit SHAs. | Pending |
| P1 | PS | Generate SBOM and release attestations. | Pending |
| P1 | PS | Remove committed binaries or formalize signed/attested release artifacts. | Pending |
| P1 | PO | Complete admin guide, user guide, release notes, changelog, and release checklist. | In progress |

## Relationship to SOC 2 and ISO/IEC 27001

SSDF alignment does not equal SOC 2 or ISO/IEC 27001 certification. It does, however, produce the engineering controls and evidence that support those frameworks later.

| Framework | How SSDF helps |
|---|---|
| SOC 2 Security | Supports secure development, change management, access controls, vulnerability management, and incident response evidence. |
| ISO/IEC 27001 | Supports secure development lifecycle controls, risk treatment, supplier/dependency management, technical vulnerability management, change management, and evidence collection. |

## Operating rule

When there is a choice between permissive compatibility and stronger security, DA Manager will select the stronger security option unless it materially breaks core app usage. Exceptions must be documented as risk acceptances.
