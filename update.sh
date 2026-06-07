#!/usr/bin/env bash
set -Eeuo pipefail

SCRIPT_DIR="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" >/dev/null 2>&1 && pwd -P)"

source "$SCRIPT_DIR/scripts/update/defaults.sh"
source "$SCRIPT_DIR/scripts/update/ui.sh"
source "$SCRIPT_DIR/scripts/update/prereqs.sh"
source "$SCRIPT_DIR/scripts/update/git.sh"
source "$SCRIPT_DIR/scripts/update/npm.sh"
source "$SCRIPT_DIR/scripts/update/backend.sh"
source "$SCRIPT_DIR/scripts/update/web.sh"

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

resolve_update_paths
cd "$REPO_ROOT"

if [[ ! -d .git ]]; then
  echo "Not a Git repository: $REPO_ROOT" >&2
  exit 1
fi

refresh_common_paths
require_cmd git "Install Git."
require_cmd go "Install Go and reopen the shell so PATH is refreshed."

GOOS_VALUE="$(go env GOOS)"
BINARY_NAME="arrakis-control-panel"
if [[ "$GOOS_VALUE" == "windows" ]]; then
  BINARY_NAME="arrakis-control-panel.exe"
fi
BUILD_OUTPUT_DIR="$(resolve_output_directory "$OUTPUT_DIR" "$GOOS_VALUE")"
BACKEND_BINARY="$BUILD_OUTPUT_DIR/$BINARY_NAME"
REPO_ROOT_BINARY="$REPO_ROOT/$BINARY_NAME"

printf 'Repo folder:    %s\n' "$REPO_ROOT"
printf 'Output folder:  %s\n' "$BUILD_OUTPUT_DIR"
printf 'Build version:  %s\n' "$VERSION"
printf 'Target GOOS:    %s\n' "$GOOS_VALUE"

invoke_git_pull_if_safe
step "Ledger size check" run bash scripts/check-ledger-size.sh
run_go_tests
run_backend_build
copy_backend_binary_and_assets
run_web_validation_and_build

cd "$REPO_ROOT"
step "Git auto-commit successful changes" invoke_auto_commit_if_needed
step "Git auto-push committed changes" invoke_auto_push_if_needed
UPDATE_SUCCEEDED=1

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
