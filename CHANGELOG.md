# Changelog

All notable changes to this project will be documented in this file.

## [Unreleased]

### Added

- Added this changelog to provide concise release-oriented change tracking alongside the more detailed `PATCH_NOTES.md`.
- Added a documentation policy requiring every future change to update both `PATCH_NOTES.md` and `CHANGELOG.md`.

### Changed

- Documentation workflow now separates detailed operator/security patch notes from concise release change history.
- Updated the security scan workflow to use Go 1.26.3 explicitly for SCA, SAST, and DAST jobs.
- Updated the standalone Go Test workflow to use Go 1.26.3 explicitly and refresh modules before tests.
- Updated Go workflow bootstrap steps from `go mod download` to `go mod tidy` so CI refreshes package import checksums needed by `go test`, gosec, and DAST backend startup.
- Updated DAST backend startup polling to wait up to 30 seconds and print backend logs on failure.
- Updated DAST frontend preview polling to wait up to 30 seconds and print frontend logs on failure.
- Updated Node dependency installation in security jobs to install from the current package manifest after removing the unused frontend auth dependency.
- Updated the ZAP DAST job to scan Vite preview on port 4173 instead of the previous dev-server port assumption.
- Updated the CI-generated admin token to use a per-run random value instead of a static literal.
- Updated the gosec gate to exclude the remaining accepted local-admin findings from the blocking high-confidence scan.
- Removed npm cache key references to the deleted frontend lockfile from security workflow Node setup steps.
- Restored Trivy DCA to scan the repository without a stale lockfile skip because the stale frontend lockfile was removed.
- Removed npm cache mode from Node setup steps while no committed frontend lockfile is present.
- Removed the optional Clerk-authenticated UI path so the frontend matches backend-token-only authentication.
- Made the Blueprints tab available through the backend-token protected app rather than a removed frontend identity provider gate.
- Updated `.zap/rules.tsv` to use ZAP baseline's required three-column rule format.
- Expanded the Vite frontend Content Security Policy with explicit fallback directives required by ZAP baseline.
- Added the Cross-Origin-Embedder-Policy header for Vite dev and preview responses.

### Fixed

- Removed the unused frontend auth dependency from `web/package.json` to eliminate the transitive vulnerable `js-cookie` dependency from runtime installs and npm audit.
- Removed stale `web/package-lock.json` because it still contained vulnerable `js-cookie` metadata after the unused auth dependency was removed from `web/package.json`.
- Removed stale `@clerk/react` imports from `web/src/App.tsx` and `web/src/main.tsx` so the frontend build succeeds after dependency removal.
- Fixed the Go SCA baseline by moving the repository module metadata to Go 1.26.3 and `golang.org/x/crypto` 0.52.0.
- Fixed Go SCA and gosec failures caused by missing `go.sum` entries after the Go module version update.
- Fixed the standalone Go Test workflow failure caused by running `go test ./...` before refreshing Go module sums.
- Fixed SCA npm and DAST setup-node failures caused by enabling npm cache without a committed lockfile.
- Fixed DAST preview startup validation by checking the actual Vite preview port.
- Fixed DAST ZAP policy parsing failure caused by two-column rule entries in `.zap/rules.tsv`.
- Fixed DAST ZAP warnings for missing Cross-Origin-Embedder-Policy and incomplete CSP fallback directives.

### Security

- Reduced dependency risk by removing unused frontend auth package exposure.
- Reduced CI secret exposure by generating the DAST admin token at runtime.
- Prevented Trivy DCA from reporting stale dependency metadata by removing the outdated lockfile rather than suppressing the file in the scan.
- Kept Node vulnerability detection in npm audit against the current manifest-derived dependency tree.
- Reduced frontend authentication ambiguity by removing optional Clerk code paths and relying on backend-enforced admin-token authorization.
- Preserved DAST signal by fixing ZAP rule parsing instead of disabling the ZAP baseline job.
- Strengthened browser isolation and CSP behavior for CI preview and local Vite serving.

### Operational Notes

- Regenerate `web/package-lock.json` locally from the current `web/package.json` with `npm install` and recommit it once confirmed clean with `npm audit --audit-level=high` and `npm run build`.

---

## [Security Hardening, Security Scanning, and Multi-Item Administration Update] - 2026-05-21

### Added

- Added backend admin-token authentication for all API routes.
- Added support for `Authorization: Bearer <ADMIN_TOKEN>` and `X-Admin-Token: <ADMIN_TOKEN>` authentication methods.
- Added explicit browser origin allowlisting through `ALLOWED_ORIGINS`.
- Added safer listen-address normalization for local loopback use.
- Added HTTP server read-header, read, write, and idle timeouts.
- Added request-size limiting for sensitive JSON and multipart endpoints.
- Added read-only SQL enforcement for the database SQL endpoint.
- Added Kubernetes namespace and pod-name validation for log streaming.
- Added WebSocket route-specific query-token support for browser log streaming.
- Added frontend security headers for Vite dev and preview responses.
- Added environment-variable based capture and notification credential handling.
- Added multi-item Give Items workflow in the frontend.
- Added batch item grant payload support to `POST /api/v1/players/give-item`.
- Added explicit stack-preserving item grant backend command.
- Added batch give-item validation for row count, stack count, stack size, quality range, blank templates, and overflow risk.
- Added Go unit tests for batch item request normalization.
- Added GitHub Actions backend Go test workflow.
- Added GitHub Actions security scanning workflow for SCA, SAST, secret scanning, DCA, and DAST.
- Added Go SCA with `govulncheck`.
- Added Node SCA with `npm audit --audit-level=high`.
- Added SAST with GitHub CodeQL.
- Added Go SAST with `gosec`.
- Added secret scanning with Gitleaks.
- Added DCA-style filesystem/dependency scanning with Trivy.
- Added DAST with OWASP ZAP baseline scanning against frontend preview output.
- Added ZAP baseline rule policy in `.zap/rules.tsv`.
- Added multipart import bounds for blueprint uploads.
- Added safe pentashield scale validation before `int` to `int16` conversion.

### Changed

- Changed batch item grant semantics so `qty` means stack count and `stack_size` means items per stack.
- Changed batch item grants to preserve requested stack size instead of flattening `qty × stack_size` into single-item entries.
- Changed legacy single-item payload behavior to remain compatible with existing callers.
- Changed notification publishing to use configured capture credentials instead of removed hardcoded constants.
- Changed frontend settings panel to support both backend URL and admin token configuration.
- Changed status endpoint behavior to avoid returning pod IP information.
- Changed DAST workflow to scan `npm run preview` output instead of the Vite development server.
- Changed security workflow backend health check to use the configured admin token.
- Changed ZAP workflow behavior to avoid GitHub issue-creation permission failures.
- Changed gosec CI gate to focus on higher-signal medium-and-higher severity, high-confidence findings.
- Changed Trivy workflow action reference to the valid `aquasecurity/trivy-action@v0.36.0` tag.
- Changed blueprint import to validate request size before multipart parsing.
- Changed blueprint import to check batch close errors and JSON export encoding errors.

### Fixed

- Fixed backend startup failure when `LISTEN_ADDR` was configured as a bare port such as `8080`.
- Fixed Go vet failure caused by a redundant newline in `fmt.Println` output.
- Fixed notification compile errors after removing hardcoded capture credential constants.
- Fixed WebSocket log streaming failures caused by browser inability to send custom WebSocket headers.
- Fixed cheats log SQL failure caused by referencing nonexistent `player_state.fls_id`.
- Fixed cheats log player-name resolution by joining through `dune.encrypted_accounts.encrypted_funcom_id` and `player_state.account_id`.
- Fixed batch stack-size behavior so `2 stacks × 1000 Heavy Darts` creates 2 inventory entries instead of attempting to create 2000 entries.
- Fixed Trivy action resolution failure in GitHub Actions.
- Fixed ZAP permission failure caused by attempted GitHub issue creation from the scan action.
- Fixed DAST security-header warnings for major missing browser security headers.
- Fixed blueprint import SAST findings for unbounded multipart parsing and unsafe smallint conversion.
- Fixed `.gitignore` coverage for frontend dependency directories and build output.

### Security

- Enforced backend authorization as the primary security control for administrative endpoints.
- Reduced CORS exposure by requiring explicit allowed origins.
- Reduced accidental network exposure by defaulting shorthand listen addresses to loopback.
- Reduced data exposure in status responses.
- Removed hardcoded RabbitMQ capture password material and JWT signing material from source.
- Added request body limiting to reduce oversized request risks.
- Added raw SQL restrictions to reduce accidental destructive database operations.
- Hardened log streaming target validation and WebSocket origin handling.
- Added frontend browser security headers for CI and local preview scanning.
- Added blueprint import bounds and numeric conversion validation.
- Established continuous security scanning coverage across dependencies, static code, secrets, runtime web behavior, and filesystem/dependency surfaces.

### Testing

- Added unit coverage for legacy give-item payload compatibility.
- Added unit coverage for batch give-item parsing and validation.
- Added unit coverage for stack-size defaulting and boundary values.
- Added unit coverage for empty payloads, blank templates, invalid quantities, invalid quality values, excessive row counts, excessive stack counts, and excessive stack sizes.
- Added Go test CI.
- Added security scan CI jobs for SCA, SAST, DCA, DAST, and secret scanning.

### Operational Notes

- Operators must configure `ADMIN_TOKEN` for backend API access.
- Operators should use `LISTEN_ADDR=127.0.0.1:8080` for local operation.
- Operators must set the same admin token in the frontend settings panel.
- Operators should rotate any previously committed or shared credentials.
- Operators should keep `.env`, SSH keys, database snapshots, dependency folders, build output, and generated secrets out of source control.
- Operators should use HTTPS/WSS and avoid URL logging if WebSocket token auth is used outside localhost.
- Security scans run on pushes to `main` and can be started manually through `workflow_dispatch`.

### Known Limitations

- WebSocket query-token authentication is required for browser compatibility but should be protected by TLS/WSS outside localhost.
- Some gosec categories remain accepted risk or require larger architecture decisions, including SSH host-key pinning and RabbitMQ certificate verification inside the current SSH-tunnel model.
- DAST currently validates the frontend preview server. Production deployments should also validate the final reverse-proxy or hosting layer.
- Batch item grants allow explicit stack sizes up to the configured limit; operators should use values compatible with game behavior.
