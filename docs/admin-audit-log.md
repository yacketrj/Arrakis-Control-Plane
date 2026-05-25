# Admin Action Audit Log

## Purpose

The Admin Action Audit Log records high-impact administrator actions so server operators can review what changed, when it changed, which safety class applied, and whether the backend reported success or failure.

This is a P0 foundation feature for Dune Admin. It supports safer expansion into mutation-heavy workflows such as Inventory Studio v2, Player 360 quick actions, guild/faction administration, journey management, teleport/rescue support, and stored procedure execution.

## Current implementation

The current foundation includes:

- an append-only JSONL audit sink
- automatic middleware capture for protected mutating HTTP requests
- protected audit event read endpoint
- status capture for success/failure classification
- risk classification from the Mutation Safety Framework
- optional or enforced operator reason capture through `X-Admin-Reason` or JSON body `reason`
- allowlisted target identifier capture
- mutation metadata fields for preview requirement, destructive status, rollback hint, operator warnings, and recommended path
- unit coverage for audit behavior, public route exclusion, body restoration, metadata capture, and reason enforcement paths

## Protected endpoint

```text
GET /api/v1/audit/events
```

This endpoint is protected by the existing admin authorization middleware.

## Audit file

Default path:

```text
admin-audit.jsonl
```

Override path:

```text
ADMIN_AUDIT_LOG=/path/to/admin-audit.jsonl
```

The file is opened append-only and created with owner-only permissions when the backend writes the first event.

## Captured event fields

Each audit event records:

- timestamp
- HTTP method
- request path
- normalized action name
- mutation risk: `low`, `medium`, `high`, or `destructive`
- optional reason text from `X-Admin-Reason` or JSON body `reason`
- allowlisted target identifiers
- HTTP status
- duration in milliseconds
- result classification: `success` or `failure`
- whether a reason is required for the action
- whether a preview should be shown before the action
- whether the action is destructive
- rollback hint, where known
- operator warnings, where known
- recommended operator path, where known

## Captured target identifiers

The audit middleware only captures this allowlist:

```text
player_id
account_id
actor_id
controller_id
item_id
faction_id
storage_id
```

The middleware intentionally does not log raw request bodies.

## Reason capture and enforcement

Reason text can be supplied through:

```text
X-Admin-Reason: support correction note
```

or as JSON:

```json
{
  "reason": "support correction note"
}
```

Reason enforcement is controlled by:

```text
ADMIN_REQUIRE_REASON=true
```

When enabled, high-risk and destructive mutations must include a reason. When disabled, reasons are still captured when provided.

## Security rules

- Do not log admin tokens.
- Do not log database credentials.
- Do not log SSH keys.
- Do not log raw request bodies by default.
- Do not expose the audit endpoint in the user portal.
- Keep audit events under protected admin routes only.
- Treat reason text as support context, not authorization.
- Review exported or copied audit content before external sharing.

## Current limitations

- Operator identity is still based on shared admin-token access rather than named operator accounts.
- Audit records are local JSONL files; rotation, retention, and export workflow are not yet automated.
- Rollback hints describe operator guidance, but the backend does not yet create automatic before-change snapshots.
- Preview requirement is metadata only until a shared frontend confirmation component is added.

## Follow-up tasks

1. Add named operator identity once authentication supports individual users.
2. Add audit export with built-in redaction review steps.
3. Add local retention and rotation guidance.
4. Add typed before-change snapshots for inventory, progression, teleport, and storage mutations.
5. Wire a shared frontend confirmation dialog to the mutation safety metadata.
6. Add operator-visible filtering by action, risk, result, target, and timestamp.

## Validation

```bash
go test ./...
cd web
npm install
npm audit --audit-level=high
npm run typecheck
npm run lint
npm run build
```
