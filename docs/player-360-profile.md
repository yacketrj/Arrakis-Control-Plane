# Player 360 Profile

## Purpose

Player 360 Profile is the next P1 operator-support feature for Dune Admin. It provides one protected read-only support view that brings together the player identity, online state, current runtime context, inventory summary, vehicles, currencies, faction status, specializations, journey state, recent events, and dungeon history.

The first implementation must be read-only. Mutating quick actions should be added only after the read-only view is stable and the shared mutation-safety confirmation pattern is ready.

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

Recommended endpoint:

```text
GET /api/v1/players/{id}/profile
```

The endpoint must be protected by the existing admin token middleware.

The endpoint should aggregate existing read-only player data and return partial results where possible. If one section fails, the response should include a section-level error instead of hiding all other available player context.

## Response shape

Recommended top-level shape:

```json
{
  "player_id": "string-or-number",
  "identity": {},
  "online_state": {},
  "location": {},
  "inventory_summary": {},
  "vehicles": [],
  "currencies": [],
  "factions": [],
  "specializations": [],
  "journey_summary": {},
  "recent_events": [],
  "dungeon_history": [],
  "section_errors": []
}
```

The exact field names may follow existing backend naming conventions, but the frontend should receive a stable structure that can render each section independently.

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
   - top stack or notable item information where available
   - link or affordance to the existing full inventory view

3. **Vehicles**
   - owned or associated vehicles
   - vehicle IDs and basic state where available

4. **Faction Status**
   - faction identity
   - reputation/tier values where available

5. **Specializations**
   - specialization rows and XP where available

6. **Journey Summary**
   - current journey indicators
   - completion/reset context where available

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
- Prefer section-level safe errors such as `inventory unavailable` over raw backend failures.

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

- Add a player profile response model.
- Add `GET /api/v1/players/{id}/profile`.
- Reuse existing read-only helpers for player details, online state, currency, inventory, vehicles, factions, specs, journey, events, and dungeons.
- Add section-level error handling.
- Add Go tests for successful partial response behavior.

### Slice 2: frontend Player 360 page

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

## Follow-up tasks

- Add shared frontend mutation confirmation component before Player 360 quick actions.
- Add safe quick actions only after read-only profile validation.
- Add before-change snapshots for any later quick action.
- Add links from Player 360 to Inventory Studio v2 once implemented.
- Add support for operator notes only after named operator identity exists.
