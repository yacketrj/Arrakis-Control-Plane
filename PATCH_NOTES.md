# Dune Admin Release Notes

## Current update: Player 360 roadmap and design foundation

### Why this update was made

The DA Manager roadmap reached a clean transition point between the P0 safety foundation and the P1 operator-support surface. Admin Audit is implemented, Mutation Safety has a backend foundation, and Player 360 Profile is now the next feature slice. The roadmap still pointed at the completed audit foundation as the next implementation step, so the planning documents needed to be corrected before code work begins.

### What changed

- Updated `docs/admin-feature-design-and-priorities.md` with the current roadmap position between Phase 1 and Phase 2.
- Marked Admin Action Audit Log as implemented foundation work in the feature roadmap.
- Clarified Mutation Safety Framework as an in-progress backend foundation with remaining frontend confirmation and typed wrapper work.
- Corrected the Player 360 roadmap entry to fold Currency and Online Status into Player Info.
- Corrected the Battlegroup Status v2 roadmap row so Prometheus/Grafana graph work remains a future diagnostic improvement.
- Added `docs/player-360-profile.md` as the design document for the next P1 read-only implementation slice.

### Security and operator impact

- Keeps the project aligned with the read-only-before-write rule.
- Prevents Player 360 from becoming another mutation surface before visibility, preview, reason flow, and audit expectations are clear.
- Defines Player 360 v1 as a protected read-only support view that aggregates existing player context without adding new high-risk actions.
- Keeps future quick actions dependent on Mutation Safety and audit controls.

### Validation

Documentation-only update. Expected validation before implementation work continues:

```bash
go test ./...
cd web
npm install
npm audit --audit-level=high
npm run typecheck
npm run lint
npm run build
```

---

## Previous update: Admin audit and mutation-safety documentation sync

### Why this update was made

The DA Manager workstream has moved beyond the original SSH tunnel foundation. The repository now includes the admin audit log and mutation-safety foundation, but several documentation and release-tracking files still described those features as future work. This update brings the operator-facing notes back in line with the code on `main`.

### What changed

- Updated `docs/admin-audit-log.md` to describe the current audit event model, protected audit endpoint, audit file behavior, captured fields, reason capture, and known limitations.
- Updated `docs/admin-implementation-tasks.md` so the Admin Action Audit Log is marked done and Mutation Safety Framework is tracked as active in-progress work.
- Clarified that Player 360 Profile is the next feature slice after the audit and mutation-safety foundation.
- Preserved the DA Manager requirement that every feature update keeps `PATCH_NOTES.md`, `CHANGELOG.md`, and relevant docs synchronized.

### Security and operator impact

- Operators now have clearer documentation for protected audit review.
- The task tracker now reflects that mutation-heavy future features should build on the landed audit foundation.
- The next Player 360 work should begin read-only and reuse audit and mutation-safety controls before adding quick actions.

### Validation

Expected validation:

```bash
go test ./...
cd web
npm install
npm audit --audit-level=high
npm run typecheck
npm run lint
npm run build
```

---

## Previous update: SSH tunnel management foundation

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
