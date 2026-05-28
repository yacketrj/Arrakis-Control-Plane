# DA Manager SSDF Control Matrix

## Purpose

This matrix maps DA Manager controls to NIST SP 800-218 SSDF practice groups and provides a working evidence model for release readiness, security review, and future SOC 2 / ISO alignment.

## Control status legend

| Status | Meaning |
|---|---|
| Implemented | Control is implemented and documented. |
| Implemented - Needs Evidence | Control exists but needs test or release evidence. |
| In Progress | Control design or implementation is active. |
| Planned | Control is accepted for backlog but not started. |
| Deferred | Control is intentionally delayed with risk acceptance. |

## Control matrix

| Control ID | SSDF Group | Control | Implementation expectation | Current status | Evidence source |
|---|---|---|---|---|---|
| PO-01 | PO | Secure-development baseline | NIST SSDF is adopted as the project secure-development baseline. | Implemented | `docs/NIST_SSDF_ALIGNMENT.md` |
| PO-02 | PO | Security posture rule | Strict secure defaults are selected unless they materially break core app usage. | Implemented | `docs/NIST_SSDF_ALIGNMENT.md`; security tracker |
| PO-03 | PO | Corporate documentation set | Maintain user guide, admin guide, release checklist, release notes, changelog, control matrix, and risk register. | In Progress | `docs/` and root docs |
| PO-04 | PO | Change/release discipline | Releases include change summary, risk, impact, validation, rollback, known issues, and approval record. | Implemented - Needs Evidence | `docs/RELEASE_CHECKLIST.md` |
| PO-05 | PO | Control ownership | Release owner, security reviewer, and operations owner are identified for each release. | Planned | Release checklist |
| PO-06 | PO | Compliance readiness | Maintain readiness position for SOC 2 and ISO/IEC 27001 mapping. | Implemented | `docs/COMPLIANCE_READINESS.md` |
| PS-01 | PS | Secret exclusion | `.env`, SSH keys, DB passwords, admin tokens, and generated secrets must not be committed. | In Progress | `.gitignore`; Gitleaks results |
| PS-02 | PS | Secret scanning | Run Gitleaks or equivalent before release. | Planned | Release checklist evidence |
| PS-03 | PS | Dependency scanning | Run `govulncheck`, `npm audit`, Trivy, and dependency review. | Planned | Release checklist evidence |
| PS-04 | PS | Static security scanning | Run `gosec ./...` and do not suppress host-key validation findings after remediation. | Planned | Release checklist evidence |
| PS-05 | PS | CI workflow hardening | Pin GitHub Actions to full commit SHAs. | Planned | Workflow review |
| PS-06 | PS | SBOM generation | Generate SBOM for release artifacts. | Planned | Release artifacts |
| PS-07 | PS | Artifact provenance | Generate checksums, signatures, or attestations for release artifacts. | Planned | Release artifacts |
| PS-08 | PS | Binary provenance | Generated binaries should not be committed unless signed/attested release artifacts are documented. | Planned | Repository review |
| PW-01 | PW | Loopback backend default | Backend defaults to `127.0.0.1` binding. | Implemented | `.env.example`; server startup behavior |
| PW-02 | PW | Fail-closed remote exposure | Non-loopback `LISTEN_ADDR` fails closed unless explicit secure exposure setting is present. | Implemented - Needs Evidence | `exposure.go`; release checklist |
| PW-03 | PW | Strict admin token format | `ADMIN_TOKEN` must be 43 base64url characters generated from 32 random bytes. | Implemented - Needs Evidence | `auth.go`; release checklist |
| PW-04 | PW | Frontend token validation | Browser Access Key rejects weak, malformed, whitespace, and placeholder values. | Implemented - Needs Evidence | `web/src/api/client.ts` |
| PW-05 | PW | Browser token persistence reduction | Frontend token is stored in `sessionStorage`, not persistent `localStorage`, as an interim control. | Implemented - Needs Evidence | `web/src/api/client.ts` |
| PW-06 | PW | Future secure session model | Replace browser-stored token with memory-only or HttpOnly secure session-cookie model. | Planned | Security tracker |
| PW-07 | PW | WebSocket token protection | Static admin token must not appear in WebSocket URLs. | Implemented - Needs Evidence | `ws_ticket.go`; `handlers_logs.go`; frontend Logs tab |
| PW-08 | PW | One-time WebSocket tickets | Log streaming uses scoped, short-lived, one-time tickets. | Implemented - Needs Evidence | `ws_ticket.go`; release checklist |
| PW-09 | PW | SSH client key restriction | DA Manager SSH client key must be Ed25519. | In Progress | `ssh.go`; `main.go`; release checklist |
| PW-10 | PW | SSH host-key validation | SSH host identity must be validated through `known_hosts`. | In Progress | `ssh.go`; release checklist |
| PW-11 | PW | Ed25519 host-key restriction | Remote SSH host key must be Ed25519. | In Progress | `ssh.go`; release checklist |
| PW-12 | PW | Mutation accountability | High/destructive mutations require reasons and audit records by default. | Planned | Mutation tests; audit log review |
| PW-13 | PW | SQL safety | SQL interface permits read-only operations by default and rejects dangerous operations. | In Progress | SQL handler/tests |
| PW-14 | PW | CORS restriction | Only configured allowed origins receive browser access. | Implemented - Needs Evidence | `server.go`; `auth.go` |
| PW-15 | PW | Request size limits | JSON request bodies are size limited. | Implemented - Needs Evidence | `auth.go`; handlers using `limitBody` |
| PW-16 | PW | Logging redaction | Credentials, tokens, and sensitive identifiers are redacted from logs. | Planned | Redaction tests; log review |
| RV-01 | RV | Security finding intake | External audit findings are converted into remediation work. | Implemented | `SECURITY_REMEDIATION_TODO.md` |
| RV-02 | RV | Vulnerability prioritization | Findings are ranked by exploitability, impact, and release risk. | Implemented | `SECURITY_REMEDIATION_TODO.md` |
| RV-03 | RV | Fix validation | Remediation requires build/test/security evidence. | Implemented - Needs Evidence | Release checklist |
| RV-04 | RV | Secret rotation | Tokens, keys, and passwords are rotated after exposure-path changes. | Planned | Release notes; admin evidence |
| RV-05 | RV | Regression prevention | Security fixes add tests or checks to prevent recurrence. | Planned | Test evidence |
| RV-06 | RV | Release communication | Release notes communicate security changes, known risks, and operator actions. | In Progress | `RELEASE_NOTES.md` |

## Minimum release evidence set

Every release should attach or reference:

- Go test output.
- Frontend typecheck/lint/build output.
- npm audit output.
- govulncheck output.
- gosec output.
- gitleaks output.
- Trivy output or documented exception.
- Manual validation notes for WebSocket ticketing, non-loopback startup failure, SSH host-key rejection, and admin-token validation.
- Release notes.
- Changelog entry.
- Risk register updates.
- Rollback plan.

## Open control priorities

| Priority | Control | Reason |
|---|---|---|
| P0 | SSH validation evidence | SSH is a privileged trust boundary. |
| P0 | Token and WebSocket validation evidence | Token leakage was a high-impact audit finding. |
| P0 | Non-loopback fail-closed evidence | Direct exposure of admin backend is a material risk. |
| P1 | Mutation reason default enforcement | Supports accountability and change traceability. |
| P1 | CI/SBOM/provenance controls | Required for software supply-chain maturity. |
| P1 | Risk register and release records | Required for corporate development discipline. |
