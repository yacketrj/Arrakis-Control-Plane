# Changelog

All notable changes to this project will be documented in this file.

## [Unreleased]

### Added

- Added comprehensive refactor and improvement review at `docs/refactor-review.md`.
- Added focused augmented item stats domain model in `item_augments.go`.
- Added augment preset catalog at `web/src/tabs/augmentPresets.ts`.
- Added manual item template refresh endpoint at `POST /api/v1/players/templates/refresh`.
- Added frontend client support for manual item template refresh.
- Added frontend quality workflow for install, audit, typecheck, lint, and build.
- Added frontend `typecheck` script for explicit TypeScript validation.
- Added backend support for augmented Give Item payloads with per-item augment name, augment grade, roll value, explicit roll arrays, roll count, and effect indices.
- Added `FAugmentedItemStats` JSON generation for newly granted augmented item stacks.
- Added augmented give-item validation tests covering normalization, invalid augment inputs, roll bounds, aligned augment arrays, and empty stats behavior.
- Added frontend API payload types for augmented item grants.
- Added and wired the augmented Give Item modal component at `web/src/tabs/GiveItemModalAugmented.tsx`.
- Added `docs/augmented-give-items.md` with request examples, stored JSON shape, validation rules, and implementation notes.
- Added `docs/item-template-source-strategy.md` documenting the recommended hybrid database-plus-JSON item template source strategy.
- Added Linux dependency installer script at `scripts/linux/install-deps.sh`.
- Added Linux development launcher at `scripts/linux/run-dev.sh`.
- Added Linux production-style build helper at `scripts/linux/build-linux.sh`.
- Added Linux systemd backend installer at `scripts/linux/install-systemd.sh`.
- Added Linux helper functional test suite at `scripts/linux/test-linux.sh`.
- Added Linux helper GitHub Actions workflow at `.github/workflows/linux-helper-tests.yml`.
- Added Linux operating guide at `docs/linux.md`.
- Added README Linux quick start, build, validation, and systemd service instructions.

### Changed

- Strengthened Go quality workflow to run formatting verification, module graph verification, vet, and tests.
- Moved augment validation and stats JSON serialization out of database command code and into `item_augments.go`.
- Updated the augmented Give Item modal with augment presets, explicit comma-separated roll arrays, and generated payload preview.
- Refreshed item templates after `/api/v1/reconnect` so the cached hybrid item list stays current after reconnects.
- Removed the legacy embedded Give Item modal from `PlayersTab.tsx` so the active UI has a single augmented item workflow.
- Changed batch Give Item backend behavior so augmented grants write item stats directly on inserted stack rows instead of creating plain `{}` stats.
- Changed client Give Item row typing so item rows can carry optional augment arrays.
- Documented that database template discovery should run on connect, reconnect, manual refresh, or low-frequency scheduled refresh rather than on every UI search.
- Documented that `item-data.json` should remain as curated metadata and unseen-template fallback while live database templates provide current observed server templates.
- Updated `.gitignore` to exclude Linux build output, frontend build output, frontend dependencies, and local Linux runtime logs.
- Updated README from the older multi-give patch-only note into the primary project quick-start document.
- Updated Vite preview and dev CSP to remove inline script execution.
- Updated Vite preview and dev CSP so style elements load from self only while style attributes remain allowed for the current React UI.
- Kept DAST as a blocking security gate and continued remediating ZAP findings instead of disabling the scan.

### Fixed

- Restored player handler endpoint coverage after the item-template refresh refactor so existing player mutation routes remain available.
- Fixed Give Item workflow drift by removing the stale embedded modal after wiring the augmented modal.
- Fixed augmented item persistence so `AppliedAugments`, `AppliedAugmentRollData`, and `AppliedAugmentQualities` are generated as aligned arrays under `FAugmentedItemStats`.
- Fixed augmented batch grants so legacy non-augmented single-item payloads remain backward compatible.
- Fixed the latest DAST warning from run `26277991316` by tightening the `style-src` and `style-src-elem` directives in `web/vite.config.ts`.
- Fixed prior DAST warnings for missing browser isolation headers, incomplete CSP fallback directives, malformed ZAP baseline rules, and inline script execution.
- Fixed prior frontend build failures after unused auth dependency cleanup.
- Fixed prior Go module scan and test workflow bootstrap failures.
- Fixed prior stale frontend lockfile dependency metadata that referenced a vulnerable transitive dependency after the unused auth package was removed.

### Security

- Increased CI coverage for both Go and frontend quality gates before changes are treated as production-ready.
- Reduced mutation risk by isolating augment validation and serialization into a focused backend model file.
- Reduced operator error by adding preset-driven augment defaults and payload preview before submission.
- Reduced UI drift risk by removing the unused legacy Give Item modal code path.
- Reduced unsafe item-edit risk by validating augment names, augment counts, grade bounds, roll bounds, and roll array limits before writing `dune.items.stats`.
- Avoided database load amplification by documenting cached template discovery instead of per-keystroke database search.
- Preserved item template correctness by recommending a hybrid live database plus curated JSON model.
- Reduced frontend script execution risk by keeping script directives self-only.
- Reduced frontend style injection surface by removing inline allowance from style elements.
- Preserved style-attribute compatibility required by the current UI while retaining DAST visibility.
- Added Linux systemd hardening defaults including `NoNewPrivileges`, `PrivateTmp`, `ProtectSystem`, and constrained write paths.
- Added automated Linux helper tests for shell syntax, documentation coverage, build helper behavior, run-dev first-run bootstrap, and run-dev process-launch behavior.
- Documented Linux operational controls for loopback binding, admin token handling, SSH key protection, and reverse-proxy/TLS requirements.

### Operational Notes

- The Players tab Give Item button now opens only the augmented Give Item modal.
- Item templates can now be refreshed through `POST /api/v1/players/templates/refresh`; database search remains cached and operator-controlled rather than per-keystroke.
- Frontend quality now runs through GitHub Actions; local validation is still `cd web && npm install && npm audit --audit-level=high && npm run typecheck && npm run lint && npm run build`.
- Linux helper scripts are committed as text files. Run `chmod +x scripts/linux/*.sh` in local clones before execution.
- Regenerate `web/package-lock.json` locally from the current `web/package.json` with `npm install` and recommit it once confirmed clean with `npm audit --audit-level=high` and `npm run build`.

---

## [Security Hardening, Security Scanning, and Multi-Item Administration Update] - 2026-05-21

### Added

- Added backend admin-token authentication for all API routes.
- Added explicit browser origin allowlisting through `ALLOWED_ORIGINS`.
- Added safer listen-address normalization for local loopback use.
- Added HTTP server timeouts and request-size limits.
- Added read-only SQL enforcement for the database SQL endpoint.
- Added Kubernetes namespace and pod-name validation for log streaming.
- Added WebSocket route-specific query-token support for browser log streaming.
- Added frontend security headers for Vite dev and preview responses.
- Added multi-item Give Items workflow and batch item grant support.
- Added Go unit tests for batch item request normalization.
- Added GitHub Actions security scanning workflow covering SCA, SAST, secret scanning, DCA, and DAST.
- Added ZAP baseline rule policy in `.zap/rules.tsv`.
- Added blueprint import bounds checks.

### Changed

- Changed batch item grant semantics so `qty` means stack count and `stack_size` means items per stack.
- Changed status endpoint behavior to avoid returning pod IP information.
- Changed DAST workflow to scan Vite preview output.
- Changed gosec CI gate to focus on higher-signal findings.

### Fixed

- Fixed backend startup for bare port listen addresses.
- Fixed WebSocket log streaming authentication.
- Fixed cheats log SQL lookup.
- Fixed stack-size behavior for batch grants.
- Fixed Trivy, ZAP, Go SCA, gosec, Node SCA, and frontend build workflow failures found during security scan stabilization.

### Security

- Enforced backend authorization as the primary control for administrative endpoints.
- Reduced CORS exposure and accidental network exposure.
- Removed hardcoded credential material from source.
- Hardened frontend security headers and CSP.
- Established continuous security scanning coverage.