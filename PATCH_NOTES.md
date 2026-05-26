# Dune Admin Release Notes

## Current update: Startup hardening and DB-down UI gating

### Why this update was made

Startup and database connectivity are release-critical. If SSH, tunnel, discovery, or DB connection fails, the application must not look operational while DB-backed tools are unusable.

### What changed

- Added `config_paths.go` with config-safe local path expansion and validation.
- Added support for documented Windows percent-variable home-directory SSH key paths.
- Rejected unsupported PowerShell-style path expressions during validation.
- Tightened runtime validation for SSH, DB, admin token, tunnel mode, and port configuration.
- Changed startup so invalid configuration exits and tells the operator to rerun setup.
- Changed connection failure messaging from ambiguous startup behavior to explicit degraded mode.
- Added frontend DB-down gating for DB-backed tabs.
- Added a DB-unavailable banner when backend status reports `db_connected=false`.
- Added a DB-required panel with SSH/DB/tunnel state and reconnect retry.
- Added frontend startup diagnostic typing.
- Expanded the production security release gate for startup reliability, config validation, injection resistance, logging redaction, encryption in transit, and secret-at-rest requirements.
- Added Battlegroup config-management and observability design notes from operator review intake.

### Security and operator impact

- Invalid configuration should fail earlier and more clearly.
- DB-backed tabs no longer present as usable when DB connectivity is down.
- Secret values are validated as opaque values for presence/control characters.
- Full encrypted secret storage remains a P0 release-gate item.
- The direct DB connection-construction patch was blocked by the connector safety filter and should be validated from the local checkout if `update.ps1` reports compile or startup issues.

### Validation

Validation required in the Windows development environment:

```powershell
.\update.ps1
```

Also manually validate startup repair, unsupported path rejection, supported path expansion, and DB-down frontend gating.

---

## Previous update: Inventory Studio v2 post-action diff panel

Inventory Studio v2 added a post-action diff panel for confirmed add, repair, and removal workflows.

---

## Previous update: Inventory Studio v2 confirmed catalog item add

Inventory Studio v2 added confirmed catalog-item add with quantity and quality inputs.

---

## Previous update: Inventory Studio v2 confirmed item removal

Inventory Studio v2 added confirmed selected-item removal with a before-action inventory snapshot.

---

## Previous update: Inventory Studio v2 confirmed item repair

Inventory Studio v2 added confirmed selected-item repair with a before-action inventory snapshot.

---

## Previous update: Inventory Studio v2 item catalog browser

Inventory Studio v2 added a read-only item catalog browser.
