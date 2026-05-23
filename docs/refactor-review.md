# Dune Admin Refactor and Improvement Review

## Scope

This review covers the current Go backend, React/Vite frontend, GitHub Actions quality gates, Linux helper scripts, augmented Give Item workflow, and operational documentation.

## Executive summary

Dune Admin now has backend authentication, security scanning, Linux support, augmented item grants, and release documentation. The next major improvement area is maintainability: split large UI files, isolate item/augment serialization logic, separate HTTP handlers from database persistence, add a cached item-template provider, and expand automated frontend/backend test coverage.

## Synchronization completed

- `PlayersTab.tsx` now opens `GiveItemModalAugmented.tsx` from the active Give Item button.
- The older embedded Give Item modal was removed from `PlayersTab.tsx` to eliminate duplicate UI behavior.
- The active UI, backend payload model, API client types, docs, patch notes, and changelog now describe the same augmented Give Item workflow.
- Augmented grants support item grade, augment grade, roll value, explicit roll arrays, roll count, and effect indices.

## Highest-priority refactors

### 1. Split `PlayersTab.tsx`

`PlayersTab.tsx` still owns player listing, inventory, vehicles, player actions, XP, journey, teleport, history, and modal orchestration. This should be split into smaller files.

Recommended structure:

```text
web/src/tabs/players/PlayersTab.tsx
web/src/tabs/players/PlayersTable.tsx
web/src/tabs/players/InventoryModal.tsx
web/src/tabs/players/PlayerActionsModal.tsx
web/src/tabs/players/GiveItemModalAugmented.tsx
web/src/tabs/players/hooks/usePlayers.ts
web/src/tabs/players/hooks/usePlayerActions.ts
web/src/tabs/players/types.ts
```

### 2. Move augment serialization into a backend domain file

`FAugmentedItemStats` generation should live outside command/handler code.

Recommended files:

```text
item_augments.go
item_augments_test.go
give_item_request.go
give_item_request_test.go
```

This keeps validation, serialization, and game-compatible JSON shape in one focused place.

### 3. Add an augment metadata catalog

Some augments require multiple roll values. The UI currently allows manual roll count and explicit roll arrays, but a catalog would provide safer defaults.

Suggested shape:

```json
{
  "T6_Augment_Magazinecapacity1": {
    "display_name": "Magazine Capacity I",
    "default_grade": 5,
    "roll_count": 3,
    "default_rolls": [1.0, 1.0, 1.0]
  }
}
```

### 4. Implement a cached hybrid item-template provider

Use live database discovery plus curated `item-data.json` metadata. Do not query the database on every UI search keystroke.

Refresh on:

- Backend startup after DB connection.
- `/api/v1/reconnect`.
- Manual operator refresh endpoint.
- Optional 15-60 minute background refresh.

### 5. Separate handlers, services, and repositories

Recommended backend shape:

```text
handlers_players.go      HTTP request/response only
players_service.go       business rules and validation
players_repository.go    SQL and persistence only
```

This will make unit tests easier and reduce regression risk.

## Backend improvements

- Replace package-level `globalDB` usage with dependency injection through an application struct.
- Pass `r.Context()` into database operations instead of using `context.Background()`.
- Add query timeouts for long-running operational queries.
- Keep augmented grants from merging into plain or differently augmented stacks.
- Add admin action audit logs for item grants, currency grants, XP changes, teleport, kick, journey wipe, and codex wipe.
- Avoid logging admin tokens, SSH keys, DB passwords, or full environment values.

## Frontend improvements

- Move each modal out of `PlayersTab.tsx`.
- Add shared clamped numeric input components.
- Split `web/src/api/client.ts` by domain.
- Add row-level validation in the augmented Give Item modal before submit.
- Add a collapsible generated-payload preview for advanced operators.
- Warn operators when granting items to online players.

## Testing improvements

### Backend

Add tests for:

- `POST /api/v1/players/give-item` handler behavior.
- Generated SQL insert behavior for augmented items.
- JSON round trips using captured augmented item examples.
- Boundary values for row count, stack size, augment count, and multi-roll augments.

### Frontend

Add Vitest and React Testing Library coverage for:

- Give Item modal row creation/removal.
- Augment creation/removal.
- Payload generation for single-roll and multi-roll augments.
- Submit disabled/enabled behavior.
- API error rendering.

### CI

- Add a frontend test workflow after Vitest is added.
- Keep Go tests, Linux helper tests, SCA, SAST, DCA, DAST, and secret scanning.
- Ensure frontend build runs on every frontend change.

## Security improvements

- Keep backend token authorization as the primary control for admin endpoints.
- Add mutation rate limits for high-impact operations.
- Add action audit logging with sanitized payload summaries.
- Continue treating SSH host-key pinning and RabbitMQ TLS verification as operator-trust-model follow-ups.
- Keep DAST enabled rather than suppressing browser security findings.

## Recommended implementation order

1. Split `PlayersTab.tsx` into smaller player components.
2. Move augmented stats serialization into focused backend model files.
3. Add augment metadata/presets and UI defaults.
4. Add cached database-plus-JSON template provider and manual refresh endpoint.
5. Add frontend unit tests.
6. Add admin action audit logging.
7. Add service/repository boundaries and reduce package-level global state.
8. Add integration-style tests for inventory mutation SQL.

## Risks and cautions

- Augment roll structure is inferred from observed examples; keep explicit roll arrays editable until enough examples exist for a reliable catalog.
- Avoid granting augmented items to online players until inventory sync behavior is better characterized.
- Database template discovery only finds templates that already appeared in the live database, so it should not replace curated fallback metadata.
