# Arrakis Control Panel Release Notes

## Current update: Audit metadata helper refactor

### Why this update was made

The Go code-quality/refactor gate requires reviewing handler boundaries, audit/mutation-safety helper boundaries, and typed/allowlisted execution surfaces before final `v0.1.0`. After extracting pure audit risk classification, the next low-risk boundary was audit metadata extraction.

This slice separates audit request-body inspection, target extraction, sanitization, remote-address extraction, and admin-token hashing from audit middleware and audit persistence.

### What changed

- Added `audit_metadata.go`.
- Moved audit metadata helper definitions out of `audit_log.go`:
  - `mutationAuditMetadata`
  - `mutationAuditTargetKeys`
  - `extractMutationAuditMetadata`
  - `payloadString`
  - `auditScalar`
  - `sanitizedAuditString`
  - `auditRemoteAddr`
  - `auditAdminTokenHash`
- Left audit middleware, audit persistence, audit status/result helpers, and audit event handling in `audit_log.go`.
- Preserved existing audit middleware call sites and metadata behavior.

### Security and operator impact

- No route behavior changed.
- No mutation behavior changed.
- No new endpoint was added.
- Audit metadata extraction behavior is intended to remain identical.
- This makes future review of request-body inspection, sensitive-text redaction, and target extraction easier.

### Validation

Validation pending.

Recommended validation:

```bash
./update.sh
```

### Remaining refactor work

- Continue Go review for mutation-safety helper boundaries.
- Continue Go review for handler boundaries and typed/allowlisted execution surfaces.
- Consider splitting audit storage into `audit_store.go` if validation remains clean.

---

## Previous update: Audit risk helper refactor

### Validation

Validated from the canonical local update path:

```bash
./update.sh
```
