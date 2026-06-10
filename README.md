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

PowerShell support exists through `update.ps1`, but the current validated release workflow is `./update.sh`.

## Development quick start

Create configuration:

```bash
cp .env.example .env
```

Generate a strict backend admin token:

```bash
python - <<'PY'
import secrets
print(secrets.token_urlsafe(32))
PY
```

Set at minimum:

```env
ADMIN_TOKEN=<43-character base64url token from 32 random bytes>
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
