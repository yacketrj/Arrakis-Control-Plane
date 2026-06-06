#!/usr/bin/env bash

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
