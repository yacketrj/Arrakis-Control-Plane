# Changelog

All notable changes to this project will be documented in this file.

This project follows a corporate change-management style informed by ITIL release/change practices and NIST SSDF secure-development expectations. Entries should distinguish security, operational, feature, documentation, validation, and known-risk changes.

## [Unreleased]

### Added

- Added `handlers_database_test.go` with database handler security coverage for bounded query parameters, control-character rejection, numeric function OID validation, unsafe SQL rejection, and redacted SQL/database result shape.
- Added `docs/database-endpoint-security.md` with the `ASEA-004` database endpoint security review state, current guardrails, added tests, and remaining work.
- Added AppSec Discord self-session route regression tests for normal Discord access to `me`, `logout`, and `/self/*` plus denial from representative admin routes.
- Added `appsec_auth_boundary_test.go` with AppSec auth-boundary regression coverage for public allowlist, self-service classification, representative admin-only routes, and WebSocket-ticket enforcement.
- Added initial `docs/appsec-endpoint-audit.md` with route inventory, auth-boundary summary, findings, remediation backlog, and manual abuse-case checklist.
- Added P0 comprehensive AppSec endpoint audit backlog task covering public and protected backend routes.
- Added future documentation requirement for `docs/appsec-endpoint-audit.md`.
- Added Inventory Studio stack-size update backend in `inventory_stack_size.go`.
- Added protected `POST /api/v1/players/item/stack-size` endpoint for selected item stack-size edits.
- Added `web/src/api/inventoryStudioMutations.ts` for Inventory Studio mutation helpers.
- Added confirmed stack-size edit control to the Inventory Studio selected-item panel.
- Added P5 Discord bot setup and usage guide task to the implementation tracker.
- Added P2 Guild Management backlog task covering create/delete guild, add/remove player membership, and guild-rank changes.
- Added P2 Player tab guild workflow task covering add/remove selected player from a guild and promote/change selected player rank from the Player tab.
- Added Inventory Studio browser-session action history for recent completed add, repair, and removal action diffs.
- Added Inventory Studio action-history export to local JSON.
- Added Discord/self-service frontend API helper at `web/src/api/discordSelfService.ts` with cookie-aware calls to `/api/v1/self/*` and admin-token support for link management.
- Added Discord Player Links admin tab at `web/src/tabs/DiscordPlayerLinksTab.tsx` for listing, creating, editing, and deleting Discord-to-player mappings.
- Added read-only My Player Card tab at `web/src/tabs/SelfPlayerCardTab.tsx` for linked Discord sessions.
- Added Discord Links and My Player Card navigation entries in `web/src/App.tsx`.
- Added Discord player link foundation in `discord_player_links.go` for admin-managed Discord ID to player actor ID mapping.
- Added protected Discord player link admin endpoints for list, upsert, and delete workflows.
- Added read-only self-service endpoints at `/api/v1/self/player-link` and `/api/v1/self/player-card` for linked Discord sessions.
- Added `discord_player_links_test.go` coverage for link validation, store helper behavior, handlers, current session link lookup, and self-service auth gating.
- Added `docs/discord-player-links.md` with storage, endpoint, auth-boundary, validation, and safety notes.
- Added Farming Requests frontend tab at `web/src/tabs/FarmingRequestsTab.tsx` for coordination-only request/order management.
- Added `web/src/api/inventoryRequests.ts` as a separate frontend API module for inventory request/order endpoints.
- Added Farming Requests navigation wiring in `web/src/App.tsx`.
- Added Discord bot command adapter skeleton in `discord_bot_adapter.go` for Discord-style personal requests, guild requests, farm orders, fill updates, and cancel updates.
- Added `discord_bot_adapter_test.go` coverage for personal request, guild request, farm-order, fill/cancel update, and unsupported-command adapter paths.
- Added inventory request/order backend coordination model in `inventory_requests.go` for personal and guild requests plus farming orders.
- Added `inventory_requests_test.go` coverage for request validation, handler lifecycle, order linking, fill propagation, and missing-request rejection.
- Added `docs/inventory-requests-orders.md` with storage, endpoint, model, frontend UI, validation, status propagation, and safety-boundary notes.
- Added Discord auth route/session coverage in `discord_auth_test.go` for route registration, role mapping, session lookup, expiry eviction, session hash generation, and logout invalidation.
- Added `docs/discord-auth.md` with runtime configuration, endpoint, role mapping, session behavior, and current limitation notes.
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

- Updated `docs/appsec-endpoint-audit.md` so `ASEA-004` is validated as partial remediation after database endpoint hardening and clean local validation.
- Updated `PATCH_NOTES.md` with validated database endpoint security hardening status.
- Updated database endpoint handlers to trim and bound query parameters, reject unsafe controls, require numeric function OIDs, trim manual SQL before validation, and redact sampled/search/manual SQL output.
- Updated `update.sh` to colorize validation output for `RUN`, `PASS`, `FAIL`, and `Update failed.` status lines.
- Updated `docs/appsec-endpoint-audit.md` so `ASEA-003` is validated as partial remediation after mutation-safety classification, audit metadata JSON, and colored-output validation.
- Updated `PATCH_NOTES.md` with validated mutation-safety classification and colored update-output status.
- Updated `docs/appsec-endpoint-audit.md` so `ASEA-002` is validated as remediated after clean local validation with a non-blocking Tailwind/Rolldown plugin timing warning.
- Updated `PATCH_NOTES.md` with validated Discord self-session route remediation status and the non-blocking plugin timing warning.
- Updated `docs/appsec-endpoint-audit.md` so `ASEA-001` is validated as partial remediation after clean local validation.
- Updated `PATCH_NOTES.md` with validated AppSec auth-boundary regression test status.
- Updated `PATCH_NOTES.md` with verified Farming Requests hook-dependency lint validation.
- Updated `docs/admin-implementation-tasks.md` so the AppSec endpoint audit is In Progress after creating the initial audit document.
- Updated `PATCH_NOTES.md` with the initial AppSec endpoint audit pass.
- Updated `docs/admin-implementation-tasks.md` with the P0 comprehensive AppSec endpoint audit backlog item and required audit scope.
- Updated `PATCH_NOTES.md` with AppSec endpoint audit backlog planning details.
- Updated `docs/inventory-studio.md` with stack-size edit endpoint, payload, UI behavior, safety behavior, and validation notes.
- Updated `docs/admin-implementation-tasks.md` so Inventory Studio stack-size edit is In Progress and quality edit is Next.
- Updated `PATCH_NOTES.md` with Inventory Studio stack-size edit workflow status.
- Updated `docs/admin-implementation-tasks.md` with Discord bot setup guide, guild management, and Player tab guild workflow backlog entries.
- Updated `PATCH_NOTES.md` with backlog planning additions for Discord bot documentation and guild-management workflows.
- Updated `docs/admin-implementation-tasks.md` so Inventory Studio browser-session action history is marked validated and done.
- Updated `PATCH_NOTES.md` with verified Inventory Studio action history validation.
- Updated `docs/inventory-studio.md` with action history behavior, export, clear, reset, and safety notes.
- Updated `docs/admin-implementation-tasks.md` so Inventory Studio post-action diff is marked done and stack-size edit is the next planned workflow after action-history validation.
- Updated `PATCH_NOTES.md` with Inventory Studio action history status.
- Updated `PATCH_NOTES.md` with manual Discord self-service frontend tab validation.
- Updated `PATCH_NOTES.md` with verified Discord self-service frontend tab validation.
- Updated `docs/discord-player-links.md` with Discord Links and My Player Card frontend behavior.
- Updated `web/src/App.tsx` tab gating so My Player Card can load with Discord session cookies while administrative tabs still require configured browser admin access.
- Updated `PATCH_NOTES.md` with manual Discord player link validation status.
- Updated `PATCH_NOTES.md` with verified Discord player link validation status.
- Updated auth middleware so normal registered Discord sessions can reach `/api/v1/self/*` only while admin-token and Discord-admin access remain required elsewhere.
- Registered Discord player link admin endpoints and read-only self-service endpoints in `routes.go`.
- Updated `PATCH_NOTES.md` with Discord player link foundation status.
- Updated `PATCH_NOTES.md` with verified Farming Requests UI validation status.
- Updated `docs/inventory-requests-orders.md` with Farming Requests frontend behavior and validation expectations.
- Updated `PATCH_NOTES.md` with the Farming Requests UI status.
- Updated `PATCH_NOTES.md` with the Discord bot command adapter status and the clean full local validation result.
- Registered protected inventory request/order endpoints in `routes.go`.
- Updated CORS middleware to allow `PATCH` for request/order update endpoints.
- Updated `PATCH_NOTES.md` with the inventory request/order backend status.
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

- Fixed `audit_log_test.go` so the audit metadata payload uses valid JSON while still testing newline sanitization.
- Fixed Farming Requests tab `react-hooks/exhaustive-deps` warning by stabilizing `load` with `useCallback` and depending on `load` in the reload effect.
- Fixed Discord Links tab type usage by importing `CSSProperties` directly from React.
- Fixed Discord-player link text validation so raw control characters are rejected before trimming.
- Fixed Farming Requests tab type usage by importing `CSSProperties` directly from React and removing an unused memoized open-request list.
- Fixed frontend lint failure from `no-control-regex` by scoping the rule exception to `web/src/api/client.ts`, where browser access-key validation intentionally rejects whitespace and control characters.
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

- Hardened database endpoints with bounded/control-character checked query parameters, numeric function OID validation, manual SQL trimming, and redaction for sampled/search/manual SQL output.
- Validated `ASEA-004` as partial remediation through clean local `./update.sh` validation; SQL timeout review, expanded read-only bypass tests, live-data redaction review, and manual abuse-case validation remain open.
- Tightened high-risk mutation classification for reconnect, database SQL, log stream ticket issuance, notify, and direct item-row edits.
- Validated `ASEA-003` as partial remediation through clean local `./update.sh` validation; full endpoint-by-endpoint audit-event assertion coverage remains open.
- Allowed registered non-admin Discord sessions to access only `GET /api/v1/auth/discord/me`, `POST /api/v1/auth/discord/logout`, and `/api/v1/self/*`; all admin review, player, database, infrastructure, and mutation routes remain admin-only.
- Validated `ASEA-002` as remediated through clean local `./update.sh` validation; the build emitted a non-blocking Tailwind/Rolldown plugin timing warning.
- Added AppSec auth-boundary regression tests for public route allowlisting, self-service path classification, representative admin-only routes, and WebSocket-ticket denial.
- Validated `ASEA-001` as partial remediation through clean local `./update.sh` validation; generated full-route auth-boundary coverage remains a future hardening follow-up.
- Added initial AppSec route inventory and auth-boundary audit document for all registered endpoints.
- Added AppSec findings ASEA-001 through ASEA-006 covering auth-boundary tests, Discord session UX, mutation audit/reason coverage, database endpoint review, infrastructure/log endpoint review, and browser token/CORS follow-up.
- Added AppSec endpoint audit planning requirement for all public and protected backend endpoints.
- Scoped the future endpoint audit to auth boundaries, input validation, request limits, CORS/session behavior, mutation safety, audit reason coverage, SQL injection, command execution, WebSocket/log-stream behavior, data exposure, rate limits, abuse cases, frontend helper behavior, and remediation evidence.
- Kept Inventory Studio stack-size edit behind shared mutation confirmation, admin reason capture, before-action snapshot export, post-action diff, and action-history capture.
- Limited the new backend stack-size mutation to `dune.items.stack_size` with item ID and `1..9999` stack-size validation.
- Added guild-management planning requirements that future guild mutations must use schema discovery, confirmation, admin reason capture, before-change review, refresh/diff where practical, and audit visibility.
- Kept Inventory Studio action history browser-local; it does not add backend persistence, rollback automation, or new mutation routes.
- Kept existing Inventory Studio confirmed add, repair, and removal workflows behind shared mutation confirmation, before-action snapshot export, and admin reason capture.
- Kept Discord self-service frontend read-only; My Player Card calls only `/api/v1/self/player-link` and `/api/v1/self/player-card`.
- Kept Discord Links as an admin-management surface for existing backend-protected link APIs.
- Added Discord player link foundation as a prerequisite for future self-service and kept it read-only for normal Discord sessions.
- Scoped normal registered Discord sessions to `/api/v1/self/*` only.
- Kept Discord player link management behind admin token or Discord admin session.
- Kept self-service player card output read-only and derived from the existing Player 360 profile builder.
- Kept Farming Requests UI coordination-only; it does not write player inventory, guild storage, claim rewards, Player 360, or game-state tables.
- Kept Farming Requests frontend API separate from the high-risk player/admin API surface.
- Kept Discord bot command adapter non-networked and dependency-free for this slice; it does not register slash commands, connect to Discord, or execute runtime actions on its own.
- Kept Discord bot command adapter mapped to the coordination-only inventory request/order model instead of any direct game-state mutation path.
- Kept inventory request/order backend coordination-only; it does not mutate player inventory, guild storage, claim rewards, Player 360, or game-state tables.
- Serialized in-process access to the local inventory request/order JSON store and retained `0600` file permissions.
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

### Validation

- Validated database endpoint security hardening from the canonical local update path:
  - `./update.sh`
- Validated mutation-safety classification coverage, audit metadata JSON fix, and colored update output from the canonical local update path:
  - `./update.sh`
- Validated Discord self-session route remediation from the canonical local update path:
  - `./update.sh`
  - Non-blocking build-performance warning observed: `[PLUGIN_TIMINGS] Your build spent significant time in plugin @tailwindcss/vite:generate:build`.
- Validated AppSec auth-boundary regression tests from the canonical local update path:
  - `./update.sh`
- Validated clean canonical update sequence after the Farming Requests hook-dependency lint fix:
  - `./update.sh`
- Validated Inventory Studio stack-size edit code/build gate from the canonical local update path:
  - `./update.sh`
- Recorded initial AppSec endpoint audit document as documentation/audit review only.
- Recorded AppSec endpoint audit backlog addition as documentation/planning review only.
- Recorded backlog planning additions for Discord bot documentation and guild-management workflows as documentation/planning review only.
- Validated Inventory Studio action-history changes from the canonical local update path:
  - `./update.sh`
- Validated Inventory Studio action-history manual release checks:
  - add, repair, and removal history append behavior
  - reset on player change
  - JSON export
  - clear-history behavior
- Validated clean local full update sequence after the scoped frontend lint fix:
  - `go test -v ./...`
  - backend Windows build
  - `npm install`
  - `npm audit --audit-level=high`
  - `npm run typecheck`
  - `npm run lint`
  - `npm run build`
- Validated Farming Requests UI frontend gates from the local checkout after the new tab/API/navigation changes:
  - `npm run typecheck`
  - `npm run lint`
  - `npm run build`
- Validated Discord player link backend gates from the local checkout after the raw control-character validation fix:
  - `go test ./...`
  - `go build ./...`
- Validated Discord player link manual release checks:
  - admin link CRUD
  - normal Discord `/api/v1/self/*` access
  - normal Discord denial from admin paths
  - unlinked Discord safe failures
  - read-only self player-card behavior
- Validated Discord self-service frontend tabs from the canonical local update path:
  - `./update.sh`
- Validated Discord self-service frontend manual release checks:
  - Discord Links tab list/create/edit/delete behavior
  - My Player Card through Discord session cookies without a Browser Access Key

### Validation still required before release

- Manually validate selected-item stack-size edit, unchanged-value guard, before-action snapshot export, required reason capture, post-action diff, action-history append, and inventory reload behavior.
- Complete comprehensive AppSec endpoint audit, including SAST, DAST, dependency review, handler-by-handler review, and manual abuse-case validation.
- Manually exercise Farming Requests UI list, create, group, fill, and cancel workflows.
- Manually exercise inventory request/order personal/guild requests, order creation, fill/cancel propagation, and `PATCH` browser preflight.
- Manually validate Discord OAuth login/callback, session context, logout, and registered-user review with configured Discord OAuth.
- Validate WebSocket ticket behavior manually.
- Validate fail-closed non-loopback startup behavior.
- Validate Ed25519-only SSH client key behavior.
- Validate SSH known-host mismatch rejection.
- Validate strict backend admin-token behavior.
- Validate Gitleaks, govulncheck, gosec, Trivy, SBOM generation, and artifact attestations before release.

### Known issues

- Discord player link storage is local JSON and not yet a durable multi-instance database-backed identity mapping.
- Local `update.sh` auto-commit can fail when Git `user.name` and `user.email` are not configured in the checkout or globally.
- Inventory request/order storage is local JSON and not yet a durable multi-instance database-backed ledger.
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
