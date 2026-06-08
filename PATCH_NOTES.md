# Arrakis Control Panel Release Notes

## Current update: Route registration grouping

### Why this update was made

`routes.go` registered all API routes in one long function. That made route review harder and increased the chance of accidental drift as new domains are added. This refactor groups routes by domain without changing any route method, path, handler, or middleware behavior.

### What changed

- Split `registerRoutes` into domain-specific helpers:
  - `registerPublicRoutes`
  - `registerDiscordAuthRoutes`
  - `registerSelfServiceRoutes`
  - `registerCoreAdminRoutes`
  - `registerBattlegroupRoutes`
  - `registerPlayerRoutes`
  - `registerInventoryCoordinationRoutes`
  - `registerDatabaseRoutes`
  - `registerLogRoutes`
  - `registerNotificationRoutes`
  - `registerStorageRoutes`
  - `registerBlueprintRoutes`
- Preserved the existing `registerRoutes(mux *http.ServeMux)` entry point.
- Preserved all existing route mappings.

### Security and operator impact

- No route behavior intentionally changed.
- No mutation behavior changed.
- No new endpoint was added.
- No endpoint was removed.
- Middleware chain remains unchanged.
- Player 360 remains read-only.

### Validation

Validation pending from the canonical local update path:

```bash
./update.sh
```

### Remaining refactor work

- Continue update-script modularization review.
- Review PowerShell script modularization or document a final `v0.1.0` deferral.
- Continue Go review for handler boundaries, audit/mutation-safety helper boundaries, and typed execution surfaces.

---

## Previous update: Refactor gate and application identity cleanup

### Validation

Validated from the canonical local update path:

```bash
./update.sh
```
