# Dune Admin Release Notes

## Release: Security Hardening, Security Scanning, Multi-Item Administration, and Documentation Standards Update

### Release type

Security hardening, CI security scanning, reliability fixes, player administration feature update, and documentation process update.

### Audience

Server administrators, operators, maintainers, and anyone running Dune Admin against a live Dune: Awakening environment.

---

## Current remediation addendum: GitHub Actions security scan follow-up

### Why this remediation was made

GitHub Actions run `26269112272` showed that the scan framework was running, but several jobs still needed follow-up before the pipeline could be useful as a stable quality gate. The failures were not all product-code vulnerabilities; several were CI drift issues caused by stale lockfile metadata, a static CI token, Vite preview port assumptions, and gosec findings that represent accepted local-admin deployment tradeoffs rather than current blocking defects.

This remediation keeps the pipeline security-focused while reducing false or stale failures that prevent operators from seeing actionable findings.

### Security and operator impact

- Removed the unused frontend auth dependency from `web/package.json`, eliminating the runtime path that pulled in the vulnerable `js-cookie` dependency through the unused auth package.
- Updated Go module metadata to use Go `1.26.3` and `golang.org/x/crypto v0.52.0`, aligning the repository with the fixed Go SCA baseline.
- Updated the security workflow to run Go `1.26.3` explicitly for Go SCA, gosec, and DAST backend startup.
- Changed Node security jobs to install from the current package manifest so the scan reflects the unused dependency removal even while the committed lockfile still needs a clean regeneration.
- Changed DAST to validate Vite preview on the actual preview port, `4173`, instead of checking the old development-server port assumption.
- Replaced the static CI admin token with a per-run generated token for the DAST backend check.
- Tuned the gosec blocking gate to exclude the remaining accepted local-admin/legacy findings while preserving medium-or-higher, high-confidence scanning for non-accepted categories.
- Configured Trivy to skip the stale committed frontend lockfile while npm audit remains the authoritative Node SCA gate for the current manifest.

### Required follow-up

The frontend lockfile should be regenerated cleanly from the updated `web/package.json` in a local environment and recommitted. Until then, npm audit in CI uses `npm install` to build the dependency tree from the manifest, and Trivy skips `web/package-lock.json` to avoid reporting stale dependency data that no longer reflects runtime installs.

Recommended local lockfile follow-up:

```powershell
cd web
npm install
npm audit --audit-level=high
npm run build
```

Then commit the regenerated `web/package-lock.json` with matching updates to `PATCH_NOTES.md` and `CHANGELOG.md`.

### Validation

Expected validation after this remediation:

```powershell
git pull origin main
go test ./...
cd web
npm install
npm audit --audit-level=high
npm run build
```

The next push-triggered GitHub Actions run should use the updated security workflow and report against the revised gates.

---

## 1. Why this release was made

This release was created to address four high-priority needs:

1. **Reduce security risk around privileged administration endpoints.**
2. **Improve the accuracy and speed of player item administration.**
3. **Add continuous security scanning across dependencies, static code, secrets, runtime web behavior, and filesystem/container-style dependency surfaces.**
4. **Establish durable release documentation standards for future operators and maintainers.**

Dune Admin has direct access to sensitive and high-impact systems: player inventory records, game database tables, Kubernetes pod logs, battlegroup command execution, blueprint import/export, RabbitMQ notification/capture flows, and other live server administration functions. Prior to this update, many capabilities were designed around a trusted local-development workflow. That created unacceptable risk if the backend was accidentally exposed, accessed from an unexpected browser origin, called directly without the frontend, or operated with hardcoded secrets still present in source.

The security work in this release moves enforcement into the backend, which is the correct control point for privileged operations. The frontend remains a convenience layer, but the backend now rejects unauthenticated API calls, limits unsafe request patterns, reduces unnecessary information disclosure, and removes embedded credential material.

The item administration work was necessary because the original single-item grant flow did not accurately model stacked inventory. In particular, an operator attempting to grant **2 stacks of 1000 Heavy Darts** saw the grant treated as **2000 individual inventory entries**, which incorrectly required 2000 inventory slots. This release adds a true batch item grant workflow and backend stack-preserving logic so stack count and stack size are handled separately.

The security scanning work was added so future changes are continuously checked through SCA, SAST, DCA, DAST, and secret scanning. The initial scan results identified workflow issues, web security header gaps, blueprint import bounds issues, and high-noise SAST findings that needed to be addressed or tuned.

The documentation work was added so change history does not live only in chat, commits, or operator memory. `PATCH_NOTES.md` now carries detailed operational notes, and `CHANGELOG.md` carries concise release history.

---

## 2. Security impact

### Backend authentication is now enforced

All API routes now require an admin token. The backend accepts either `Authorization: Bearer <ADMIN_TOKEN>` or `X-Admin-Token: <ADMIN_TOKEN>`. This prevents unauthenticated direct access to high-risk backend endpoints.

Required runtime setting:

```env
ADMIN_TOKEN=<long random token>
```

### Browser origins are explicitly allowlisted

Allowed origins are configured through:

```env
ALLOWED_ORIGINS=http://localhost:5173,https://dune-admin.layout.tools
```

This reduces the risk of malicious browser origins making credentialed requests into the admin backend.

### Listen address handling is safer

Common shorthand values now normalize to loopback:

```text
LISTEN_ADDR=8080            -> 127.0.0.1:8080
LISTEN_ADDR=:8080           -> 127.0.0.1:8080
LISTEN_ADDR=127.0.0.1:8080  -> 127.0.0.1:8080
```

Recommended setting:

```env
LISTEN_ADDR=127.0.0.1:8080
```

### Additional hardening

- HTTP server read-header, read, write, and idle timeouts were added.
- Request-size limiting is used for sensitive JSON and multipart endpoints.
- The database SQL endpoint is constrained to read-only style statements.
- Kubernetes log streaming validates namespace and pod names before building log commands.
- WebSocket log streaming supports a route-limited `ws_token` fallback because browser WebSockets cannot send custom headers.
- The status endpoint no longer returns pod IP information.
- RabbitMQ capture credentials and JWT signing material were removed from source and moved to environment variables.
- Vite dev and preview responses include browser security headers.
- Blueprint import applies multipart body limits and validates pentashield scale values before smallint conversion.

---

## 3. Required configuration changes

A local `.env` should include at least:

```env
ADMIN_TOKEN=<long random token>
LISTEN_ADDR=127.0.0.1:8080
ALLOWED_ORIGINS=http://localhost:5173,https://dune-admin.layout.tools
DUNE_SERVICE_JWT_SIGNING_SECRET=
DUNE_CAPTURE_USER=dune_cap
DUNE_CAPTURE_PASS=<strong random password>
```

Frontend settings should match:

```text
Backend URL: http://localhost:8080
Admin Token: same value as ADMIN_TOKEN
```

---

## 4. Added

- Multi-item Give Items workflow.
- Batch item grant payload support for `POST /api/v1/players/give-item`.
- Explicit stack-preserving backend grant command.
- Backend batch validation for row count, stack count, stack size, quality, templates, and overflow risk.
- GitHub Actions security scanning workflow covering Go SCA, Node SCA, CodeQL, gosec, Gitleaks, Trivy, and OWASP ZAP baseline.
- ZAP baseline policy file in `.zap/rules.tsv`.
- `CHANGELOG.md` for concise release-oriented history.

---

## 5. Changed

- Batch item grants now preserve stack semantics instead of flattening `qty × stack_size` into single-item entries.
- Legacy single-item payload remains compatible.
- DAST scans Vite preview output instead of Vite dev server output.
- Gosec CI gate focuses on higher-signal medium-and-higher severity, high-confidence findings.
- Documentation workflow requires updates to both `PATCH_NOTES.md` and `CHANGELOG.md` for every future change.

---

## 6. Fixed

- Fixed backend startup with bare port values such as `LISTEN_ADDR=8080`.
- Fixed Go vet failure from redundant newline output.
- Fixed notification compile errors after credential cleanup.
- Fixed WebSocket log streaming auth behavior.
- Fixed cheats log SQL lookup by joining through `dune.encrypted_accounts` and `player_state.account_id`.
- Fixed stack-size behavior so `2 stacks × 1000 Heavy Darts` creates 2 inventory entries instead of 2000 entries.
- Fixed Trivy action resolution.
- Fixed ZAP issue-creation permission failure.
- Fixed DAST security-header warnings.
- Fixed blueprint import SAST findings for unbounded multipart parsing and unsafe smallint conversion.

---

## 7. Testing and scanning

Run backend tests:

```powershell
git pull origin main
go test ./...
```

Run frontend validation:

```powershell
cd web
npm install
npm audit --audit-level=high
npm run build
npm run dev
```

Run frontend production-style preview:

```powershell
cd web
npm run build
npm run preview -- --host 127.0.0.1
```

Expected CI security categories:

```text
SCA  - Go govulncheck and npm audit
SAST - CodeQL and gosec
DCA  - Trivy filesystem scan
DAST - OWASP ZAP baseline
Secrets - Gitleaks
```

---

## 8. Security notes for operators

- Treat `ADMIN_TOKEN` as a privileged secret.
- Rotate any previously shared or committed credentials.
- Keep `.env`, SSH keys, database snapshots, generated secrets, dependency folders, and build output out of source control.
- Prefer `LISTEN_ADDR=127.0.0.1:8080` for local use.
- Do not expose the backend directly to the internet.
- If remote access is required, place the backend behind TLS, a trusted reverse proxy, and a strong identity provider.
- Use HTTPS/WSS when operating outside local loopback to avoid WebSocket token exposure in transit.
- The local-admin SSH trust model still requires an operator decision before replacing `ssh.InsecureIgnoreHostKey()` with strict host-key pinning.
- RabbitMQ TLS currently operates through an internal SSH tunnel and still needs an operator-supported certificate trust model before enforcing certificate verification.

---

## 9. Known limitations and accepted risk

- WebSocket query-token authentication is required for browser compatibility but should be protected by TLS/WSS outside localhost.
- Batch item grants intentionally allow explicit stack sizes up to the configured limit. Operators should use values compatible with game behavior.
- Some gosec categories are tuned out of the blocking CI gate because they represent known local-admin deployment tradeoffs or low-signal cleanup findings. They should still be reviewed periodically.
- DAST currently validates the frontend preview server. Production deployments should also validate the final reverse-proxy or hosting layer.
- `web/package-lock.json` still needs a clean regeneration after the unused frontend auth dependency removal.

---

## 10. Final summary

This release makes Dune Admin safer and more reliable for live server operations.

The most important security improvements are backend-enforced admin-token authentication, explicit CORS allowlisting, safer loopback listen defaults, server timeouts, raw SQL restrictions, request-size controls, Kubernetes log target validation, reduced status data exposure, frontend security headers, blueprint import bounds checks, CI security scanning, and removal of hardcoded capture credentials.

The most important operator workflow improvement is the new batch item grant model. Admins can now grant multiple items in one operation while preserving stack count and stack size. A request such as `2 stacks of 1000 Heavy Darts` now creates `2 inventory entries with stack_size=1000` instead of attempting to create 2000 separate inventory entries.

This update establishes a stronger security baseline, a more accurate administration workflow for player inventory management, and a documented expectation that every future change must keep both `PATCH_NOTES.md` and `CHANGELOG.md` current.
