# Changelog

All notable changes to this project will be documented in this file.

This project follows a corporate change-management style informed by ITIL release/change practices and NIST SSDF secure-development expectations. Entries should distinguish security, operational, feature, documentation, validation, and known-risk changes.

## [Unreleased]

### Added

- Added Discord auth route/session coverage in `discord_auth_test.go` for route registration, role mapping, session lookup, expiry eviction, session hash generation, and logout invalidation.
- Added `docs/discord-auth.md` with runtime configuration, endpoint, role mapping, session behavior, validation, and current limitation notes.
- Added NIST SP 800-218 SSDF as the primary secure-development baseline for DA Manager.
- Added `docs/NIST_SSDF_ALIGNMENT.md` with SSDF practice-group mapping, release gates, and control priorities.
- Added `docs/COMPLIANCE_READINESS.md` with SOC 2 and ISO/IEC 27001 readiness assessment.
- Added `docs/ADMIN_GUIDE.md` with deployment, SSH, token, network exposure, diagnostics, validation, backup, recovery, incident response, and change-management guidance.
- Added `docs/RELEASE_CHECKLIST.md` with ITIL-style release metadata, risk/impact assessment, validation gates, rollback plan, secret-rotation checklist, known-issue tracking, and approval record.
- Added `docs/CONTROL_MATRIX.md` to map DA Manager controls to NIST SSDF practice groups and evidence sources.
- Added `docs/RISK_REGISTER.md` to track active security, operational, release, and compliance risks.
- Added `SECURITY_REMEDIATION_TODO.md` to convert external security audit findings into a prioritized remediation backlog.
- Added one-time WebSocket log stream tickets through `ws_ticket.go` and `POST /api/v1/logs/stream-ticket`.
- Added fail-closed backend exposure validation for non-loopback `LISTEN_ADDR` values.
- Added extracted confirmed player admin-actions modal at `web/src/tabs/PlayerAdminActionsModal.tsx`.
- Added dedicated confirmed player move modal at `web/src/tabs/PlayerTeleportModal.tsx`.
- Added confirmed player resource, XP, specialization, and journey node actions modal at `web/src/tabs/PlayerActionsModalConfirmed.tsx`.
- Added extracted confirmed inventory modal at `web/src/tabs/InventoryModal.tsx`.
- Added shared frontend mutation confirmation hook at `web/src/hooks/useMutationConfirmation.tsx`.
- Added mutation confirmation support for backend classification lookup, conservative fallback classification, risk/warning display, target context, rollback guidance, and admin reason capture.
- Added validated Player 360 v1 read-only profile workflow with backend endpoint, standalone tab, and Players-table launcher.
- Added Player 360 row launcher from the Players table through `web/src/tabs/PlayersTabWith360Launcher.tsx`.
- Added automatic Player 360 profile loading when opened from a selected Players table row.
- Added standalone Player 360 frontend tab and main app navigation entry for read-only profile lookup by PlayerCharacter actor ID.
- Added `web/src/api/playerProfile.ts` frontend response types and protected fetch helper for `GET /api/v1/players/{id}/profile`.
- Added Player 360 backend profile foundation with `GET /api/v1/players/{id}/profile`.
- Added `player_profile.go` for protected read-only Player 360 aggregation, response modeling, section-level safe errors, inventory summaries, journey summaries, online-state matching, and ID helper behavior.
- Added `routes.go` to centralize backend HTTP route registration and include the Player 360 profile endpoint.
- Added `player_profile_test.go` coverage for Player 360 summary helpers, online-state matching, ID helper behavior, and safe section-error wording.
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

- Registered Discord auth endpoints in `routes.go` for login, callback, current auth context, logout, and registered-user review.
- Updated `PATCH_NOTES.md` with the Discord auth route/session validation status.
- Updated project governance posture to treat DA Manager as a corporate development effort with NIST SSDF control alignment and ITIL-style release/change records.
- Updated `.env.example` for hardened defaults: Ed25519 SSH key, SSH known_hosts path, strict admin token guidance, loopback backend binding, and explicit remote exposure override documentation.
- Updated frontend admin-token handling to reject malformed Browser Access Keys and migrate valid legacy `localStorage` tokens into `sessionStorage`.
- Updated log streaming to request a one-time WebSocket ticket before opening the stream.
- Updated backend startup to refuse non-loopback binding unless `DUNE_ADMIN_REMOTE_EXPOSURE=reverse-proxy-tls` is explicitly configured.
- Updated SSH key auto-detection to prefer DA Manager Ed25519 key names and remove RSA auto-detection from the app path.
- Updated the active Players-table wrapper so Give Item, Inventory, Actions, Move, Admin, and Player 360 workflows route through extracted modal surfaces rather than expanding `PlayersTab.tsx`.
- Updated Give Item, Inventory repair/delete, player resource/spec actions, journey node actions, player move, and player admin actions to use shared mutation confirmation and required admin reason capture.
- Updated `docs/mutation-safety-framework.md` with the shared frontend confirmation hook, integration pattern, limitations, and follow-up migration tasks.
- Updated `docs/admin-implementation-tasks.md` so the active focus is migrating existing high-risk workflows to the shared confirmation hook before Player 360 quick actions.
- Updated `PATCH_NOTES.md` with the shared frontend mutation confirmation foundation status.
- Marked Player 360 Profile as validated and Done in implementation tracking after clean compile validation.
- Updated `PATCH_NOTES.md` with the validated Player 360 v1 status.
- Updated `web/src/App.tsx` so the Players tab uses the Player 360 launcher wrapper.
- Updated `web/src/tabs/Player360Tab.tsx` to read the selected player actor ID from the launcher and auto-load that profile.
- Updated `docs/player-360-profile.md` and `PATCH_NOTES.md` with the Player 360 launcher status.
- Updated Player 360 to fetch specialization tracks by controller ID when available, matching the existing player action flow.
- Updated `docs/player-360-profile.md` with current frontend tab status, validation requirements, and follow-up work.
- Updated `docs/admin-implementation-tasks.md` so Player 360 validation is the active implementation focus.
- Updated `PATCH_NOTES.md` with the Player 360 read-only frontend tab update.
- Refactored `server.go` to call shared route registration through `registerRoutes`.
- Updated `PATCH_NOTES.md` with the Player 360 backend profile foundation update.
- Updated `docs/admin-feature-design-and-priorities.md` so the next implementation slice is Player 360 Profile read-only foundation instead of the already-completed audit foundation.
- Updated the Player 360 roadmap entry to fold Currency and Online Status into Player Info.
- Clarified Battlegroup Status v2 as future Prometheus/Grafana graph and diagnostic improvement work.
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

- Fixed static admin-token exposure through WebSocket URLs by replacing `ws_token` with one-time stream tickets.
- Fixed unsafe backend exposure behavior by failing closed on non-loopback `LISTEN_ADDR` unless explicitly acknowledged for reverse-proxy/TLS deployment.
- Fixed frontend legacy token persistence by removing valid legacy tokens from `localStorage` after migration to `sessionStorage` and deleting invalid legacy values.
- Fixed Player 360 specialization lookup so it uses the controller ID when the player identity is available.
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

- Kept Discord login/callback as the only public Discord auth routes while session identity, logout, and registered-user review remain behind the normal backend auth path.
- Adopted NIST SSDF as the active secure-development control baseline.
- Added governance and evidence documentation needed to support future SOC 2 / ISO-style mappings.
- Replaced static WebSocket admin-token query usage with short-lived one-time scoped stream tickets.
- Added fail-closed backend binding behavior for unsafe non-loopback exposure.
- Tightened frontend Browser Access Key validation and moved current interim storage to session scope.
- Documented the future target to replace browser-stored tokens with memory-only or HttpOnly secure session-cookie authentication.
- Documented SSH hardening expectations: Ed25519 client keys, Ed25519 host keys, and mandatory known_hosts validation.
- Routed the active Players-table mutation workflows through shared confirmation and required admin reason capture: Give Item, Inventory repair/delete, resource/spec actions, journey node actions, player move, and player admin actions.
- Kept Player 360 read-only while the high-risk player mutation surface was migrated into extracted confirmed modals.
- Added shared frontend confirmation support so future high-risk UI actions can display mutation safety metadata and capture an admin reason before sending the mutation request.
- Kept Player 360 quick actions blocked until existing high-risk workflows are migrated to shared confirmation and reason capture.
- Validated Player 360 v1 as a protected read-only support surface with no new player mutation paths.
- Added a read-only Player 360 launcher from Players without changing existing mutation workflows.
- Added a read-only Player 360 frontend path without adding new player mutation workflows.
- Added safe section-level error wording to the read-only Player 360 backend response to avoid exposing raw backend details.
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

### Validation required

- Run `go test -v ./...`.
- Run `go build` or `./update.ps1` / `./update.sh`.
- Run frontend typecheck, lint, and build.
- Validate Discord auth route/session tests through `go test ./...` and manually validate OAuth login/callback, session context, logout, and registered-user review with configured Discord OAuth.
- Validate WebSocket ticket behavior manually.
- Validate fail-closed non-loopback startup behavior.
- Validate Ed25519-only SSH client key behavior.
- Validate SSH known-host mismatch rejection.
- Validate strict backend admin-token behavior.
- Validate Gitleaks, govulncheck, gosec, Trivy, npm audit, and SBOM generation before release.

### Known issues

- Browser token storage remains JavaScript-readable during the active session. Future target is memory-only token handling or HttpOnly secure session-cookie authentication.
- Mutation reason enforcement is not yet defaulted to enabled for all high/destructive actions.
- CI workflow hardening, SBOM generation, and artifact attestations remain open.
- Committed binary provenance requires cleanup or formal signed/attested release process.

## Change entry template

```markdown
## [Version] - YYYY-MM-DD

### Security
- 

### Added
- 

### Changed
- 

### Fixed
- 

### Operations
- 

### Documentation
- 

### Validation
- 

### Known issues
- 
```
