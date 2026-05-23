# User Portal and Protected Admin Portal Design

## Purpose

Dune Admin needs two clearly separated surfaces:

1. an unprotected, player-safe user portal
2. a protected administrator portal

This separation keeps public server information available to players while preventing accidental exposure of sensitive operational, database, account, inventory, moderation, and infrastructure data.

## Portal model

### User Portal

The user portal is for information that is safe to expose without administrator credentials.

Allowed user-facing content:

- redacted server status
- maintenance announcements
- server rules
- support links
- event schedule
- claim/reward instructions
- known issue notices

Disallowed user-facing content:

- account IDs
- FLS IDs
- actor IDs
- controller IDs
- inventory details
- currency balances
- faction reputation details
- database schema/table/routine inspection
- logs
- cheat records
- battlegroup execution tools
- storage and blueprint tools
- mutation actions
- SSH host, namespace, pod, queue, or infrastructure details

### Protected Admin Portal

The protected admin portal is the operations console.

Admin-only content:

- player search and Player 360 workflows
- inventory and item mutation
- Claim Rewards Queue tooling
- database table inspection
- database routine discovery and inspection
- battlegroup operations
- logs and cheat records
- notifications and broadcasts
- storage and blueprint workflows
- teleport/rescue workflows
- journey, progression, guild, faction, and currency tools
- audit logs

## Backend route pattern

Use this route split:

```text
/api/v1/public/*    reviewed player-safe endpoints
/api/v1/*           protected administrator endpoints
```

Only deliberately reviewed endpoints should be added under `/api/v1/public/*`.

## Current public endpoint

```text
GET /api/v1/public/status
```

This endpoint must remain redacted. It should not include SSH host, namespace, pod names, database details, player counts, player names, or operational internals unless each field has been reviewed as public-safe.

## Security rules

- Public endpoints must be read-only.
- Public endpoints must not proxy arbitrary backend data.
- Public endpoints must not expose raw errors containing infrastructure details.
- Public endpoints should return stable, minimal JSON.
- Admin endpoints must continue to require administrator authorization.
- Any new public endpoint requires documentation and review before implementation.

## Frontend behavior

The frontend should make the separation visible:

```text
User Portal | Protected Admin Portal
```

The user portal can load without an admin token. The admin portal should prompt for backend settings and an admin token when needed.

## Implementation sequence

1. Add redacted public status endpoint.
2. Add User Portal landing page.
3. Keep existing admin tabs in the protected Admin Portal.
4. Move database routine discovery and inspection into Admin Portal only.
5. Add tests ensuring public status does not expose sensitive fields.
6. Add audit logging before implementing additional mutation-heavy admin workflows.

## Validation

Expected validation commands:

```bash
go test ./...
cd web
npm install
npm audit --audit-level=high
npm run typecheck
npm run lint
npm run build
```
