# Mutation Safety Framework v1

## Purpose

Mutation Safety Framework v1 adds shared safety metadata around high-impact administrator operations. It builds on the Admin Action Audit Log by classifying mutating requests, capturing operator reason text when supplied, and recording a small allowlist of target identifiers.

This is a P0 foundation layer for future features such as Player 360 quick actions, Inventory Studio v2, guild/faction administration, journey management, safe teleport/rescue, and allowlisted routine execution.

## What v1 provides

- Risk classification for protected mutating requests.
- Protected classification endpoint for frontend preview workflows.
- Optional `reason` capture from JSON request bodies.
- Optional `X-Admin-Reason` header capture for admin workflows.
- Environment-controlled reason enforcement for high-risk and destructive actions.
- Allowlisted target ID capture.
- Request-body restoration after audit inspection so handlers can still decode the request normally.
- Audit event fields for risk, reason, target metadata, preview requirement, destructive status, rollback hint, operator warnings, and recommended path.
- Unit tests for classification, handler behavior, reason enforcement, body restoration, target capture, and oversized-body handling.

## Protected classification endpoint

```text
GET /api/v1/mutation-safety/classify?method=POST&path=/api/v1/players/give-item
```

The response includes the normalized action name, risk level, reason requirement, preview requirement, destructive flag, and any available guidance fields.

## Risk levels

| Risk | Meaning |
|---|---|
| low | Reserved for future low-impact mutations |
| medium | General protected POST/PUT/PATCH mutation |
| high | High-impact player/admin mutation such as item grants, live grants, teleport, journey changes, and faction changes |
| destructive | Operations that can remove or replace important server or player state |

## Captured metadata

The audit middleware only captures these allowlisted fields:

- `reason`
- `player_id`
- `account_id`
- `actor_id`
- `controller_id`
- `item_id`
- `faction_id`
- `storage_id`

The middleware intentionally does not log full request bodies.

## Reason capture and enforcement

Reason text can be supplied in the `X-Admin-Reason` header or in a JSON request body field named `reason`.

Reason enforcement is controlled by `ADMIN_REQUIRE_REASON`. When enabled, high-risk and destructive requests must include a reason. When disabled, reason text is still captured when supplied.

## Security rules

- Do not log admin tokens.
- Do not log database credentials, SSH keys, or environment values.
- Do not log arbitrary request bodies.
- Keep public routes outside audit capture.
- Keep audit records available only from protected admin routes.
- Treat reason and target metadata as support metadata, not authorization.
- Keep mutation preview metadata separate from actual authorization checks.

## Current limitations

- Operator identity is still based on shared admin-token access rather than named operator accounts.
- Preview requirement is metadata only until a shared frontend confirmation component is added.
- Rollback hints describe operator guidance, but the backend does not yet create automatic before-change snapshots.
- Reason enforcement is environment-controlled and not yet configurable from the UI.
- Typed mutation wrappers are still needed for workflow-specific before/after metadata.

## Follow-up tasks

1. Add shared frontend mutation confirmation component.
2. Add typed backend mutation wrappers per high-risk endpoint.
3. Add before-change snapshot helpers for inventory, journey/progression, teleport, and storage operations.
4. Add named operator identity when authentication supports individual users.
5. Add audit export and filtering support.
6. Add UI visibility for reason-enforcement state.

## Validation

```bash
gofmt -w *.go
go test ./...
cd web
npm install
npm audit --audit-level=high
npm run typecheck
npm run lint
npm run build
```
