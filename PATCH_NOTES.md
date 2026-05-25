# Dune Admin Release Notes

## Current update: Player 360 launcher from Players table

### Why this update was made

Player 360 compiled cleanly as a protected read-only backend and standalone frontend tab. This slice adds the operator shortcut that was planned after validation: a row-level `360` launcher from the existing Players table.

### What changed

- Added `web/src/tabs/PlayersTabWith360Launcher.tsx` as a wrapper around the existing Players tab.
- Added a read-only `360` button beside existing player row actions.
- Updated `web/src/App.tsx` so the Players tab uses the launcher wrapper.
- Updated `web/src/tabs/Player360Tab.tsx` so Player 360 reads the selected player ID and auto-loads the profile.
- Updated `docs/player-360-profile.md` with the launcher status.

### Security and operator impact

- Player 360 remains read-only.
- No existing Inventory, Give Item, or Actions flows were changed.
- No new player mutation paths were added.
- The launcher stores only the selected player actor ID in browser local storage.

### Validation

Expected validation after this launcher slice:

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

## Previous update: Player 360 read-only frontend tab

Added the standalone read-only Player 360 frontend tab and navigation entry.

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
