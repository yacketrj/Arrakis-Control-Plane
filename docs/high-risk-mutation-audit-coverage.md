# High-Risk Mutation Audit Coverage

## Purpose

This note tracks the AppSec follow-up for endpoint-by-endpoint audit-event assertion coverage on high-risk and destructive mutation routes.

The goal is to make sure every high-risk/destructive mutation route emits an audit event with the expected mutation-safety classification and core accountability metadata.

## Coverage added

`audit_log_test.go` now includes `TestAuditMiddlewareHighRiskAndDestructiveRouteCoverage`.

The test iterates across the current high-risk/destructive mutation route set and verifies that each route emits exactly one audit event with:

- request method
- request path
- mutation-safety action
- mutation-safety risk
- destructive flag
- requires-reason flag
- requires-preview flag
- HTTP status
- result outcome
- sanitized reason
- common target metadata
- request ID

## Covered routes

The test covers representative concrete paths for:

- Discord player-link deletion
- reconnect
- Battlegroup exec
- player item/currency/faction/scrip/XP/intel/live grants
- player kick
- item deletion
- item stack-size update
- specialization reset/set
- faction-tier mutation
- journey complete/reset/wipe
- tutorial/codex destructive mutations
- item repair
- teleport
- storage give-item
- manual database SQL
- log-stream ticket issuance
- notification broadcast
- blueprint import

## Security impact

- No route behavior changed.
- No mutation behavior changed.
- No new endpoint was added.
- Player 360 remains read-only.
- The change adds regression evidence that high-risk/destructive mutation routes produce expected audit accountability fields.

## Remaining work

This is partial coverage for audit-event assertions. Remaining audit work includes:

- handler-specific pre/post-change snapshot assertions where practical
- route-specific target metadata assertions beyond the shared common body fields
- negative-path audit assertions for blocked mutations
- SAST/DAST/dependency evidence
- manual abuse-case validation

## Validation

Required from the canonical local update path:

```bash
./update.sh
```
