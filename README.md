# Dune Admin Fork

Dune Admin is a Go backend plus React/Vite frontend for administering a Dune: Awakening server environment through SSH, database, log, blueprint, storage, player, and battlegroup workflows.

## Quick start on Linux

```bash
git pull origin main
chmod +x scripts/linux/*.sh
./scripts/linux/install-deps.sh
cp .env.example .env
nano .env
./scripts/linux/run-dev.sh
```

Open:

```text
http://127.0.0.1:5173
```

In the gear/settings panel, set:

```text
Backend URL: http://127.0.0.1:8080
Admin Token: same value as ADMIN_TOKEN in .env
```

Full Linux instructions are in [`docs/linux.md`](docs/linux.md).

## Linux helper scripts

```text
scripts/linux/install-deps.sh     Install/check Linux dependencies.
scripts/linux/run-dev.sh          Start backend and Vite frontend locally.
scripts/linux/build-linux.sh      Build frontend and Linux backend binary.
scripts/linux/install-systemd.sh  Install backend as a systemd service.
```

## Required runtime configuration

Copy `.env.example` to `.env` and set at minimum:

```env
ADMIN_TOKEN=<long random token>
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

Generate an admin token with:

```bash
openssl rand -base64 32
```

The setup wizard can also populate discovered connection settings:

```bash
go run . -setup
```

## Build and validate

Backend:

```bash
go mod tidy
go test ./...
```

Frontend:

```bash
cd web
npm install
npm audit --audit-level=high
npm run build
```

Linux production-style build:

```bash
./scripts/linux/build-linux.sh
```

Output:

```text
dist/linux/dune-admin
```

## Systemd backend service

```bash
./scripts/linux/build-linux.sh
sudo ./scripts/linux/install-systemd.sh
sudo nano /opt/dune-admin/.env
sudo systemctl start dune-admin
sudo systemctl status dune-admin
journalctl -u dune-admin -f
```

## Security notes

- Treat `ADMIN_TOKEN` as a privileged administrative secret.
- Keep `LISTEN_ADDR=127.0.0.1:8080` unless the backend is protected by a trusted reverse proxy, TLS, and identity controls.
- Never commit `.env`, SSH keys, logs, database snapshots, `dist/`, `web/node_modules/`, or generated credentials.
- Use HTTPS/WSS when operating outside local loopback.
- Keep `PATCH_NOTES.md` and `CHANGELOG.md` updated for every change.
