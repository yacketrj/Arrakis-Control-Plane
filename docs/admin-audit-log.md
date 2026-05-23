# Admin Action Audit Log

## Purpose

The Admin Action Audit Log records high-impact administrator actions so server operators can review what changed, when it changed, and whether the backend reported success or failure.

This is a P0 foundation feature. It must exist before Dune Admin expands into more mutation-heavy workflows such as Inventory Studio v2, Player 360 quick actions, guild/faction administration, journey management, and stored procedure execution.

## Current implementation

The current foundation includes:

- an append-only JSONL audit sink
- automatic middleware capture for protected mutating HTTP requests
- protected audit event read endpoint
- status capture for success/failure classification
- redacted, minimal event fields

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

## Captured event fields

Each audit event records:

- timestamp
- HTTP method
- path
- normalized action name
- HTTP status
- duration in milliseconds
- result classification: `success` or `failure`

## Current limitations

The first implementation is deliberately conservative. It does not yet record:

- operator identity beyond the shared admin-token model
- target player/account/controller IDs
- sanitized payload summaries
- rollback hints
- reason text
- typed mutation metadata

Those fields should be added by the Mutation Safety Framework.

## Security rules

- Do not log admin tokens.
- Do not log database credentials.
- Do not log SSH keys.
- Do not log raw request bodies by default.
- Do not expose the audit endpoint in the user portal.
- Keep audit events under protected admin routes only.

## Follow-up tasks

1. Add frontend Audit tab.
2. Add tests for audit middleware success and failure behavior.
3. Add tests proving `/api/v1/public/*` routes are not audited.
4. Add typed audit metadata helpers for high-risk workflows.
5. Add required reason fields for destructive actions.
6. Add rollback hints where feasible.
7. Add operator identity once authentication supports named operators instead of a shared token.

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
