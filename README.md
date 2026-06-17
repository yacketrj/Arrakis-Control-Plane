# Arrakis Control Panel

Arrakis Control Panel is a Go backend plus React/Vite frontend for administering a Dune: Awakening self-hosted server environment through SSH, database, log, blueprint, storage, player, inventory, Discord, and battlegroup workflows.

The current product name is **Arrakis Control Panel**. Older internal references to `DA Manager` are deprecated and should be replaced when found.

## Upstream attribution

Arrakis Control Panel is a fork of Icehunter's `dune-admin` project by Ryan Wilson:

```text
https://github.com/Icehunter/dune-admin
```

This fork preserves upstream attribution and builds on that original Dune: Awakening server-admin work. Future RMQ/live-admin work should also preserve the upstream acknowledgement that the original `dune-admin` README gives to `@adainrivers` and the `dune-dedicated-server-manager` project for RabbitMQ server-command research.

## Current release

Current release-candidate label:

```text
0.1.0-rc.1
```

Release evidence is tracked in:

```text
docs/releases/v0.1.0-rc.1.md
```

Roadmap and feature priorities are tracked in:

```text
docs/roadmap.md
```

Release policy and release-train goals are tracked in:

```text
docs/release-versioning.md
```

Release deviations are tracked in:

```text
docs/release-deviation-log.md
```

Final release-gate disposition is tracked in:

```text
docs/final-v0.1.0-gate-status.md
```

## Prerequisites and required tooling

Install the local applications and command-line tools below before running the backend, frontend, or validation scripts.

### Required local applications

| Requirement | Purpose | Notes |
|---|---|---|
| Git | Clone, pull, commit, and push the repository. | Required by `./update.sh` and normal development flow. |
| Go 1.26.3 or compatible newer Go toolchain | Build and test the backend. | The active module declares `go 1.26.3` in `go.mod`. |
| Node.js with npm | Install frontend dependencies and run Vite/TypeScript/ESLint builds. | Use a current LTS or newer version compatible with the checked-in frontend dependency lockfile. |
| Bash | Run the canonical `./update.sh` validation/build workflow. | Linux/macOS shells work directly. On Windows, use Git Bash / MINGW64 or WSL. |
| PowerShell 7+ | Run the Windows validation helper. | Required only when using `update.ps1`; Windows operators should prefer PowerShell 7 or newer. |
| OpenSSH client tools | Connect to the Dune server host and manage known-host verification. | `ssh` and `ssh-keyscan` are needed for Ed25519 SSH key and host-key setup. |
| Modern browser | Use the React/Vite operator UI. | Chrome, Edge, Firefox, or equivalent current browser. |

### Required server-side access

Arrakis Control Panel is an admin application. You need access to the target self-hosted Dune: Awakening environment before the app is useful.

| Requirement | Purpose | Notes |
|---|---|---|
| Dune: Awakening self-hosted server host | Remote administration target. | The backend reaches this host over SSH. |
| SSH user and Ed25519 private key | Authenticated backend access to the server host. | Configure `SSH_HOST`, `SSH_USER`, `SSH_KEY`, and `SSH_KNOWN_HOSTS` in `.env`. |
| PostgreSQL-compatible game database access | Player, inventory, storage, and server data access. | Configure `DB_PORT`, `DB_USER`, `DB_PASS`, `DB_NAME`, and `DB_SCHEMA`. The app can use SSH tunnel settings from `.env`. |
| Admin token | Backend API protection. | `ADMIN_TOKEN` must be a strict 43-character base64url token generated from 32 random bytes. See `SECURITY.md`. |
| Allowed frontend origins | Browser API access control. | Configure `ALLOWED_ORIGINS` for local dev and preview ports. |

### Optional integrations and release tools

| Tool or service | Purpose | Notes |
|---|---|---|
| Discord application credentials | Optional Discord OAuth/session and role mapping. | Keep disabled until configured with least-privilege settings. |
| Discord bot token | Optional Discord bot command plane. | Keep disabled until configured; do not expose raw server-command execution. |
| GitHub CLI (`gh`) | Optional local release, issue, and PR workflow helper. | The Bash prerequisite helper can attempt to install it when needed. |
| `govulncheck`, `gosec`, `gitleaks`, `trivy`, `syft` | Optional security/release evidence tools. | Useful before final/stable releases for vulnerability, static-analysis, secret-scan, filesystem/container scan, and SBOM evidence. |
| Reverse proxy, VPN, and TLS certificate tooling | Optional protected remote access. | Do not expose the backend directly to the public internet. Remote exposure requires trusted TLS/identity controls. |

### Quick prerequisite check

From the repository root, these commands should resolve before a normal validation run:

```bash
git --version
go version
node --version
npm --version
ssh -V
```

On Windows, also check:

```powershell
pwsh --version
```

## Canonical update, validation, and build workflow

Use `./update.sh` from the repository root as the canonical validation/build path.

```bash
git pull --ff-only origin main
./update.sh
```

The update script performs the project validation sequence, including:

- safe fast-forward pull when the worktree allows it
- ledger-size check
- Go tests
- Go backend build
- backend binary copy to the repository root unless skipped
- frontend dependency install
- npm audit
- TypeScript typecheck
- frontend lint
- frontend production build
- optional auto-commit/auto-push behavior when validation succeeds

Common options:

```bash
./update.sh --help
./update.sh --skip-git-pull
./update.sh --skip-web-install
./update.sh --skip-auto-commit
./update.sh --skip-auto-push
./update.sh --clean-web-dependencies
```

Run the script from wherever you cloned the repository:

```bash
cd /path/to/Arrakis-Control-Plane
./update.sh
```

On Windows, use Git Bash / MINGW64 or another shell that can execute Bash scripts. The drive letter and mount path depend on the local machine and should not be assumed.

PowerShell support exists through `update.ps1` and is validated for the current hardening slice:

```powershell
.\update.ps1 -SkipAutoPush
```

## Development quick start

Create configuration:

```bash
cp .env.example .env
```

Set required runtime values in `.env`; see `SECURITY.md` for strict admin-token guidance.

Minimum local configuration fields:

```env
ADMIN_TOKEN=<strict 43-character base64url token>
LISTEN_ADDR=127.0.0.1:8080
ALLOWED_ORIGINS=http://127.0.0.1:5173,http://localhost:5173,http://127.0.0.1:4173,http://localhost:4173
SSH_HOST=<dune-host>:22
SSH_USER=dune
SSH_KEY=./sshKey
DB_PORT=15432
DB_USER=dune
DB_PASS=<database password>
DB_NAME=dune
DB_SCHEMA=dune
```
