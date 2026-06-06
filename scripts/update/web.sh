#!/usr/bin/env bash

run_web_validation_and_build() {
  if [[ ! -d "$WEB_ROOT" ]]; then
    echo "Web folder not found; skipping web build: $WEB_ROOT"
    return 0
  fi

  cd "$WEB_ROOT"
  echo "Web folder:     $(pwd)"

  if [[ ! -f package.json ]]; then
    echo "package.json not found; skipping web build"
    return 0
  fi

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
}
