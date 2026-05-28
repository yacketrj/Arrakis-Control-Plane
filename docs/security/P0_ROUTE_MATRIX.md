# P0 Route Risk Matrix

This matrix classifies the current admin API surface before new orchestration, world/Sietch configuration, metrics, announcements, or player mutation features are expanded.

## Enforcement model

All non-public API routes are protected by backend admin-token authentication. Mutating routes are audited. High-risk and destructive mutations require an admin reason when `ADMIN_REQUIRE_REASON=true`.

Risk levels:

- `low`: read-only or status-only endpoint.
- `medium`: state-changing endpoint with limited blast radius or operational reconnect behavior.
- `high`: player/game-state mutation, live operation, or runtime command.
- `destructive`: deletion, wipe, import, restore, or bulk irreversible mutation.

Required controls for all future high-risk or destructive routes:

1. Typed request model.
2. Input allowlist and target validation.
3. Admin reason.
4. Audit event.
5. Redaction of secrets in errors/logs.
6. Preview or dry-run where practical.
7. Backup/snapshot before destructive writes where practical.

## Public routes

| Method | Route | Risk | Auth | Mutation | Notes |
|---|---|---:|---:|---:|---|
| GET | `/api/v1/public/status` | low | no | no | Public liveness endpoint only. Must not expose host, DB, SSH, runtime, pod, account, or token detail. |

## Core/admin routes

| Method | Route | Risk | Auth | Mutation | Reason | Notes |
|---|---|---:|---:|---:|---:|---|
| GET | `/api/v1/status` | low | yes | no | no | Safe status payload only; do not expose secrets. |
| POST | `/api/v1/reconnect` | medium | yes | yes | optional | Reopens SSH/DB/tunnels and refreshes templates. Audit required. |
| GET | `/api/v1/connectivity/diagnostics` | low | yes | no | no | Diagnostics output must remain redacted. |
| GET | `/api/v1/audit/events` | low | yes | no | no | Returns recent audit records. |
| GET | `/api/v1/mutation-safety/classify` | low | yes | no | no | Classification helper. |
| POST | `/api/v1/notify` | medium | yes | yes | optional | Notification path; keep payload bounded and redacted. |

## Battlegroup/runtime routes

| Method | Route | Risk | Auth | Mutation | Reason | Notes |
|---|---|---:|---:|---:|---:|---|
| GET | `/api/v1/battlegroup/status` | low | yes | no | no | Runtime status only. |
| GET | `/api/v1/battlegroup/health` | low | yes | no | no | Runtime health diagnostics. |
| GET | `/api/v1/battlegroup/pods` | low | yes | no | no | Container/pod discovery. Validate target names. |
| POST | `/api/v1/battlegroup/exec` | high | yes | yes | yes | Runtime command allowlist only. No arbitrary shell. Docker support must remain disabled until provider hardening is complete. |

## Player read routes

| Method | Route | Risk | Auth | Mutation | Reason | Notes |
|---|---|---:|---:|---:|---:|---|
| GET | `/api/v1/players` | low | yes | no | no | Player list. |
| GET | `/api/v1/players/online` | low | yes | no | no | Online state. |
| GET | `/api/v1/players/currency` | low | yes | no | no | Currency view. |
| GET | `/api/v1/players/factions` | low | yes | no | no | Faction view. |
| GET | `/api/v1/players/specs` | low | yes | no | no | Specs view. |
| GET | `/api/v1/players/templates` | low | yes | no | no | Template catalog. |
| GET | `/api/v1/players/{id}/profile` | low | yes | no | no | Player profile. |
| GET | `/api/v1/players/{id}/inventory` | low | yes | no | no | Inventory read. |
| GET | `/api/v1/players/{id}/journey` | low | yes | no | no | Journey read. |
| GET | `/api/v1/players/{id}/char-xp` | low | yes | no | no | XP read. |
| GET | `/api/v1/players/{id}/specs` | low | yes | no | no | Player specs read. |
| GET | `/api/v1/players/{id}/vehicles` | low | yes | no | no | Vehicle read. |
| GET | `/api/v1/players/partitions` | low | yes | no | no | Partition read. |
| GET | `/api/v1/players/{id}/events` | low | yes | no | no | Event read. |
| GET | `/api/v1/players/{id}/dungeons` | low | yes | no | no | Dungeon read. |

## Player mutation routes

| Method | Route | Risk | Auth | Mutation | Reason | Notes |
|---|---|---:|---:|---:|---:|---|
| POST | `/api/v1/players/templates/refresh` | medium | yes | yes | optional | Refreshes server-derived template catalog. |
| POST | `/api/v1/players/give-item` | high | yes | yes | yes | Direct inventory write. May require relog. Validate item catalog and target. |
| POST | `/api/v1/players/give-currency` | high | yes | yes | yes | Player economy mutation. |
| POST | `/api/v1/players/grant-live` | high | yes | yes | yes | Claim/reward path only for validated plain grants. |
| POST | `/api/v1/players/give-faction-rep` | high | yes | yes | yes | Faction mutation. |
| POST | `/api/v1/players/give-scrip` | high | yes | yes | yes | Currency mutation. |
| POST | `/api/v1/players/award-xp` | high | yes | yes | yes | Progression mutation. |
| POST | `/api/v1/players/award-char-xp` | high | yes | yes | yes | Character progression mutation. |
| POST | `/api/v1/players/award-intel` | high | yes | yes | yes | Progression/resource mutation. |
| POST | `/api/v1/players/kick` | high | yes | yes | yes | Live player action. |
| DELETE | `/api/v1/players/item/{id}` | destructive | yes | yes | yes | Capture item row before deletion. |
| POST | `/api/v1/players/reset-spec` | destructive | yes | yes | yes | Progression reset. Snapshot affected rows. |
| POST | `/api/v1/players/set-faction-tier` | high | yes | yes | yes | Faction progression mutation. |
| POST | `/api/v1/players/journey/complete` | high | yes | yes | yes | Journey mutation. |
| POST | `/api/v1/players/journey/reset` | destructive | yes | yes | yes | Journey reset. Snapshot first. |
| POST | `/api/v1/players/journey/wipe` | destructive | yes | yes | yes | Journey wipe. Snapshot first. |
| POST | `/api/v1/players/delete-tutorials` | destructive | yes | yes | yes | Tutorial deletion. Snapshot first. |
| POST | `/api/v1/players/wipe-codex` | destructive | yes | yes | yes | Codex wipe. Snapshot first. |
| POST | `/api/v1/players/set-spec-xp` | high | yes | yes | yes | Progression mutation. |
| POST | `/api/v1/players/repair-item` | high | yes | yes | yes | Inventory mutation. |
| POST | `/api/v1/players/teleport` | high | yes | yes | yes | Record prior partition/location before teleport. |

## Database routes

| Method | Route | Risk | Auth | Mutation | Reason | Notes |
|---|---|---:|---:|---:|---:|---|
| GET | `/api/v1/database/tables` | low | yes | no | no | Schema/table listing. |
| GET | `/api/v1/database/describe` | low | yes | no | no | Table description. Validate identifiers. |
| GET | `/api/v1/database/sample` | low | yes | no | no | Bounded sample. |
| GET | `/api/v1/database/search` | low | yes | no | no | Read-only search. |
| GET | `/api/v1/database/functions` | low | yes | no | no | Function listing. |
| GET | `/api/v1/database/functions/inspect` | low | yes | no | no | Function inspection. Redact sensitive function text if needed. |
| POST | `/api/v1/database/sql` | medium | yes | controlled read-only | optional | Single-statement read-only SQL only. No mutation SQL. |

## Logs routes

| Method | Route | Risk | Auth | Mutation | Reason | Notes |
|---|---|---:|---:|---:|---:|---|
| GET | `/api/v1/logs/pods` | low | yes | no | no | Log target discovery. Validate names/IDs. |
| POST | `/api/v1/logs/stream-ticket` | medium | yes | yes | optional | Short-lived ticket issuance. Audit required. |
| GET | `/api/v1/logs/stream` | medium | ticket | no | no | WebSocket stream. Origin and one-time ticket required. Redact streamed lines. |
| GET | `/api/v1/logs/cheats` | low | yes | no | no | Cheat log read. Redact errors. |

## Storage routes

| Method | Route | Risk | Auth | Mutation | Reason | Notes |
|---|---|---:|---:|---:|---:|---|
| GET | `/api/v1/storage` | low | yes | no | no | Storage list. |
| GET | `/api/v1/storage/{id}/items` | low | yes | no | no | Storage inventory read. |
| POST | `/api/v1/storage/{id}/give-item` | high | yes | yes | yes | Storage mutation may require zone/server restart to appear. |

## Blueprint routes

| Method | Route | Risk | Auth | Mutation | Reason | Notes |
|---|---|---:|---:|---:|---:|---|
| GET | `/api/v1/blueprints` | low | yes | no | no | Blueprint list. |
| GET | `/api/v1/blueprints/{id}/export` | low | yes | no | no | Blueprint export. Avoid leaking secrets in export metadata. |
| POST | `/api/v1/blueprints/import` | destructive | yes | yes | yes | Import can overwrite/create game data. Require reason and backup when implemented. |
