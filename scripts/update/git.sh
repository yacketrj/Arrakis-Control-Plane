#!/usr/bin/env bash

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
