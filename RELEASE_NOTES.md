# DA Manager Release Notes

## Purpose

Release notes communicate operationally relevant changes to administrators, operators, and security reviewers. They should be updated for every release and should summarize what changed, why it matters, validation performed, known risks, and required operator actions.

## Current release: Unreleased Security Governance and Hardening Baseline

### Release classification

| Field | Value |
|---|---|
| Change type | Normal change |
| Risk level | High until validation completes |
| Release theme | NIST SSDF alignment, security hardening, corporate release discipline |
| Primary audience | Administrators, operators, security reviewers |
| Production readiness | Not production-ready until validation checklist is complete |

### Executive summary

This release establishes DA Manager as a corporate-style secure-development effort aligned to NIST SP 800-218 SSDF. It adds formal security governance documentation, release evidence expectations, risk tracking, and several high-priority security hardening changes.

The release addresses major audit-driven concerns around SSH trust, token exposure, WebSocket token leakage, unsafe backend exposure, browser token persistence, and release evidence maturity.

### Security changes

| Area | Change | Operator impact |
|---|---|---|
| Secure-development baseline | Adopted NIST SSDF as the primary control baseline. | Future releases must include control/evidence updates. |
| WebSocket log streaming | Replaced static `ws_token` URL usage with short-lived one-time stream tickets. | Browser Network tab should no longer show static admin token in WebSocket URLs. |
| Backend exposure | Non-loopback `LISTEN_ADDR` now fails closed unless explicitly acknowledged. | Direct public backend exposure is blocked by default. |
| Frontend token handling | Browser Access Key validation is stricter and token storage is session-scoped. | Operators may need to re-enter token after browser/tab restart. |
| SSH hardening | App direction is Ed25519-only SSH client keys and known_hosts validation. | Operators must configure Ed25519 keys and host-key trust. |
| Compliance posture | Added clear statement that SOC 2/ISO are not yet met; SSDF readiness is the working target. | Avoid overclaiming compliance readiness. |

### Documentation changes

Added or updated:

- `docs/NIST_SSDF_ALIGNMENT.md`
- `docs/COMPLIANCE_READINESS.md`
- `docs/ADMIN_GUIDE.md`
- `docs/RELEASE_CHECKLIST.md`
- `docs/CONTROL_MATRIX.md`
- `docs/RISK_REGISTER.md`
- `SECURITY_REMEDIATION_TODO.md`
- `CHANGELOG.md`
- `.env.example`

### Required operator actions

Before using this release in a real administrative environment:

1. Pull the latest repository state.
2. Generate or confirm Ed25519 SSH key usage.
3. Add the remote host's Ed25519 host key to `known_hosts`.
4. Generate a strict 43-character admin token.
5. Keep backend bound to `127.0.0.1:8080` unless behind approved secure exposure controls.
6. Run the release checklist.
7. Rotate credentials after token/transport validation is complete.

### Required validation

Minimum validation before release acceptance:

```bash
go test -v ./...
govulncheck ./...
gosec ./...
gitleaks detect --source .
trivy fs --severity HIGH,CRITICAL .
cd web && npm audit --audit-level=high && npm run typecheck && npm run lint && npm run build
```

Manual validation:

- Confirm WebSocket URL contains a one-time ticket, not `ws_token` or `ADMIN_TOKEN`.
- Confirm consumed WebSocket ticket cannot be reused.
- Confirm expired WebSocket ticket fails.
- Confirm `LISTEN_ADDR=0.0.0.0:8080` fails startup without explicit secure exposure setting.
- Confirm weak/malformed `ADMIN_TOKEN` fails validation.
- Confirm non-Ed25519 SSH client key is rejected.
- Confirm changed or unknown SSH host key is rejected.

### Known risks

| Risk | Status |
|---|---|
| Browser token remains JavaScript-readable during active session. | Accepted interim state; future memory-only or HttpOnly secure session-cookie model planned. |
| Mutation reason enforcement is not yet default-on for all high/destructive actions. | Planned P1 remediation. |
| CI action pinning, SBOM, and artifact attestations are not complete. | Planned P1 remediation. |
| Generated binary provenance requires cleanup or formal release artifact process. | Planned P1 remediation. |
| Full SOC 2 / ISO readiness is not achieved. | Documented; SSDF readiness is the target. |

### Rollback notes

If startup or authentication fails after this release:

1. Confirm `.env` values against `.env.example`.
2. Confirm `ADMIN_TOKEN` is 43 base64url characters.
3. Confirm `LISTEN_ADDR=127.0.0.1:8080`.
4. Confirm `SSH_KEY` points to an Ed25519 private key.
5. Confirm `SSH_KNOWN_HOSTS` exists and contains the remote Ed25519 host key.
6. Revert to the prior known-good commit only if required after collecting logs and validation output.

### Acceptance criteria

This release should not be considered accepted until:

- Release checklist is completed.
- High-risk validation gates pass or have documented risk acceptance.
- Release notes and changelog are current.
- Risk register is reviewed.
- Security remediation tracker is updated.
