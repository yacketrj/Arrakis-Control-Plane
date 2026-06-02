# Inventory Requests and Farming Orders

The inventory request/order backend is a coordination surface for farming requests and fulfillment tracking. It does **not** mutate player inventory, guild storage, claim rewards, item stacks, currency, XP, or any other game-state table.

## Purpose

This slice supports two coordination workflows:

1. **Personal requests**: a Discord-linked player can request items or materials for personal use.
2. **Guild requests**: a guild or community can request materials for shared crafting, stocking, or farming operations.

Farming orders group one or more requests into an assignable work item. Filling an order marks linked requests fulfilled, but it does not deliver items in-game.

## Storage

The backend uses a local JSON store:

```env
INVENTORY_REQUEST_STORE=inventory-requests.json
```

If unset, the backend writes `inventory-requests.json` in the working directory. The file is written with `0600` permissions. Access is serialized with an in-process mutex to reduce clobber risk from simultaneous handler calls.

Current limitation: this is process-local file storage, not a multi-node database. A future production version should move this to a durable table if DA Manager is deployed with multiple backend instances.

## Request model

Inventory requests include:

- `scope`: `personal` or `guild`
- `requester_discord_id`
- `requester_name`
- `guild_id`
- `player_id`
- `item_name`
- `item_template_id`
- `quantity`
- `notes`
- `status`: `open`, `ordered`, `fulfilled`, or `cancelled`
- `order_id`
- timestamps

Validation rules:

- `scope` is required and must be `personal` or `guild`.
- `item_name` is required.
- `quantity` must be between `1` and `999999`.
- Personal requests require `requester_discord_id`.
- Guild requests require `guild_id`.
- Text fields are trimmed, length-limited, and reject unsupported control characters.

## Order model

Inventory orders include:

- `scope`: `personal` or `guild`
- `guild_id`
- `requester_discord_id`
- `assignee_discord_id`
- `assignee_name`
- `request_ids`
- `status`: `open`, `filled`, or `cancelled`
- `notes`
- timestamps

Validation rules:

- `request_ids` must include at least one existing request.
- Linked requests must match the order scope.
- Guild orders can only link requests from the same `guild_id`.
- Fulfilled or cancelled requests cannot be added to a new order.
- Duplicate request IDs are deduplicated before validation.

## Endpoints

| Method | Path | Purpose |
| --- | --- | --- |
| `GET` | `/api/v1/inventory/requests` | List requests, optionally filtered by `scope`, `status`, `guild_id`, or `requester_discord_id`. |
| `POST` | `/api/v1/inventory/requests` | Create a personal or guild inventory request. |
| `PATCH` | `/api/v1/inventory/requests/{id}` | Update request status, linked order ID, or notes. |
| `GET` | `/api/v1/inventory/orders` | List orders, optionally filtered by `scope`, `status`, or `guild_id`. |
| `POST` | `/api/v1/inventory/orders` | Create a farming order from one or more existing request IDs. |
| `PATCH` | `/api/v1/inventory/orders/{id}` | Update order status, assignee, or notes. |

All endpoints are protected by the existing backend auth middleware when served normally. They are not public Discord OAuth paths.

## Frontend UI

The frontend adds a protected **Farming Requests** tab.

The tab supports:

- Filtering requests and orders by scope/status.
- Creating personal or guild inventory requests.
- Selecting open requests and grouping them into a farming order.
- Marking open orders as `filled` or `cancelled`.
- Refreshing request/order state from the backend.

The tab uses `web/src/api/inventoryRequests.ts` rather than the high-risk player/admin API client. This keeps the coordination workflow separate from player inventory mutation, Player 360, and direct admin actions.

## Status propagation

Creating an order marks linked requests as `ordered` and sets their `order_id`.

Marking an order `filled` sets `completed_at` on the order and marks linked requests as `fulfilled`.

Marking an order `cancelled` clears linked request `order_id` values and returns linked requests to `open`.

## Browser/API behavior

The backend CORS middleware allows `PATCH` so browser clients can call the update endpoints after a successful preflight.

## Validation

Run backend tests from the local checkout or CI:

```bash
go test ./...
```

Run frontend checks from the local checkout:

```bash
cd web
npm run typecheck
npm run lint
npm run build
```

Manual validation should confirm:

1. Creating a personal request requires `requester_discord_id`.
2. Creating a guild request requires `guild_id`.
3. Creating an order with a missing request ID is rejected.
4. Creating an order marks linked requests `ordered`.
5. Filling an order marks linked requests `fulfilled`.
6. Cancelling an order returns linked requests to `open`.
7. Browser preflight allows `PATCH` for the update endpoints.
8. The Farming Requests tab can list, create, group, fill, and cancel request/order records without touching player inventory.

## Safety boundary

This feature is intentionally non-mutating. It is a request/order ledger only. In-game delivery, player inventory writes, guild storage writes, and Player 360 self-service actions must remain separate and must not be enabled until Discord identity-to-player mapping and the mutation-safety workflow are explicitly implemented.
