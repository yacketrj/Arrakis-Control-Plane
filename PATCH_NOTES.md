# Arrakis Control Panel Release Notes

## Current update: Audit risk helper refactor

### Why this update was made

The Go code-quality/refactor gate requires reviewing handler boundaries, audit/mutation-safety helper boundaries, and typed/allowlisted execution surfaces before final `v0.1.0`. The first low-risk Go refactor target was `audit_log.go`, which mixed audit middleware, persistence, metadata extraction, status/result helpers, and action/risk classification.

This slice extracts pure action/risk classification into its own file without changing route behavior, mutation behavior, audit persistence, or handler registration.

### What changed

- Added `audit_risk.go`.
- Moved pure audit/mutation classification helpers out of `audit_log.go`:
  - `auditActionName`
  - `mutationRiskForRequest`
  - `highRiskMutationPathMarkers`
- Left audit middleware, audit persistence, audit metadata extraction, and audit event handling in `audit_log.go`.
- Preserved the existing `mutationSafetyForPath` and audit middleware call sites.

### Security and operator impact

- No route behavior changed.
- No mutation behavior changed.
- No new endpoint was added.
- Audit risk classification rules are intended to remain identical.
- This makes future review of high-risk/destructive operation classification easier.

### Validation

Validation pending.

Recommended validation:

```bash
./update.sh
```

### Remaining refactor work

- Continue Go review for audit metadata extraction boundaries.
- Continue Go review for mutation-safety helper boundaries.
- Continue Go review for handler boundaries and typed/allowlisted execution surfaces.

---

## Previous update: Roadmap discoverability update

### Validation

Validated from the canonical local update path:

```bash
./update.sh
```
