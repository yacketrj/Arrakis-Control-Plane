# Compliance Readiness Assessment

## Executive assessment

DA Manager is not currently SOC 2 ready or ISO/IEC 27001 ready as a complete organizational program.

The application has made meaningful progress toward secure-product controls, but SOC 2 and ISO/IEC 27001 are not achieved by code hardening alone. Both require documented governance, risk management, control ownership, operating evidence, change management, incident response, access review, vendor management, backup/recovery evidence, and repeatable audit artifacts.

Current posture: **pre-audit hardening / control design in progress**.

## Audit perspective

If reviewed today by an external auditor, the likely result would be:

| Framework | Likely result | Reason |
|---|---|---|
| SOC 2 Type I | Not ready | Some technical controls exist, but formal control design, policies, control owners, scope, evidence, and operating procedures are incomplete. |
| SOC 2 Type II | Not ready | Type II requires operating effectiveness evidence over time. Current controls are still being designed and implemented. |
| ISO/IEC 27001:2022 certification | Not ready | There is no documented ISMS scope, risk assessment, risk treatment plan, Statement of Applicability, internal audit program, management review, or corrective-action process. |

## Current strengths

| Area | Current strength |
|---|---|
| Secure defaults | Backend defaults to loopback binding. Non-loopback exposure fails closed unless explicitly acknowledged. |
| SSH trust | Direction is strong: Ed25519-only SSH keys and mandatory known_hosts validation. |
| Token handling | Admin token format is being tightened to 32 random bytes encoded as 43 base64url characters. |
| Browser token storage | Frontend moved from persistent localStorage to stricter sessionStorage as an interim measure. |
| WebSocket exposure | Static admin token has been removed from WebSocket URLs and replaced with one-time stream tickets. |
| Auditability | Security remediation tracker exists and is being maintained. |
| Release process | Update scripts run build and validation gates before committing. |
| Change discipline | Work is being prioritized by severity and release-blocking risk. |

## Primary gaps

| Gap | Impact | Priority |
|---|---|---|
| No formal ISMS scope | ISO/IEC 27001 cannot be assessed without defined scope and boundaries. | P0 for certification readiness |
| No control matrix | SOC 2 and ISO evidence cannot be mapped cleanly to implemented controls. | P0 |
| No evidence register | Implemented controls are not yet tied to repeatable audit evidence. | P0 |
| No formal risk register | Risks are discussed and tracked, but not in an auditable risk-management format. | P0 |
| No access review process | Admin access and key/token handling need periodic review and approval evidence. | P0 |
| No incident response procedure | Security incidents require documented triage, containment, eradication, recovery, and post-incident review. | P0 |
| No formal change/release management record | Commits exist, but ITIL-style change classification, approval, risk, validation, and rollback evidence are incomplete. | P0 |
| No backup/restore evidence | Backups and restore tests must be documented and periodically validated. | P0 |
| No vulnerability management SLA | Scanning and remediation need defined severity targets and evidence. | P1 |
| No vendor/dependency governance | Dependency and toolchain risk need documented review and acceptance. | P1 |
| No security training/acceptable-use evidence | Required for many organizational audit programs. | P2 |

## SOC 2 readiness view

SOC 2 evaluates service organization controls using Trust Services Criteria. The most relevant initial scope for DA Manager would be Security, with Availability and Confidentiality as likely future criteria.

| SOC 2 area | Current state | Readiness |
|---|---|---|
| Control environment | Informal project discipline exists; formal policies and accountability do not. | Low |
| Communication and information | Documentation is being created; evidence trail is incomplete. | Medium-low |
| Risk assessment | Security findings are prioritized, but no formal risk register exists. | Medium-low |
| Monitoring activities | Some validation gates exist; control monitoring is incomplete. | Medium-low |
| Logical access | Improving through strict tokens, session-scoped browser storage, SSH key controls, and loopback defaults. | Medium |
| Change management | Git history and update scripts exist; ITIL-style change records and release approvals are incomplete. | Medium-low |
| System operations | Diagnostics and operational runbooks are being built. | Medium-low |
| Incident response | Not yet sufficiently documented. | Low |
| Vendor/dependency management | npm audit exists; SBOM, attestations, dependency review, and action pinning remain open. | Low-medium |

## ISO/IEC 27001 readiness view

ISO/IEC 27001 requires an information security management system, not just secure code.

| ISO/IEC 27001 area | Current state | Readiness |
|---|---|---|
| Context and scope | Not formally defined. | Low |
| Leadership and responsibilities | Informal ownership exists; formal roles and responsibilities are not documented. | Low |
| Planning and risk treatment | Security backlog exists; formal risk assessment and risk treatment plan are missing. | Low-medium |
| Support and competence | Development expertise exists; training/evidence process is not documented. | Low |
| Operations | Deployment and operations guidance is being created. | Medium-low |
| Performance evaluation | Validation scripts exist; internal audit and metrics review do not. | Low |
| Improvement | Remediation tracker exists; corrective-action process needs formalization. | Medium-low |
| Annex A controls | Several technical controls are being implemented, but the Statement of Applicability is missing. | Low-medium |

## Minimum path to audit readiness

### Phase 1: Governance foundation

1. Define system and audit scope.
2. Create control matrix mapped to SOC 2 Security and ISO/IEC 27001:2022 clauses/Annex A.
3. Create risk register and risk treatment plan.
4. Assign control owners.
5. Create evidence register.
6. Create formal change/release management procedure.
7. Create incident response procedure.
8. Create access review procedure.

### Phase 2: Technical control completion

1. Complete SSH known_hosts and Ed25519-only enforcement validation.
2. Complete strict backend admin-token enforcement validation.
3. Remove static tokens from all URLs and logs.
4. Add security headers and CSP.
5. Add SBOM generation and artifact attestations.
6. Pin CI actions to immutable commit SHAs.
7. Remove committed binaries or move them to signed/attested release artifacts.
8. Add backup and restore validation evidence.

### Phase 3: Operating evidence

1. Maintain release notes and changelog for every release.
2. Maintain change records with risk, validation, and rollback notes.
3. Record access reviews.
4. Record vulnerability scans and remediation decisions.
5. Record backup/restore tests.
6. Record incident response tabletop or real incident reviews.
7. Perform internal audit against the control matrix.
8. Perform management review and approve residual risk.

## Recommended audit target

The most realistic near-term target is not certification. It is **SOC 2 / ISO-aligned control readiness**.

Recommended target state:

```text
Security controls designed, documented, implemented, and producing repeatable evidence.
```

Once that state is stable for at least one release cycle, the project can decide whether to pursue formal SOC 2 or ISO/IEC 27001 work.

## Bottom line

DA Manager is moving in the right direction technically, but it would not pass SOC 2 or ISO/IEC 27001 today. The secure-development posture is improving quickly; the main deficit is governance, evidence, and operating process maturity.
