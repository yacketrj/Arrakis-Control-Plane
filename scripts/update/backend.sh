#!/usr/bin/env bash

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

run_go_tests() {
  if [[ "$SKIP_GO_TESTS" -eq 0 ]]; then
    step "Go tests" run go test -v ./...
  else
    echo "Skipping Go tests because --skip-go-tests was supplied."
  fi
}

run_backend_build() {
  step "Go backend build" run mkdir -p "$BUILD_OUTPUT_DIR"
  local ldflags="-s -w -X main.version=$VERSION"
  step "Go backend compile" run go build -trimpath -ldflags "$ldflags" -o "$BACKEND_BINARY" .

  if [[ ! -f "$BACKEND_BINARY" ]]; then
    echo "Backend build completed, but expected binary was not found: $BACKEND_BINARY" >&2
    exit 1
  fi

  if command -v stat >/dev/null 2>&1; then
    echo "Backend build output: $BACKEND_BINARY ($(stat -c%s "$BACKEND_BINARY" 2>/dev/null || stat -f%z "$BACKEND_BINARY") bytes)"
  else
    echo "Backend build output: $BACKEND_BINARY"
  fi
}

copy_backend_binary_and_assets() {
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
}
