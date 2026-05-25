# Dune Admin Release Notes

## Current update: Player 360 backend profile foundation

### Why this update was made

Player 360 is moving from roadmap planning into implementation. This slice adds the protected read-only backend foundation so operators can eventually review one consolidated player support profile before any new high-risk quick actions are introduced.

### What changed

- Added `player_profile.go` with the Player 360 response model, aggregation handler, section-level safe error handling, and helper functions.
- Added `GET /api/v1/players/{id}/profile` through shared route registration.
- Added `routes.go` to centralize HTTP route registration.
- Refactored `server.go` to call shared route registration.
- Added `player_profile_test.go` coverage for summary helpers, online-state matching, ID helper behavior, and safe error wording.
- Updated `docs/player-360-profile.md` and `docs/admin-implementation-tasks.md` with backend status.

### Security and operator impact

- Player 360 remains read-only in this slice.
- No new mutation paths were added.
- The route remains protected by the existing admin middleware.
- Section errors use safe wording instead of exposing raw backend details.

### Validation

Expected validation:

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

---

## Previous update: Player 360 roadmap and design foundation

The roadmap and design documents were updated so Player 360 is the next P1 read-only implementation slice after the P0 safety foundation.

---

## Previous update: Admin audit and mutation-safety documentation sync

Documentation and tracking were synced with the landed audit foundation and in-progress mutation-safety foundation.

---

## Previous update: SSH tunnel management foundation

Managed SSH tunnel behavior was added for protected game-management database access.

---

## Previous update: Battlegroup Health Diagnostics

Read-only Battlegroup Health Diagnostics and support-bundle export were added for operator troubleshooting.
