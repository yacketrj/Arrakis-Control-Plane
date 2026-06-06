#!/usr/bin/env bash

sudo_run() {
  if [[ "$(id -u)" -eq 0 ]]; then
    run "$@"
  elif command -v sudo >/dev/null 2>&1; then
    run sudo "$@"
  else
    echo "sudo is required to install missing prerequisites. Install $* manually or rerun as root." >&2
    exit 1
  fi
}

refresh_common_paths() {
  export PATH="$PATH:/usr/local/go/bin:/usr/local/bin:/opt/homebrew/bin:$HOME/go/bin"
}

install_with_package_manager() {
  local logical_name="$1"
  shift

  if [[ "$SKIP_PREREQ_INSTALL" -eq 1 ]]; then
    echo "$logical_name was not found on PATH and prerequisite auto-install is disabled." >&2
    exit 1
  fi

  section "Install prerequisite: $logical_name"

  if command -v apt-get >/dev/null 2>&1; then
    sudo_run apt-get update
    sudo_run apt-get install -y "$@"
  elif command -v dnf >/dev/null 2>&1; then
    sudo_run dnf install -y "$@"
  elif command -v yum >/dev/null 2>&1; then
    sudo_run yum install -y "$@"
  elif command -v pacman >/dev/null 2>&1; then
    sudo_run pacman -Sy --needed --noconfirm "$@"
  elif command -v apk >/dev/null 2>&1; then
    sudo_run apk add --no-cache "$@"
  elif command -v brew >/dev/null 2>&1; then
    run brew install "$@"
  else
    echo "No supported package manager was found. Install $logical_name manually." >&2
    exit 1
  fi

  refresh_common_paths
}

install_prerequisite_for_command() {
  local name="$1"
  case "$name" in
    git) install_with_package_manager "Git" git ;;
    go) install_with_package_manager "Go" golang-go ;;
    node) install_with_package_manager "Node.js" nodejs npm ;;
    npm) install_with_package_manager "npm" nodejs npm ;;
    gh) install_with_package_manager "GitHub CLI" gh ;;
    *) echo "No auto-install mapping exists for missing command: $name" >&2; exit 1 ;;
  esac
}

require_cmd() {
  local name="$1"
  local hint="${2:-}"
  if command -v "$name" >/dev/null 2>&1; then
    return 0
  fi

  echo "$name was not found on PATH. Attempting prerequisite auto-install."
  install_prerequisite_for_command "$name"

  if ! command -v "$name" >/dev/null 2>&1; then
    if [[ -n "$hint" ]]; then
      echo "$name is still unavailable after auto-install. $hint" >&2
    else
      echo "$name is still unavailable after auto-install." >&2
    fi
    exit 1
  fi
}
