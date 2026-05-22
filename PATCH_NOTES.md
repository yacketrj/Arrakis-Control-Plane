# Dune Admin Release Notes

## Current remediation addendum: GitHub Actions security scan follow-up

### Why this remediation was made

GitHub Actions run `26277991316` passed every security job except DAST. The remaining DAST finding was a ZAP CSP warning for inline style allowance in the Vite preview response.

This remediation keeps the DAST gate enabled and tightens the frontend policy instead of suppressing the scan.

### Security and operator impact

- Updated `web/vite.config.ts` so `style-src` and `style-src-elem` load styles from `self` only.
- Kept `style-src-attr` compatible with the current React UI because the app still uses style attributes extensively.
- Preserved the earlier self-only script policy, browser isolation headers, ZAP baseline rules, and Vite preview security headers.
- Continued using the full security pipeline for SCA, SAST, DCA, DAST, and secret scanning.

### Validation

Expected local validation:

```powershell
git pull origin main
go mod tidy
go test ./...
cd web
npm install
npm audit --audit-level=high
npm run build
npm run preview -- --host 127.0.0.1
```

The next push-triggered workflow should rerun the complete security scan suite.

### Required follow-up

Regenerate `web/package-lock.json` from the current manifest when working locally:

```powershell
cd web
npm install
npm audit --audit-level=high
npm run build
```

Then commit the regenerated lockfile with matching updates to `PATCH_NOTES.md` and `CHANGELOG.md`.

---

## Release: Security Hardening, Security Scanning, Multi-Item Administration, and Documentation Standards Update

### Release type

Security hardening, CI security scanning, reliability fixes, player administration feature update, and documentation process update.

### Audience

Server administrators, operators, maintainers, and anyone running Dune Admin against a live Dune: Awakening environment.

---

## Summary

This release makes Dune Admin safer and more reliable for live server operations. It establishes backend-enforced admin-token authentication, explicit CORS allowlisting, safer loopback listen defaults, server timeouts, raw SQL restrictions, request-size controls, Kubernetes log target validation, reduced status data exposure, hardened frontend security headers, blueprint import bounds checks, CI security scanning, removal of hardcoded capture credentials, and the requirement that every future change keep both `PATCH_NOTES.md` and `CHANGELOG.md` current.

---

## Added

- Multi-item Give Items workflow.
- Batch item grant payload support for `POST /api/v1/players/give-item`.
- Explicit stack-preserving backend grant command.
- Backend batch validation for row count, stack count, stack size, quality, templates, and overflow risk.
- GitHub Actions security scanning workflow covering Go SCA, Node SCA, CodeQL, gosec, Gitleaks, Trivy, and OWASP ZAP baseline.
- ZAP baseline policy file in `.zap/rules.tsv`.
- `CHANGELOG.md` for concise release-oriented history.

---

## Changed

- Batch item grants now preserve stack semantics instead of flattening `qty × stack_size` into single-item entries.
- Legacy single-item payload remains compatible.
- DAST scans Vite preview output instead of Vite dev server output.
- Gosec CI gate focuses on higher-signal medium-and-higher severity, high-confidence findings.
- Documentation workflow requires updates to both `PATCH_NOTES.md` and `CHANGELOG.md` for every future change.

---

## Fixed

- Fixed backend startup with bare port values such as `LISTEN_ADDR=8080`.
- Fixed Go vet failure from redundant newline output.
- Fixed notification compile errors after credential cleanup.
- Fixed WebSocket log streaming auth behavior.
- Fixed cheats log SQL lookup by joining through `dune.encrypted_accounts` and `player_state.account_id`.
- Fixed stack-size behavior so stacked item grants use the correct number of inventory entries.
- Fixed Trivy action resolution.
- Fixed ZAP issue-creation permission failure.
- Fixed DAST security-header warnings.
- Fixed blueprint import SAST findings for unbounded multipart parsing and unsafe smallint conversion.

---

## Security notes for operators

- Treat `ADMIN_TOKEN` as a privileged secret.
- Rotate any previously shared or committed credentials.
- Keep `.env`, SSH keys, database snapshots, generated secrets, dependency folders, and build output out of source control.
- Prefer `LISTEN_ADDR=127.0.0.1:8080` for local use.
- Do not expose the backend directly to the internet.
- Place remote access behind TLS, a trusted reverse proxy, and a strong identity provider.
