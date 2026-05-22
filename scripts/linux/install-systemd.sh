#!/usr/bin/env bash
set -Eeuo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
SERVICE_NAME="${SERVICE_NAME:-dune-admin}"
INSTALL_DIR="${INSTALL_DIR:-/opt/dune-admin}"
SERVICE_USER="${SERVICE_USER:-dune-admin}"
BINARY_PATH="${ROOT_DIR}/dist/linux/dune-admin"

log() { printf '[linux-systemd] %s\n' "$*"; }
fail() { printf '[linux-systemd] %s\n' "$*" >&2; exit 1; }

require_root() {
  if [[ "${EUID}" -ne 0 ]]; then
    fail "Run this script as root or through sudo."
  fi
}

main() {
  require_root

  if ! command -v systemctl >/dev/null 2>&1; then
    fail "systemctl was not found. This installer targets systemd-based Linux hosts."
  fi

  if [[ ! -x "${BINARY_PATH}" ]]; then
    fail "Linux binary not found at ${BINARY_PATH}. Run scripts/linux/build-linux.sh first."
  fi

  log "Creating service user ${SERVICE_USER} if needed"
  if ! id "${SERVICE_USER}" >/dev/null 2>&1; then
    useradd --system --home-dir "${INSTALL_DIR}" --shell /usr/sbin/nologin "${SERVICE_USER}"
  fi

  log "Installing files into ${INSTALL_DIR}"
  mkdir -p "${INSTALL_DIR}"
  install -m 0755 "${BINARY_PATH}" "${INSTALL_DIR}/dune-admin"
  install -m 0644 "${ROOT_DIR}/.env.example" "${INSTALL_DIR}/.env.example"

  if [[ ! -f "${INSTALL_DIR}/.env" ]]; then
    install -m 0600 "${ROOT_DIR}/.env.example" "${INSTALL_DIR}/.env"
    log "Created ${INSTALL_DIR}/.env from .env.example. Edit it before starting the service."
  fi

  chown -R "${SERVICE_USER}:${SERVICE_USER}" "${INSTALL_DIR}"

  log "Writing /etc/systemd/system/${SERVICE_NAME}.service"
  cat > "/etc/systemd/system/${SERVICE_NAME}.service" <<SERVICE
[Unit]
Description=Dune Admin backend
After=network-online.target
Wants=network-online.target

[Service]
Type=simple
User=${SERVICE_USER}
Group=${SERVICE_USER}
WorkingDirectory=${INSTALL_DIR}
EnvironmentFile=${INSTALL_DIR}/.env
ExecStart=${INSTALL_DIR}/dune-admin
Restart=on-failure
RestartSec=5
NoNewPrivileges=true
PrivateTmp=true
ProtectSystem=full
ProtectHome=read-only
ReadWritePaths=${INSTALL_DIR}

[Install]
WantedBy=multi-user.target
SERVICE

  systemctl daemon-reload
  systemctl enable "${SERVICE_NAME}.service"

  log "Installed ${SERVICE_NAME}.service"
  log "Edit ${INSTALL_DIR}/.env, then start with: sudo systemctl start ${SERVICE_NAME}"
  log "Check status with: sudo systemctl status ${SERVICE_NAME}"
}

main "$@"
