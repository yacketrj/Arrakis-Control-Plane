# Dune Admin Release Notes

## Current update: Admin audit log test coverage

### Why this update was made

The Admin Action Audit Log is a P0 foundation feature. After adding the audit sink, middleware, protected endpoint, and UI tab, the next step was to add unit coverage for the behavior that protects the audit model from regressions.

### What changed

- Added `audit_log_test.go`.
- Added coverage confirming protected mutating requests create audit events.
- Added coverage confirming failure statuses are recorded as failures.
- Added coverage confirming read-only and public-safe routes are not audited.
- Added coverage confirming audit reads are sorted newest-first and respect caller limits.

### Security and operator impact

- Improves confidence that public user-portal routes stay out of protected admin audit capture.
- Improves confidence that failed admin mutations are still recorded.
- Keeps the audit foundation minimal and secret-safe while preserving the path for typed metadata in the upcoming Mutation Safety Framework.

### Validation

Expected validation:

```bash
go test ./...
cd web
npm install
npm audit --audit-level=high
npm run typecheck
npm run lint
npm run build
```

---

## Previous update: Admin action audit log foundation

### Why this update was made

The feature roadmap identifies the Admin Action Audit Log as a P0 foundation item. As Dune Admin gains more powerful workflows, protected mutating requests need an append-only record that shows what endpoint was called, when it was called, whether it succeeded, and how long it took.

### What changed

- Added `audit_log.go` with an append-only JSONL audit sink.
- Added audit middleware that records protected mutating HTTP requests.
- Added a protected audit read endpoint at `GET /api/v1/audit/events`.
- Excluded public-safe routes from audit capture.
- Added `docs/admin-audit-log.md` describing the audit model, security rules, current limitations, and follow-up tasks.

### Security and operator impact

- Audit events are intentionally minimal and do not record request bodies, admin tokens, database credentials, SSH keys, or other secrets.
- Audit records are written to `admin-audit.jsonl` by default.
- Operators can override the path with `ADMIN_AUDIT_LOG`.
- This is the first audit foundation; typed payload summaries, operator names, reason fields, and rollback hints remain part of the upcoming Mutation Safety Framework.

### Validation

Expected validation:

```bash
go test ./...
cd web
npm install
npm audit --audit-level=high
npm run typecheck
npm run lint
npm run build
```

---

## Previous update: Public status endpoint and DB routine wiring

### Why this update was made

The admin console is being split into player-safe user portal surfaces and protected admin tools. The user portal needs only reviewed, redacted public APIs. The DB routine inspection feature also needed final UI wiring.

### What changed

- Added `/api/v1/public/status` as a redacted unauthenticated status endpoint.
- Kept sensitive status fields behind the protected `/api/v1/status` endpoint.
- Wired the DB Routines admin tab into the main protected admin interface.
- Documented the separation between public-safe portal data and protected admin operations.

### Validation

```bash
go test ./...
cd web
npm install
npm audit --audit-level=high
npm run typecheck
npm run lint
npm run build
```

---

## Previous update: Live Claim Rewards delivery mode

### Why this update was made

Direct inventory database writes support the richest Give Item workflow, including stack placement, grades, and augmented item stats. Online players may not see those direct inventory changes until they relog because the game client/server inventory state is already loaded.

This update exposes the existing Claim Rewards delivery path directly in the Give Item modal as a separate operator-selectable delivery mode.

### What changed

- Added a delivery mode selector to `GiveItemModalAugmented.tsx`.
- Added `Inventory Write` mode for direct database item creation with item grade, stack size, and augment support.
- Added `Live Claim Rewards` mode for online-friendly plain item delivery through the existing live reward API path.
- Added delivery-mode-aware payload preview.
- Blocked Live Claim Rewards mode when selected rows contain item grade greater than zero or augments, because the live reward path only supports plain template-and-amount entries.
- Disabled grade and augment controls while Live Claim Rewards mode is selected.

### Operator guidance

Use `Live Claim Rewards` when the player is online and you need them to receive plain items without logging out. The player should see a Claim Rewards prompt.

Use `Inventory Write` when you need exact inventory placement, item grade, stack control, or augmented item stats. For online players, this may still require logout/login before the game refreshes client-visible inventory state.

### Validation

Expected validation:

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

---

## Previous update: Go Quality run 26326549538 remediation

### Why this update was made

GitHub Actions run `26326549538` failed in the expanded template merge test because the test expected three merged templates, while the hybrid merge correctly returns four: two live database templates, one curated name-only template, and one curated item-rule template.

### Remediation

- Updated the test to expect four merged templates.
- Added explicit verification that the curated `item_rule_template` entry is present in the serialized template response.
- Preserved the lower-case friendly-name lookup assertion for `DB_TEMPLATE`.
