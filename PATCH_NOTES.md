# Dune Admin Release Notes

## Current update: Discord bot command adapter skeleton

### Why this update was made

The farming request/order backend needs a Discord-facing adapter layer before an actual bot runtime is introduced. This slice adds command-shape normalization and payload conversion without adding a networked bot process, Discord gateway client, or new runtime dependency.

### What changed

- Added `discord_bot_adapter.go` with a non-network command adapter for Discord-style farming request/order commands.
- Added adapter support for personal item requests, guild item requests, farm-order creation, filled-order updates, and cancelled-order updates.
- Reused existing inventory request/order validation and normalization instead of creating a second validation path.
- Added `discord_bot_adapter_test.go` coverage for personal requests, guild requests, farm orders, fill updates, cancel updates, and unsupported command rejection.

### Security and operator impact

- This is an adapter skeleton only. It does not connect to Discord, register slash commands, open a gateway connection, or execute bot actions on its own.
- The adapter does not mutate player inventory, guild storage, claim rewards, currency, XP, Player 360, or any game-state table.
- Player 360 remains read-only. Self-service player-card actions remain blocked until Discord identity-to-player mapping exists and is explicitly enforced.
- The adapter deliberately maps bot-style inputs into the existing request/order coordination model so future bot runtime work can stay thin and testable.

### Validation

User-provided local validation completed cleanly after the scoped frontend lint fix:

```bash
go test -v ./...
go build ./...
npm install
npm audit --audit-level=high
npm run typecheck
npm run lint
npm run build
```

The local `update.sh` run failed only at the Git auto-commit step because Git author identity was not configured in the local checkout.

---

## Previous update: Inventory request/order backend coordination model

### Why this update was made

The next safe community-support slice needs a way to collect personal and guild farming requests without writing to player inventory or guild storage. This creates a backend coordination ledger for requests and farming orders while preserving the current rule that Player 360 remains read-only and self-service inventory changes are not enabled.

### What changed

- Added `inventory_requests.go` with a local JSON-backed request/order store.
- Added personal and guild inventory request modeling with validation for scope, requester, guild, item name, quantity, notes, and status.
- Added farming order modeling that groups one or more requests and tracks assignee, status, completion timestamp, and notes.
- Added status propagation so order creation marks linked requests `ordered`, filled orders mark linked requests `fulfilled`, and cancelled orders return linked requests to `open`.
- Registered protected backend endpoints:
  - `GET /api/v1/inventory/requests`
  - `POST /api/v1/inventory/requests`
  - `PATCH /api/v1/inventory/requests/{id}`
  - `GET /api/v1/inventory/orders`
  - `POST /api/v1/inventory/orders`
  - `PATCH /api/v1/inventory/orders/{id}`
- Added `inventory_requests_test.go` coverage for request normalization, invalid request payload rejection, create/list/order/fill lifecycle, and missing-request order rejection.
- Added in-process mutex serialization around the local JSON store to reduce concurrent write clobber risk.
- Updated CORS to allow `PATCH` for browser-based update endpoints.
- Added `docs/inventory-requests-orders.md` with endpoint, model, validation, storage, and safety-boundary notes.

### Security and operator impact

- This feature is coordination-only. It does not mutate player inventory, guild storage, claim rewards, currency, XP, Player 360, or any game-state table.
- Inventory request/order endpoints are protected by the normal backend auth middleware and are not public Discord OAuth paths.
- The default `inventory-requests.json` store is local file storage written with `0600` permissions.
- The JSON store is appropriate for the current backend slice, but multi-instance or production use should move this ledger to a durable database table.

### Validation

Validated clean from the local checkout after the strict admin-token fixture repair:

```bash
go test ./...
go build ./...
```

Operator/browser validation still recommended before release: personal request creation, guild request creation, request filtering, order creation, order fill/cancel status propagation, and browser preflight behavior for `PATCH` update endpoints.

---

## Previous update: Discord auth route registration and session tests

### Why this update was made

Discord OAuth and session handlers existed in the backend, but the centralized route registry did not expose the Discord auth endpoints. The login, callback, session identity, logout, and registered-user inspection paths must be reachable before Discord-backed operator identity and later safe self-service workflows can be built on top of them.

### What changed

- Registered Discord auth endpoints in `routes.go`:
  - `GET /api/v1/auth/discord/login`
  - `GET /api/v1/auth/discord/callback`
  - `GET /api/v1/auth/discord/me`
  - `POST /api/v1/auth/discord/logout`
  - `GET /api/v1/auth/discord/users`
- Added `discord_auth_test.go` coverage for route registration, Discord role mapping, default normal-role behavior, session lookup helpers, expired-session eviction, session hash generation, and logout invalidation.
- Isolated Discord route tests from any local `discord-users.json` file by using a temporary test user-store path.

### Security and operator impact

- Discord login and callback remain the only Discord auth routes intended to be public through the existing `isPublicPath` allowlist.
- Session identity, logout, and registered-user inspection are now registered backend routes and remain subject to the server's existing auth middleware behavior when served normally.
- Player 360 remains read-only. This change does not add player mutation or self-service player-card behavior.
- Discord role/session behavior now has durable backend test coverage before additional Discord bot and self-service work is layered on top.

### Validation

Validation required from the local checkout or CI:

```powershell
.\update.ps1
```

or:

```bash
go test ./...
```

Also manually validate that Discord OAuth login redirects when configured, callback rejects missing/invalid state, logout clears the session cookie, and `/api/v1/auth/discord/me` reports the expected auth context for an authenticated Discord session or admin token.

---

## Previous update: Startup hardening and DB-down UI gating

### Why this update was made

Startup and database connectivity are release-critical. If SSH, tunnel, discovery, or DB connection fails, the application must not look operational while DB-backed tools are unusable.

### What changed

- Added `config_paths.go` with config-safe local path expansion and validation.
- Added support for documented Windows percent-variable home-directory SSH key paths.
- Rejected unsupported PowerShell-style path expressions during validation.
- Tightened runtime validation for SSH, DB, admin token, tunnel mode, and port configuration.
- Changed startup so invalid configuration exits and tells the operator to rerun setup.
- Changed connection failure messaging from ambiguous startup behavior to explicit degraded mode.
- Added frontend DB-down gating for DB-backed tabs.
- Added a DB-unavailable banner when backend status reports `db_connected=false`.
- Added a DB-required panel with SSH/DB/tunnel state and reconnect retry.
- Added frontend startup diagnostic typing.
- Expanded the production security release gate for startup reliability, config validation, injection resistance, logging redaction, encryption in transit, and secret-at-rest requirements.
- Added Battlegroup config-management and observability design notes from operator review intake.

### Security and operator impact

- Invalid configuration should fail earlier and more clearly.
- DB-backed tabs no longer present as usable when DB connectivity is down.
- Secret values are validated as opaque values for presence/control characters.
- Full encrypted secret storage remains a P0 release-gate item.
- The direct DB connection-construction patch was blocked by the connector safety filter and should be validated from the local checkout if `update.ps1` reports compile or startup issues.

### Validation

Validation required in the Windows development environment:

```powershell
.\update.ps1
```

Also manually validate startup repair, unsupported path rejection, supported path expansion, and DB-down frontend gating.

---

## Previous update: Inventory Studio v2 post-action diff panel

Inventory Studio v2 added a post-action diff panel for confirmed add, repair, and removal workflows.

---

## Previous update: Inventory Studio v2 confirmed catalog item add

Inventory Studio v2 added confirmed catalog-item add with quantity and quality inputs.

---

## Previous update: Inventory Studio v2 confirmed item removal

Inventory Studio v2 added confirmed selected-item removal with a before-action inventory snapshot.

---

## Previous update: Inventory Studio v2 confirmed item repair

Inventory Studio v2 added confirmed selected-item repair with a before-action inventory snapshot.

---

## Previous update: Inventory Studio v2 item catalog browser

Inventory Studio v2 added a read-only item catalog browser.
