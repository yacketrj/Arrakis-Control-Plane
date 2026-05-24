# Mutation Safety Framework v1

## Purpose

Mutation Safety Framework v1 adds shared safety metadata around high-impact administrator operations. It builds on the Admin Action Audit Log by classifying mutating requests, capturing optional operator reason text, and recording a small allowlist of target identifiers.

This is a P0 foundation layer for future features such as Player 360 quick actions, Inventory Studio v2, guild/faction administration, journey management, safe teleport/rescue, and any future allowlisted routine execution.

## What v1 provides

- Risk classification for protected mutating requests.
- Optional `reason` capture from JSON request bodies.
- Allowlisted target ID capture.
- Request-body restoration after audit inspection so handlers can still decode the request normally.
- Audit UI columns for risk, reason, and target metadata.
- Unit tests for success/failure audit behavior, public route exclusion, body restoration, target capture, and oversized-body handling.

## Risk levels

| Risk | Meaning |
|---|---|
| low | Reserved for future low-impact mutations |
| medium | General protected POST/PUT/PATCH mutation |
| high | High-impact player/admin mutation such as item grants, live grants, teleport, journey changes, and faction changes |
| destructive | DELETE, wipe, delete, or import-style operations that can remove or replace state |

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

## Security rules

- Do not log admin tokens.
- Do not log passwords, DB credentials, SSH keys, or environment values.
- Do not log arbitrary request bodies.
- Keep public routes outside audit capture.
- Keep audit records available only from protected admin routes.
- Treat reason and target metadata as support metadata, not authorization.

## Current limitations

- Reason fields are optional in v1.
- Operator identity is still based on shared admin-token access rather than named operator accounts.
- Rollback hints are not yet implemented.
- Payload summaries are intentionally minimal until typed mutation wrappers are added.

## Follow-up tasks

1. Add shared frontend mutation confirmation component.
2. Require reason text for destructive operations.
3. Add typed backend mutation wrappers per high-risk endpoint.
4. Add rollback hints where feasible.
5. Add named operator identity when authentication supports it.
6. Add export support for audit events.

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
