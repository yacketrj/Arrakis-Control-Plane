# Discord Auth

Discord auth is the identity foundation for future operator and safe self-service workflows. Player 360 remains read-only until Discord identity-to-player mapping exists and is explicitly enforced.

## Runtime configuration

Set these values when enabling Discord OAuth:

```env
DISCORD_AUTH_ENABLED=true
DISCORD_CLIENT_ID=<discord application client id>
DISCORD_CLIENT_SECRET=<discord application client secret>
DISCORD_REDIRECT_URI=<backend callback url>/api/v1/auth/discord/callback
DISCORD_GUILD_ID=<discord guild id>
DISCORD_ADMIN_ROLE_IDS=<comma-separated admin role ids>
DISCORD_NORMAL_ROLE_IDS=<comma-separated normal role ids>
DISCORD_USER_STORE=discord-users.json
DISCORD_POST_LOGIN_REDIRECT=/
SESSION_COOKIE_SECURE=true
```

`SESSION_COOKIE_SECURE=true` is required for non-local HTTPS deployments. The backend also sets secure cookies automatically when `DISCORD_REDIRECT_URI` starts with `https://`.

## Registered endpoints

| Method | Path | Purpose | Public path | Normal Discord session access |
| --- | --- | --- | --- | --- |
| `GET` | `/api/v1/auth/discord/login` | Start Discord OAuth and set the state cookie. | Yes | N/A |
| `GET` | `/api/v1/auth/discord/callback` | Complete Discord OAuth, map roles, upsert the registered-user store, and issue a session cookie. | Yes | N/A |
| `GET` | `/api/v1/auth/discord/me` | Return current auth context for an authenticated Discord session or admin-token request. | No | Yes |
| `POST` | `/api/v1/auth/discord/logout` | Clear the session cookie and remove the in-memory Discord session. | No | Yes |
| `GET` | `/api/v1/auth/discord/users` | Return the registered Discord user store for administrative review. | No | No |

Only login and callback are included in the public-path allowlist.

Registered non-admin Discord sessions are intentionally limited to:

- `GET /api/v1/auth/discord/me`
- `POST /api/v1/auth/discord/logout`
- `/api/v1/self/*`

All other protected routes require either a valid backend admin token or a Discord admin session.

## Role mapping

Role mapping is intentionally simple:

1. Any Discord role ID listed in `DISCORD_ADMIN_ROLE_IDS` maps to `admin`.
2. If no admin role matches and `DISCORD_NORMAL_ROLE_IDS` is configured, a matching role maps to `normal`.
3. If `DISCORD_NORMAL_ROLE_IDS` is empty, any guild member without an admin role maps to `normal`.
4. If no rule matches, the user maps to `none` and the callback rejects access.

Admin role mapping wins over normal role mapping.

## Session behavior

The backend stores Discord sessions in memory with a 12-hour TTL. The `dune_admin_session` cookie is HttpOnly, SameSite=Lax, path-scoped to `/`, and uses the secure flag when configured or when the Discord redirect URI is HTTPS.

Expired sessions are rejected and evicted on lookup. Logout removes the in-memory session and clears the browser cookie.

## Validation

Run backend tests from the local checkout or CI:

```bash
go test ./...
```

Manual OAuth validation should confirm:

1. `GET /api/v1/auth/discord/login?json=1` returns an authorize URL and sets a state cookie when Discord auth is configured.
2. `GET /api/v1/auth/discord/callback` rejects missing or invalid `code` / `state` values.
3. A valid callback for a guild member with an allowed role writes/updates the registered-user store and issues `dune_admin_session`.
4. `/api/v1/auth/discord/me` returns `auth_type: "discord"` and the expected role for a valid Discord session.
5. `/api/v1/auth/discord/me` returns `auth_type: "admin-token"` when reached through an authenticated admin-token request.
6. `/api/v1/auth/discord/logout` clears the cookie and invalidates the session.
7. `/api/v1/auth/discord/users` is not reachable without administrative authentication when served through the normal backend middleware.
8. A normal registered Discord session can reach only `me`, `logout`, and `/api/v1/self/*` routes.
9. A normal registered Discord session cannot reach admin review, player, database, infrastructure, or mutation routes.

## Current limitation

The session store is process-local memory. A backend restart invalidates active Discord sessions. Production deployments that need restart-stable sessions should add a durable, encrypted session store before relying on Discord sessions as the only operator identity mechanism.
