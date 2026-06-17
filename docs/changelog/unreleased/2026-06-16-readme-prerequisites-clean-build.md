# README prerequisites and clean build validation

## Date

2026-06-16

## Scope

Documentation-only update for the README prerequisites and required tooling section, followed by operator-reported clean build validation.

## Files changed

- `README.md`
- `PATCH_NOTES.md`
- `CHANGELOG.md`

## README additions

The README now includes an explicit `Prerequisites and required tooling` section.

Covered local applications and tools:

- Git
- Go 1.26.3 or compatible newer Go toolchain
- Node.js with npm
- Bash
- PowerShell 7+
- OpenSSH client tools
- modern browser

Covered server-side access requirements:

- Dune: Awakening self-hosted server host
- SSH user and Ed25519 private key
- PostgreSQL-compatible game database access
- admin token
- allowed frontend origins

Covered optional integrations and release tools:

- Discord OAuth credentials
- Discord bot credentials
- GitHub CLI
- optional security and release evidence tools, including `govulncheck`, `gosec`, `gitleaks`, `trivy`, and `syft`
- reverse proxy, VPN, and TLS tooling

The README also now links to `docs/final-v0.1.0-gate-status.md` from the release-reference list.

## Validation

The release owner reported clean local build validation after the README prerequisites update.

Operator-reported validation command:

```bash
./update.sh
```

Windows validation path remains:

```powershell
.\update.ps1 -SkipAutoPush
```

## Runtime impact

No runtime behavior changed.

No backend routes, mutation behavior, endpoints, or higher-risk future admin surfaces were added as part of this documentation-only update.

## Remaining gate

Post-release verification remains pending until tag/artifact install or launch evidence is recorded.
