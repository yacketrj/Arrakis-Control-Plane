#!/usr/bin/env bash
set -Eeuo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
TMP_DIR="$(mktemp -d)"
trap 'rm -rf "${TMP_DIR}"' EXIT

log() { printf '[linux-test] %s\n' "$*"; }
fail() { printf '[linux-test] %s\n' "$*" >&2; exit 1; }
assert_file() { [[ -f "$1" ]] || fail "Expected file not found: $1"; }
assert_executable() { [[ -x "$1" ]] || fail "Expected executable file not found: $1"; }
assert_contains() { grep -Fq "$2" "$1" || fail "Expected '$2' in $1"; }

write_fake_tools() {
  local fakebin="$1"
  mkdir -p "${fakebin}"

  cat > "${fakebin}/node" <<'EOF_NODE'
#!/usr/bin/env bash
if [[ "${1:-}" == "-p" ]]; then
  echo "22"
else
  echo "v22.0.0"
fi
EOF_NODE

  cat > "${fakebin}/npm" <<'EOF_NPM'
#!/usr/bin/env bash
set -Eeuo pipefail
case "${1:-}" in
  --version) echo "10.0.0" ;;
  install|audit) exit 0 ;;
  run)
    if [[ "${2:-}" == "build" ]]; then exit 0; fi
    if [[ "${2:-}" == "dev" ]]; then sleep 30; fi
    ;;
esac
exit 0
EOF_NPM

  cat > "${fakebin}/go" <<'EOF_GO'
#!/usr/bin/env bash
set -Eeuo pipefail
if [[ "${1:-}" == "version" ]]; then
  echo "go version go1.26.3 linux/amd64"
  exit 0
fi
if [[ "${1:-}" == "run" ]]; then
  sleep 30
  exit 0
fi
if [[ "${1:-}" == "build" ]]; then
  out=""
  prev=""
  for arg in "$@"; do
    if [[ "${prev}" == "-o" ]]; then
      out="${arg}"
      break
    fi
    prev="${arg}"
  done
  [[ -n "${out}" ]] || { echo "missing -o" >&2; exit 1; }
  mkdir -p "$(dirname "${out}")"
  printf '#!/usr/bin/env bash\necho dune-admin test binary\n' > "${out}"
  chmod 0755 "${out}"
  exit 0
fi
exit 0
EOF_GO

  chmod 0755 "${fakebin}/node" "${fakebin}/npm" "${fakebin}/go"
}

test_shell_syntax() {
  log "Checking Linux shell script syntax"
  bash -n "${ROOT_DIR}/scripts/linux/install-deps.sh"
  bash -n "${ROOT_DIR}/scripts/linux/build-linux.sh"
  bash -n "${ROOT_DIR}/scripts/linux/run-dev.sh"
  bash -n "${ROOT_DIR}/scripts/linux/install-systemd.sh"
}

test_documentation_coverage() {
  log "Checking Linux documentation coverage"
  assert_file "${ROOT_DIR}/docs/linux.md"
  assert_contains "${ROOT_DIR}/README.md" "Quick start on Linux"
  assert_contains "${ROOT_DIR}/docs/linux.md" "systemd"
  assert_contains "${ROOT_DIR}/docs/linux.md" "ADMIN_TOKEN"
  assert_contains "${ROOT_DIR}/.gitignore" "/dist/"
  assert_contains "${ROOT_DIR}/.gitignore" "web/node_modules/"
}

test_build_helper_with_fake_tools() {
  log "Testing build helper with fake toolchain"
  local fakebin="${TMP_DIR}/fakebin-build"
  write_fake_tools "${fakebin}"
  rm -rf "${ROOT_DIR}/dist/linux"
  PATH="${fakebin}:${PATH}" VERSION=test-linux "${ROOT_DIR}/scripts/linux/build-linux.sh"
  assert_executable "${ROOT_DIR}/dist/linux/dune-admin"
  assert_file "${ROOT_DIR}/dist/linux/.env.example"
  assert_file "${ROOT_DIR}/dist/linux/README.md"
  rm -rf "${ROOT_DIR}/dist/linux"
}

test_run_dev_bootstrap_without_env() {
  log "Testing run-dev first-run bootstrap behavior"
  local sandbox="${TMP_DIR}/run-dev-bootstrap"
  mkdir -p "${sandbox}/scripts/linux"
  cp "${ROOT_DIR}/scripts/linux/run-dev.sh" "${sandbox}/scripts/linux/run-dev.sh"
  cp "${ROOT_DIR}/.env.example" "${sandbox}/.env.example"
  chmod 0755 "${sandbox}/scripts/linux/run-dev.sh"
  (cd "${sandbox}" && ./scripts/linux/run-dev.sh)
  assert_file "${sandbox}/.env"
}

test_run_dev_launches_processes_with_fake_tools() {
  log "Testing run-dev process launch with fake toolchain"
  local sandbox="${TMP_DIR}/run-dev-launch"
  local fakebin="${TMP_DIR}/fakebin-run"
  mkdir -p "${sandbox}/scripts/linux" "${sandbox}/web"
  cp "${ROOT_DIR}/scripts/linux/run-dev.sh" "${sandbox}/scripts/linux/run-dev.sh"
  cp "${ROOT_DIR}/.env.example" "${sandbox}/.env.example"
  cp "${ROOT_DIR}/.env.example" "${sandbox}/.env"
  chmod 0755 "${sandbox}/scripts/linux/run-dev.sh"
  write_fake_tools "${fakebin}"
  set +e
  PATH="${fakebin}:${PATH}" timeout 3 "${sandbox}/scripts/linux/run-dev.sh" > "${sandbox}/run-dev.out" 2>&1
  local status=$?
  set -e
  [[ "${status}" -eq 124 ]] || fail "Expected run-dev to keep running until timeout; got status ${status}. Output: $(cat "${sandbox}/run-dev.out")"
  assert_file "${sandbox}/.dune-admin-backend.log"
  assert_file "${sandbox}/.dune-admin-frontend.log"
  assert_contains "${sandbox}/run-dev.out" "Open http://127.0.0.1:5173"
}

main() {
  test_shell_syntax
  test_documentation_coverage
  test_build_helper_with_fake_tools
  test_run_dev_bootstrap_without_env
  test_run_dev_launches_processes_with_fake_tools
  log "Linux helper tests passed"
}

main "$@"
