# SSH Tunnel Management

## Purpose

Dune Admin supports a managed SSH tunnel path for game-management data traffic. The default mode creates local loopback tunnels through the configured SSH host and routes supported game-management traffic through those tunnels instead of requiring direct access to discovered cluster addresses.

Remote operational commands still execute over the authenticated SSH session. Database operations use the managed tunnel, RabbitMQ capture can use the same tunnel policy, and managed tunnels are closed during shutdown or reconnect.

## Tunnel modes

| Mode | Behavior |
|---|---|
| auto | Creates managed local tunnels for supported game-management traffic. This is the default. |
| existing | Uses an operator-controlled path. For database access, configure the local host and local port to match the external tunnel. |
| off | Direct connection mode for troubleshooting only. |

## Configuration

```text
SSH_TUNNEL_MODE=auto
SSH_TUNNEL_LOCAL_HOST=127.0.0.1
DB_TUNNEL_LOCAL_PORT=0
```

`DB_TUNNEL_LOCAL_PORT=0` lets the operating system choose an available local port for the managed PostgreSQL tunnel. Set a specific local port only when you need a predictable port or when using existing tunnel mode.

## Runtime behavior

On connect or reconnect:

1. Dune Admin opens the SSH session.
2. Dune Admin discovers the database pod and pod IP through Kubernetes.
3. In managed tunnel mode, Dune Admin binds a local loopback listener.
4. PostgreSQL traffic connects to the local tunnel.
5. The tunnel forwards traffic through the SSH session to the discovered internal database address.
6. On reconnect or shutdown, Dune Admin closes database connections, closes managed tunnels, and then closes SSH.

## RabbitMQ capture behavior

RabbitMQ capture now follows the normalized tunnel policy:

| Mode | Capture behavior |
|---|---|
| auto | Opens a short-lived managed tunnel connection for AMQP traffic. Closing the AMQP connection also closes that ephemeral tunnel. |
| existing | Uses the authenticated SSH connection path for the internal AMQP address. |
| off | Uses direct TCP dialing. This is intended only for troubleshooting from trusted network paths. |

## Status API

The protected status endpoint includes tunnel state for operator visibility:

```text
GET /api/v1/status
```

Relevant fields:

| Field | Description |
|---|---|
| tunnel_mode | Normalized tunnel mode: `auto`, `existing`, or `off`. |
| tunnels | Active managed tunnels. Each entry includes name, local address, and remote address. |

This status is protected and should remain part of the admin portal only.

## Operational guidance

Use `auto` for normal operations. Use `existing` when a controlled external tunnel is required by a local operating model. Use `off` only when debugging connectivity from a trusted network path.

When using `existing`, make sure the external local database tunnel is already established before starting Dune Admin or before running reconnect.

## Validation

```bash
go mod tidy
git diff -- go.mod go.sum
go test ./...
cd web
npm run lint
npm run build
```
