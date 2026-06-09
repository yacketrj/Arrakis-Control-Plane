function Get-GitStatusLines {
  $status = @(& git status --porcelain)
  if ($LASTEXITCODE -ne 0) { throw 'Unable to check git working tree status.' }
  return $status
}

function Write-GitStatusPreview {
  param([string[]]$StatusLines)
  if (-not $StatusLines -or $StatusLines.Count -eq 0) { return }
  Write-Host 'Changed files:' -ForegroundColor Yellow
  $StatusLines | Select-Object -First 12 | ForEach-Object { Write-Host "  $_" -ForegroundColor Yellow }
  if ($StatusLines.Count -gt 12) { Write-Host "  ... $($StatusLines.Count - 12) more" -ForegroundColor Yellow }
}

function Invoke-GitPullIfSafe {
  if ($SkipGitPull) { Write-Host 'Skipping git pull because -SkipGitPull was supplied.' -ForegroundColor Yellow; return }
  $statusLines = @(Get-GitStatusLines)
  if ($statusLines.Count -gt 0 -and -not $AllowDirtyWorktree) {
    Write-Host 'Local changes detected; skipping git pull to avoid merging over uncommitted work.' -ForegroundColor Yellow
    Write-Host 'These changes will be validated and auto-committed if all gates pass.' -ForegroundColor Yellow
    Write-GitStatusPreview -StatusLines $statusLines
    return
  }
  if ($statusLines.Count -gt 0 -and $AllowDirtyWorktree) {
    Write-Host 'Dirty worktree pull allowed because -AllowDirtyWorktree was supplied.' -ForegroundColor Yellow
    Write-GitStatusPreview -StatusLines $statusLines
  }
  Invoke-Step 'Git pull --ff-only' { Invoke-Native 'git' @('pull', '--ff-only') }
}

function Invoke-AutoCommitIfNeeded {
  if ($SkipAutoCommit) { Write-Host 'Skipping auto-commit because -SkipAutoCommit was supplied.' -ForegroundColor Yellow; return }
  Set-Location $RepoRoot
  $statusLines = @(Get-GitStatusLines)
  if ($statusLines.Count -eq 0) { Write-Host 'No repository changes detected; no auto-commit created.' -ForegroundColor Yellow; return }
  Write-GitStatusPreview -StatusLines $statusLines
  Invoke-Native 'git' @('add', '-A')
  $stagedFiles = @(& git diff --cached --name-only)
  if ($LASTEXITCODE -ne 0) { throw 'Unable to inspect staged changes before auto-commit.' }
  if ($stagedFiles.Count -eq 0) { Write-Host 'No staged changes detected after git add; no auto-commit created.' -ForegroundColor Yellow; return }
  Invoke-Native 'git' @('commit', '-m', $CommitMessage)
  $script:AutoCommitSha = (& git rev-parse --short HEAD).Trim()
  if ($LASTEXITCODE -ne 0) { throw 'Auto-commit was created, but the commit SHA could not be read.' }
  Write-Host "Auto-commit created: $script:AutoCommitSha" -ForegroundColor Green
}

function Invoke-AutoPushIfNeeded {
  if ($SkipAutoPush) { Write-Host 'Skipping auto-push because -SkipAutoPush was supplied.' -ForegroundColor Yellow; return }
  Set-Location $RepoRoot
  $branchState = (& git status --short --branch | Select-Object -First 1)
  if ($LASTEXITCODE -ne 0) { throw 'Unable to inspect branch state before auto-push.' }
  if ($branchState -match 'behind') { throw 'Refusing auto-push because the local branch is behind its upstream. Pull first, then rerun.' }
  if ($branchState -notmatch 'ahead') { Write-Host 'Branch is not ahead of upstream; no auto-push needed.' -ForegroundColor Yellow; return }
  Invoke-Native 'git' @('push')
  Write-Host 'Auto-push completed.' -ForegroundColor Green
}
