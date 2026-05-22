#!/usr/bin/env bash
set -Eeuo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
BACKEND_LOG="${ROOT_DIR}/.dune-admin-backend.log"
FRONTEND_LOG="${ROOT_DIR}/.dune-admin-frontend.log"

log() { printf '[linux-dev] %s\n' "$*"; }
fail() { printf '[linux-dev] %s\n' "$*" >&2; exit 1; }

command -v go >/dev/null 2>&1 || fail "go was not found. Run scripts/linux/install-deps.sh first."
command -v npm >/dev/null 2>&1 || fail "npm was not found. Install Node.js 22+."

if [[ ! -f "${ROOT_DIR}/.env" ]]; then
  log "No .env found. Copying .env.example to .env."
  cp "${ROOT_DIR}/.env.example" "${ROOT_DIR}/.env"
  log "Edit ${ROOT_DIR}/.env, then rerun this script."
  log "You can also run: go run . -setup"
  exit 0
fi

cleanup() {
  if [[ -n "${BACKEND_PID:-}" ]] && kill -0 "${BACKEND_PID}" 2>/dev/null; then
    kill "${BACKEND_PID}" 2>/dev/null || true
  fi
  if [[ -n "${FRONTEND_PID:-}" ]] && kill -0 "${FRONTEND_PID}" 2>/dev/null; then
    kill "${FRONTEND_PID}" 2>/dev/null || true
  fi
}
trap cleanup EXIT INT TERM

log "Starting backend on configured LISTEN_ADDR"
cd "${ROOT_DIR}"
go run . > "${BACKEND_LOG}" 2>&1 &
BACKEND_PID="$!"

log "Installing frontend dependencies"
cd "${ROOT_DIR}/web"
npm install

log "Starting frontend dev server"
npm run dev -- --host 127.0.0.1 > "${FRONTEND_LOG}" 2>&1 &
FRONTEND_PID="$!"

log "Backend log: ${BACKEND_LOG}"
log "Frontend log: ${FRONTEND_LOG}"
log "Open http://127.0.0.1:5173"
log "Press Ctrl+C to stop."

wait -n "${BACKEND_PID}" "${FRONTEND_PID}"
