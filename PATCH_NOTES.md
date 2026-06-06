# Dune Admin Release Notes

## Current update: Route-specific audit target assertions

### Why this update was made

The AppSec hardening track is continuing before new Live Admin / RMQ / Welcome Kit features. The previous slices proved that high-risk/destructive mutations emit audit events and that blocked mutations are still auditable. This slice tightens target accountability by ensuring audit events capture route-specific identifiers needed to investigate player, item, server-command, vehicle, and guild mutations.

### What changed

- Expanded audit target metadata extraction in `audit_log.go`.
- Added `audit_log_target_test.go`.
- Added coverage for expanded target fields:
  - `player_id`
  - `account_id`
  - `actor_id`
  - `controller_id`
  - `fls_id`
  - `item_id`
  - `item_template`
  - `item_template_id`
  - `template_id`
  - `quantity`
  - `amount`
  - `quality`
  - `vehicle_id`
  - `guild_id`
  - `rank`
  - `command`
  - `command_path`
- Added redaction coverage for expanded sensitive target fields.
- Added `docs/changelog/unreleased/2026-06-route-specific-audit-targets.md` as the durable per-slice record.

### Security and operator impact

- No route behavior changed.
- No mutation behavior changed.
- No new endpoint was added.
- Player 360 remains read-only.
- This improves audit investigation quality for high-risk mutation attempts without exposing arbitrary raw command publishing.

### Validation

Validated from the canonical local update path:

```bash
./update.sh
```

Non-blocking build-performance warning observed:

```text
[PLUGIN_TIMINGS] Your build spent significant time in plugin `@tailwindcss/vite:generate:build`.
```

### Remaining AppSec work

- pre/post-change review verification where practical
- SAST/DAST/dependency evidence
- manual abuse-case validation

---

## Previous update: Blocked mutation audit coverage

### Why this update was made

The AppSec hardening track is continuing before new Live Admin / RMQ / Welcome Kit features. The previous audit slice verified successful high-risk/destructive mutation audit events. This slice adds negative-path coverage so blocked high-risk/destructive mutations are also auditable when admin-reason enforcement rejects them before the downstream handler runs.

### What changed

- Added `audit_log_negative_test.go`.
- Added table-driven coverage for high-risk and destructive mutation routes blocked by missing admin reason.
- Verified blocked mutations:
  - do not reach the downstream handler
  - return `400 Bad Request`
  - still emit exactly one audit event
  - record the expected method and path
  - record mutation-safety action, risk, destructive flag, reason flag, and preview flag
  - record failure result and `400` status
  - preserve request ID
  - preserve common target metadata such as `player_id`, `account_id`, and `actor_id`
- Added oversized-body negative-path coverage for reason inspection.
- Added `docs/changelog/unreleased/2026-06-blocked-mutation-audit-coverage.md` as the durable per-slice record.

### Security and operator impact

- No route behavior changed.
- No mutation behavior changed.
- No new endpoint was added.
- Player 360 remains read-only.
- This adds regression evidence that blocked high-risk/destructive mutations remain visible in the admin audit trail.

### Validation

Validated from the canonical local update path:

```bash
./update.sh
```
