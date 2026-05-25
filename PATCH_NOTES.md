# Dune Admin Release Notes

## Current update: Player 360 read-only frontend tab

### Why this update was made

Player 360 now has a protected read-only backend profile endpoint, so this slice adds the first frontend surface for operators. The goal is to let operators load one consolidated support view without changing the existing Players, Inventory, Give Item, or Actions workflows and without introducing any new mutation paths.

### What changed

- Added `web/src/api/playerProfile.ts` with Player 360 response types and a protected profile fetch helper.
- Added `web/src/tabs/Player360Tab.tsx` as a standalone read-only Player 360 page.
- Added `Player 360` to the main app navigation in `web/src/App.tsx`.
- Removed the unused Player 360 modal prototype so the supported frontend path is the standalone tab.
- Updated `docs/player-360-profile.md` and `docs/admin-implementation-tasks.md` with the frontend status and validation focus.

### Security and operator impact

- Player 360 remains read-only.
- No new player mutation paths were added.
- Existing player workflows remain unchanged until the standalone tab validates cleanly.
- The frontend surfaces section-level partial-data warnings without exposing raw backend internals.

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

## Previous update: Player 360 backend profile foundation

Added the protected read-only Player 360 backend profile endpoint, route registration, helper tests, and backend documentation.

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
