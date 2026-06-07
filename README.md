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

On Windows, run the canonical path from Git Bash / MINGW64:

```bash
cd /e/dune-admin-fork
./update.sh
```

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

The setup wizard can also populate discovered connection settings:

```bash
go run . -setup
```

Start a local development session manually when needed:

```bash
go run .
```

In another shell:

```bash
cd web
npm install
npm run dev
```

Open:

```text
http://127.0.0.1:5173
```

In the frontend settings panel, set:

```text
Backend URL: http://127.0.0.1:8080
Browser Access Key: same value as ADMIN_TOKEN in .env
```

## Linux helper scripts

The Linux helpers remain available for Linux-only development or deployment tasks:

```text
scripts/linux/install-deps.sh     Install/check Linux dependencies.
scripts/linux/run-dev.sh          Start backend and Vite frontend locally.
scripts/linux/build-linux.sh      Build frontend and Linux backend binary.
scripts/linux/install-systemd.sh  Install backend as a systemd service.
```

Full Linux instructions are in:

```text
docs/linux.md
```

## Manual validation commands

`./update.sh` is preferred. Use these only when debugging a specific gate.

Backend:

```bash
go test -v ./...
go build ./...
```

Frontend:

```bash
cd web
npm install
npm audit --audit-level=high
npm run typecheck
npm run lint
npm run build
```

Ledger check:

```bash
bash scripts/check-ledger-size.sh
```

## Build outputs

`./update.sh` writes build output under:

```text
dist/<goos>/
```

Examples:

```text
dist/windows/arrakis-control-panel.exe
dist/linux/arrakis-control-panel
```

The root repository copy produced by `./update.sh` also uses the project-aligned executable name:

```text
arrakis-control-panel.exe
arrakis-control-panel
```

## Systemd backend service

For Linux service installs:

```bash
./scripts/linux/build-linux.sh
sudo ./scripts/linux/install-systemd.sh
sudo nano /opt/dune-admin/.env
sudo systemctl start dune-admin
sudo systemctl status dune-admin
journalctl -u dune-admin -f
```

The Linux systemd installer may still use the legacy `dune-admin` service path/unit name until the service migration is completed and validated. The compiled executable name is now `arrakis-control-panel`.

## Release workflow

Release candidates use Semantic Versioning-style labels:

```text
v0.1.0-rc.1
v0.1.0-rc.2
v0.1.0
```

Before tagging a release candidate:

1. Update `VERSION`.
2. Update `CHANGELOG.md`.
3. Update `PATCH_NOTES.md`.
4. Create or update `docs/releases/<version>.md`.
5. Record any scope/process deviations in `docs/release-deviation-log.md`.
6. Run `./update.sh`.
7. Tag only after validation evidence is recorded.

Create an annotated tag only from a validated checkout:

```bash
git tag -a v0.1.0-rc.1 -m "Arrakis Control Panel v0.1.0-rc.1"
git push origin v0.1.0-rc.1
```

## Security notes

- Treat `ADMIN_TOKEN` as a privileged administrative secret.
- `ADMIN_TOKEN` must be exactly 43 base64url characters generated from 32 random bytes.
- Keep `LISTEN_ADDR=127.0.0.1:8080` unless the backend is protected by a trusted reverse proxy, TLS, and identity controls.
- Never commit `.env`, SSH keys, logs, database snapshots, `dist/`, `web/node_modules/`, or generated credentials.
- Use HTTPS/WSS when operating outside local loopback.
- Keep Player 360 read-only until self-service identity mapping and mutation safety are fully validated.
- Do not expose arbitrary raw server-command publishing through the UI or Discord.
- Keep `PATCH_NOTES.md`, `CHANGELOG.md`, release checklist files, and deviation logs updated for every durable change.
