#!/usr/bin/env bash
set -Eeuo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
OUT_DIR="${ROOT_DIR}/dist/linux"
APP_NAME="arrakis-control-panel"
VERSION="${VERSION:-dev}"

log() { printf '[linux-build] %s\n' "$*"; }
fail() { printf '[linux-build] %s\n' "$*" >&2; exit 1; }

command -v go >/dev/null 2>&1 || fail "go was not found. Run scripts/linux/install-deps.sh first."
command -v node >/dev/null 2>&1 || fail "node was not found. Install Node.js 22+."
command -v npm >/dev/null 2>&1 || fail "npm was not found. Install npm."

log "Building frontend"
cd "${ROOT_DIR}/web"
npm install
npm audit --audit-level=high
npm run build

log "Building backend for Linux"
cd "${ROOT_DIR}"
mkdir -p "${OUT_DIR}"
GOOS=linux GOARCH="${GOARCH:-amd64}" CGO_ENABLED="${CGO_ENABLED:-0}" go build -trimpath -ldflags "-s -w -X main.version=${VERSION}" -o "${OUT_DIR}/${APP_NAME}" .
chmod 0755 "${OUT_DIR}/${APP_NAME}"

cp -f .env.example "${OUT_DIR}/.env.example"
cp -f README.md "${OUT_DIR}/README.md"

log "Build complete: ${OUT_DIR}/${APP_NAME}"
