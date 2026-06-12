# Security Audit Remediation Tracker

This tracker converts the external security audit into implementation work for Arrakis Control Panel.

## Release policy

The application is a high-privilege administrative control plane. Treat every item marked **P0** as release-blocking for any production or remotely exposed deployment.

## Priority order

| Rank | Priority | Severity | Area | Finding | Required outcome | Validation gate |
|---:|---|---|---|---|---|---|
| 1 | P0 | High | SSH trust | Backend SSH client must validate the remote host key. | Replace permissive SSH host-key behavior with `knownhosts`, pinned host keys, or an explicit first-use pinning flow that fails closed on mismatch. Re-enable `gosec` G106 after remediation. | `gosec ./...` runs without suppressing G106 and no unsafe SSH host-key finding remains. Add a unit/integration test that rejects an unknown or changed host key. |
| 2 | P0 | High | Admin token handling | Static admin token must not be stored long-term in browser `localStorage` or placed into WebSocket URLs. | Remove token query parameters from log-stream WebSocket flow. Use `sessionStorage` only as the current incremental hardening step; replace with short-lived one-time WS tickets or secure cookie/session flow later. Rotate existing tokens after patch. | Inspect WebSocket handshake: no token appears in URL. Browser storage does not contain long-lived privileged tokens. Auth tests cover status, HTTP routes, and WebSocket stream. |
| 3 | P0 | Medium | Transport exposure | Non-loopback backend exposure must fail closed unless an explicit secure deployment mode is configured. | If `LISTEN_ADDR` is not loopback, require an explicit remote/TLS/proxy acknowledgment setting. Keep loopback as default. Document reverse-proxy/TLS requirement. | Starting with `LISTEN_ADDR=0.0.0.0:8080` and no explicit secure-mode setting aborts startup. HTTPS/WSS reverse-proxy mode is documented and tested. |
| 4 | P0 | Medium | Secrets and logs | Rotate credentials that may have been exposed during testing and prevent future leakage. | Rotate `ADMIN_TOKEN`, DB password, and SSH key if any were present in logs, browser URLs, screenshots, crash reports, or reverse-proxy logs. Ensure redaction covers token, DB password, SSH paths where needed, account IDs, and player identifiers. | `gitleaks detect --source .` passes. Review runtime logs and browser network traces for token/credential leakage. |
| 5 | P1 | Low/Medium | Mutation accountability | Mutation reason enforcement should be enabled by default for high/destructive actions. | Default `ADMIN_REQUIRE_REASON=true`. Server rejects empty/whitespace reasons for high/destructive mutation classes. UI requires a reason before submitting those actions. | Tests assert sensitive mutation routes reject missing reason and accept non-empty reason. Audit records include reason. |
| 6 | P1 | Medium | CI supply chain | GitHub Actions and release provenance need hardening. | Pin third-party Actions to full commit SHAs, add SBOM generation, add artifact attestations, and restore excluded security rules after code fixes. | CI workflow uses immutable action refs. SBOM and attestation artifacts are published for releases. `govulncheck`, `npm audit`, `gosec`, `gitleaks`, `trivy`, and CodeQL pass. |
| 7 | P1 | Medium | Binary provenance | Committed binaries should be removed or treated as signed release artifacts. | Remove `arrakis-control-panel.exe` from source control unless there is a documented signed/attested binary release process. Keep generated binaries ignored. | `git ls-files` does not include generated binaries, or release docs include checksums/signatures/attestations. |
| 8 | P1 | Architecture | Browser token storage | `sessionStorage` is an incremental hardening step, not final secure storage. | Replace browser-stored admin token with memory-only token handling or a short-lived HttpOnly secure session-cookie model. Prefer HttpOnly cookie-backed sessions with CSRF protections when remote UI deployment is supported. | Refresh/tab-close behavior is documented. Browser dev tools must not show a long-lived privileged token in `localStorage` or `sessionStorage` after the final model is implemented. |
| 9 | P2 | Architecture | Auth model | Static shared admin-token auth should be replaced long term. | Move toward short-lived server-issued sessions, OIDC/identity-provider integration, or secure cookie-backed admin sessions with CSRF protections where applicable. | Threat model updated. End-to-end auth tests cover login/session expiry/revocation. |
| 10 | P2 | Dependency assurance | Frontend dependency evidence should be SBOM-backed. | Generate full Go/npm SBOM and dependency review output instead of relying on manifest-level review. | SBOM checked into release artifacts; dependency review runs in CI. |

## Implementation sequencing

### Phase 1: Immediate release blockers

1. Fix SSH host-key validation.
2. Remove admin token from WebSocket query strings.
3. Stop persisting long-lived admin tokens in browser `localStorage`; current interim state is `sessionStorage` with strict token validation.
4. Enforce fail-closed startup for non-loopback backend binding without explicit secure deployment mode.
5. Rotate operational credentials after token/transport changes land.

### Phase 2: Near-term hardening

1. Default mutation reason enforcement to enabled.
2. Add server and UI tests for high/destructive mutation reason requirements.
3. Replace `sessionStorage` admin token handling with memory-only or HttpOnly secure session-cookie auth.
4. Pin GitHub Actions to commit SHAs.
5. Add SBOM and artifact attestation workflow steps.
6. Remove committed binary artifacts or move them to signed release artifacts.

### Phase 3: Longer-term architecture

1. Replace static token authentication with session-based or identity-provider-backed authentication.
2. Add admin role separation if the app grows beyond a single trusted operator model.
3. Add release threat model documentation for local-only, private-network, and internet-exposed deployment modes.

## Validation checklist

Run these after the corresponding fixes land:

```bash
govulncheck ./...
gosec ./...
gitleaks detect --source .
trivy fs --severity HIGH,CRITICAL .
cd web && npm audit --audit-level=high && npm run typecheck && npm run lint && npm run build
```

Manual checks:

- Open browser dev tools and verify no privileged token appears in WebSocket URLs.
- Start backend with non-loopback listen address and no secure deployment override; startup must abort.
- Send high/destructive mutation without a reason; request must be rejected.
- Connect to a host with a missing or mismatched SSH host key; connection must be rejected.
- Confirm the interim browser token is not stored in `localStorage`; final architecture should remove long-lived browser-accessible token storage entirely.

## Notes

- The audit did not identify confirmed vulnerable direct Go dependency versions in the visible dependency set.
- The highest-risk issues are repository-specific design/implementation issues rather than NVD-tracked dependency CVEs.
- Treat token rotation as mandatory after removing URL/local-storage token exposure paths.
- Current frontend token state is intentionally interim: strict-format `sessionStorage` improves posture over `localStorage`, but memory-only or HttpOnly secure session-cookie auth remains the desired end state.
