# Dune Admin Release Notes

## Current update: SSH tunnel management foundation

### Why this update was made

Dune Admin game-management workflows need to consistently use SSH-controlled paths instead of assuming direct access to in-cluster game services. This improves operator safety for local, cloud, and remote deployments where Kubernetes service and database endpoints should not be exposed directly.

### What changed

- Added SSH tunnel configuration flags and environment-driven settings.
- Added managed tunnel lifecycle support for local loopback forwarding.
- Updated database connection logic so PostgreSQL access can route through a managed SSH tunnel by default.
- Added support for `auto`, `existing`, and `off` tunnel modes.
- Added cleanup of managed tunnels during shutdown and reconnect paths.
- Added `docs/ssh-tunnel-management.md` with tunnel modes, configuration, runtime behavior, and validation guidance.

### Security and operator impact

- Keeps database access bound to local loopback in managed tunnel mode.
- Avoids requiring direct exposure of internal Kubernetes database addresses to the operator workstation.
- Allows externally managed tunnels when the operator environment already controls tunnel creation.
- Keeps direct connection mode available only as a troubleshooting escape hatch.

### Validation

Expected validation:

```bash
go mod tidy
git diff -- go.mod go.sum
go test ./...
cd web
npm run lint
npm run build
```

---

## Previous update: Battlegroup Health Diagnostics

### Why this update was made

Operators need a fast, read-only way to collect the Kubernetes and service-state signals that usually explain Dune server availability issues. The previous Battlegroup view showed pod status and allowed server control actions, but it did not provide a broader health bundle for services, stateful workloads, persistent volumes, node state, recent events, or pod metrics.

This update moves the Server Health Command Center roadmap forward by adding protected Battlegroup Health Diagnostics while preserving the security rule that diagnostics must not mutate infrastructure or expose data through the public portal.

### What changed

- Added protected backend endpoint `GET /api/v1/battlegroup/health`.
- Added read-only diagnostic collection for pods, services, statefulsets, deployments, PVCs, recent events, nodes, and pod metrics.
- Added namespace validation before building namespace-scoped Kubernetes commands.
- Added frontend API support for the health diagnostic response.
- Added a `Health Diagnostics` view in the Battlegroup tab.
- Added diagnostic cards showing section name, description, command, output, and command error.
- Added raw support-bundle export for exact local diagnostic capture.
- Added redacted support-bundle export that masks common infrastructure identifiers before download.
- Added `docs/battlegroup-health-diagnostics.md` with endpoint, UI, security, validation, troubleshooting, export, and redaction guidance.
- Fixed the Battlegroup view buttons to use supported HeroUI variants.

### Security and operator impact

- Diagnostics are protected by admin authorization.
- Commands are fixed server-side; operators cannot submit arbitrary commands through the health endpoint.
- The diagnostic bundle is read-only.
- Metrics-server errors are surfaced per section and do not block the rest of the health result.
- Raw bundles remain available for local operations review.
- Redacted bundles provide a safer handoff starting point by masking IPv4 addresses, IPv6 addresses, UUIDs, and common cloud/internal hostname patterns.
- Redacted export is a helper, not a full data-loss-prevention system; operators should still review exported bundles before external sharing.
- The feature improves supportability without adding a new mutation path.

### Validation

Expected validation:

```bash
go mod tidy
git diff -- go.mod go.sum
go test ./...
cd web
npm run lint
npm run build
```

---
