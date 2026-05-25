# Dune Admin Release Notes

## Current update: Shared frontend mutation confirmation foundation

### Why this update was made

The next DA Manager slice prepares the frontend safety layer needed before adding Player 360 quick actions or expanding high-risk operator workflows. Player 360 remains read-only; this update adds the reusable confirmation surface that future mutating flows must use.

### What changed

- Added `web/src/hooks/useMutationConfirmation.tsx` as a shared frontend confirmation hook.
- The hook classifies the target mutation through `/api/v1/mutation-safety/classify` before showing confirmation.
- Added conservative local fallback classification so the UI still prompts for high-risk or destructive-looking paths if the backend classification request fails.
- Added support for displaying risk, action name, operator warnings, recommended path, rollback hint, target context, extra details, and admin reason capture.
- Updated `docs/mutation-safety-framework.md` with the frontend integration pattern and current limitations.
- Updated `docs/admin-implementation-tasks.md` so the active focus is migrating existing high-risk workflows to the shared confirmation hook.

### Security and operator impact

- Player 360 remains a read-only support view.
- No new player mutation paths were added.
- Existing Players, Inventory, Give Item, and Actions workflows were not behaviorally changed in this slice.
- Future Player 360 quick actions remain blocked until relevant mutation flows are wired through the shared confirmation hook and reason capture.

### Validation

Validation still required in the target development environment:

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

## Previous update: Player 360 validated read-only profile

Player 360 v1 was validated as a protected read-only profile with a standalone tab and Players-table launcher.

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
