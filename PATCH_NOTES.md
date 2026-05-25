# Dune Admin Release Notes

## Current update: Player 360 validated read-only profile

### Why this update was made

Player 360 is now implemented and compiled cleanly as the first P1 operator-support surface after the P0 audit and mutation-safety foundation. This update closes the Player 360 v1 slice as a protected read-only profile with a standalone tab and Players-table launcher.

### What changed

- Marked Player 360 v1 as validated in `docs/player-360-profile.md`.
- Marked Player 360 Profile as Done in `docs/admin-implementation-tasks.md`.
- Confirmed the backend profile endpoint, standalone frontend tab, auto-load behavior, and Players-table launcher have compiled cleanly.
- Preserved the rule that Player 360 remains read-only and does not add quick actions yet.

### Security and operator impact

- Player 360 remains a read-only support view.
- No existing Inventory, Give Item, or Actions flows were changed.
- No new player mutation paths were added.
- Future Player 360 quick actions remain blocked on shared mutation-safety confirmation.

### Validation

Validated clean:

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

## Previous update: Player 360 launcher from Players table

Added the read-only Players-table `360` launcher and Player 360 auto-load behavior.

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
