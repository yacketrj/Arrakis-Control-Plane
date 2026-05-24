# SSH Tunnel Management

## Purpose

Dune Admin now supports a managed SSH tunnel path for game-management data traffic. The default mode creates a local loopback tunnel through the configured SSH host and routes PostgreSQL access through that tunnel instead of relying on direct access to discovered cluster addresses.

Remote operational commands still execute over the authenticated SSH session. Database operations use the managed tunnel, and the tunnel is closed during shutdown or reconnect.

## Tunnel modes

| Mode | Behavior |
|---|---|
| auto | Creates a managed local tunnel for supported game-management traffic. This is the default. |
| existing | Uses an operator-created local tunnel. Configure the local host and local port to match the external tunnel. |
| off | Direct connection mode for troubleshooting only. |

## Configuration

```text
SSH_TUNNEL_MODE=auto
SSH_TUNNEL_LOCAL_HOST=127.0.0.1
DB_TUNNEL_LOCAL_PORT=0
```

`DB_TUNNEL_LOCAL_PORT=0` lets the operating system choose an available local port. Set a specific local port only when you need a predictable port or when using existing tunnel mode.

## Runtime behavior

On connect or reconnect:

1. Dune Admin opens the SSH session.
2. Dune Admin discovers the database pod and pod IP through Kubernetes.
3. In managed tunnel mode, Dune Admin binds a local loopback listener.
4. PostgreSQL traffic connects to the local tunnel.
5. The tunnel forwards traffic through the SSH session to the discovered internal database address.
6. On reconnect or shutdown, Dune Admin closes database connections, closes managed tunnels, and then closes SSH.

## Operational guidance

Use `auto` for normal operations. Use `existing` when a controlled external tunnel is required by a local operating model. Use `off` only when debugging connectivity from a trusted network path.

## Validation

```bash
go mod tidy
git diff -- go.mod go.sum
go test ./...
cd web
npm run lint
npm run build
```
