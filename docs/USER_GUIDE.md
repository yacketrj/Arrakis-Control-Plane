# DA Manager User Guide

## Purpose

DA Manager is a local-first administrative interface for inspecting and managing a self-hosted Dune: Awakening environment. It provides operational workflows for runtime health, player investigation, inventory actions, storage, blueprints, logs, database inspection, and audit review.

This guide is intended for operators using the web UI.

## Security overview

DA Manager is a high-privilege administrative control plane. Treat access to the UI, backend, admin token, SSH key, database credentials, logs, and audit records as privileged administrative access.

Current security posture:

- Backend defaults to loopback access: `127.0.0.1:8080`.
- Non-loopback backend binding fails closed unless explicitly configured for secure reverse-proxy/VPN exposure.
- Browser access requires a strict 43-character admin token.
- The frontend stores the access token in session-scoped browser storage as an interim control.
- Log streaming uses short-lived one-time WebSocket tickets.
- SSH access requires Ed25519 client keys.
- SSH host trust requires the remote Ed25519 host key in `known_hosts`.

## Start the backend

From the repository root:

```powershell
.\dune-admin.exe
```

Expected startup pattern:

```text
dune-admin listening on 127.0.0.1:8080
runtime: docker
ssh: <user>@<host>:22
ssh tunnel mode: auto
ssh tunnel postgres: 127.0.0.1:<local-port> -> <remote-db-endpoint>:5432
```

If startup is degraded, DB-backed tabs remain unavailable until connectivity is restored.

## Run diagnostics

```powershell
.\dune-admin.exe -diagnose
```

Healthy stages:

```text
[OK] config
[OK] ssh_key
[OK] ssh_dial
[OK] remote_runtime_tools
[OK] db_discovery
[OK] remote_db_tcp
[OK] local_tunnel_bind
```

Fix the first failing stage before troubleshooting later stages.

## Start the frontend

Development mode:

```powershell
cd web
npm run dev
```

Open:

```text
http://localhost:5173
```

Production-like local preview:

```powershell
npm run build
npm run preview -- --host 127.0.0.1 --port 4173
```

Open:

```text
http://localhost:4173
```

## Configure browser access

1. Open the frontend.
2. Open settings.
3. Set Backend URL to `http://localhost:8080` unless changed intentionally.
4. Paste the backend `ADMIN_TOKEN` into **Browser Access Key**.
5. Save and reload.

The Browser Access Key must be exactly 43 base64url characters:

```text
A-Z a-z 0-9 _ -
```

The UI rejects weak, malformed, whitespace-containing, and placeholder values.

## Token handling

The browser access token is session-scoped. Closing the tab or browser may require re-entry. This is intentional and reduces long-term browser exposure compared with persistent storage.

Do not paste the admin token into:

- Discord.
- Tickets.
- Screenshots.
- Logs.
- Chat windows.
- Shared terminal output.
- Documentation.

## Reconnect workflow

When DB-backed tools are unavailable:

1. Review the blocked panel.
2. Confirm SSH, tunnel, and DB state.
3. Select **Retry reconnect**.
4. Reload when reconnect succeeds.

Reconnect closes stale DB, tunnel, and SSH state, then attempts a fresh connection.

## Primary tabs

| Tab | Purpose | Operator caution |
|---|---|---|
| Audit | Review administrative actions and mutation history. | Use for accountability and change review. |
| Battlegroup | Runtime status, health, and controlled command execution. | Commands can affect live services. |
| Players | Player search/listing and online state. | Use to identify correct IDs before mutation. |
| Player 360 | Consolidated per-player view. | Confirm identity before actions. |
| Inventory Studio | Inventory investigation and mutation workflows. | Mutations require clear reason capture when enforced. |
| Database | Read-only inspection and controlled SQL tooling. | Only read-only SQL should be permitted by default. |
| Logs | Runtime log target discovery and streaming. | Log streaming uses one-time WebSocket tickets. |
| Storage | Storage/container inspection and item workflows. | Confirm target container before mutation. |
| Blueprints | Blueprint inspection/export/import. | Imports can materially alter player/server state. |

## Log streaming

Log streaming uses a one-time ticket flow:

1. Select a log target.
2. The UI requests a one-time ticket from the backend.
3. The UI opens the WebSocket stream using that ticket.
4. The backend consumes the ticket and streams logs.

Legacy `ws_token` behavior is not accepted. If log streaming fails, refresh log targets and reconnect.

## Safe operating practices

- Verify the target player, account, container, or runtime before mutation.
- Record a clear reason for high-risk or destructive operations.
- Avoid running the UI on untrusted machines.
- Do not expose the backend directly to the public internet.
- Use VPN or a trusted reverse proxy if remote access is required.
- Review Audit after sensitive operations.
- Export and share redacted support bundles only after operator review.

## Common issues

| Symptom | Likely cause | Action |
|---|---|---|
| Unauthorized | Missing or invalid browser access key. | Re-enter the 43-character token. |
| DB unavailable | SSH tunnel or database discovery failed. | Run `-diagnose` and retry reconnect. |
| Logs fail to stream | Ticket issuance or WebSocket connection failed. | Refresh log targets and reconnect. |
| Backend refuses startup | Unsafe non-loopback bind. | Use `127.0.0.1:8080` or configure approved secure exposure mode. |
| SSH dial fails | Missing key, non-Ed25519 key, or unknown host key. | Confirm Ed25519 key and `known_hosts`. |
| Browser token disappears | Session-scoped token storage cleared. | Re-enter Browser Access Key. |

## Escalation notes

Escalate to the administrator or release owner when:

- SSH host-key mismatch occurs.
- Admin token may have been disclosed.
- DB password may have been disclosed.
- Backend was exposed beyond loopback unexpectedly.
- Unauthorized or unexplained admin actions appear in Audit.
- Logs show repeated authentication failures.
