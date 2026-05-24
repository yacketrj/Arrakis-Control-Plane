# Battlegroup Health Diagnostics

## Purpose

Battlegroup Health Diagnostics gives operators a read-only view of the Kubernetes and service state that most often explains server availability problems.

This feature is part of the Server Health Command Center roadmap. It is intentionally diagnostic-only and does not mutate Kubernetes resources, database state, RabbitMQ state, or player data.

## Endpoint

```text
GET /api/v1/battlegroup/health
```

The endpoint is protected by the administrator authorization middleware.

## UI location

```text
Battlegroup > Health Diagnostics
```

The Battlegroup tab now has two views:

- `Pods`
- `Health Diagnostics`

The `Health Diagnostics` view runs the read-only diagnostic bundle and renders each section as a card with:

- section name
- section description
- command executed
- command output
- command error, when present

## Diagnostic sections

The backend currently collects:

| Section | Purpose |
|---|---|
| pods | Shows namespace pod readiness, restarts, age, node placement, and pod IPs. |
| services | Shows service exposure and service discovery state. |
| statefulsets | Shows stateful workload readiness for services such as PostgreSQL and RabbitMQ. |
| deployments | Shows deployment rollout readiness for BGD, gateway, text router, and related workloads. |
| persistent_volumes | Shows PVC state for database and stateful service storage. |
| recent_events | Shows namespace events sorted by timestamp. |
| nodes | Shows cluster node readiness and placement information. |
| pod_metrics | Shows pod CPU and memory usage when metrics-server is available. |

## Security model

- The endpoint is protected and must not be exposed in the public user portal.
- Commands are fixed server-side and are not operator-supplied.
- The backend validates the configured Kubernetes namespace before building namespace-scoped commands.
- The command bundle is read-only.
- Metrics-server failures are informational; they should not block the rest of the diagnostic result.

## Operator workflow

1. Open the protected admin portal.
2. Open `Battlegroup`.
3. Select `Health Diagnostics`.
4. Click `Run Diagnostics`.
5. Review the cards in this order:
   - pods
   - recent_events
   - services
   - statefulsets
   - deployments
   - persistent_volumes
   - nodes
   - pod_metrics

## Troubleshooting guidance

### Pods not ready

Check:

- pod status
- restart count
- recent events
- node placement
- image pull or scheduling errors

### Services not reachable

Check:

- service type
- exposed ports
- service selectors
- pod readiness

### Database or RabbitMQ instability

Check:

- statefulsets
- PVCs
- pod restarts
- recent events
- pod metrics, when available

### Metrics unavailable

`pod_metrics` may fail when metrics-server is not installed or is unavailable. This is not treated as a fatal health endpoint failure.

## Validation

```bash
go mod tidy
git diff -- go.mod go.sum
go test ./...
cd web
npm run lint
npm run build
```

## Follow-up tasks

- Add structured parsing for the diagnostic sections.
- Add health badges for unhealthy sections.
- Add export/download support for a support bundle.
- Add a redact-before-copy helper for external bug reports.
- Add a server health landing page that combines DB, SSH, Kubernetes, RabbitMQ, and game process state.
