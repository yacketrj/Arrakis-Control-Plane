#!/usr/bin/env bash

get_node_process_summary() {
  if command -v pgrep >/dev/null 2>&1; then
    local pids
    pids="$(pgrep -f 'node|npm|vite' || true)"
    if [[ -z "$pids" ]]; then
      echo "No running node/npm/vite processes were detected."
      return 0
    fi
    echo "Running node/npm/vite processes detected:"
    ps -p $pids -o pid=,comm=,args= 2>/dev/null | sed -n '1,12p' | sed 's/^/  /' || true
  else
    echo "pgrep is unavailable; skipping process summary."
  fi
}

remove_node_modules_for_repair() {
  local web_root_path="$1"
  local node_modules="$web_root_path/node_modules"
  if [[ ! -d "$node_modules" ]]; then
    echo "node_modules is not present; repair will perform a fresh install."
    return 0
  fi

  echo "Removing web node_modules for automatic dependency repair..."
  rm -rf -- "$node_modules" || {
    get_node_process_summary
    echo "Unable to remove web node_modules automatically. A process or scanner may be locking files under node_modules." >&2
    exit 1
  }
}

invoke_npm_install_with_repair() {
  local clean="$1"
  local install_label="NPM install"
  local install_args=(install)
  if [[ "$clean" -eq 1 ]]; then
    install_label="NPM clean install"
    install_args=(ci)
  fi

  if step "$install_label" run npm "${install_args[@]}"; then
    return 0
  fi

  if [[ "$SKIP_NPM_REPAIR" -eq 1 ]]; then
    echo "$install_label failed and automatic npm repair is disabled by --skip-npm-repair." >&2
    exit 1
  fi

  echo "Attempting npm recovery: cache verify, retry, then dependency repair if needed."
  if step "NPM cache verify" run npm cache verify && step "$install_label retry" run npm "${install_args[@]}"; then
    return 0
  fi

  get_node_process_summary
  remove_node_modules_for_repair "$WEB_ROOT"
  step "$install_label after dependency repair" run npm "${install_args[@]}"
}

show_npm_lock_help() {
  cat >&2 <<'EOF'

NPM dependency update failed after automatic recovery attempts.
The script already retried npm and attempted node_modules repair unless --skip-npm-repair was supplied.
If removal failed, a running process or security scanner is probably locking web/node_modules files.
Close the listed process or pause the scanner and rerun the script.
You can bypass dependency installation only when dependencies are already valid by running: ./update.sh --skip-web-install
EOF
}
