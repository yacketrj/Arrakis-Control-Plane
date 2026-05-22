# Dune Admin Release Notes

## Release: Security Hardening and Multi-Item Administration Update

### Release type

Security hardening, reliability fixes, and player administration feature update.

### Audience

Server administrators, operators, maintainers, and anyone running Dune Admin against a live Dune: Awakening environment.

---

## 1. Why this release was made

This release was created to address two high-priority needs:

1. **Reduce security risk around privileged administration endpoints.**
2. **Improve the accuracy and speed of player item administration.**

Dune Admin has direct access to sensitive and high-impact systems: player inventory records, game database tables, Kubernetes pod logs, battlegroup command execution, blueprint import/export, RabbitMQ notification/capture flows, and other live server administration functions. Prior to this update, many of those capabilities were designed around a trusted local-development workflow. That created unacceptable risk if the backend was accidentally exposed, accessed from an unexpected browser origin, called directly without the frontend, or operated with hardcoded secrets still present in source.

The security work in this release moves enforcement into the backend, which is the correct control point for privileged operations. The frontend remains a convenience layer, but the backend now rejects unauthenticated API calls, limits unsafe request patterns, reduces unnecessary information disclosure, and removes embedded credential material.

The item administration work was also necessary because the original single-item grant flow did not accurately model stacked inventory. In particular, an operator attempting to grant **2 stacks of 1000 Heavy Darts** saw the grant treated as **2000 individual inventory entries**, which incorrectly required 2000 inventory slots. This release adds a true batch item grant workflow and backend stack-preserving logic so stack count and stack size are handled separately.

---

## 2. Security impact

### Backend authentication is now enforced

All API routes now require an admin token. The backend accepts either of the following authentication methods:

```http
Authorization: Bearer <ADMIN_TOKEN>
X-Admin-Token: <ADMIN_TOKEN>
```

This prevents unauthenticated direct access to high-risk backend endpoints.

Required runtime setting:

```env
ADMIN_TOKEN=<long random token>
```

### Browser origins are explicitly allowlisted

The backend no longer relies on permissive wildcard CORS behavior. Allowed origins are configured through:

```env
ALLOWED_ORIGINS=http://localhost:5173,https://dune-admin.layout.tools
```

This reduces the risk of malicious browser origins making credentialed requests into the admin backend.

### Listen address handling is safer

The backend now normalizes local listen-address values so common shorthand values stay bound to loopback.

Examples:

```text
LISTEN_ADDR=8080            -> 127.0.0.1:8080
LISTEN_ADDR=:8080           -> 127.0.0.1:8080
LISTEN_ADDR=127.0.0.1:8080  -> 127.0.0.1:8080
```

Recommended setting:

```env
LISTEN_ADDR=127.0.0.1:8080
```

### HTTP server timeouts were added

The backend HTTP server now configures explicit read-header, read, write, and idle timeouts. This lowers exposure to slow-client and connection-resource exhaustion behavior.

### Request body limits were added

A JSON request-size limiting helper was introduced and applied to the batch item grant path. This protects sensitive mutation endpoints from oversized request bodies.

### Raw SQL is restricted

The database SQL endpoint is now constrained to read-only style statements. Allowed prefixes include:

```text
SELECT
WITH
SHOW
EXPLAIN
```

The backend rejects semicolon-separated statements and destructive SQL keywords such as:

```text
INSERT, UPDATE, DELETE, DROP, ALTER, TRUNCATE, CREATE, GRANT, REVOKE,
COPY, CALL, DO, EXECUTE, MERGE, VACUUM, ANALYZE
```

This does not replace proper database permissions, but it greatly reduces the chance of destructive SQL being executed through the UI by mistake.

### Kubernetes log streaming was hardened

Log streaming now validates namespace and pod names before building a `kubectl logs` command. WebSocket origin checks also use the configured origin allowlist.

### WebSocket authentication was fixed

Browser WebSocket connections cannot set custom request headers like `X-Admin-Token`. The backend now supports a WebSocket-only query-token fallback for the log stream endpoint:

```text
/api/v1/logs/stream?...&ws_token=<ADMIN_TOKEN>
```

This fallback is limited to the log stream WebSocket route and is intended to support browser-based log streaming while keeping backend authentication enforced.

### Status data was reduced

The status endpoint no longer returns pod IP information. It still reports basic connectivity state, but avoids returning unnecessary internal network details.

### Hardcoded capture credentials were removed

RabbitMQ capture credentials and JWT signing material were removed from source and moved to environment variables:

```env
DUNE_SERVICE_JWT_SIGNING_SECRET=
DUNE_CAPTURE_USER=dune_cap
DUNE_CAPTURE_PASS=<strong random password>
```

Operators should rotate any previously committed or shared credential values.

---

## 3. Required configuration changes

A local `.env` should include at least:

```env
ADMIN_TOKEN=<long random token>
LISTEN_ADDR=127.0.0.1:8080
ALLOWED_ORIGINS=http://localhost:5173,https://dune-admin.layout.tools
DUNE_SERVICE_JWT_SIGNING_SECRET=
DUNE_CAPTURE_USER=dune_cap
DUNE_CAPTURE_PASS=<strong random password>
```

Frontend settings should match:

```text
Backend URL: http://localhost:8080
Admin Token: same value as ADMIN_TOKEN
```

For local development, run the backend and frontend separately:

```powershell
# Terminal 1
cd "Z:\Unreal Projects\Icarus\dune-admin-fork"
go run .
```

```powershell
# Terminal 2
cd "Z:\Unreal Projects\Icarus\dune-admin-fork\web"
npm run dev
```

Open:

```text
http://localhost:5173
```

---

## 4. Added

### Multi-item Give Items workflow

The player Give Item workflow now supports multiple item rows in one operation. Each row can specify:

- Item template
- Stack count
- Grade / quality
- Stack size

The frontend submits batch grants through:

```text
api.players.giveItems(playerID, items)
```

### Batch item grant payload

The existing endpoint remains:

```http
POST /api/v1/players/give-item
```

It now supports a batch payload:

```json
{
  "player_id": 123,
  "items": [
    {
      "template": "ItemTemplateHeavyDarts",
      "qty": 2,
      "quality": 1,
      "stack_size": 1000
    }
  ]
}
```

In the batch payload:

```text
qty        = number of stacks
stack_size = number of items per stack
```

Example:

```text
qty=2, stack_size=1000 -> 2 inventory slots, 2000 total items
```

### Explicit stack grant backend command

A new backend command was added for stack-preserving item grants:

```text
cmdGiveItemStacks(playerID, template, stacks, stackSize, quality)
```

This command creates the requested number of item rows and assigns the requested `stack_size` to each row.

### Backend batch validation

Batch grants now validate:

```text
maximum rows per request: 100
maximum stack count per row: 9999
maximum stack size per row: 9999
quality range: 0 through 5
```

The backend rejects blank templates, empty item lists, non-positive quantities, out-of-range quality values, excessive stack counts, excessive stack sizes, and total quantity overflow risks.

### Go unit tests for batch normalization

Unit tests were added for the batch item request normalizer, including legacy payload compatibility, batch payload parsing, stack-size defaulting, validation failures, and boundary limits.

### Go test workflow

A GitHub Actions workflow was added to run:

```bash
go test ./...
```

on pushes to `main`.

---

## 5. Changed

### Give Item semantics for batch payloads

Batch item grants now preserve stack semantics instead of flattening `qty × stack_size` into a single quantity.

Before:

```text
2 stacks × 1000 Heavy Darts -> 2000 inventory slots requested
```

After:

```text
2 stacks × 1000 Heavy Darts -> 2 inventory slots requested
```

Volume / weight checks still use total item count, so stack grants continue to respect inventory capacity rules beyond slot count.

### Legacy single-item payload remains compatible

The old single-item payload still works:

```json
{
  "player_id": 123,
  "template": "ItemTemplateHeavyDarts",
  "qty": 10,
  "quality": 1
}
```

Legacy behavior remains flat quantity-based for compatibility with existing callers.

### Notification and capture credential usage

Notification publishing now uses configured capture credentials instead of removed hardcoded constants.

### Frontend settings panel

The settings panel now supports both backend URL and admin token configuration for local operation.

---

## 6. Fixed

### Fixed backend startup with bare port values

`LISTEN_ADDR=8080` previously caused:

```text
listen tcp: address 8080: missing port in address
```

The backend now normalizes bare ports to loopback host-port form.

### Fixed Go vet failure in capture output

A redundant newline in a `fmt.Println` call caused `go test ./...` to fail during vet checks. The output was cleaned up.

### Fixed notification compile error after credential cleanup

After removing hardcoded capture constants, notification publishing still referenced the old names. It now uses the configured capture credential helpers.

### Fixed WebSocket log streaming authentication

Log streaming previously failed after backend authentication was introduced because browser WebSockets could not send the admin token header. The backend now supports the route-limited `ws_token` fallback.

### Fixed cheats log SQL error

The cheats log query previously joined against a non-existent column:

```sql
ps.fls_id
```

This caused:

```text
ERROR: column ps.fls_id does not exist (SQLSTATE 42703)
```

The query now joins through `dune.encrypted_accounts` using `encrypted_funcom_id`, then resolves player names through `player_state.account_id`.

### Fixed stack-size behavior for batch grants

Batch item grants no longer require one inventory slot per item when the operator specifies stack size. They now create the requested number of stacks with the requested stack size.

---

## 7. Security notes for operators

- Treat `ADMIN_TOKEN` as a privileged secret.
- Rotate any previously shared or committed credentials.
- Keep `.env`, SSH keys, database snapshots, and generated secrets out of source control.
- Prefer `LISTEN_ADDR=127.0.0.1:8080` for local use.
- Do not expose the backend directly to the internet.
- If remote access is required, place the backend behind TLS, a trusted reverse proxy, and a strong identity provider.
- The WebSocket `ws_token` is placed in the URL because browsers cannot send custom WebSocket headers. Use HTTPS/WSS when operating outside local loopback to avoid token exposure in transit.

---

## 8. Validation

Run backend tests:

```powershell
git pull origin main
go test ./...
```

Run backend locally:

```powershell
go run .
```

Run frontend locally:

```powershell
cd web
npm install
npm run build
npm run dev
```

Recommended manual validation:

1. Confirm backend starts on `127.0.0.1:8080`.
2. Confirm unauthenticated API requests are rejected.
3. Confirm the frontend works after saving `ADMIN_TOKEN` in the settings panel.
4. Confirm Logs -> Cheats (7d) loads without the `ps.fls_id` SQL error.
5. Confirm pod log streaming connects without WebSocket auth errors.
6. Confirm granting 2 stacks of 1000 Heavy Darts creates 2 inventory slots, not 2000.
7. Confirm invalid batch item rows are rejected with clear validation errors.

---

## 9. Known limitations

- The frontend build workflow was not added through the connector because the workflow file write was blocked during this patch series. Local validation with `npm run build` is still recommended.
- WebSocket query-token authentication is a practical browser compatibility solution, but deployments outside localhost should use TLS/WSS and avoid logging full URLs.
- Batch item grants intentionally allow explicit stack sizes up to the configured limit. Operators should use reasonable values that are compatible with game behavior.

---

## 10. Final summary

This release makes Dune Admin safer and more reliable for live server operations.

The most important security improvements are backend-enforced admin-token authentication, explicit CORS allowlisting, safer loopback listen defaults, server timeouts, raw SQL restrictions, request-size controls, Kubernetes log target validation, reduced status data exposure, and removal of hardcoded capture credentials.

The most important operator workflow improvement is the new batch item grant model. Admins can now grant multiple items in one operation while preserving stack count and stack size. A request such as:

```text
2 stacks of 1000 Heavy Darts
```

now creates:

```text
2 inventory entries with stack_size=1000
```

instead of attempting to create 2000 separate inventory entries.

This update establishes a stronger security baseline and a more accurate administration workflow for player inventory management, log review, and local operator use.
