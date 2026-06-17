# Arrakis Control Panel Release Notes

## Current update: Frontend typecheck fix

### Why this update was made

Frontend validation failed in the Battlegroup tab.

### What changed

- Removed an unsupported UI prop from the server-control button.
- Kept the visible runtime explanation text below the disabled controls.
- Durable detail is archived in:

```text
docs/changelog/unreleased/2026-06-16-heroui-button-title-typecheck.md
```

### Impact

- Frontend typecheck should no longer fail on that button prop.
- Runtime behavior did not change.

### Validation

Run the normal update script.

### Remaining final-v0.1.0 gate

- Post-release verification after tag/artifact install or launch.
