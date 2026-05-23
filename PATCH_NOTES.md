# Dune Admin Release Notes

## Current update: Augmented Give Items and item template source strategy

### Why this update was made

The Give Item workflow needed support for augmented items, including augment names, augment grades, and normalized roll data. Operators also needed a clear decision on whether item templates should come from the live database, a JSON file, or both.

### Security and operator impact

- Added backend support for augmented Give Item payloads on `POST /api/v1/players/give-item`.
- Added per-item augment definitions with augment name, augment grade, roll value, explicit roll arrays, roll count, and effect indices.
- Added validation for augment name presence, maximum augment count, grade bounds, roll bounds, roll count bounds, and explicit roll-array bounds before writing item stats.
- Added `FAugmentedItemStats` generation using the observed game-compatible `[[], payload]` wrapper shape.
- Preserved alignment between `AppliedAugments`, `AppliedAugmentRollData`, and `AppliedAugmentQualities`.
- Preserved legacy single-item non-augmented payload behavior for existing callers.
- Augmented grants create new stack rows rather than topping up existing stacks, avoiding accidental merging between augmented and plain items.
- Added frontend API types for augmented Give Item payloads.
- Added a ready-to-wire augmented Give Item modal component at `web/src/tabs/GiveItemModalAugmented.tsx`.

### Current UI status

The active `PlayersTab.tsx` Give Item button now opens `GiveItemModalAugmented.tsx`. The prior embedded modal has been renamed to `LegacyGiveItemModal` as a short-term rollback/reference component until the larger player tab is split into smaller files.

### Augmented payload example

```json
{
  "player_id": 123,
  "items": [
    {
      "template": "ItemTemplateWeaponExample",
      "qty": 1,
      "quality": 5,
      "stack_size": 1,
      "augments": [
        {
          "name": "T6_Augment_Damage1",
          "grade": 5,
          "roll": 1.0,
          "roll_count": 1,
          "effect_indices": []
        },
        {
          "name": "T6_Augment_Magazinecapacity1",
          "grade": 5,
          "rolls": [1.0, 1.0, 1.0],
          "effect_indices": []
        }
      ]
    }
  ]
}
```

### Item template source strategy

Use a hybrid template source:

1. Pull live observed item templates from the database at backend connect, reconnect, manual refresh, or low-frequency scheduled refresh.
2. Keep `item-data.json` as curated fallback metadata for display names, stack defaults, volume defaults, aliases, and templates that have not appeared in the live database.
3. Serve the frontend from a merged in-memory cache instead of querying the database on every search keystroke.

A `WITH` query can improve readability and can combine observed template IDs with observed stack or volume data, but the query shape is not the primary optimization. The main performance win is caching the result and avoiding repeated typeahead queries against `dune.items`.

Recommended refresh cadence:

- Backend startup after DB connection.
- `/api/v1/reconnect`.
- Manual operator refresh endpoint.
- Optional background refresh every 15 to 60 minutes for long-running admin sessions.

Do not create indexes automatically from the admin app. If a large live database makes template refresh slow, operators can evaluate a controlled database migration such as a `template_id` index.

### Testing added

Added Go unit tests for:

- Augmented Give Item request normalization.
- Legacy single-item payload with augments.
- Invalid augment name, grade, roll, and roll-array validation.
- `FAugmentedItemStats` JSON generation.
- Empty stats behavior when no augments are supplied.

### Validation

```bash
go mod tidy
go test ./...
cd web
npm install
npm audit --audit-level=high
npm run build
```

### Known limitations and follow-up

- The legacy embedded Give Item modal remains exported as `LegacyGiveItemModal` for short-term rollback/reference and should be removed after the player tab is split into smaller files.
- Some augments appear to require multiple stat rolls. Until an augment metadata catalog exists, operators should provide explicit `rolls` arrays for augments known to need multiple values.
- Database template refresh strategy is documented, but a manual refresh endpoint and scheduled refresh loop have not yet been implemented.
- `web/package-lock.json` needs clean regeneration from the current frontend manifest.

---

## Previous update: Linux version support and helper tests

### Why this update was made

The Windows-oriented workflow is now complemented by a Linux version so operators can install dependencies, run the app locally, build a Linux backend binary, and install the backend as a systemd service without translating Windows steps manually.

After the initial Linux scripts were added, automated tests were added so Linux helper behavior is validated continuously instead of relying only on manual review.

### Security and operator impact

- Added Linux helper scripts under `scripts/linux/` for dependency setup, local development, Linux builds, and systemd installation.
- Added Linux helper functional tests in `scripts/linux/test-linux.sh`.
- Added `.github/workflows/linux-helper-tests.yml` so GitHub Actions runs Linux helper tests on pushes to `main` and manual workflow dispatch.
- Added a Linux operating guide at `docs/linux.md` with configuration, build, run, service, validation, and security steps.
- Updated the README with Linux quick-start instructions and runtime configuration guidance.
- Updated `.gitignore` so Linux build output, frontend build output, frontend dependencies, and local runtime logs are not committed.
- The systemd installer creates a dedicated service user and enables service hardening options including no new privileges, private tmp, protected system paths, and constrained write access.
- Linux guidance continues to require loopback backend binding by default, strong admin tokens, SSH key protection, and TLS/reverse-proxy controls for any remote exposure.

### Validation

```bash
go mod tidy
go test ./...
cd web
npm install
npm audit --audit-level=high
npm run build
```

---

## Release: Security Hardening, Security Scanning, Multi-Item Administration, and Documentation Standards Update

### Summary

This release makes Dune Admin safer and more reliable for live server operations. It establishes backend-enforced admin-token authentication, explicit CORS allowlisting, safer loopback listen defaults, server timeouts, raw SQL restrictions, request-size controls, Kubernetes log target validation, reduced status data exposure, hardened frontend security headers, blueprint import bounds checks, CI security scanning, removal of hardcoded capture credentials, and the requirement that every future change keep both `PATCH_NOTES.md` and `CHANGELOG.md` current.

### Security notes for operators

- Treat `ADMIN_TOKEN` as a privileged secret.
- Rotate any previously shared or committed credentials.
- Keep `.env`, SSH keys, database snapshots, generated secrets, dependency folders, and build output out of source control.
- Prefer `LISTEN_ADDR=127.0.0.1:8080` for local use.
- Do not expose the backend directly to the internet.
- Place remote access behind TLS, a trusted reverse proxy, and a strong identity provider.
