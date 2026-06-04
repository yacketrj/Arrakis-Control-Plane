# Browser Token and CORS Security Notes

## Purpose

This note tracks the AppSec review state for DA Manager browser-token handling and CORS behavior under `ASEA-006`.

The current admin UI supports an interim Browser Access Key flow. The token is no longer persisted in `localStorage`; valid legacy values are migrated to `sessionStorage` and removed from legacy storage. This is safer than persistent local storage, but still JavaScript-readable. The long-term target remains memory-only token handling or HttpOnly secure session-cookie authentication with an explicit CSRF design.

## Reviewed surface

| Area | Files | Current posture |
|---|---|---|
| Backend admin token validation | `auth.go`, `auth_test.go` | Strict 43-character base64url token validation, forbidden placeholder rejection, control-character rejection, and constant-time comparison. |
| CORS origin parsing | `auth.go`, `server.go`, `auth_test.go` | Allowed origins are exact-match only after validation. Wildcard, `null`, control-character, non-HTTP(S), userinfo, path, query, and fragment origins are rejected/ignored. |
| CORS headers | `server.go`, `auth_test.go` | Allowed origins are reflected exactly with `Vary: Origin`; disallowed origins are not reflected. |
| Browser token storage | `web/src/api/client.ts` | Uses `sessionStorage` for current Browser Access Key and removes valid/invalid legacy `localStorage` values. |
| Frontend token validation | `web/src/api/client.ts` | Enforces same 43-character base64url shape and rejects whitespace/control characters/placeholders before storage. |

## Current guardrails

- Backend `ADMIN_TOKEN` must be exactly 43 base64url characters generated from 32 random bytes.
- Backend rejects missing, malformed, placeholder, whitespace, and control-character tokens.
- Browser Access Key validation mirrors the strict backend token shape.
- Current frontend storage uses `sessionStorage`, not durable `localStorage`.
- Legacy `localStorage` values are migrated only when valid, then removed.
- Invalid legacy `localStorage` values are removed.
- CORS allowed origins are exact-match only.
- Misconfigured CORS values such as `*`, `null`, `file://`, `javascript:`, origins with userinfo, origins with path/query/fragment, and origins containing controls are rejected/ignored.
- Disallowed browser origins do not receive `Access-Control-Allow-Origin`.

## Added regression tests

`auth_test.go` now covers:

- unsafe allowed-origin value rejection
- safe origin value acceptance
- parser behavior when mixed safe and unsafe allowed origins are configured
- CORS preflight behavior for disallowed origins
- CORS preflight behavior for allowed origins and `Vary: Origin`

## Remaining ASEA-006 work

This is partial remediation. Before `ASEA-006` can be closed, finish:

- replace JavaScript-readable Browser Access Key storage with memory-only handling or HttpOnly secure session-cookie authentication
- define and test CSRF behavior before adding cookie-based admin auth
- decide whether `X-Admin-Token` should remain in browser CORS headers after the auth model changes
- add frontend unit tests for `adminTokenValidationError`, legacy-token migration, invalid legacy-token removal, and session storage behavior
- manually validate CORS behavior in browser against configured local and reverse-proxy origins
- verify reverse-proxy/TLS deployment guidance aligns with the final auth model

## Validation

Required from the canonical local update path:

```bash
./update.sh
```
