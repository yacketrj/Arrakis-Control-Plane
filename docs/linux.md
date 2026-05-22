# Dune Admin Linux Guide

This guide covers running Dune Admin on Linux for local development, production-style builds, and systemd service operation.

## Supported layout

Dune Admin is a two-part application:

- Go backend API and SSH/database tunnel process.
- React/Vite frontend in `web/`.

The backend reads `.env` from the repository root or from the service working directory. The frontend stores the backend URL and admin token in the browser settings panel.

## Requirements

- Linux x86_64 or arm64.
- Go 1.26.3.
- Node.js 22 or newer.
- npm.
- SSH access to the Dune host.
- PostgreSQL access through the app's SSH tunnel workflow.
- A strong `ADMIN_TOKEN`.

## Install Linux dependencies

From the repository root:

```bash
chmod +x scripts/linux/*.sh
./scripts/linux/install-deps.sh
```

The dependency installer supports common Linux package managers including `apt`, `dnf`, `yum`, `pacman`, and `zypper`. It installs base build dependencies and installs Go 1.26.3 into `/usr/local/go` by default when needed.

Node.js installation is not forced by default because package-manager Node versions vary. Install Node.js 22+ through your preferred distro, package manager, or node version manager. To let the helper try the distro package manager anyway:

```bash
INSTALL_NODE=1 ./scripts/linux/install-deps.sh
```

## Configure `.env`

Create a local environment file:

```bash
cp .env.example .env
openssl rand -base64 32
```

Set the generated value as `ADMIN_TOKEN` in `.env`.

Minimum local settings:

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

The setup wizard can also populate connection settings:

```bash
go run . -setup
```

## Development run

From the repository root:

```bash
chmod +x scripts/linux/*.sh
./scripts/linux/run-dev.sh
```

Open:

```text
http://127.0.0.1:5173
```

In the gear/settings panel, set:

```text
Backend URL: http://127.0.0.1:8080
Admin Token: same value as ADMIN_TOKEN
```

## Production-style Linux build

From the repository root:

```bash
chmod +x scripts/linux/*.sh
./scripts/linux/build-linux.sh
```

The build output is written to:

```text
dist/linux/dune-admin
```

The build script also runs:

```bash
npm install
npm audit --audit-level=high
npm run build
go build
```

## Run built backend manually

```bash
cd dist/linux
cp .env.example .env
nano .env
./dune-admin
```

Serve the frontend separately through Vite preview or a static web server. For local validation:

```bash
cd web
npm install
npm run build
npm run preview -- --host 127.0.0.1
```

## Install backend as a systemd service

Build first:

```bash
./scripts/linux/build-linux.sh
```

Install the backend service:

```bash
sudo ./scripts/linux/install-systemd.sh
```

The service installer:

- Creates a `dune-admin` system user when needed.
- Installs the backend binary into `/opt/dune-admin`.
- Creates `/opt/dune-admin/.env` from `.env.example` if one does not exist.
- Writes `/etc/systemd/system/dune-admin.service`.
- Enables the service.

Edit service configuration before starting:

```bash
sudo nano /opt/dune-admin/.env
```

Start and inspect:

```bash
sudo systemctl start dune-admin
sudo systemctl status dune-admin
journalctl -u dune-admin -f
```

## Security notes

- Keep `LISTEN_ADDR=127.0.0.1:8080` unless the backend is behind a trusted reverse proxy with TLS and identity controls.
- Never commit `.env`, SSH keys, generated tokens, logs, database snapshots, `dist/`, or `web/node_modules/`.
- Treat `ADMIN_TOKEN` as a privileged administrative secret.
- Regenerate `web/package-lock.json` locally after `npm install`, then run `npm audit --audit-level=high` and `npm run build` before committing it.
- Use HTTPS/WSS when operating outside local loopback.

## Validation

```bash
go mod tidy
go test ./...
cd web
npm install
npm audit --audit-level=high
npm run build
```

The GitHub Actions security workflow validates SCA, SAST, DCA, DAST, and secret scanning on pushes to `main`.
