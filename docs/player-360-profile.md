# Player 360 Profile

## Purpose

Player 360 Profile is a protected read-only support view for Dune Admin. It brings together player identity, online state, current map context, inventory summary, vehicles, currencies, faction status, specializations, character XP, journey state, recent events, and dungeon history.

The first version must remain read-only. Mutating quick actions belong to later work after validation and after shared mutation-safety confirmation is ready.

## Current implementation status

Player 360 is in progress on `main` and the backend plus standalone frontend tab have compiled cleanly.

Implemented so far:

- `player_profile.go` backend response model, handler, aggregation helpers, and safe section-level errors.
- `GET /api/v1/players/{id}/profile` registered through `routes.go`.
- `server.go` refactored to call shared route registration.
- `player_profile_test.go` helper coverage for summaries, ID matching, online-state matching, and safe error wording.
- `web/src/api/playerProfile.ts` frontend profile fetch helper and response types.
- `web/src/tabs/Player360Tab.tsx` standalone read-only Player 360 tab.
- `web/src/App.tsx` navigation entry for `Player 360`.
- `web/src/tabs/PlayersTabWith360Launcher.tsx` wrapper that adds a read-only `360` launcher beside existing player row actions.
- Player 360 auto-loads the selected player ID when opened from the Players table.
- Existing Players, Inventory, Give Item, and Actions flows remain unchanged.

## Scope for v1

Player 360 v1 includes:

- protected read-only backend overview endpoint
- standalone frontend Player 360 tab
- Player Info section that folds in Currency and Online Status context
- read-only inventory, journey, vehicle, faction, specialization, event, and dungeon summaries
- partial-data behavior through `section_errors`
- player-row launch path from the existing Players table
- no new player mutation paths

## Out of scope for v1

Do not include these until later slices:

- inventory editing
- item deletion or repair from Player 360
- item grants
- currency grants
- teleport or rescue mutations
- journey reset or completion actions
- faction mutation actions
- arbitrary SQL or stored procedure execution

## Protected backend route

```text
GET /api/v1/players/{id}/profile
```

The route is protected by the existing admin middleware because it is registered under the protected API mux and is not listed as a public route.

## Current response shape

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

Each frontend section should render independently. If one section fails, the profile should still display other available sections.

## Frontend behavior

The current frontend surface is a standalone `Player 360` tab. Operators can either enter a PlayerCharacter actor ID directly or use the `360` button beside a row in the existing Players table.

The tab currently renders:

- Player Info
- Online and Location
- Character Progression
- Currencies
- Faction Status
- Inventory Summary
- Inventory Preview
- Journey Summary
- Vehicles
- Recent Events
- Dungeon History
- Partial-data warnings when a section is unavailable

Existing Players, Inventory, Give Item, and Actions workflows remain unchanged.

## Security controls

- Keep Player 360 protected by admin authorization.
- Do not expose Player 360 through public portal routes.
- Do not add new mutation paths in v1.
- Do not include admin tokens, database credentials, SSH keys, or raw environment values in responses.
- Do not expose raw SQL, connection strings, or internal infrastructure details in section errors.
- Keep later quick actions dependent on audit and mutation-safety controls.

## Audit requirements

The read-only Player 360 profile endpoint does not create mutation audit events because it does not change state.

Later Player 360 quick actions must use the Mutation Safety Framework, require preview where appropriate, capture operator reason where required, create audit events, and include rollback guidance where practical.

## Validation steps

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

## Next work

1. Re-run validation after the launcher wrapper change.
2. Clean up any frontend typecheck/lint/build findings.
3. Improve Player 360 section display and add links to existing full inventory/actions views after validation.
4. Add shared frontend mutation confirmation before any Player 360 quick actions.
5. Link Player 360 to Inventory Studio v2 after Inventory Studio exists.
