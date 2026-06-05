# Generated Route Auth Boundary Coverage

## Purpose

This note tracks the generated full-route auth-boundary coverage added after the initial AppSec endpoint audit.

The goal is to prevent route drift: any route registered in `routes.go` must have an explicit AppSec auth-boundary expectation in tests. If a new route is added without a matching expectation, the test fails and forces the route's intended exposure model to be reviewed.

## Coverage added

`appsec_route_inventory_test.go` parses `routes.go` with Go's AST parser and extracts every `mux.HandleFunc("METHOD /path", handler)` route pattern. It then compares the discovered registered route list to `appsecExpectedRouteAuth`.

The test fails when:

- a registered route has no AppSec auth-boundary expectation
- an expectation references a route no longer registered in `routes.go`
- a public route does not bypass auth as expected
- a self-service route is not denied without auth, allowed with admin token, and allowed with a normal registered Discord session
- an admin route is not denied without auth, allowed with admin token, and denied to a normal registered Discord session
- the WebSocket log-stream route does not require a one-time ticket for WebSocket upgrade requests before admin-token fallback

## Auth-boundary classes

| Class | Meaning |
|---|---|
| `public` | Explicit public route. Must bypass backend admin-token validation. |
| `self-service` | Requires admin token or registered Discord self-service/session access. |
| `admin` | Requires admin token or Discord admin session. Normal Discord sessions must fail. |
| `websocket-ticket` | Normal non-upgrade GET remains admin-protected; WebSocket upgrade requires a one-time scoped ticket. |

## Current expected routes

The expectation map covers all routes currently registered in `routes.go`, including:

- public health and Discord OAuth initiation/callback routes
- Discord self-session routes
- Discord player-link admin routes
- self-service player-card routes
- status, diagnostics, audit, and mutation-safety routes
- Battlegroup infrastructure routes
- player read and mutation routes
- inventory request/order coordination routes
- database routes
- log and WebSocket ticket routes
- notification route
- storage routes
- blueprint routes

## Security impact

- No route behavior changed.
- No new route was added.
- No mutation path was added.
- Player 360 remains read-only.
- This is regression coverage only: it makes future route additions fail closed at test time until auth expectations are explicitly reviewed.

## Remaining work

This closes the generated route inventory/auth-boundary coverage gap for currently registered routes after local validation. Remaining AppSec audit work includes:

- endpoint-by-endpoint audit-event assertion coverage
- mutation reason/pre-post review verification for every high-risk mutation
- handler-level input validation tests beyond currently covered slices
- SAST/DAST/dependency evidence
- manual abuse-case validation

## Validation

Required from the canonical local update path:

```bash
./update.sh
```
