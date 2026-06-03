# AppSec Endpoint Audit

## Purpose

This document is the durable application-security audit record for DA Manager backend endpoints.

The audit covers public, Discord-session, admin-token, WebSocket, infrastructure, database, and high-risk mutation surfaces. It is intended to be updated as endpoint behavior changes.

## Audit status

| Area | Status | Notes |
|---|---|---|
| Route inventory | In Progress | Initial inventory taken from `routes.go`. |
| Auth boundary review | In Progress | Initial middleware review complete; `appsec_auth_boundary_test.go` covers public, self-service, Discord self-session, representative admin, and WebSocket-ticket boundaries. Generated full-route coverage remains a follow-up. |
| Input validation review | In Progress | Full handler-by-handler review still required. |
| Mutation reason coverage | In Progress | Needs endpoint-by-endpoint confirmation. |
| SAST | Pending | Run and record tool/version/result. |
| DAST | Pending | Run and record tool/version/result. |
| Dependency review | Pending | Run and record tool/version/result. |
| Manual abuse-case review | Pending | Endpoint-by-endpoint manual cases still required. |

## Reviewed source files

Initial static review used:

- `routes.go`
- `auth.go`
- `server.go`
- `discord_auth.go`
- `appsec_auth_boundary_test.go`
- existing roadmap and release-tracking documents

## Global middleware and boundary summary

The backend server wraps registered routes in this order:

```text
corsMiddleware(auditMiddleware(mutationSafetyMiddleware(authMiddleware(mux))))
```

Security-relevant observed behavior:

- `OPTIONS` requests are handled by CORS/preflight logic.
- Public paths bypass normal auth only when listed in `isPublicPath`.
- Current public paths are:
  - `GET /api/v1/public/status`
  - `GET /api/v1/auth/discord/login`
  - `GET /api/v1/auth/discord/callback`
- `GET /api/v1/logs/stream` with WebSocket upgrade uses a one-time log-stream ticket instead of static admin-token query usage.
- Normal admin access accepts either `Authorization: Bearer <token>` or `X-Admin-Token` when the configured backend admin token is strict-valid and matches in constant time.
- Discord admin sessions can reach protected routes when Discord auth is enabled.
- Registered non-admin Discord sessions can reach only these self-session/self-service surfaces:
  - `GET /api/v1/auth/discord/me`
  - `POST /api/v1/auth/discord/logout`
  - `/api/v1/self/*`
- Backend startup normalizes loopback listen addresses and fails closed for unsafe non-loopback exposure unless explicitly configured elsewhere.
- Error responses are passed through sensitive-text redaction before being returned.

## Endpoint classification legend

| Classification | Meaning |
|---|---|
| Public | Explicitly allowed by `isPublicPath`; does not require admin token. |
| Self-service | Requires registered Discord session and path under `/api/v1/self/*`, or admin access. |
| Discord self-session | Requires registered Discord session for the caller's own auth/session context, or admin access. |
| Admin | Requires strict admin token or Discord admin session. |
| WebSocket ticket | Requires valid one-time stream ticket for WebSocket log stream. |
| Mutation | Changes runtime, database, storage, player, request/order, or infrastructure state. |
| High risk | Can expose sensitive data, run commands, change player/game state, or affect infrastructure. |

## Endpoint inventory

### Public endpoints

| Method | Path | Handler | Classification | Initial notes |
|---|---|---|---|---|
| GET | `/api/v1/public/status` | `handlePublicStatus` | Public | Redacted public health only. |
| GET | `/api/v1/auth/discord/login` | `handleDiscordLogin` | Public | OAuth initiation. Review state/nonce/session behavior. |
| GET | `/api/v1/auth/discord/callback` | `handleDiscordCallback` | Public | OAuth callback. Review state validation, cookie flags, redirect behavior, error handling. |

### Discord/self-service endpoints

| Method | Path | Handler | Classification | Initial notes |
|---|---|---|---|---|
| GET | `/api/v1/self/player-link` | `handleSelfPlayerLink` | Self-service | Must enforce current Discord session to linked player only. |
| GET | `/api/v1/self/player-card` | `handleSelfPlayerCard` | Self-service | Read-only self player card. Confirm no unrelated player lookup. |
| GET | `/api/v1/auth/discord/me` | `handleDiscordMe` | Discord self-session | Normal registered Discord sessions can inspect their own auth context. |
| POST | `/api/v1/auth/discord/logout` | `handleDiscordLogout` | Discord self-session mutation | Normal registered Discord sessions can clear their own session cookie and in-memory session. |
| GET | `/api/v1/auth/discord/users` | `handleDiscordUsers` | Admin | Registered-user review surface. |
| GET | `/api/v1/auth/discord/player-links` | `handleListDiscordPlayerLinks` | Admin | Identity mapping list. Sensitive association data. |
| POST | `/api/v1/auth/discord/player-links` | `handleUpsertDiscordPlayerLink` | Admin mutation | Identity mapping mutation. Needs validation, reason/audit review. |
| DELETE | `/api/v1/auth/discord/player-links/{discord_id}` | `handleDeleteDiscordPlayerLink` | Admin mutation | Identity mapping mutation. Needs validation, reason/audit review. |

### Core status, diagnostics, audit, and safety endpoints

| Method | Path | Handler | Classification | Initial notes |
|---|---|---|---|---|
| GET | `/api/v1/status` | `handleStatus` | Admin | May expose runtime, DB, SSH, or tunnel status. Review redaction. |
| POST | `/api/v1/reconnect` | `handleReconnect` | Admin mutation/high risk | Reopens DB/SSH/tunnels. Needs audit/reason review and abuse-case coverage. |
| GET | `/api/v1/connectivity/diagnostics` | `handleConnectivityDiagnostics` | Admin/high risk | Review host, port, username, and secret redaction. |
| GET | `/api/v1/diagnostics/export` | `handleDiagnosticExport` | Admin/high risk | Review raw vs redacted content and external-sharing warnings. |
| GET | `/api/v1/audit/events` | `handleAdminAuditEvents` | Admin | Audit log review endpoint. Review pagination, filtering, and data exposure. |
| GET | `/api/v1/mutation-safety/classify` | `handleMutationSafetyClassify` | Admin | Classification helper. Review fallback behavior and coverage. |

### Battlegroup endpoints

| Method | Path | Handler | Classification | Initial notes |
|---|---|---|---|---|
| GET | `/api/v1/battlegroup/status` | `handleBGStatus` | Admin/high risk | Infrastructure status exposure. Review redaction. |
| GET | `/api/v1/battlegroup/health` | `handleBGHealth` | Admin/high risk | Kubernetes diagnostic exposure. Review fixed-command constraints and redaction. |
| POST | `/api/v1/battlegroup/exec` | `handleBGExec` | Admin mutation/high risk | Server-control command path. Must require confirmation, reason, audit, allowlist. |
| GET | `/api/v1/battlegroup/pods` | `handleBGPods` | Admin/high risk | Pod metadata exposure. Review namespace and name validation. |

### Player read endpoints

| Method | Path | Handler | Classification | Initial notes |
|---|---|---|---|---|
| GET | `/api/v1/players` | `handleGetPlayers` | Admin | Player list. PII/game data exposure review. |
| GET | `/api/v1/players/online` | `handleGetOnlineState` | Admin | Online-state exposure. |
| GET | `/api/v1/players/currency` | `handleGetCurrency` | Admin | Currency exposure. |
| GET | `/api/v1/players/factions` | `handleGetFactions` | Admin | Faction exposure. |
| GET | `/api/v1/players/specs` | `handleGetSpecs` | Admin | Specialization exposure. |
| GET | `/api/v1/players/templates` | `handleGetTemplates` | Admin | Item template exposure. |
| GET | `/api/v1/players/{id}/profile` | `handleGetPlayerProfile` | Admin | Player 360 read-only profile. Confirm safe section-level errors. |
| GET | `/api/v1/players/{id}/inventory` | `handleGetInventory` | Admin | Inventory exposure. Path ID validation review required. |
| GET | `/api/v1/players/{id}/journey` | `handleGetJourney` | Admin | Progression exposure. |
| GET | `/api/v1/players/{id}/char-xp` | `handleGetCharXP` | Admin | XP exposure. |
| GET | `/api/v1/players/{id}/specs` | `handleGetPlayerSpecs` | Admin | Player specialization exposure. |
| GET | `/api/v1/players/{id}/vehicles` | `handleGetPlayerVehicles` | Admin | Vehicle exposure. |
| GET | `/api/v1/players/partitions` | `handleGetPartitions` | Admin | Partition data exposure. |
| GET | `/api/v1/players/{id}/events` | `handleGetPlayerEvents` | Admin | Event exposure. |
| GET | `/api/v1/players/{id}/dungeons` | `handleGetPlayerDungeons` | Admin | Dungeon/progression exposure. |

### Player mutation endpoints

| Method | Path | Handler | Classification | Initial notes |
|---|---|---|---|---|
| POST | `/api/v1/players/templates/refresh` | `handleRefreshTemplates` | Admin mutation | Refreshes template cache. Review audit/reason coverage. |
| POST | `/api/v1/players/give-item` | `handleGiveItems` | Admin mutation/high risk | Direct inventory write. Must require reason/audit/validation. |
| POST | `/api/v1/players/give-currency` | `handleGiveCurrency` | Admin mutation/high risk | Economy mutation. |
| POST | `/api/v1/players/grant-live` | `handleGrantLive` | Admin mutation/high risk | Claim reward/live grant path. |
| POST | `/api/v1/players/give-faction-rep` | `handleGiveFactionRep` | Admin mutation/high risk | Faction mutation. |
| POST | `/api/v1/players/give-scrip` | `handleGiveScrip` | Admin mutation/high risk | Currency mutation. |
| POST | `/api/v1/players/award-xp` | `handleAwardXP` | Admin mutation/high risk | XP mutation. |
| POST | `/api/v1/players/award-char-xp` | `handleAwardCharXP` | Admin mutation/high risk | Character XP mutation. |
| POST | `/api/v1/players/award-intel` | `handleAwardIntel` | Admin mutation/high risk | Intel mutation. |
| POST | `/api/v1/players/kick` | `handleKick` | Admin mutation/high risk | Active session/player disruption. |
| DELETE | `/api/v1/players/item/{id}` | `handleDeleteItem` | Admin mutation/high risk | Inventory deletion. |
| POST | `/api/v1/players/item/stack-size` | `handleSetItemStackSize` | Admin mutation/high risk | Direct item stack-size update. New route requires validation. |
| POST | `/api/v1/players/reset-spec` | `handleResetSpec` | Admin mutation/high risk | Progression mutation. |
| POST | `/api/v1/players/set-faction-tier` | `handleSetFactionTier` | Admin mutation/high risk | Faction mutation. |
| POST | `/api/v1/players/journey/complete` | `handleJourneyComplete` | Admin mutation/high risk | Journey mutation. |
| POST | `/api/v1/players/journey/reset` | `handleJourneyReset` | Admin mutation/high risk | Journey mutation. |
| POST | `/api/v1/players/journey/wipe` | `handleJourneyWipe` | Admin mutation/high risk | Destructive journey mutation. |
| POST | `/api/v1/players/delete-tutorials` | `handleDeleteTutorials` | Admin mutation/high risk | Destructive/tutorial state mutation. |
| POST | `/api/v1/players/wipe-codex` | `handleWipeCodex` | Admin mutation/high risk | Destructive codex mutation. |
| POST | `/api/v1/players/set-spec-xp` | `handleSetSpecXP` | Admin mutation/high risk | Specialization XP mutation. |
| POST | `/api/v1/players/repair-item` | `handleRepairItem` | Admin mutation/high risk | Inventory durability mutation. |
| POST | `/api/v1/players/teleport` | `handleTeleportPlayer` | Admin mutation/high risk | Movement/rescue mutation. |

### Inventory request/order endpoints

| Method | Path | Handler | Classification | Initial notes |
|---|---|---|---|---|
| GET | `/api/v1/inventory/requests` | `handleListInventoryRequests` | Admin | Coordination data. Review tenant/user visibility. |
| POST | `/api/v1/inventory/requests` | `handleCreateInventoryRequest` | Admin mutation | Local coordination store mutation. |
| PATCH | `/api/v1/inventory/requests/{id}` | `handleUpdateInventoryRequest` | Admin mutation | Local coordination store mutation. |
| GET | `/api/v1/inventory/orders` | `handleListInventoryOrders` | Admin | Coordination data. |
| POST | `/api/v1/inventory/orders` | `handleCreateInventoryOrder` | Admin mutation | Local coordination store mutation. |
| PATCH | `/api/v1/inventory/orders/{id}` | `handleUpdateInventoryOrder` | Admin mutation | Local coordination store mutation. |

### Database endpoints

| Method | Path | Handler | Classification | Initial notes |
|---|---|---|---|---|
| GET | `/api/v1/database/tables` | `handleDBTables` | Admin/high risk | Schema exposure. |
| GET | `/api/v1/database/describe` | `handleDBDescribe` | Admin/high risk | Schema exposure. Review identifier validation. |
| GET | `/api/v1/database/sample` | `handleDBSample` | Admin/high risk | Data exposure. Review sampling limits/redaction. |
| GET | `/api/v1/database/search` | `handleDBSearch` | Admin/high risk | Dynamic search. Review SQL injection and result limits. |
| GET | `/api/v1/database/functions` | `handleDBFunctions` | Admin/high risk | Function metadata exposure. |
| GET | `/api/v1/database/functions/inspect` | `handleDBFunctionInspect` | Admin/high risk | Function definition exposure. |
| POST | `/api/v1/database/sql` | `handleDBSQL` | Admin mutation/high risk | Manual SQL runner. Confirm read-only guard, dangerous keyword detection, reason/audit behavior, and timeout/result limits. |

### Log endpoints

| Method | Path | Handler | Classification | Initial notes |
|---|---|---|---|---|
| GET | `/api/v1/logs/pods` | `handleLogPods` | Admin/high risk | Pod/log target listing. Review target validation. |
| POST | `/api/v1/logs/stream-ticket` | `handleIssueLogStreamTicket` | Admin mutation/high risk | Issues one-time scoped ticket. Review TTL, scope, replay, and audit behavior. |
| GET | `/api/v1/logs/stream` | `handleLogStream` | WebSocket ticket/high risk | Requires ticket on WebSocket upgrade. Review non-upgrade behavior. |
| GET | `/api/v1/logs/cheats` | `handleGetCheatLog` | Admin/high risk | Sensitive log exposure. Review redaction and result limits. |

### Notification endpoint

| Method | Path | Handler | Classification | Initial notes |
|---|---|---|---|---|
| POST | `/api/v1/notify` | `handleNotify` | Admin mutation/high risk | Broadcast/notification path. Review message validation, size limits, audit reason, and abuse cases. |

### Storage endpoints

| Method | Path | Handler | Classification | Initial notes |
|---|---|---|---|---|
| GET | `/api/v1/storage` | `handleListStorage` | Admin | Storage list exposure. |
| GET | `/api/v1/storage/{id}/items` | `handleGetStorageItems` | Admin | Storage item exposure. |
| POST | `/api/v1/storage/{id}/give-item` | `handleGiveItemToStorage` | Admin mutation/high risk | Storage inventory mutation. Requires confirmation, reason, validation, audit review. |

### Blueprint endpoints

| Method | Path | Handler | Classification | Initial notes |
|---|---|---|---|---|
| GET | `/api/v1/blueprints` | `handleListBlueprints` | Admin | Blueprint metadata exposure. |
| GET | `/api/v1/blueprints/{id}/export` | `handleExportBlueprint` | Admin/high risk | Blueprint export. Review output size and data exposure. |
| POST | `/api/v1/blueprints/import` | `handleImportBlueprint` | Admin mutation/high risk | Import mutation. Requires confirmation, reason, payload limits, validation, audit review. |

## Initial findings and remediation backlog

| ID | Severity | Status | Finding | Recommended action | Validation evidence |
|---|---|---|---|---|---|
| ASEA-001 | High | Validated partial remediation | No generated endpoint inventory/auth-boundary regression test existed for the full `routes.go` surface. Initial regression coverage now exists for public allowlist behavior, self-service path classification, admin-only representative routes, and WebSocket-ticket denial. Generated full-route coverage remains an open hardening follow-up. | Expand `appsec_auth_boundary_test.go` to generated full-route coverage in a future pass. | `appsec_auth_boundary_test.go` added and validated clean through `./update.sh`. |
| ASEA-002 | Medium | Remediated pending validation | Discord `me` and `logout` handlers supported session-cookie behavior, but middleware limited registered non-admin Discord sessions to `/api/v1/self/*`. | Added a narrow Discord self-session middleware route allowlist for `GET /api/v1/auth/discord/me` and `POST /api/v1/auth/discord/logout`; added regression tests to confirm normal Discord sessions can use those two routes but not admin routes. | Local validation pending. |
| ASEA-003 | High | Open | High-risk mutation endpoints require endpoint-by-endpoint verification for `X-Admin-Reason`, audit logging, mutation-safety classification, request-size limits, and pre/post-change review behavior. | Create a mutation endpoint checklist and add automated tests for reason/audit coverage where feasible. | Pending. |
| ASEA-004 | High | Open | Database search/manual SQL endpoints need dedicated SQL injection, read-only guard, timeout, result-limit, and redaction review. | Perform handler-specific review and add abuse-case tests for dangerous SQL, multi-statement attempts, identifier injection, and large result sets. | Pending. |
| ASEA-005 | High | Open | Infrastructure command/log endpoints need dedicated command allowlist, target validation, ticket replay, TTL, and data-redaction review. | Review battlegroup exec, diagnostics, log pod list, stream-ticket, stream, and cheat-log handlers; add tests for invalid targets and replay attempts. | Pending. |
| ASEA-006 | Medium | Open | CORS allows configured origins and `X-Admin-Token`/`X-Admin-Reason` headers. Browser-stored admin token remains JavaScript-readable per known issue. | Continue migration plan toward memory-only or HttpOnly secure session-cookie auth; document CSRF approach before cookie-based admin auth. | Pending. |

## Manual abuse-case checklist

Use this checklist during the full audit pass:

- Send each public route without admin token and confirm only intended public routes respond.
- Send each admin route without token and confirm `401` or expected denial.
- Send each admin route with malformed backend token configuration and confirm fail-closed behavior.
- Send each self-service route with no Discord session, unregistered Discord session, registered unlinked Discord session, and linked Discord session.
- Send every mutation with missing, empty, and oversized JSON bodies.
- Send every path parameter route with non-numeric, negative, zero, oversized, encoded slash, and traversal-like values.
- Send every query parameter route with missing, blank, oversized, SQL metacharacters, wildcard-heavy, and Unicode control inputs.
- Attempt SQL multi-statements, comments, dangerous keywords, function calls, copy/export operations, and long-running queries against database endpoints.
- Attempt command target injection and invalid Kubernetes target names against infrastructure endpoints.
- Attempt WebSocket log-stream reuse, expired ticket use, wrong target use, non-upgrade GET behavior, and cross-origin use.
- Confirm all high-risk mutations produce audit events with action, target, reason, and outcome.
- Confirm all sensitive errors are redacted.

## Validation commands to record when complete

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

Canonical local validation:

```bash
./update.sh
```

Security validation to record when run:

```bash
gitleaks detect
# govulncheck ./...
# gosec ./...
# trivy fs .
# DAST baseline against local dev/preview deployment
```
