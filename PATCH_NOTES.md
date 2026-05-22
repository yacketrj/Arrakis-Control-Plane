# Dune Admin Release Notes

## Release: Security Hardening, Security Scanning, Multi-Item Administration, and Documentation Standards Update

### Release type

Security hardening, CI security scanning, reliability fixes, player administration feature update, and documentation process update.

### Audience

Server administrators, operators, maintainers, and anyone running Dune Admin against a live Dune: Awakening environment.

### Documentation standard going forward

Every future code change, configuration change, security remediation, workflow update, behavior change, bug fix, or operator-facing update must include matching updates to both:

- `PATCH_NOTES.md`
- `CHANGELOG.md`

`PATCH_NOTES.md` must explain why the change was made, the security/operator impact, required configuration changes, validation steps, and known limitations or accepted risk.

`CHANGELOG.md` must provide concise release-oriented change tracking using standard categories such as:

- `Added`
- `Changed`
- `Fixed`
- `Security`
- `Testing`
- `Operational Notes`
- `Known Limitations`

This release adds `CHANGELOG.md` so future updates have both detailed patch notes and concise historical release tracking.

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

All API routes now require an admin token. The backend accepts either of the following authentication methods:

```http
Authorization: Bearer <ADMIN_TOKEN>
X-Admin-Token: <ADMIN_TOKEN>
```

This prevents unauthenticated direct access to high-risk backend endpoints.

Required runtime setting:

```env
ADMIN_TOKEN=<long random token>
```

### Browser origins are explicitly allowlisted

The backend no longer relies on permissive wildcard CORS behavior. Allowed origins are configured through:

```env
ALLOWED_ORIGINS=http://localhost:5173,https://dune-admin.layout.tools
```

This reduces the risk of malicious browser origins making credentialed requests into the admin backend.

### Listen address handling is safer

The backend now normalizes local listen-address values so common shorthand values stay bound to loopback.

```text
LISTEN_ADDR=8080            -> 127.0.0.1:8080
LISTEN_ADDR=:8080           -> 127.0.0.1:8080
LISTEN_ADDR=127.0.0.1:8080  -> 127.0.0.1:8080
```

Recommended setting:

```env
LISTEN_ADDR=127.0.0.1:8080
```

### HTTP server timeouts were added

The backend HTTP server now configures explicit read-header, read, write, and idle timeouts. This lowers exposure to slow-client and connection-resource exhaustion behavior.

### Request body limits were added

Request-size limiting is now used for sensitive JSON and multipart endpoints. This reduces exposure to oversized request bodies and memory exhaustion behavior.

### Raw SQL is restricted

The database SQL endpoint is constrained to read-only style statements. Allowed prefixes include:

```text
SELECT
WITH
SHOW
EXPLAIN
```

The backend rejects semicolon-separated statements and destructive SQL keywords such as:

```text
INSERT, UPDATE, DELETE, DROP, ALTER, TRUNCATE, CREATE, GRANT, REVOKE,
COPY, CALL, DO, EXECUTE, MERGE, VACUUM, ANALYZE
```

This does not replace proper database permissions, but it reduces the chance of destructive SQL being executed through the UI by mistake.

### Kubernetes log streaming was hardened

Log streaming now validates namespace and pod names before building a `kubectl logs` command. WebSocket origin checks also use the configured origin allowlist.

### WebSocket authentication was fixed

Browser WebSocket connections cannot set custom request headers like `X-Admin-Token`. The backend now supports a WebSocket-only query-token fallback for the log stream endpoint:

```text
/api/v1/logs/stream?...&ws_token=<ADMIN_TOKEN>
```

This fallback is limited to the log stream WebSocket route and is intended to support browser-based log streaming while keeping backend authentication enforced.

### Status data was reduced

The status endpoint no longer returns pod IP information. It still reports basic connectivity state, but avoids returning unnecessary internal network details.

### Hardcoded capture credentials were removed

RabbitMQ capture credentials and JWT signing material were removed from source and moved to environment variables:

```env
DUNE_SERVICE_JWT_SIGNING_SECRET=
DUNE_CAPTURE_USER=dune_cap
DUNE_CAPTURE_PASS=<strong random password>
```

Operators should rotate any previously committed or shared credential values.

### Frontend security headers were added

Vite dev and preview now set security headers to reduce browser-side risk during local and CI validation:

```text
X-Frame-Options
X-Content-Type-Options
Referrer-Policy
Permissions-Policy
Cross-Origin-Opener-Policy
Cross-Origin-Resource-Policy
Content-Security-Policy
Cache-Control
```

These headers specifically address DAST findings for clickjacking protection, content-type sniffing, CSP, and browser permission scoping.

### Blueprint import bounds were hardened

Blueprint import now applies multipart body limits before parsing and validates pentashield scale values before converting to PostgreSQL `smallint[]`. This addresses memory-exhaustion risk and integer-conversion overflow risk identified by SAST.

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

For local development, run the backend and frontend separately:

```powershell
# Terminal 1
cd "Z:\Unreal Projects\Icarus\dune-admin-fork"
go run .
```

```powershell
# Terminal 2
cd "Z:\Unreal Projects\Icarus\dune-admin-fork\web"
npm run dev
```

Open:

```text
http://localhost:5173
```

---

## 4. Added

### Multi-item Give Items workflow

The player Give Item workflow now supports multiple item rows in one operation. Each row can specify item template, stack count, grade/quality, and stack size.

### Batch item grant payload

The existing endpoint remains:

```http
POST /api/v1/players/give-item
```

It now supports a batch payload:

```json
{
  "player_id": 123,
  "items": [
    {
      "template": "ItemTemplateHeavyDarts",
      "qty": 2,
      "quality": 1,
      "stack_size": 1000
    }
  ]
}
```

In the batch payload:

```text
qty        = number of stacks
stack_size = number of items per stack
```

Example:

```text
qty=2, stack_size=1000 -> 2 inventory slots, 2000 total items
```

### Explicit stack grant backend command

A backend command was added for stack-preserving item grants:

```text
cmdGiveItemStacks(playerID, template, stacks, stackSize, quality)
```

### Backend batch validation

Batch grants now validate maximum row count, maximum stack count, maximum stack size, quality range, blank templates, empty item lists, non-positive quantities, and total quantity overflow risks.

### Security scanning workflow

A GitHub Actions workflow was added for continuous security validation:

```text
.github/workflows/security-scans.yml
```

The workflow includes:

- Go SCA through `govulncheck`.
- Node SCA through `npm audit --audit-level=high`.
- SAST through GitHub CodeQL.
- Go SAST through `gosec`.
- Secret scanning through Gitleaks.
- DCA-style filesystem/dependency scanning through Trivy.
- DAST through OWASP ZAP baseline scanning against the built frontend preview server.

### ZAP baseline policy

A ZAP rule policy file was added:

```text
.zap/rules.tsv
```

The policy suppresses known local-preview noise while keeping meaningful browser security-header checks active.

### Changelog

A new `CHANGELOG.md` was added to provide concise release-oriented historical tracking alongside detailed patch notes.

---

## 5. Changed

### Give Item semantics for batch payloads

Batch item grants now preserve stack semantics instead of flattening `qty × stack_size` into a single quantity.

Before:

```text
2 stacks × 1000 Heavy Darts -> 2000 inventory slots requested
```

After:

```text
2 stacks × 1000 Heavy Darts -> 2 inventory slots requested
```

Volume / weight checks still use total item count, so stack grants continue to respect inventory capacity rules beyond slot count.

### Legacy single-item payload remains compatible

The old single-item payload still works and remains flat quantity-based for compatibility with existing callers.

### DAST now scans frontend preview instead of dev server

The security workflow now starts the built frontend with:

```powershell
npm run preview -- --host 127.0.0.1
```

This better represents production build output than scanning the Vite development server.

### Gosec gate tuned for high-signal findings

The gosec job now focuses on medium-and-higher severity, high-confidence findings and excludes accepted legacy/local-admin noise categories from the blocking gate.

### Documentation workflow

Future changes must update both `PATCH_NOTES.md` and `CHANGELOG.md`. Patch notes remain detailed and operator-focused; the changelog is concise and release-history-focused.

---

## 6. Fixed

### Fixed backend startup with bare port values

`LISTEN_ADDR=8080` previously caused:

```text
listen tcp: address 8080: missing port in address
```

The backend now normalizes bare ports to loopback host-port form.

### Fixed Go vet failure in capture output

A redundant newline in a `fmt.Println` call caused `go test ./...` to fail during vet checks. The output was cleaned up.

### Fixed notification compile error after credential cleanup

After removing hardcoded capture constants, notification publishing still referenced the old names. It now uses the configured capture credential helpers.

### Fixed WebSocket log streaming authentication

Log streaming previously failed after backend authentication was introduced because browser WebSockets could not send the admin token header. The backend now supports the route-limited `ws_token` fallback.

### Fixed cheats log SQL error

The cheats log query previously joined against a non-existent column:

```sql
ps.fls_id
```

This caused:

```text
ERROR: column ps.fls_id does not exist (SQLSTATE 42703)
```

The query now joins through `dune.encrypted_accounts` using `encrypted_funcom_id`, then resolves player names through `player_state.account_id`.

### Fixed stack-size behavior for batch grants

Batch item grants no longer require one inventory slot per item when the operator specifies stack size. They now create the requested number of stacks with the requested stack size.

### Fixed Trivy action resolution

The Trivy action reference was corrected to:

```yaml
uses: aquasecurity/trivy-action@v0.36.0
```

### Fixed ZAP issue-creation permission failure

The ZAP action no longer attempts to create GitHub issues during CI. This avoids `Resource not accessible by integration` failures while preserving scan output and failure behavior.

### Fixed DAST security-header warnings

Frontend preview responses now include security headers that address the major OWASP ZAP baseline warnings for missing anti-clickjacking, content-type, CSP, and permissions policy headers.

### Fixed blueprint import SAST findings

Blueprint import now limits request bodies before multipart parsing and validates pentashield scale values before converting `int` to `int16`.

---

## 7. Testing and scanning

Run backend tests:

```powershell
git pull origin main
go test ./...
```

Run backend locally:

```powershell
go run .
```

Run frontend locally:

```powershell
cd web
npm install
npm run build
npm run dev
```

Run frontend production-style preview locally:

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

Recommended manual validation:

1. Confirm backend starts on `127.0.0.1:8080`.
2. Confirm unauthenticated API requests are rejected.
3. Confirm authenticated `/api/v1/status` succeeds with `X-Admin-Token`.
4. Confirm the frontend works after saving `ADMIN_TOKEN` in the settings panel.
5. Confirm Logs -> Cheats (7d) loads without the `ps.fls_id` SQL error.
6. Confirm pod log streaming connects without WebSocket auth errors.
7. Confirm granting 2 stacks of 1000 Heavy Darts creates 2 inventory slots, not 2000.
8. Confirm invalid batch item rows are rejected with clear validation errors.
9. Confirm blueprint import rejects oversized uploads and out-of-range pentashield scale values.
10. Confirm GitHub Actions security workflow resolves all actions and starts all scan jobs.
11. Confirm `CHANGELOG.md` and `PATCH_NOTES.md` are both updated for every future change.

---

## 8. Security notes for operators

- Treat `ADMIN_TOKEN` as a privileged secret.
- Rotate any previously shared or committed credentials.
- Keep `.env`, SSH keys, database snapshots, generated secrets, dependency folders, and build output out of source control.
- Prefer `LISTEN_ADDR=127.0.0.1:8080` for local use.
- Do not expose the backend directly to the internet.
- If remote access is required, place the backend behind TLS, a trusted reverse proxy, and a strong identity provider.
- The WebSocket `ws_token` is placed in the URL because browsers cannot send custom WebSocket headers. Use HTTPS/WSS when operating outside local loopback to avoid token exposure in transit.
- The local-admin SSH trust model still requires an operator decision before replacing `ssh.InsecureIgnoreHostKey()` with strict host-key pinning.
- RabbitMQ TLS currently operates through an internal SSH tunnel and still needs an operator-supported certificate trust model before enforcing certificate verification.

---

## 9. Known limitations and accepted risk

- WebSocket query-token authentication is required for browser compatibility but should be protected by TLS/WSS outside localhost.
- Batch item grants intentionally allow explicit stack sizes up to the configured limit. Operators should use values compatible with game behavior.
- Some gosec categories are tuned out of the blocking CI gate because they represent known local-admin deployment tradeoffs or low-signal cleanup findings. They should still be reviewed periodically.
- DAST currently validates the frontend preview server. Production deployments should also validate the final reverse-proxy or hosting layer.

---

## 10. Final summary

This release makes Dune Admin safer and more reliable for live server operations.

The most important security improvements are backend-enforced admin-token authentication, explicit CORS allowlisting, safer loopback listen defaults, server timeouts, raw SQL restrictions, request-size controls, Kubernetes log target validation, reduced status data exposure, frontend security headers, blueprint import bounds checks, CI security scanning, and removal of hardcoded capture credentials.

The most important operator workflow improvement is the new batch item grant model. Admins can now grant multiple items in one operation while preserving stack count and stack size. A request such as:

```text
2 stacks of 1000 Heavy Darts
```

now creates:

```text
2 inventory entries with stack_size=1000
```

instead of attempting to create 2000 separate inventory entries.

This update establishes a stronger security baseline, a more accurate administration workflow for player inventory management, and a documented expectation that every future change must keep both `PATCH_NOTES.md` and `CHANGELOG.md` current.