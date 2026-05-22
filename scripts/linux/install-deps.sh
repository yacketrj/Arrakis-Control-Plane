#!/usr/bin/env bash
set -Eeuo pipefail

GO_VERSION="${GO_VERSION:-1.26.3}"
NODE_MAJOR="${NODE_MAJOR:-22}"
INSTALL_GO="${INSTALL_GO:-1}"
INSTALL_NODE="${INSTALL_NODE:-0}"

log() { printf '\033[1;32m[linux-install]\033[0m %s\n' "$*"; }
warn() { printf '\033[1;33m[linux-install]\033[0m %s\n' "$*"; }
err() { printf '\033[1;31m[linux-install]\033[0m %s\n' "$*" >&2; }

run_root() {
  if [[ "${EUID}" -eq 0 ]]; then
    "$@"
  elif command -v sudo >/dev/null 2>&1; then
    sudo "$@"
  else
    err "This step requires root privileges. Install sudo or run as root."
    exit 1
  fi
}

install_os_packages() {
  if command -v apt-get >/dev/null 2>&1; then
    log "Installing base packages with apt-get"
    run_root apt-get update
    run_root apt-get install -y ca-certificates curl git openssl build-essential tar gzip
  elif command -v dnf >/dev/null 2>&1; then
    log "Installing base packages with dnf"
    run_root dnf install -y ca-certificates curl git openssl make gcc gcc-c++ tar gzip
  elif command -v yum >/dev/null 2>&1; then
    log "Installing base packages with yum"
    run_root yum install -y ca-certificates curl git openssl make gcc gcc-c++ tar gzip
  elif command -v pacman >/dev/null 2>&1; then
    log "Installing base packages with pacman"
    run_root pacman -Sy --needed --noconfirm ca-certificates curl git openssl base-devel tar gzip
  elif command -v zypper >/dev/null 2>&1; then
    log "Installing base packages with zypper"
    run_root zypper --non-interactive install ca-certificates curl git openssl make gcc gcc-c++ tar gzip
  else
    warn "No supported package manager found. Ensure ca-certificates, curl, git, openssl, tar, gzip, make, and gcc are installed."
  fi
}

arch_name() {
  case "$(uname -m)" in
    x86_64|amd64) echo "amd64" ;;
    aarch64|arm64) echo "arm64" ;;
    *) err "Unsupported CPU architecture: $(uname -m)"; exit 1 ;;
  esac
}

install_go_toolchain() {
  if command -v go >/dev/null 2>&1; then
    local current
    current="$(go version | awk '{print $3}' | sed 's/^go//')"
    log "Found Go ${current}"
    if [[ "${current}" == "${GO_VERSION}" ]]; then
      return
    fi
    warn "Repository targets Go ${GO_VERSION}; current Go is ${current}."
  fi

  if [[ "${INSTALL_GO}" != "1" ]]; then
    warn "Skipping Go install because INSTALL_GO=${INSTALL_GO}."
    return
  fi

  local arch tarball url tmp
  arch="$(arch_name)"
  tarball="go${GO_VERSION}.linux-${arch}.tar.gz"
  url="https://go.dev/dl/${tarball}"
  tmp="/tmp/${tarball}"

  log "Installing Go ${GO_VERSION} for linux-${arch}"
  curl -fsSL "${url}" -o "${tmp}"
  run_root rm -rf /usr/local/go
  run_root tar -C /usr/local -xzf "${tmp}"
  rm -f "${tmp}"

  if ! grep -q '/usr/local/go/bin' "${HOME}/.profile" 2>/dev/null; then
    printf '\n# Go toolchain\nexport PATH=/usr/local/go/bin:$PATH\n' >> "${HOME}/.profile"
  fi
  export PATH="/usr/local/go/bin:${PATH}"
  log "Go installed: $(go version)"
}

check_node() {
  if command -v node >/dev/null 2>&1; then
    local major
    major="$(node -p 'process.versions.node.split(".")[0]')"
    log "Found Node.js $(node --version)"
    if [[ "${major}" -lt "${NODE_MAJOR}" ]]; then
      warn "Node.js ${NODE_MAJOR}+ is recommended; current major is ${major}."
    fi
  else
    warn "Node.js was not found. Install Node.js ${NODE_MAJOR}+ before running the frontend."
    if [[ "${INSTALL_NODE}" == "1" ]]; then
      install_node_from_package_manager
    fi
  fi

  if command -v npm >/dev/null 2>&1; then
    log "Found npm $(npm --version)"
  else
    warn "npm was not found. Install npm with Node.js ${NODE_MAJOR}+."
  fi
}

install_node_from_package_manager() {
  if command -v apt-get >/dev/null 2>&1; then
    run_root apt-get update
    run_root apt-get install -y nodejs npm
  elif command -v dnf >/dev/null 2>&1; then
    run_root dnf install -y nodejs npm
  elif command -v yum >/dev/null 2>&1; then
    run_root yum install -y nodejs npm
  elif command -v pacman >/dev/null 2>&1; then
    run_root pacman -Sy --needed --noconfirm nodejs npm
  elif command -v zypper >/dev/null 2>&1; then
    run_root zypper --non-interactive install nodejs npm
  else
    warn "No supported package manager found for Node.js installation."
  fi
}

main() {
  log "Preparing Linux dependencies for Dune Admin"
  install_os_packages
  install_go_toolchain
  check_node
  log "Dependency check complete"
  log "Open a new shell or run: export PATH=/usr/local/go/bin:\$PATH"
}

main "$@"
