# Player 360 Profile

## Purpose

Player 360 Profile is the next P1 operator-support feature for Dune Admin. It provides one protected read-only support view that brings together the player identity, online state, current runtime context, inventory summary, vehicles, currencies, faction status, specializations, journey state, recent events, and dungeon history.

The first implementation must be read-only. Mutating quick actions should be added only after the read-only view is stable and the shared mutation-safety confirmation pattern is ready.

## Current implementation status

Slice 1 backend foundation has started:

- Added `player_profile.go` with the Player 360 response model, read-only aggregation handler, section-level error handling, and helper functions.
- Registered `GET /api/v1/players/{id}/profile` through shared route registration in `routes.go`.
- Refactored `server.go` to call shared route registration instead of carrying the full route list inline.
- Added `player_profile_test.go` coverage for inventory summary, journey summary, online-state matching, ID helper behavior, and safe section-error wording.

Remaining Slice 1 work before frontend implementation:

- Run Go validation in a connected/local development environment.
- Confirm the profile endpoint compiles against all current files.
- Add any missing endpoint integration tests if needed after local validation.

## User/operator problem

Current player support workflows require operators to move between separate player, online-status, currency, inventory, vehicle, faction, spec, journey, event, and dungeon views. This makes support slower and increases the chance that an operator acts without seeing the full player context.

Player 360 solves this by making the player detail page the starting point for support triage.

## Scope for v1

Player 360 v1 should:

- provide a protected read-only backend overview endpoint
- add a frontend Player 360 detail view
- fold Currency and Online Status into Player Info
- reuse existing player APIs and database helpers where practical
- avoid adding new mutating operations
- surface missing/partial data clearly instead of failing the entire page
- keep audit and mutation-safety foundations available for later quick actions

## Out of scope for v1

Do not include these in the first read-only slice:

- inventory editing
- item deletion
- item repair
- item grant quick actions
- currency grants
- teleport/rescue mutations
- journey reset/complete actions
- faction mutation actions
- arbitrary SQL or stored procedure execution

These belong to later slices after the read-only support view is stable.

## Protected backend route

Implemented endpoint:

```text
GET /api/v1/players/{id}/profile
```

The endpoint is protected by the existing admin token middleware because it is registered under the protected API mux and is not a public route.

The endpoint aggregates existing read-only player data and returns partial results where possible. If one section fails, the response includes a section-level error instead of hiding all other available player context.

## Response shape

Current top-level shape:

```json
{
  "player_id": 123,
  "identity": {},
  "online_state": {},
  "location": {},
  "inventory_summary": {},
  "vehicles": [],
  "currencies": [],
  "factions": [],
  "specializations": [],
  "character_xp": {},
  "journey_summary": {},
  "recent_events": [],
  "dungeon_history": [],
  "section_errors": []
}
```

Each section is intended to render independently. A failed section should add a safe entry to `section_errors` rather than exposing raw database connection strings, SQL text, or internal details.

## Frontend behavior

The Player 360 page should include these sections:

1. **Player Info**
   - identity
   - player/account/controller IDs
   - online/offline state
   - current map, partition, or last-known location where available
   - currency summary

2. **Inventory Summary**
   - total item count
   - total stack count
   - unique template count
   - limited preview list
   - link or affordance to the existing full inventory view

3. **Vehicles**
   - owned or associated vehicles
   - vehicle IDs and basic state where available

4. **Faction Status**
   - faction identity
   - reputation/tier values where available

5. **Specializations and Character XP**
   - specialization rows and XP where available
   - character XP and level where available

6. **Journey Summary**
   - total journey nodes
   - completed nodes
   - revealed nodes
   - pending rewards
   - limited preview list

7. **Recent Events**
   - recent player events already exposed by existing APIs

8. **Dungeon History**
   - recent dungeon or activity history already exposed by existing APIs

## Security controls

- Keep the route protected by admin authorization.
- Do not expose Player 360 through public portal routes.
- Do not add new mutation paths in v1.
- Do not include admin tokens, database credentials, SSH keys, or raw environment values in responses.
- Do not expose raw SQL query text or internal connection strings in frontend section errors.
- Prefer section-level safe errors such as `inventory unavailable` or `section unavailable` over raw backend failures.

## Audit requirements

The read-only Player 360 profile endpoint does not need to create mutation audit events because it does not change state.

Later Player 360 quick actions must:

- use the Mutation Safety Framework
- require preview for high-risk actions
- capture operator reason where required
- create audit events
- include rollback guidance where practical

## Validation steps

Expected validation after implementation:

```bash
go test ./...
cd web
npm install
npm audit --audit-level=high
npm run typecheck
npm run lint
npm run build
```

## Implementation plan

### Slice 1: read-only backend profile endpoint

Status: in progress.

- Add a player profile response model. Done.
- Add `GET /api/v1/players/{id}/profile`. Done.
- Reuse existing read-only helpers for player details, online state, currency, inventory, vehicles, factions, specs, journey, events, and dungeons. Done.
- Add section-level error handling. Done.
- Add Go tests for summary/helper behavior. Done.
- Run full local Go validation. Pending.

### Slice 2: frontend Player 360 page

Status: next after backend validation.

- Add frontend API client support for the profile endpoint.
- Add a Player 360 detail page or panel from the existing Players view.
- Fold Currency and Online Status into Player Info.
- Keep existing pages available until the new page is validated.
- Show section loading/error states clearly.

### Slice 3: cleanup and navigation consolidation

- Remove duplicate navigation only after operators can access the same information through Player 360.
- Update screenshots or docs if a UI guide is added later.
- Confirm support workflows remain clear for new and veteran operators.

## Known limitations

- Some sections depend on current database visibility and may be unavailable if the backend is disconnected.
- Online state may differ from database state during travel, login/logout, or partition transitions.
- Inventory writes may not be reflected instantly for online players; v1 is read-only and should not imply mutation safety.
- This view does not replace Inventory Studio v2.
- Backend endpoint validation still needs to be run in a local development environment with Go available.

## Follow-up tasks

- Add frontend Player 360 page and API client support.
- Add shared frontend mutation confirmation component before Player 360 quick actions.
- Add safe quick actions only after read-only profile validation.
- Add before-change snapshots for any later quick action.
- Add links from Player 360 to Inventory Studio v2 once implemented.
- Add support for operator notes only after named operator identity exists.
