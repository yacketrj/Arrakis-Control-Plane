# AppSec Endpoint Audit

## Purpose

This document is the durable application-security audit record for DA Manager backend endpoints.

The audit covers public, Discord-session, admin-token, WebSocket, infrastructure, database, and high-risk mutation surfaces. It is intended to be updated as endpoint behavior changes.

## Audit status

| Area | Status | Notes |
|---|---|---|
| Route inventory | In Progress | Initial inventory taken from `routes.go`. |
| Auth boundary review | In Progress | Initial middleware review complete; `appsec_auth_boundary_test.go` covers public, self-service, Discord self-session, representative admin, and WebSocket-ticket boundaries. Generated full-route coverage remains a follow-up. |
| Input validation review | In Progress | Database handler parameter bounds/control-character checks, numeric function OID validation, unsafe SQL rejection, database output redaction, infrastructure namespace/command checks, log target checks, and log-ticket replay tests are in progress/validated by slice. Full handler-by-handler review still required. |
| Mutation reason coverage | In Progress | High-risk/destructive classification, oversized-body reason-enforcement, audit metadata parsing, and colorized validation-output changes are validated. Full endpoint-by-endpoint audit-event assertion coverage still required. |
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
- `audit_log.go`
- `mutation_safety.go`
- `mutation_safety_handler.go`
- `mutation_safety_test.go`
- `handlers_database.go`
- `handlers_database_test.go`
- `docs/database-endpoint-security.md`
- `handlers_battlegroup.go`
- `handlers_logs.go`
- `ws_ticket.go`
- `runtime_commands.go`
- `infrastructure_security_test.go`
- `docs/infrastructure-log-endpoint-security.md`
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
- Mutation safety classification now treats reconnect, Battlegroup exec, database SQL, log stream ticket issuance, notify, direct item-row edits, inventory writes, player state edits, storage writes, and destructive reset/wipe/delete/import paths as high-risk or destructive as appropriate.
- Database endpoint handlers now trim and bound query parameters, reject unsafe control characters, require numeric function OIDs, redact sampled/search rows, redact manual SQL output, and keep existing database row-limit guardrails.
- Infrastructure/log handlers now normalize Battlegroup command input, enforce a static command allowlist, validate runtime namespace before Kubernetes command construction, redact returned runtime/log/cheat outputs, and test log-stream ticket single-use, wrong-target, and expiry behavior.
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
| POST | `/api/v1/reconnect` | `handleReconnect` | Admin mutation/high risk | Reopens DB/SSH/tunnels. Mutation-safety classification now marks this high risk. |
| GET | `/api/v1/connectivity/diagnostics` | `handleConnectivityDiagnostics` | Admin/high risk | Review host, port, username, and secret redaction. |
| GET | `/api/v1/diagnostics/export` | `handleDiagnosticExport` | Admin/high risk | Review raw vs redacted content and external-sharing warnings. |
| GET | `/api/v1/audit/events` | `handleAdminAuditEvents` | Admin | Audit log review endpoint. Review pagination, filtering, and data exposure. |
| GET | `/api/v1/mutation-safety/classify` | `handleMutationSafetyClassify` | Admin | Classification helper. Review fallback behavior and coverage. |

### Battlegroup endpoints

| Method | Path | Handler | Classification | Initial notes |
|---|---|---|---|---|
| GET | `/api/v1/battlegroup/status` | `handleBGStatus` | Admin/high risk | Validates Kubernetes namespace before command construction and redacts command output. |
| GET | `/api/v1/battlegroup/health` | `handleBGHealth` | Admin/high risk | Validates Kubernetes namespace before fixed health commands and redacts section output/errors. |
| POST | `/api/v1/battlegroup/exec` | `handleBGExec` | Admin mutation/high risk | Server-control command path. Normalizes command input and enforces static allowlist. |
| GET | `/api/v1/battlegroup/pods` | `handleBGPods` | Admin/high risk | Validates Kubernetes namespace before command construction and redacts returned lines. |

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
| POST | `/api/v1/players/give-item` | `handleGiveItems` | Admin mutation/high risk | Direct inventory write. Requires reason/audit/validation. |
| POST | `/api/v1/players/give-currency` | `handleGiveCurrency` | Admin mutation/high risk | Economy mutation. |
| POST | `/api/v1/players/grant-live` | `handleGrantLive` | Admin mutation/high risk | Claim reward/live grant path. |
| POST | `/api/v1/players/give-faction-rep` | `handleGiveFactionRep` | Admin mutation/high risk | Faction mutation. |
| POST | `/api/v1/players/give-scrip` | `handleGiveScrip` | Admin mutation/high risk | Currency mutation. |
| POST | `/api/v1/players/award-xp` | `handleAwardXP` | Admin mutation/high risk | XP mutation. |
| POST | `/api/v1/players/award-char-xp` | `handleAwardCharXP` | Admin mutation/high risk | Character XP mutation. |
| POST | `/api/v1/players/award-intel` | `handleAwardIntel` | Admin mutation/high risk | Intel mutation. |
| POST | `/api/v1/players/kick` | `handleKick` | Admin mutation/high risk | Active session/player disruption. |
| DELETE | `/api/v1/players/item/{id}` | `handleDeleteItem` | Admin mutation/destructive | Inventory deletion. |
| POST | `/api/v1/players/item/stack-size` | `handleSetItemStackSize` | Admin mutation/high risk | Direct item stack-size update. Mutation-safety classification now marks this high risk. |
| POST | `/api/v1/players/reset-spec` | `handleResetSpec` | Admin mutation/destructive | Progression mutation. |
| POST | `/api/v1/players/set-faction-tier` | `handleSetFactionTier` | Admin mutation/high risk | Faction mutation. |
| POST | `/api/v1/players/journey/complete` | `handleJourneyComplete` | Admin mutation/high risk | Journey mutation. |
| POST | `/api/v1/players/journey/reset` | `handleJourneyReset` | Admin mutation/destructive | Journey mutation. |
| POST | `/api/v1/players/journey/wipe` | `handleJourneyWipe` | Admin mutation/destructive | Journey mutation. |
| POST | `/api/v1/players/delete-tutorials` | `handleDeleteTutorials` | Admin mutation/destructive | Destructive/tutorial state mutation. |
| POST | `/api/v1/players/wipe-codex` | `handleWipeCodex` | Admin mutation/destructive | Destructive codex mutation. |
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
| GET | `/api/v1/database/describe` | `handleDBDescribe` | Admin/high risk | Schema exposure. Parameter is now bounded and control-character checked; command uses parameterized information-schema lookup. |
| GET | `/api/v1/database/sample` | `handleDBSample` | Admin/high risk | Data exposure. Parameter is now bounded/control-character checked; limit clamps to 200; response rows are redacted. |
| GET | `/api/v1/database/search` | `handleDBSearch` | Admin/high risk | Dynamic search. Term is now bounded/control-character checked; command uses parameterized search; response rows are redacted. |
| GET | `/api/v1/database/functions` | `handleDBFunctions` | Admin/high risk | Function metadata exposure. Optional term/category parameters are bounded and control-character checked. |
| GET | `/api/v1/database/functions/inspect` | `handleDBFunctionInspect` | Admin/high risk | Function definition exposure. OID is now bounded and numeric-only. |
| POST | `/api/v1/database/sql` | `handleDBSQL` | Admin mutation/high risk | Manual SQL runner. SQL is trimmed before read-only validation, output is redacted, and existing single-statement/read-only/result-limit guardrails remain. SQL timeout review remains open. |

### Log endpoints

| Method | Path | Handler | Classification | Initial notes |
|---|---|---|---|---|
| GET | `/api/v1/logs/pods` | `handleLogPods` | Admin/high risk | Validates runtime namespace/targets and redacts Docker display names. |
| POST | `/api/v1/logs/stream-ticket` | `handleIssueLogStreamTicket` | Admin mutation/high risk | Issues scoped, one-time, 60-second tickets after target validation; regression tests cover replay, wrong target, expiry, and invalid targets. |
| GET | `/api/v1/logs/stream` | `handleLogStream` | WebSocket ticket/high risk | Rejects legacy `ws_token`, validates target, and streams redacted log lines. |
| GET | `/api/v1/logs/cheats` | `handleGetCheatLog` | Admin/high risk | Redacts returned cheat-log fields before returning rows. |

### Notification endpoint

| Method | Path | Handler | Classification | Initial notes |
|---|---|---|---|---|
| POST | `/api/v1/notify` | `handleNotify` | Admin mutation/high risk | Broadcast/notification path. Mutation-safety classification now marks this high risk. Review message validation, size limits, audit reason, and abuse cases. |

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
| POST | `/api/v1/blueprints/import` | `handleImportBlueprint` | Admin mutation/destructive | Import mutation. Requires confirmation, reason, payload limits, validation, audit review. |

## Initial findings and remediation backlog

| ID | Severity | Status | Finding | Recommended action | Validation evidence |
|---|---|---|---|---|---|
| ASEA-001 | High | Validated partial remediation | No generated endpoint inventory/auth-boundary regression test existed for the full `routes.go` surface. Initial regression coverage now exists for public allowlist behavior, self-service path classification, admin-only representative routes, and WebSocket-ticket denial. Generated full-route coverage remains an open hardening follow-up. | Expand `appsec_auth_boundary_test.go` to generated full-route coverage in a future pass. | `appsec_auth_boundary_test.go` added and validated clean through `./update.sh`. |
| ASEA-002 | Medium | Validated remediation | Discord `me` and `logout` handlers supported session-cookie behavior, but middleware limited registered non-admin Discord sessions to `/api/v1/self/*`. | Added a narrow Discord self-session middleware route allowlist for `GET /api/v1/auth/discord/me` and `POST /api/v1/auth/discord/logout`; added regression tests to confirm normal Discord sessions can use those two routes but not admin routes. | `./update.sh` passed clean; build emitted non-blocking `[PLUGIN_TIMINGS]` warning for `@tailwindcss/vite:generate:build`. |
| ASEA-003 | High | Validated partial remediation | High-risk mutation endpoints require endpoint-by-endpoint verification for `X-Admin-Reason`, audit logging, mutation-safety classification, request-size limits, and pre/post-change review behavior. Initial review found some high-risk mutation paths were under-classified as medium. | Tightened mutation-safety classification for reconnect, database SQL, log stream ticket issuance, notify, and direct item-row edits. Added high-risk/destructive route coverage tests and oversized-body reason-enforcement test. Fixed audit metadata JSON test payload and added colorized validation output for `RUN`, `PASS`, and `FAIL`. Full endpoint-by-endpoint audit-event assertion coverage remains required. | `./update.sh` passed clean after the audit metadata JSON fix and update-script color-output change. |
| ASEA-004 | High | Validated partial remediation | Database search/manual SQL endpoints needed dedicated SQL injection, read-only guard, result-limit, and redaction review. Initial review confirmed parameterization/safe identifier quoting in key commands, then added handler-level parameter bounds, control-character checks, numeric OID validation, SQL trimming, and output redaction. | Added database handler security tests and `docs/database-endpoint-security.md`. SQL timeout review, expanded read-only bypass tests, live-data redaction review, and manual abuse-case validation remain open. | `./update.sh` passed clean after database handler hardening and tests. |
| ASEA-005 | High | Partially remediated pending validation | Infrastructure command/log endpoints needed dedicated command allowlist, target validation, ticket replay, TTL, and data-redaction review. Initial review found namespace validation and output redaction should be applied more consistently across runtime command/log paths. | Added Battlegroup command normalization/allowlist checks, shared runtime namespace validation, infrastructure output redaction, cheat-log field redaction, log-ticket replay/wrong-target/expiry tests, and `docs/infrastructure-log-endpoint-security.md`. Handler-level SSH/database-stub tests, command timeout review, WebSocket origin review, live runtime/manual validation, and real-output redaction review remain open. | Local validation pending. |
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
