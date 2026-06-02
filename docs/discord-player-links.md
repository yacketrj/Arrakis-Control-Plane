# Discord Player Links

Discord player links are the identity-to-player mapping foundation for future player-safe self-service. They connect a Discord user ID to one DA Manager player actor ID.

This feature is intentionally read-only for self-service. It does not enable player inventory writes, guild storage writes, claim rewards, currency changes, XP changes, Player 360 mutation, or any other game-state mutation.

## Storage

The backend uses a local JSON store:

```env
DISCORD_PLAYER_LINK_STORE=discord-player-links.json
```

If unset, the backend writes `discord-player-links.json` in the working directory. The file is written with `0600` permissions. Access is serialized with an in-process mutex.

Current limitation: this is process-local file storage, not a multi-node database. A future production version should move the mapping to a durable table if DA Manager runs more than one backend instance.

## Link model

Each link contains:

- `discord_id`: stable Discord user ID.
- `player_id`: DA Manager player actor ID.
- `player_name`: optional display name for operator review.
- `notes`: optional operator notes.
- `linked_by_discord_id`: Discord ID of the operator who wrote the link when available.
- `linked_by_auth_type`: `discord`, `admin-token`, or `unknown`.
- `created_at` / `updated_at`: UTC timestamps.

Validation rules:

- `discord_id` is required.
- `player_id` must be greater than zero.
- Text fields are trimmed, length-limited, and reject unsupported control characters.

## Admin endpoints

| Method | Path | Purpose |
| --- | --- | --- |
| `GET` | `/api/v1/auth/discord/player-links` | List current Discord-to-player links. |
| `POST` | `/api/v1/auth/discord/player-links` | Create or update a Discord-to-player link. |
| `DELETE` | `/api/v1/auth/discord/player-links/{discord_id}` | Delete a link by Discord ID. |

These endpoints require normal DA Manager administrative authentication: admin token or Discord admin session.

## Self-service endpoints

| Method | Path | Purpose |
| --- | --- | --- |
| `GET` | `/api/v1/self/player-link` | Return the current Discord session user's player link. |
| `GET` | `/api/v1/self/player-card` | Return a read-only player card summary for the linked player. |

Normal Discord sessions are allowed only under `/api/v1/self/*`. They are not allowed to access admin endpoints unless they also have the configured Discord admin role.

## Self-service player card contents

The self-service player card is derived from the existing read-only Player 360 profile builder and returns a reduced summary:

- Discord ID
- player ID
- player name
- class
- map
- online status
- location
- inventory summary
- vehicle count
- currencies
- factions
- specializations
- character XP
- journey summary
- section errors

The response intentionally omits direct mutation affordances.

## Auth boundary

The auth middleware now permits normal registered Discord sessions only when the path starts with `/api/v1/self/`.

All other protected paths still require one of:

- valid admin token / Browser Access Key
- Discord admin session
- special-case WebSocket log stream ticket for the log stream path

## Validation

Run backend validation from the local checkout:

```bash
go test ./...
go build ./...
```

Manual validation should confirm:

1. Admin token or Discord admin can create/list/delete player links.
2. Normal Discord sessions cannot access `/api/v1/status` or admin paths.
3. Normal Discord sessions can access `/api/v1/self/player-link` only when linked.
4. Normal Discord sessions can access `/api/v1/self/player-card` only when linked.
5. Unlinked Discord sessions receive a safe not-found error for self player link/card endpoints.
6. Self player card returns only read-only support data and exposes no player mutation action.

## Safety boundary

This mapping is a prerequisite for future self-service, not a self-service mutation feature. Player 360 remains read-only. Any future self-service action must explicitly verify the Discord session, enforce the mapped player ID, classify mutation risk, capture an audit reason when needed, and use the mutation-safety workflow.
