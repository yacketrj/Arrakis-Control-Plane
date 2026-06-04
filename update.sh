#!/usr/bin/env bash
set -Eeuo pipefail

REPO_ROOT=""
OUTPUT_DIR=""
VERSION="${VERSION:-dev}"
COMMIT_MESSAGE="${COMMIT_MESSAGE:-Automated successful update}"
CLEAN_WEB_DEPENDENCIES=0
SKIP_WEB_INSTALL=0
SKIP_GIT_PULL=0
SKIP_GO_TESTS=0
SKIP_WEB_AUDIT=0
SKIP_WEB_TYPECHECK=0
SKIP_WEB_LINT=0
SKIP_WEB_BUILD=0
SKIP_NPM_REPAIR=0
SKIP_AUTO_COMMIT=0
SKIP_AUTO_PUSH=0
ALLOW_DIRTY_WORKTREE=0
SKIP_ROOT_BINARY_COPY=0
SKIP_PREREQ_INSTALL=0

INITIAL_DIR="$(pwd)"
UPDATE_SUCCEEDED=0
AUTO_COMMIT_SHA=""
SCRIPT_DIR="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" >/dev/null 2>&1 && pwd -P)"

COLOR_RESET=$'\033[0m'
COLOR_DIM=$'\033[90m'
COLOR_CYAN=$'\033[36m'
COLOR_GREEN=$'\033[32m'
COLOR_RED=$'\033[31m'
COLOR_YELLOW=$'\033[33m'

usage() {
  cat <<'EOF'
Usage: ./update.sh [options]

Options:
  --repo-root PATH             Repository root. Defaults to this script's directory.
  --output-dir PATH            Build output directory. Defaults to dist/<goos>.
  --version VALUE              Build version. Defaults to VERSION env var or dev.
  --commit-message VALUE       Auto-commit message.
  --clean-web-dependencies     Run npm ci instead of npm install.
  --skip-web-install           Skip npm install/npm ci.
  --skip-git-pull              Skip git pull --ff-only.
  --skip-go-tests              Skip go test -v ./...
  --skip-web-audit             Skip npm audit --audit-level=high.
  --skip-web-typecheck         Skip npm run typecheck.
  --skip-web-lint              Skip npm run lint.
  --skip-web-build             Skip npm run build.
  --skip-npm-repair            Do not attempt npm cache/node_modules repair.
  --skip-auto-commit           Do not auto-commit successful validated changes.
  --skip-auto-push             Do not auto-push when branch is ahead.
  --allow-dirty-worktree       Allow git pull with local changes present.
  --skip-root-binary-copy      Do not copy dist binary back to repo root.
  --skip-prereq-install        Do not auto-install missing Git/Go/Node/npm tools.
  -h, --help                   Show help.
EOF
}

while [[ $# -gt 0 ]]; do
  case "$1" in
    --repo-root) REPO_ROOT="${2:-}"; shift 2 ;;
    --output-dir) OUTPUT_DIR="${2:-}"; shift 2 ;;
    --version) VERSION="${2:-}"; shift 2 ;;
    --commit-message) COMMIT_MESSAGE="${2:-}"; shift 2 ;;
    --clean-web-dependencies) CLEAN_WEB_DEPENDENCIES=1; shift ;;
    --skip-web-install) SKIP_WEB_INSTALL=1; shift ;;
    --skip-git-pull) SKIP_GIT_PULL=1; shift ;;
    --skip-go-tests) SKIP_GO_TESTS=1; shift ;;
    --skip-web-audit) SKIP_WEB_AUDIT=1; shift ;;
    --skip-web-typecheck) SKIP_WEB_TYPECHECK=1; shift ;;
    --skip-web-lint) SKIP_WEB_LINT=1; shift ;;
    --skip-web-build) SKIP_WEB_BUILD=1; shift ;;
    --skip-npm-repair) SKIP_NPM_REPAIR=1; shift ;;
    --skip-auto-commit) SKIP_AUTO_COMMIT=1; shift ;;
    --skip-auto-push) SKIP_AUTO_PUSH=1; shift ;;
    --allow-dirty-worktree) ALLOW_DIRTY_WORKTREE=1; shift ;;
    --skip-root-binary-copy) SKIP_ROOT_BINARY_COPY=1; shift ;;
    --skip-prereq-install) SKIP_PREREQ_INSTALL=1; shift ;;
    -h|--help) usage; exit 0 ;;
    *) echo "Unknown option: $1" >&2; usage >&2; exit 2 ;;
  esac
done

if [[ -z "$REPO_ROOT" ]]; then
  REPO_ROOT="$SCRIPT_DIR"
fi
REPO_ROOT="$(cd "$REPO_ROOT" && pwd -P)"
WEB_ROOT="$REPO_ROOT/web"

section() {
  printf '\n%s=== %s ===%s\n' "$COLOR_CYAN" "$1" "$COLOR_RESET"
}

colorize_output() {
  sed -E \
    -e "s/(^|[[:space:]])(=== RUN[[:space:]])/\1${COLOR_CYAN}\2${COLOR_RESET}/g" \
    -e "s/(^|[[:space:]])(--- PASS:|PASS)([[:space:]]|$)/\1${COLOR_GREEN}\2${COLOR_RESET}\3/g" \
    -e "s/(^|[[:space:]])(--- FAIL:|FAIL)([[:space:]]|$)/\1${COLOR_RED}\2${COLOR_RESET}\3/g" \
    -e "s/(^|[[:space:]])(Update failed\.)([[:space:]]|$)/\1${COLOR_RED}\2${COLOR_RESET}\3/g"
}

run() {
  printf '%s>>> ' "$COLOR_DIM"
  printf '%q ' "$@"
  printf '%s\n' "$COLOR_RESET"
  set +e
  "$@" 2>&1 | colorize_output
  local status=${PIPESTATUS[0]}
  set -e
  return "$status"
}

step() {
  local name="$1"
  shift
  section "$name"
  "$@"
}

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

get_git_status_lines() {
  git status --porcelain
}

write_git_status_preview() {
  local lines
  lines="$(get_git_status_lines || true)"
  [[ -z "$lines" ]] && return 0
  echo "Changed files:"
  echo "$lines" | sed -n '1,12p' | sed 's/^/  /'
  local count
  count="$(echo "$lines" | wc -l | tr -d ' ')"
  if [[ "$count" -gt 12 ]]; then
    echo "  ... $((count - 12)) more"
  fi
}

invoke_git_pull_if_safe() {
  if [[ "$SKIP_GIT_PULL" -eq 1 ]]; then
    echo "Skipping git pull because --skip-git-pull was supplied."
    return 0
  fi

  local status_lines
  status_lines="$(get_git_status_lines)"
  if [[ -n "$status_lines" && "$ALLOW_DIRTY_WORKTREE" -eq 0 ]]; then
    echo "Local changes detected; skipping git pull to avoid merging over uncommitted work."
    echo "These changes will be validated and auto-committed if all gates pass."
    write_git_status_preview
    return 0
  fi

  if [[ -n "$status_lines" && "$ALLOW_DIRTY_WORKTREE" -eq 1 ]]; then
    echo "Dirty worktree pull allowed because --allow-dirty-worktree was supplied."
    write_git_status_preview
  fi

  step "Git pull --ff-only" run git pull --ff-only origin main
}

invoke_auto_commit_if_needed() {
  if [[ "$SKIP_AUTO_COMMIT" -eq 1 ]]; then
    echo "Skipping auto-commit because --skip-auto-commit was supplied."
    return 0
  fi

  cd "$REPO_ROOT"
  local status_lines
  status_lines="$(get_git_status_lines)"
  if [[ -z "$status_lines" ]]; then
    echo "No repository changes detected; no auto-commit created."
    return 0
  fi

  write_git_status_preview
  run git add -A

  local staged_files
  staged_files="$(git diff --cached --name-only)"
  if [[ -z "$staged_files" ]]; then
    echo "No staged changes detected after git add; no auto-commit created."
    return 0
  fi

  run git commit -m "$COMMIT_MESSAGE"
  AUTO_COMMIT_SHA="$(git rev-parse --short HEAD)"
  echo "Auto-commit created: $AUTO_COMMIT_SHA"
}

invoke_auto_push_if_needed() {
  if [[ "$SKIP_AUTO_PUSH" -eq 1 ]]; then
    echo "Skipping auto-push because --skip-auto-push was supplied."
    return 0
  fi

  cd "$REPO_ROOT"
  local branch_state
  branch_state="$(git status --short --branch | head -n 1)"

  if [[ "$branch_state" == *behind* ]]; then
    echo "Refusing auto-push because the local branch is behind its upstream. Pull first, then rerun." >&2
    exit 1
  fi

  if [[ "$branch_state" != *ahead* ]]; then
    echo "Branch is not ahead of upstream; no auto-push needed."
    return 0
  fi

  run git push
  echo "Auto-push completed."
}

resolve_output_directory() {
  local requested="$1"
  local goos="$2"
  if [[ -z "$requested" ]]; then
    echo "$REPO_ROOT/dist/$goos"
    return 0
  fi
  if [[ "$requested" = /* ]]; then
    echo "$requested"
  else
    echo "$REPO_ROOT/$requested"
  fi
}

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

on_error() {
  local exit_code=$?
  echo ""
  printf '%sUpdate failed.%s\n' "$COLOR_RED" "$COLOR_RESET" >&2
  case "${BASH_COMMAND:-}" in
    *npm*|*node_modules*|*rm\ -rf*) show_npm_lock_help ;;
  esac
  cd "$REPO_ROOT" 2>/dev/null || cd "$INITIAL_DIR" 2>/dev/null || true
  echo ""
  echo "Current folder: $(pwd)"
  exit "$exit_code"
}
trap on_error ERR

cd "$REPO_ROOT"

if [[ ! -d .git ]]; then
  echo "Not a Git repository: $REPO_ROOT" >&2
  exit 1
fi

refresh_common_paths
require_cmd git "Install Git."
require_cmd go "Install Go and reopen the shell so PATH is refreshed."

GOOS_VALUE="$(go env GOOS)"
BINARY_NAME="dune-admin"
if [[ "$GOOS_VALUE" == "windows" ]]; then
  BINARY_NAME="dune-admin.exe"
fi
BUILD_OUTPUT_DIR="$(resolve_output_directory "$OUTPUT_DIR" "$GOOS_VALUE")"
BACKEND_BINARY="$BUILD_OUTPUT_DIR/$BINARY_NAME"
REPO_ROOT_BINARY="$REPO_ROOT/$BINARY_NAME"

printf 'Repo folder:    %s\n' "$REPO_ROOT"
printf 'Output folder:  %s\n' "$BUILD_OUTPUT_DIR"
printf 'Build version:  %s\n' "$VERSION"
printf 'Target GOOS:    %s\n' "$GOOS_VALUE"

invoke_git_pull_if_safe

if [[ "$SKIP_GO_TESTS" -eq 0 ]]; then
  step "Go tests" run go test -v ./...
else
  echo "Skipping Go tests because --skip-go-tests was supplied."
fi

step "Go backend build" bash -c '
  set -Eeuo pipefail
  mkdir -p "$1"
  ldflags="-s -w -X main.version=$2"
  go build -trimpath -ldflags "$ldflags" -o "$3" .
' bash "$BUILD_OUTPUT_DIR" "$VERSION" "$BACKEND_BINARY"

if [[ ! -f "$BACKEND_BINARY" ]]; then
  echo "Backend build completed, but expected binary was not found: $BACKEND_BINARY" >&2
  exit 1
fi

if command -v stat >/dev/null 2>&1; then
  echo "Backend build output: $BACKEND_BINARY ($(stat -c%s "$BACKEND_BINARY" 2>/dev/null || stat -f%z "$BACKEND_BINARY") bytes)"
else
  echo "Backend build output: $BACKEND_BINARY"
fi

if [[ "$SKIP_ROOT_BINARY_COPY" -eq 0 ]]; then
  step "Copy backend binary to repo root" run cp -f "$BACKEND_BINARY" "$REPO_ROOT_BINARY"
  chmod +x "$REPO_ROOT_BINARY" || true
  echo "Repo root binary:    $REPO_ROOT_BINARY"
else
  echo "Skipping repo-root binary copy because --skip-root-binary-copy was supplied."
fi

for asset in .env.example README.md; do
  if [[ -f "$REPO_ROOT/$asset" ]]; then
    cp -f "$REPO_ROOT/$asset" "$BUILD_OUTPUT_DIR/$asset"
  fi
done

if [[ -d "$WEB_ROOT" ]]; then
  cd "$WEB_ROOT"
  echo "Web folder:     $(pwd)"

  if [[ -f package.json ]]; then
    require_cmd node "Install Node.js 22+."
    require_cmd npm "Install npm or repair the Node.js installation."

    step "Node version" run node --version
    step "NPM version" run npm --version

    if [[ "$SKIP_WEB_INSTALL" -eq 0 ]]; then
      invoke_npm_install_with_repair "$CLEAN_WEB_DEPENDENCIES"
    else
      echo "Skipping npm install because --skip-web-install was supplied."
    fi

    if [[ "$SKIP_WEB_AUDIT" -eq 0 ]]; then
      step "NPM audit" run npm audit --audit-level=high
    else
      echo "Skipping npm audit because --skip-web-audit was supplied."
    fi

    if [[ "$SKIP_WEB_TYPECHECK" -eq 0 ]]; then
      step "Web typecheck" run npm run typecheck
    else
      echo "Skipping web typecheck because --skip-web-typecheck was supplied."
    fi

    if [[ "$SKIP_WEB_LINT" -eq 0 ]]; then
      step "Web lint" run npm run lint
    else
      echo "Skipping web lint because --skip-web-lint was supplied."
    fi

    if [[ "$SKIP_WEB_BUILD" -eq 0 ]]; then
      step "Web build" run npm run build
    else
      echo "Skipping web build because --skip-web-build was supplied."
    fi
  else
    echo "package.json not found; skipping web build"
  fi
else
  echo "Web folder not found; skipping web build: $WEB_ROOT"
fi

cd "$REPO_ROOT"
step "Git auto-commit successful changes" invoke_auto_commit_if_needed
step "Git auto-push committed changes" invoke_auto_push_if_needed
UPDATE_SUCCEEDED=1
trap - ERR

cd "$REPO_ROOT"
echo ""
echo "Update complete."
echo "Backend binary: $BACKEND_BINARY"
echo "Repo root exe:   $REPO_ROOT_BINARY"
if [[ "$SKIP_WEB_BUILD" -eq 0 ]]; then
  echo "Frontend build:  $WEB_ROOT/dist"
fi
echo "Copied assets:   $BUILD_OUTPUT_DIR/.env.example, $BUILD_OUTPUT_DIR/README.md"
if [[ -n "$AUTO_COMMIT_SHA" ]]; then
  echo "Auto-commit:     $AUTO_COMMIT_SHA"
fi
