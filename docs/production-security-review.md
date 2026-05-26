# Production Security Review and Release Gate

## Release posture

```text
NOT PRODUCTION READY
```

Dune Admin is a high-risk administrative tool. It can expose player identity, inventory, progression, database state, server operations, logs, and mutation workflows. Production release is blocked until the security gates in this document are closed or formally risk-accepted.

## Release-blocking security gates

### P0: Configuration validation

Required outcomes:

- Setup validates all required values before writing `.env`.
- Runtime validates all required values read from `.env`.
- SSH private key paths support documented formats only.
- Windows `%USERPROFILE%\.ssh\id_rsa` style paths are supported.
- Unsupported PowerShell expression paths such as `$env:USERPROFILE\.ssh\id_rsa` are rejected with a clear error.
- Missing or unreadable SSH keys block startup when SSH is required.
- Invalid ports and invalid tunnel modes block startup.

### P0: Secret-at-rest handling

Required outcomes:

- `.env` must not store plaintext DB passwords in production mode.
- `.env` must not store plaintext admin tokens in production mode.
- Existing plaintext `.env` values may remain supported for development with clear warnings.
- A documented migration path must exist from plaintext `.env` to encrypted or OS-protected secret storage.

Acceptable options include Windows DPAPI, OS keychain/credential manager integration, or operator-provided external secret environment variables.

Not acceptable for production:

- base64 treated as encryption
- reversible encryption with the key stored next to the encrypted value
- logging decrypted secrets

### P0: SSH tunnel and startup behavior

Required outcomes:

- If `SSH_TUNNEL_MODE=auto`, startup establishes the managed SSH tunnel before DB-dependent features are marked available.
- If `SSH_TUNNEL_MODE=existing`, startup validates that the configured local tunnel endpoint is reachable.
- If `SSH_TUNNEL_MODE=off`, direct DB access is explicit and documented.
- Failed required tunnel setup must block or clearly gate DB-dependent functionality.
- UI must show a blocking state when `db_connected=false`.

### P0: Database discovery safety

Required outcomes:

- Document the discovery flow: SSH to configured host, detect runtime, locate DB service/pod/container, establish tunnel or direct endpoint, connect PostgreSQL.
- Discovery errors identify whether failure happened at SSH, runtime detection, DB service discovery, tunnel setup, or DB authentication.
- `database pod not found in cluster` must be supplemented with actionable diagnostics.
- Kubernetes and Docker discovery paths must both be reviewed.

### P0: Authentication and authorization

Required outcomes:

- Admin token authentication is consistently enforced on protected API routes.
- Public-safe routes are explicitly identified and reviewed.
- Default admin token behavior is safe.
- Setup generates a strong token when missing.
- Frontend token storage is documented and reviewed.

### P0: Mutation safety and auditability

Required outcomes:

- Every mutating route is inventoried.
- Every high-risk or destructive route requires backend-side admin reason enforcement in production mode.
- Every high-risk or destructive route writes audit records.
- Frontend confirmation is treated as UX only; backend enforcement remains authoritative.

### P0: Raw SQL execution controls

Required outcomes:

- Raw SQL execution is disabled by default in production mode or replaced with approved typed routines.
- If enabled, raw SQL requires explicit configuration, admin reason, audit logging, and operator warning.
- Dangerous statements are blocked or require break-glass mode.

### P0: Battlegroup command/config controls

Required outcomes:

- Battlegroup command execution is disabled by default in production mode or constrained to allowlisted commands.
- Shell command input is reviewed for injection and expansion risk.
- UserEngine.ini and UserGame.ini modification is treated as high-risk game/economy configuration management.
- Config editing must include discovery, backup, diff, validation, confirmed deployment, rollback, and audit logging.

### P1: Frontend security and UX gating

Required outcomes:

- Protected views clearly show disconnected/degraded state.
- DB-dependent tabs are gated when DB is unavailable.
- Player 360 and Inventory Studio use searchable player selection; raw actor ID remains fallback only.
- Raw JSON displays are reviewed for sensitive data exposure.
- Browser-local snapshots are documented as sensitive files.

### P1: CORS and network exposure

Required outcomes:

- Allowed origins are explicit in production mode.
- Wildcard CORS is not allowed in production mode.
- Listen address defaults and non-loopback exposure are reviewed.

### P1: Logging and data leakage

Required outcomes:

- Logs do not print DB passwords, admin tokens, SSH private key contents, decrypted secrets, or full connection strings.
- Errors remain actionable without exposing secrets.
- Log viewer is reviewed for sensitive server/operator data.

### P1: Supply-chain and CI/CD

Required outcomes:

- CI runs Go tests, Go vet, frontend typecheck, lint, build, and dependency audit.
- High/critical dependency vulnerabilities are reviewed.
- Release artifacts are generated reproducibly.
- Build scripts do not embed secrets.

### P1: Vulnerability intake workflow

Required outcomes:

- Vulnerability emails are triaged with severity, affected component, reproduction status, exploitability, fix owner, and target date.
- Security fixes are tracked separately from feature work.
- Release is blocked by open critical/high findings unless explicitly accepted.

## Required review artifacts

Before production release, produce:

- threat model
- protected route inventory
- mutating route inventory
- authn/authz review notes
- secret-handling review notes
- tunnel/startup review notes
- SQL and command-exec review notes
- frontend data-exposure review notes
- vulnerability triage register
- release sign-off checklist

## Release decision rule

Production release is blocked until every release-blocking gate is either fixed and validated or explicitly risk-accepted with owner, rationale, compensating controls, and expiration date.
