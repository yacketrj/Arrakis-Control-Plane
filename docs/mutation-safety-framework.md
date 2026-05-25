# Mutation Safety Framework v1

## Purpose

Mutation Safety Framework v1 adds shared safety metadata around high-impact administrator operations. It builds on the Admin Action Audit Log by classifying mutating requests, capturing operator reason text when supplied, and recording a small allowlist of target identifiers.

This is a P0 foundation layer for features such as Player 360 quick actions, Inventory Studio v2, guild/faction administration, journey management, safe teleport/rescue, and allowlisted routine execution.

## What v1 provides

- Risk classification for protected mutating requests.
- Protected classification endpoint for frontend preview workflows.
- Shared frontend mutation confirmation hook at `web/src/hooks/useMutationConfirmation.tsx`.
- Frontend confirmation support for risk display, warnings, recommended path, rollback hint, target summary, operator details, and admin reason capture.
- Confirmed extracted player workflow modals for the active Players-table mutation surface:
  - `web/src/tabs/GiveItemModalAugmented.tsx`
  - `web/src/tabs/InventoryModal.tsx`
  - `web/src/tabs/PlayerActionsModalConfirmed.tsx`
  - `web/src/tabs/PlayerTeleportModal.tsx`
  - `web/src/tabs/PlayerAdminActionsModal.tsx`
- Confirmed non-player mutation surfaces for current high-risk admin workflows:
  - Storage add/remove item flows in `web/src/tabs/StorageTab.tsx`
  - Database SQL execution in `web/src/tabs/DatabaseTab.tsx`
  - Battlegroup server control in `web/src/tabs/BattlegroupTab.tsx`
  - Blueprint import in `web/src/tabs/BlueprintsTab.tsx`
- Optional `reason` capture from JSON request bodies.
- Optional `X-Admin-Reason` header capture for admin workflows.
- Environment-controlled reason enforcement for high-risk and destructive actions.
- Allowlisted target ID capture.
- Request-body restoration after audit inspection so handlers can still decode the request normally.
- Audit event fields for risk, reason, target metadata, preview requirement, destructive status, rollback hint, operator warnings, and recommended path.
- Unit tests for classification, handler behavior, reason enforcement, body restoration, target capture, and oversized-body handling.

## Protected classification endpoint

```text
GET /api/v1/mutation-safety/classify?method=POST&path=/api/v1/players/give-item
```

The response includes the normalized action name, risk level, reason requirement, preview requirement, destructive flag, and any available guidance fields.

## Shared frontend confirmation hook

`useMutationConfirmation` is the shared frontend integration point for high-impact operator actions.

Expected use pattern:

```tsx
const { confirmMutation, confirmationDialog } = useMutationConfirmation()

const reason = await confirmMutation({
  method: 'POST',
  path: '/players/give-currency',
  title: 'Give currency',
  summary: 'Give 100 Solari to the selected player.',
  target: 'actor:12345',
  details: ['Player must be the intended support target.'],
  forceReason: true,
})

await api.players.giveCurrency(controllerId, 100, reason)

return <>{confirmationDialog}</>
```

The hook calls the protected classification endpoint before displaying the dialog. If classification is unavailable, it falls back to conservative local classification so operators still see a confirmation and reason prompt for high-risk or destructive-looking paths.

## Active Players-table coverage

The active Players-table wrapper, `web/src/tabs/PlayersTabWith360Launcher.tsx`, routes high-risk player workflows through extracted confirmed modals instead of expanding the legacy `PlayersTab.tsx` file.

Currently covered by shared confirmation and required reason capture:

- Give Item and Live Claim Rewards.
- Inventory repair and inventory delete.
- Give Currency.
- Give Scrip.
- Award Intel.
- Award Character XP.
- Give Faction Reputation.
- Set Specialization XP.
- Journey node complete.
- Journey node reset.
- Player move.
- Clear all journey progress.
- Remove tutorial records.
- Clear codex discoveries.
- Disconnect player session.

Player 360 remains read-only. No Player 360 quick actions have been added.

## Current non-player mutation coverage

Currently covered by shared confirmation and required reason capture:

- Storage add item to container.
- Storage remove item from container.
- Database Run SQL.
- Battlegroup Server Control:
  - Start
  - Stop
  - Restart
  - Update
  - Backup
  - Restore
- Blueprint import.

Read-only views remain frictionless where possible, including database table/describe/sample/search views, Battlegroup pod/status/health diagnostics, Battlegroup support-bundle export, and Blueprint export.

## Risk levels

| Risk | Meaning |
|---|---|
| low | Reserved for future low-impact mutations |
| medium | General protected POST/PUT/PATCH mutation |
| high | High-impact player/admin mutation such as item grants, live grants, teleport, journey changes, and faction changes |
| destructive | Operations that can remove or replace important server or player state |

## Captured metadata

The audit middleware only captures these allowlisted fields:

- `reason`
- `player_id`
- `account_id`
- `actor_id`
- `controller_id`
- `item_id`
- `faction_id`
- `storage_id`

The middleware intentionally does not log full request bodies.

## Reason capture and enforcement

Reason text can be supplied in the `X-Admin-Reason` header or in a JSON request body field named `reason`.

Reason enforcement is controlled by `ADMIN_REQUIRE_REASON`. When enabled, high-risk and destructive requests must include a reason. When disabled, reason text is still captured when supplied.

Frontend confirmed modals force reason collection for active high-risk workflows before the request is sent.

## Security rules

- Do not log admin tokens.
- Do not log database credentials, SSH keys, or environment values.
- Do not log arbitrary request bodies.
- Keep public routes outside audit capture.
- Keep audit records available only from protected admin routes.
- Treat reason and target metadata as support metadata, not authorization.
- Keep mutation preview metadata separate from actual authorization checks.
- Keep Player 360 read-only until quick actions are explicitly implemented as new confirmed workflows.
- Prefer extracted workflow modals over further expansion of legacy `PlayersTab.tsx`.
- Preserve low-friction read-only views while confirming write, exec, import, or destructive flows.

## Current limitations

- Operator identity is still based on shared admin-token access rather than named operator accounts.
- Rollback hints describe operator guidance, but the backend does not yet create automatic before-change snapshots.
- Reason enforcement is environment-controlled and not yet configurable from the UI.
- Typed mutation wrappers are still needed for workflow-specific before/after metadata.
- Future Inventory Studio v2 workflows and future Player 360 quick actions still need shared confirmation review as they are implemented.

## Follow-up tasks

1. Add Player 360 quick actions only as new confirmed workflows with reason capture and target metadata.
2. Add shared confirmation to future Inventory Studio v2 workflows as they are implemented.
3. Add typed backend mutation wrappers per high-risk endpoint.
4. Add before-change snapshot helpers for inventory, journey/progression, teleport, storage, SQL, battlegroup exec, and blueprint import operations.
5. Add named operator identity when authentication supports individual users.
6. Add audit export and filtering support.
7. Add UI visibility for reason-enforcement state.

## Validation

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
