# Dune Admin Release Notes

## Current update: Frontend Quality workflow stabilization

### Why this update was made

GitHub Actions run `26325532482` failed before frontend dependency installation because `actions/setup-node` was still configured to cache npm dependencies using `web/package-lock.json`. That lockfile is intentionally not committed right now, so the cache dependency path could not be resolved and the job stopped during setup.

### Security and operator impact

- Removed the frontend workflow cache dependency on the missing `web/package-lock.json`.
- Kept frontend validation active: install, high-severity npm audit, TypeScript typecheck, lint, and build.
- Added `FORCE_JAVASCRIPT_ACTIONS_TO_NODE24=true` to the frontend workflow to opt into the newer JavaScript action runtime and avoid Node 20 deprecation drift on GitHub-hosted runners.
- Preserved the lockfile follow-up requirement: regenerate `web/package-lock.json` locally from the current manifest and recommit it once clean.

### Validation

The next Frontend Quality run should proceed past `actions/setup-node` and reach dependency install/typecheck/lint/build.

Local validation remains:

```bash
cd web
npm install
npm audit --audit-level=high
npm run typecheck
npm run lint
npm run build
```

---

## Previous update: Refactor implementation progress

### Why this update was made

The refactor review identified several immediate code improvements that could be completed without waiting for a larger architectural rewrite. This update implements the first set of improvements: isolate augmented item stat serialization, add augment presets, add a generated payload preview, remove stale UI drift, and add a manual item template refresh path.

### Security and operator impact

- Moved augment validation and `FAugmentedItemStats` JSON serialization into `item_augments.go` so high-risk item stat mutation logic is centralized.
- Added an augment preset catalog at `web/src/tabs/augmentPresets.ts` with default grades and roll counts for common T6 augments.
- Updated the augmented Give Item modal with preset buttons, explicit comma-separated roll arrays, and a generated payload preview.
- Added `POST /api/v1/players/templates/refresh` so operators can refresh cached live database templates without restarting the backend.
- Updated reconnect handling so item templates refresh after successful reconnect.
- Restored player handler endpoint coverage after the template refresh refactor so existing player mutation routes remain available.
- Removed the stale embedded legacy Give Item modal from `PlayersTab.tsx`; the active Players tab now has one augmented Give Item workflow.

### Current UI status

The active `PlayersTab.tsx` Give Item button opens `GiveItemModalAugmented.tsx`. Operators can select items, set stack count, item grade, stack size, add augment presets, enter custom augment names, set augment grade, set roll values, set explicit roll arrays, and preview the outgoing payload before submission.

### Item template source strategy

The implemented path is cached and operator-controlled:

- Backend startup loads and merges database-observed templates with `item-data.json`.
- Successful `/api/v1/reconnect` refreshes templates.
- `POST /api/v1/players/templates/refresh` refreshes templates on demand.
- The frontend still searches the returned cached list locally rather than querying the database on every keystroke.

This keeps the database load low while improving correctness over JSON-only item lists.

### Testing impact

Existing Go tests continue to cover:

- Augmented Give Item request normalization.
- Legacy single-item payload with augments.
- Invalid augment name, grade, roll, and roll-array validation.
- `FAugmentedItemStats` JSON generation.
- Empty stats behavior when no augments are supplied.

Additional recommended next tests:

- Handler coverage for the new template refresh endpoint.
- Frontend tests for augment preset selection and payload preview.
- Integration-style tests around augmented item insert SQL.

### Validation

```bash
go mod tidy
go test ./...
cd web
npm install
npm audit --audit-level=high
npm run typecheck
npm run lint
npm run build
```

### Known limitations and follow-up

- `PlayersTab.tsx` is still large and should be split into smaller player table, inventory, action, and modal components.
- The augment preset catalog is intentionally small and should grow as more verified augmented item examples are captured.
- Some augments appear to require multiple stat rolls. Operators should still use explicit `rolls` arrays when the preset catalog is incomplete.
- `web/package-lock.json` needs clean regeneration from the current frontend manifest.

---

## Previous update: Augmented Give Items, UI synchronization, and refactor review

### Why this update was made

The Give Item workflow needed support for augmented items, including augment names, augment grades, and normalized roll data. Operators also needed a clear decision on whether item templates should come from the live database, a JSON file, or both.

After the augmented modal was wired into the active Players tab, the repository still had stale legacy modal code and documentation that described a temporary rollback state. This update synchronized the code and documentation so there is one active augmented Give Item workflow.

### Security and operator impact

- Added backend support for augmented Give Item payloads on `POST /api/v1/players/give-item`.
- Added per-item augment definitions with augment name, augment grade, roll value, explicit roll arrays, roll count, and effect indices.
- Added validation for augment name presence, maximum augment count, grade bounds, roll bounds, roll count bounds, and explicit roll-array bounds before writing item stats.
- Added `FAugmentedItemStats` generation using the observed game-compatible `[[], payload]` wrapper shape.
- Preserved alignment between `AppliedAugments`, `AppliedAugmentRollData`, and `AppliedAugmentQualities`.
- Preserved legacy single-item non-augmented payload behavior for existing callers.
- Augmented grants create new stack rows rather than topping up existing stacks, avoiding accidental merging between augmented and plain items.
- Added frontend API types for augmented Give Item payloads.
- Wired the active Players tab Give Item button to `web/src/tabs/GiveItemModalAugmented.tsx`.
- Removed the stale embedded legacy Give Item modal from `PlayersTab.tsx` to prevent UI drift and duplicate behavior.

---

## Previous update: Linux version support and helper tests

### Why this update was made

The Windows-oriented workflow is now complemented by a Linux version so operators can install dependencies, run the app locally, build a Linux backend binary, and install the backend as a systemd service without translating Windows steps manually.

After the initial Linux scripts were added, automated tests were added so Linux helper behavior is validated continuously instead of relying only on manual review.

### Validation

```bash
go mod tidy
go test ./...
cd web
npm install
npm audit --audit-level=high
npm run build
```

---

## Release: Security Hardening, Security Scanning, Multi-Item Administration, and Documentation Standards Update

### Summary

This release makes Dune Admin safer and more reliable for live server operations. It establishes backend-enforced admin-token authentication, explicit CORS allowlisting, safer loopback listen defaults, server timeouts, raw SQL restrictions, request-size controls, Kubernetes log target validation, reduced status data exposure, hardened frontend security headers, blueprint import bounds checks, CI security scanning, removal of hardcoded capture credentials, and the requirement that every future change keep both `PATCH_NOTES.md` and `CHANGELOG.md` current.

### Security notes for operators

- Treat `ADMIN_TOKEN` as a privileged secret.
- Rotate any previously shared or committed credentials.
- Keep `.env`, SSH keys, database snapshots, generated secrets, dependency folders, and build output out of source control.
- Prefer `LISTEN_ADDR=127.0.0.1:8080` for local use.
- Do not expose the backend directly to the internet.
- Place remote access behind TLS, a trusted reverse proxy, and a strong identity provider.
