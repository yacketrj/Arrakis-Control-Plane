#!/usr/bin/env bash
set -Eeuo pipefail

repo_root="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")/.." >/dev/null 2>&1 && pwd -P)"
cd "$repo_root"

fail=0

check_max_lines() {
  local file="$1"
  local max_lines="$2"
  local guidance="$3"

  [[ -f "$file" ]] || return 0

  local lines
  lines="$(wc -l < "$file" | tr -d ' ')"
  if [[ "$lines" -gt "$max_lines" ]]; then
    echo "Ledger size violation: $file has $lines lines; max is $max_lines." >&2
    echo "  $guidance" >&2
    fail=1
  fi
}

warn_large_markdown() {
  local warn_lines="${LEDGER_WARN_MARKDOWN_LINES:-350}"
  local files=()

  while IFS= read -r -d '' file; do
    files+=("$file")
  done < <(find . \
    -path './.git' -prune -o \
    -path './web/node_modules' -prune -o \
    -path './web/dist' -prune -o \
    -path './dist' -prune -o \
    -name '*.md' -print0)

  for file in "${files[@]}"; do
    file="${file#./}"
    local lines
    lines="$(wc -l < "$file" | tr -d ' ')"
    if [[ "$lines" -gt "$warn_lines" ]]; then
      echo "Large Markdown notice: $file has $lines lines. Consider splitting if it is a mutable ledger." >&2
    fi
  done
}

check_max_lines "CHANGELOG.md" "${CHANGELOG_MAX_LINES:-180}" \
  "Keep CHANGELOG.md as a compact index. Move details to docs/changelog/unreleased/ or docs/changelog/archive/."

check_max_lines "PATCH_NOTES.md" "${PATCH_NOTES_MAX_LINES:-220}" \
  "Keep PATCH_NOTES.md focused on the current operator-facing update only. Archive durable detail under docs/changelog/unreleased/."

check_max_lines "docs/appsec-endpoint-audit.md" "${APPSEC_AUDIT_MAX_LINES:-360}" \
  "Keep the AppSec audit file as an index/status document. Move detailed findings to dedicated docs/appsec/ or docs/changelog/ records."

if [[ -d docs/changelog/unreleased ]]; then
  while IFS= read -r -d '' file; do
    file="${file#./}"
    check_max_lines "$file" "${CHANGELOG_RECORD_MAX_LINES:-180}" \
      "Keep per-slice changelog records focused. Split oversized records by topic."
  done < <(find docs/changelog/unreleased -type f -name '*.md' -print0)
fi

warn_large_markdown

if [[ "$fail" -ne 0 ]]; then
  echo "Ledger size validation failed." >&2
  exit 1
fi

echo "Ledger size validation passed."
