# Changelog

All notable changes to this project will be documented in this file.

## [Unreleased]

### Added

- Added `docs/player-360-profile.md` as the design foundation for the next P1 read-only Player 360 implementation slice.
- Added roadmap status notes showing the project has moved from P0 safety foundation into the Player 360 planning step for Phase 2 operator support.
- Added admin audit documentation sync covering the current audit event model, protected event review endpoint, local JSONL audit path, captured fields, reason capture, and operational limitations.
- Added implementation tracker updates so Admin Action Audit Log is marked complete and Player 360 Profile is identified as the next feature slice.
- Added SSH tunnel management foundation with managed local forwarding for game-management database access.
- Added SSH tunnel mode configuration for managed, existing, and direct troubleshooting paths.
- Added managed tunnel lifecycle cleanup during reconnect and shutdown.
- Added SSH tunnel normalization tests and empty-status test coverage.
- Added `docs/ssh-tunnel-management.md` with tunnel modes, runtime behavior, operator guidance, and validation steps.
- Added Battlegroup Health Diagnostics with protected read-only Kubernetes diagnostics for pods, services, statefulsets, deployments, PVCs, recent events, nodes, and pod metrics.
- Added Battlegroup Health Diagnostics UI cards with diagnostic descriptions, executed commands, output, and per-section errors.
- Added raw and redacted Battlegroup Health support-bundle exports for local review and safer developer/support handoff.
- Added `docs/battlegroup-health-diagnostics.md` with endpoint, UI, security, validation, troubleshooting, export, and redaction guidance.
- Added Go Module Integrity workflow to enforce committed `go.sum` and prevent missing or stale module checksum drift.
- Added redacted unauthenticated public status endpoint at `/api/v1/public/status` for the future player-safe user portal.
- Added feature design and priority roadmap at `docs/admin-feature-design-and-priorities.md`, including the item delivery architecture distinction between gameplay inventory, direct inventory writes, and claim reward queue grants.
- Added Live Claim Rewards delivery mode to the augmented Give Item modal for online-friendly plain item grants through the existing live grant endpoint.
- Added delivery-mode payload preview so operators can distinguish direct inventory writes from live claim rewards before submitting.
- Added comprehensive refactor and improvement review at `docs/refactor-review.md`.
- Added focused augmented item stats domain model in `item_augments.go`.
- Added augment preset catalog at `web/src/tabs/augmentPresets.ts`.
- Added reusable frontend Give Item payload helper module at `web/src/tabs/giveItemPayload.ts`.
- Added manual item template refresh endpoint at `POST /api/v1/players/templates/refresh`.
- Added frontend client support for manual item template refresh.
- Added frontend quality workflow for install, audit, typecheck, lint, and build.
- Added frontend `typecheck` script for explicit TypeScript validation.
- Added backend support for augmented Give Item payloads with per-item augment name, augment grade, roll value, explicit roll arrays, roll count, and effect indices.
- Added `FAugmentedItemStats` JSON generation for newly granted augmented item stacks.
- Added expanded augmented give-item tests covering normalization, invalid augment inputs, quality alias handling, roll defaulting, explicit roll precedence, aligned augment arrays, empty stats behavior, template merging, and template response serialization.
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

- Updated `docs/admin-feature-design-and-priorities.md` so the next implementation slice is Player 360 Profile read-only foundation instead of the already-completed audit foundation.
- Updated the Player 360 roadmap entry to fold Currency and Online Status into Player Info.
- Clarified Battlegroup Status v2 as future Prometheus/Grafana graph and diagnostic improvement work.
- Updated `PATCH_NOTES.md` with the Player 360 planning and roadmap correction update.
- Updated `docs/admin-implementation-tasks.md` to reflect landed audit work, active mutation-safety work, and Player 360 Profile as the next feature slice.
- Routed PostgreSQL game-management access through the configured SSH tunnel policy instead of the previous inline SSH dial function.
- Updated Battlegroup tab navigation with separate `Pods` and `Health Diagnostics` views.
- Updated Battlegroup support-bundle workflow with separate raw and redacted export actions.
- Clarified that Claim Rewards Queue grants are not the same as crafting, looting, finding an item, or authoritative gameplay inventory mutation.
- Updated `GiveItemModalAugmented.tsx` so Inventory Write remains the full-featured grade/augment/stat path while Live Claim Rewards is clearly limited to plain template-and-amount grants.
- Updated `GiveItemModalAugmented.tsx` to use the shared Give Item payload helper module instead of duplicating clamping, preset, roll parsing, and payload mapping logic inline.
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

- Fixed stale roadmap text that still pointed to the completed audit log as the next implementation slice.
- Fixed the Player 360 roadmap typo from Currency wording and corrected the malformed Battlegroup Status v2 roadmap row.
- Fixed stale admin implementation tracker statuses that still described Admin Action Audit Log and Mutation Safety Framework as future-only work.
- Fixed managed SSH tunnel cleanup so repeated close paths are idempotent.
- Fixed Battlegroup view button variants to use supported HeroUI values.
- Fixed missing/stale `go.sum` risk by adding CI validation that runs `go mod tidy`, verifies `go.mod`/`go.sum` have no diff, verifies module checksums, and runs Go tests.
- Fixed Give Item operator confusion by documenting that the live reward path is a claim queue, not a universal instant inventory mutation system.
- Fixed Give Item operator confusion by exposing the existing live grant mechanism as an explicit delivery mode rather than hiding it behind a separate action.
- Fixed Go Quality run `26326549538` by aligning template merge test expectations with the hybrid database-plus-JSON merge behavior.
- Fixed frontend Give Item helper drift by making the active modal use `giveItemPayload.ts` directly.
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

- Reinforced that Player 360 Profile v1 must start as a protected read-only support surface with no new mutation paths.
- Reinforced that future Player 360 quick actions should reuse the audit and mutation-safety foundation before adding new mutation workflows.
- Improved operator accountability by documenting the current audit event capture fields, reason capture, and protected review path.
- Reduced direct infrastructure exposure by routing supported game-management database access through managed SSH tunnels by default.
- Reduced reconnect risk by cleaning managed tunnels before reopening SSH and database sessions.
- Reduced diagnostic handoff risk by adding redacted Battlegroup Health bundle export for common infrastructure identifiers.
- Reduced infrastructure mutation risk by keeping Battlegroup Health Diagnostics read-only with fixed server-side commands.
- Reduced support exposure risk by documenting redaction limits and requiring operator review before external sharing.
- Reduced supply-chain drift risk by making `go.sum` integrity a dedicated CI gate.
- Reduced design risk by distinguishing full-fidelity inventory writes from plain claim-queue grants in the feature roadmap.
- Reduced live-operation risk by blocking Live Claim Rewards mode for graded or augmented rows that require direct inventory/stat writes.
- Increased test coverage around augmented item roll defaults, explicit roll arrays, grade aliases, and template serialization before generated item stats are written.
- Increased CI coverage for both Go and frontend quality gates before changes are treated as production-ready.
- Reduced mutation risk by isolating augment validation and serialization into a focused backend model file.
- Reduced operator error by adding preset-driven augment defaults and payload preview before submission.
