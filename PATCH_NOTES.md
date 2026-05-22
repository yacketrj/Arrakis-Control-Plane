# Dune Admin Patch Notes

## Why these changes were made

This update was made to move Dune Admin from a trusted-local development tool toward a safer administrative application for managing a live Dune: Awakening server environment.

Before this patch series, the backend exposed powerful administrative capabilities with very little server-side protection. These capabilities included database access, player inventory changes, battlegroup commands, log streaming, notification publishing, blueprint import/export, and RabbitMQ capture/publish workflows. Because those actions can affect live game state, player inventories, database records, Kubernetes workloads, and internal service credentials, the security model needed to be tightened before continuing feature development.

The second major goal was to improve the item grant workflow. The original Give Item flow only supported a single item at a time and treated quantity as a flat item count. That made bulk admin work slow and created incorrect behavior for stacked items. For example, trying to grant 2 stacks of 1000 Heavy Darts was interpreted as 2000 individual inventory entries, which caused inventory-slot errors. This update adds a batch item grant model that preserves explicit stack counts and stack sizes.

Overall, this patch focuses on three priorities:

1. Harden the backend so administrative actions require an explicit admin token.
2. Reduce accidental or malicious exposure of privileged endpoints.
3. Improve the player item grant workflow so admins can grant multiple stacked items accurately.

---

## Security changes

### Backend admin-token authentication

All backend API routes are now protected by an admin-token middleware.

The backend accepts either:

- `Authorization: Bearer <ADMIN_TOKEN>`
- `X-Admin-Token: <ADMIN_TOKEN>`

This ensures that the React frontend is no longer the only access-control boundary. Direct calls to backend endpoints now require the configured server-side token.

Required `.env` setting:

```env
ADMIN_TOKEN=<long random token>
```

The frontend settings panel stores the matching token in browser `localStorage` and sends it to the backend as `X-Admin-Token`.

### CORS allowlisting

Wildcard CORS behavior was replaced with an explicit origin allowlist.

Required `.env` setting:

```env
ALLOWED_ORIGINS=http://localhost:5173,https://dune-admin.layout.tools
```

This prevents arbitrary browser origins from freely calling the backend API from a malicious site.

### Safer backend listen address

The backend now normalizes listen addresses more safely.

Examples:

```text
LISTEN_ADDR=8080       -> 127.0.0.1:8080
LISTEN_ADDR=:8080      -> 127.0.0.1:8080
LISTEN_ADDR=127.0.0.1:8080 -> 127.0.0.1:8080
```

This prevents accidental binding to all interfaces when the backend is intended for local use.

Recommended `.env` setting:

```env
LISTEN_ADDR=127.0.0.1:8080
```

### Server timeout hardening

The HTTP server now uses explicit timeouts:

- Read header timeout
- Read timeout
- Write timeout
- Idle timeout

This reduces exposure to slow-client and resource-exhaustion behavior.

### Request body limiting helper

A request body limiting helper was added for JSON endpoints. The batch item grant endpoint uses this limit to prevent oversized JSON request bodies.

### Raw SQL restrictions

The raw SQL endpoint was restricted to read-only style statements.

Allowed prefixes include:

- `SELECT`
- `WITH`
- `SHOW`
- `EXPLAIN`

The backend rejects semicolon-separated statements and destructive SQL keywords such as:

- `INSERT`
- `UPDATE`
- `DELETE`
- `DROP`
- `ALTER`
- `TRUNCATE`
- `CREATE`
- `GRANT`
- `REVOKE`
- `COPY`
- `CALL`
- `DO`
- `EXECUTE`
- `MERGE`
- `VACUUM`
- `ANALYZE`

This reduces the chance that an admin accidentally runs destructive SQL through the database tab.

### WebSocket and Kubernetes log stream hardening

Log streaming now validates Kubernetes namespace and pod names before using them as log targets.

The WebSocket upgrader also uses the same origin allowlist logic as the normal API CORS middleware.

This reduces the risk of browser-origin abuse and unsafe log target input.

### Status endpoint data reduction

The status response no longer exposes pod IP data.

The endpoint still reports whether SSH and database connections are active, but it avoids returning unnecessary internal network details.

### Capture and notification credential cleanup

Hardcoded RabbitMQ capture credentials and JWT signing material were removed from the source code.

Capture/notification-related settings now come from environment variables:

```env
DUNE_SERVICE_JWT_SIGNING_SECRET=
DUNE_CAPTURE_USER=dune_cap
DUNE_CAPTURE_PASS=<strong random password>
```

Notes:

- `DUNE_CAPTURE_PASS` should be generated fresh and treated as secret.
- `DUNE_SERVICE_JWT_SIGNING_SECRET` should only be set if you know the actual game service signing secret.
- If the signing secret is blank, capture mode falls back to the configured capture user/password path.

### Frontend admin-token support

The frontend settings gear now includes an Admin Token field.

The frontend stores the token locally and sends it to the backend with API requests, including blueprint imports.

This keeps local development convenient while still requiring the backend to enforce authentication.

---

## Player item grant changes

### Batch item grant endpoint

The existing endpoint remains:

```http
POST /api/v1/players/give-item
```

It now supports both the old single-item payload and the new batch payload.

Legacy payload still works:

```json
{
  "player_id": 123,
  "template": "ItemTemplateHeavyDarts",
  "qty": 10,
  "quality": 1
}
```

New batch payload:

```json
{
  "player_id": 123,
  "items": [
    {
      "template": "ItemTemplateHeavyDarts",
      "qty": 2,
      "quality": 1,
      "stack_size": 1000
    },
    {
      "template": "ItemTemplateOtherResource",
      "qty": 5,
      "quality": 3,
      "stack_size": 250
    }
  ]
}
```

In the batch payload:

- `qty` means number of stacks.
- `stack_size` means how many items are in each stack.
- Total granted items equals `qty × stack_size`.
- Inventory slots required equals `qty`, not `qty × stack_size`.

### Explicit stack creation

A new backend command was added to preserve stack intent:

```text
cmdGiveItemStacks(playerID, template, stacks, stackSize, quality)
```

This fixes the issue where granting 2 stacks of 1000 Heavy Darts was incorrectly interpreted as 2000 separate inventory slots.

Expected behavior now:

```text
2 stacks × 1000 Heavy Darts = 2 inventory slots, 2000 total darts
```

### Batch validation

The batch item grant handler validates each row before making changes.

Current limits:

```text
max item rows per request: 100
max stack count per row: 9999
max stack size per row: 9999
quality range: 0 through 5
```

The backend rejects:

- Empty item lists.
- Blank item templates.
- Zero or negative stack counts.
- Stack counts above the configured maximum.
- Stack sizes above the configured maximum.
- Quality below 0 or above 5.
- Requests that could overflow total quantity math.

### Inventory capacity behavior

The explicit stack grant path checks inventory capacity using stack count.

For example:

```text
qty = 2
stack_size = 1000
```

requires 2 free slots, not 2000 free slots.

The command still respects inventory volume/weight checks by calculating total item count as:

```text
qty × stack_size
```

This means the inventory slot count and item volume behavior are both handled correctly.

### Frontend multi-item UI

The Give Item modal now supports multiple rows.

Each row includes:

- Item template selection/search.
- Quantity / stack count.
- Grade / quality.
- Stack size.
- Total item count display.

The modal submits all valid rows in one request through:

```text
api.players.giveItems(playerID, items)
```

---

## Testing changes

### Backend unit tests

New unit tests were added for batch item request normalization.

Coverage includes:

- Legacy single-item payload compatibility.
- New batch payload parsing.
- Template trimming.
- Stack size defaulting.
- Empty payload rejection.
- Blank template rejection.
- Zero and negative quantity rejection.
- Quantity upper-bound rejection.
- Stack-size upper-bound rejection.
- Quality range rejection.
- More-than-100-row rejection.
- Maximum allowed quantity and stack-size acceptance.

### Go test workflow

A GitHub Actions workflow was added to run:

```bash
go test ./...
```

on pushes to `main`.

### Local validation commands

Recommended backend validation:

```powershell
git pull origin main
go test ./...
go run .
```

Recommended frontend validation:

```powershell
cd web
npm install
npm run build
npm run dev
```

---

## Operational notes

### Required local backend `.env`

A minimal local development `.env` should include:

```env
ADMIN_TOKEN=<long random token>
LISTEN_ADDR=127.0.0.1:8080
ALLOWED_ORIGINS=http://localhost:5173,https://dune-admin.layout.tools
DUNE_CAPTURE_USER=dune_cap
DUNE_CAPTURE_PASS=<strong random password>
DUNE_SERVICE_JWT_SIGNING_SECRET=
```

### Running locally

Use two terminals.

Backend:

```powershell
cd "Z:\Unreal Projects\Icarus\dune-admin-fork"
go run .
```

Frontend:

```powershell
cd "Z:\Unreal Projects\Icarus\dune-admin-fork\web"
npm run dev
```

Open:

```text
http://localhost:5173
```

Use the frontend gear icon to configure:

```text
Backend URL: http://localhost:8080
Admin Token: same value as ADMIN_TOKEN in .env
```

---

## Final summary

This patch series significantly improves the safety and usability of Dune Admin.

The backend now enforces authentication, restricts browser origins, normalizes local listen addresses, applies server timeouts, limits request bodies, reduces status-data exposure, restricts raw SQL, validates log-stream targets, and removes hardcoded capture secrets.

The player item grant workflow now supports true multi-item batch grants. Admins can grant several item types at once, set the number of stacks, set per-stack size, and set quality/grade per row. The backend now preserves explicit stack sizing instead of flattening stack grants into thousands of individual inventory entries.

The most important functional fix is that granting something like:

```text
2 stacks of 1000 Heavy Darts
```

now creates:

```text
2 inventory slots, each with stack_size = 1000
```

rather than requiring 2000 inventory slots.

The result is a safer admin tool with a faster and more accurate item-grant workflow for live server administration.
