# DA Manager Admin Guide

## Purpose

This guide defines how DA Manager should be deployed, configured, secured, operated, and recovered. It is written for system administrators, release owners, and security reviewers.

DA Manager is a high-privilege administrative control plane. Treat access to the backend, frontend session, SSH key, database credentials, and audit logs as privileged administrative access.

## Control baseline

DA Manager aligns secure-development and operational practices to NIST SP 800-218, Secure Software Development Framework (SSDF). SOC 2 and ISO/IEC 27001 mappings are future compliance targets, but SSDF is the active engineering baseline.

Primary documents:

- `docs/NIST_SSDF_ALIGNMENT.md`
- `SECURITY_REMEDIATION_TODO.md`
- `docs/COMPLIANCE_READINESS.md`
- `docs/RELEASE_CHECKLIST.md`
- `CHANGELOG.md`
- `RELEASE_NOTES.md`

## Supported deployment posture

The default supported deployment posture is local-first:

```text
Browser UI on localhost
  -> DA Manager backend on 127.0.0.1:8080
  -> SSH tunnel to remote host
  -> PostgreSQL/database endpoint reachable only through SSH tunnel
```

Do not expose the backend directly to the public internet.

## Required tools

Required for development and release validation:

- Git
- Go
- Node.js/npm
- OpenSSH client
- PowerShell on Windows
- Bash on Linux/macOS/WSL

Optional but recommended:

- GitHub CLI
- Gitleaks
- govulncheck
- gosec
- Trivy
- CodeQL or GitHub code scanning
- SBOM tooling such as Syft

## SSH security requirements

DA Manager uses SSH as a privileged trust boundary. SSH must fail closed when host identity cannot be verified.

Required configuration:

```text
SSH_HOST=<host>:22
SSH_USER=<admin-user>
SSH_KEY=~/.ssh/dune_admin_ed25519
SSH_KNOWN_HOSTS=~/.ssh/known_hosts
```

Rules:

- `SSH_KEY` must be an Ed25519 private key.
- RSA keys are not accepted by the app path.
- `SSH_KNOWN_HOSTS` must contain the remote host's Ed25519 host key.
- Host-key mismatch must be treated as a security event.
- Emergency cloud-console or serial-console keys may exist outside DA Manager, but should not be configured as the DA Manager SSH key.

Add the remote Ed25519 host key:

```bash
ssh-keyscan -t ed25519 -H <host> >> ~/.ssh/known_hosts
```

Validate SSH outside the app:

```bash
ssh -i ~/.ssh/dune_admin_ed25519 <user>@<host>
```

On Windows PowerShell:

```powershell
ssh -i "$env:USERPROFILE\.ssh\dune_admin_ed25519" <user>@<host>
```

## Admin token requirements

`ADMIN_TOKEN` must be generated from 32 random bytes and encoded as 43 base64url characters.

Allowed characters:

```text
A-Z a-z 0-9 _ -
```

Generate a token in PowerShell:

```powershell
[Convert]::ToBase64String((1..32 | ForEach-Object { Get-Random -Maximum 256 })).TrimEnd('=').Replace('+','-').Replace('/','_')
```

Generate a token in Bash:

```bash
python3 - <<'PY'
import base64, os
print(base64.urlsafe_b64encode(os.urandom(32)).decode().rstrip('='))
PY
```

Operational rules:

- Do not reuse tokens across environments.
- Rotate tokens after any suspected disclosure.
- Do not paste tokens into chat, screenshots, issue tickets, logs, or documentation.
- Browser-side token storage is currently session-scoped and remains an interim hardening step.
- The long-term target is memory-only token handling or HttpOnly secure session-cookie auth.

## Network exposure controls

`LISTEN_ADDR` must remain loopback unless explicitly deploying behind a secure boundary.

Default:

```text
LISTEN_ADDR=127.0.0.1:8080
```

Unsafe direct binds such as the following must fail closed:

```text
LISTEN_ADDR=0.0.0.0:8080
```

The only approved override is:

```text
DUNE_ADMIN_REMOTE_EXPOSURE=reverse-proxy-tls
```

Use that only when all of the following are true:

- HTTPS/WSS is terminated by a trusted reverse proxy or VPN.
- Backend access is restricted by firewall or allow-list.
- Admin authentication is enabled and validated.
- Logs do not expose tokens or credentials.
- Operator access path is documented.

## WebSocket log streaming

Static admin tokens must never appear in WebSocket URLs.

Current model:

1. Browser sends authenticated HTTP request to `POST /api/v1/logs/stream-ticket`.
2. Backend issues a short-lived one-time ticket scoped to the selected log target.
3. Browser opens `/api/v1/logs/stream?ns=...&pod=...&ticket=...`.
4. Backend consumes the ticket once and starts the stream.

Legacy `ws_token` query parameters are not accepted.

## Database connectivity

The database should not be exposed directly to the client or public network.

Preferred mode:

```text
SSH_TUNNEL_MODE=auto
SSH_TUNNEL_LOCAL_HOST=127.0.0.1
DB_TUNNEL_LOCAL_PORT=0
```

This lets DA Manager create a local tunnel to the discovered database endpoint.

## Setup and diagnostics

Run interactive setup:

```powershell
.\dune-admin.exe -setup
```

Run connectivity diagnostics:

```powershell
.\dune-admin.exe -diagnose
```

Expected diagnostic stages:

```text
config
ssh_key
ssh_dial
remote_runtime_tools
db_discovery
remote_db_tcp
local_tunnel_bind
```

Always fix the first failing stage first.

## Release validation

Use the platform update script:

Windows:

```powershell
.\update.ps1
```

Linux/macOS/WSL:

```bash
chmod +x ./update.sh
./update.sh
```

Required release evidence:

```bash
go test -v ./...
govulncheck ./...
gosec ./...
gitleaks detect --source .
trivy fs --severity HIGH,CRITICAL .
cd web && npm audit --audit-level=high && npm run typecheck && npm run lint && npm run build
```

If a tool is unavailable, document the exception in the release checklist and do not treat the release as fully validated.

## Backup and recovery

Minimum operational expectations:

- Back up `.env` securely outside the repo.
- Back up SSH private keys in an encrypted password manager or secure vault.
- Maintain a separate emergency access path such as cloud console, serial console, or break-glass account.
- Back up database state before destructive operations.
- Periodically test restore procedures.

Recovery evidence should include:

- Restore date.
- Restore operator.
- Source backup identifier.
- Target environment.
- Result.
- Issues found.
- Follow-up actions.

## Incident response

Security incidents include:

- Admin token disclosure.
- SSH key disclosure or loss.
- Host-key mismatch.
- Unauthorized admin action.
- Unexpected backend public exposure.
- Secrets committed to Git.
- Suspicious mutation activity.

Minimum response process:

1. Triage and classify severity.
2. Contain access path.
3. Rotate affected secrets.
4. Preserve evidence.
5. Remediate root cause.
6. Validate fix.
7. Document incident and corrective action.
8. Update controls to prevent recurrence.

## Administrative change management

Every release or high-risk change should include:

- Change summary.
- Business/operational reason.
- Risk and impact.
- Validation evidence.
- Rollback plan.
- Known issues.
- Operator action items.
- Approval or documented acceptance.

Use `docs/RELEASE_CHECKLIST.md`, `RELEASE_NOTES.md`, and `CHANGELOG.md` for release records.

## Secure operating rules

- Prefer the strictest secure option unless it materially breaks core app usage.
- Keep backend bound to loopback by default.
- Use Ed25519 SSH keys only for DA Manager.
- Keep host-key validation mandatory.
- Keep secrets out of Git and logs.
- Treat player/account identifiers and audit logs as sensitive operational data.
- Require clear reasons for high-risk and destructive mutations.
- Run validation gates before release.
- Keep documentation current with implemented behavior.
