param(
  [string]$RepoRoot,
  [string]$OutputDir,
  [string]$Version,
  [string]$CommitMessage,
  [switch]$CleanWebDependencies,
  [switch]$SkipWebInstall,
  [switch]$SkipGitPull,
  [switch]$SkipGoTests,
  [switch]$SkipWebAudit,
  [switch]$SkipWebTypecheck,
  [switch]$SkipWebLint,
  [switch]$SkipWebBuild,
  [switch]$SkipAutoCommit,
  [switch]$AllowDirtyWorktree
)

Set-StrictMode -Version Latest
$ErrorActionPreference = 'Stop'
$ProgressPreference = 'SilentlyContinue'

if ([string]::IsNullOrWhiteSpace($RepoRoot)) {
  $RepoRoot = $PSScriptRoot
}
if ([string]::IsNullOrWhiteSpace($Version)) {
  $Version = if ($env:VERSION) { $env:VERSION } else { 'dev' }
}
if ([string]::IsNullOrWhiteSpace($CommitMessage)) {
  $CommitMessage = if ($env:COMMIT_MESSAGE) { $env:COMMIT_MESSAGE } else { 'Automated successful update' }
}

$InitialLocation = Get-Location
$UpdateSucceeded = $false
$ExitCode = 0
$AutoCommitSha = ''

function Write-Section {
  param([Parameter(Mandatory = $true)][string]$Name)

  Write-Host ""
  Write-Host "=== $Name ===" -ForegroundColor Cyan
}

function Invoke-Step {
  param(
    [Parameter(Mandatory = $true)][string]$Name,
    [Parameter(Mandatory = $true)][scriptblock]$Command
  )

  Write-Section $Name
  & $Command
}

function Invoke-Native {
  param(
    [Parameter(Mandatory = $true)][string]$FilePath,
    [string[]]$Arguments = @()
  )

  Write-Host (">>> {0} {1}" -f $FilePath, ($Arguments -join ' ')) -ForegroundColor DarkGray
  & $FilePath @Arguments
  $nativeExitCode = $LASTEXITCODE
  if ($null -ne $nativeExitCode -and $nativeExitCode -ne 0) {
    throw "$FilePath $($Arguments -join ' ') failed with exit code $nativeExitCode"
  }
}

function Assert-CommandAvailable {
  param(
    [Parameter(Mandatory = $true)][string]$Name,
    [string]$InstallHint = ''
  )

  if (-not (Get-Command $Name -ErrorAction SilentlyContinue)) {
    $message = "$Name was not found on PATH."
    if (-not [string]::IsNullOrWhiteSpace($InstallHint)) {
      $message = "$message $InstallHint"
    }
    throw $message
  }
}

function Get-GitStatusLines {
  $status = @(& git status --porcelain)
  if ($LASTEXITCODE -ne 0) {
    throw 'Unable to check git working tree status.'
  }
  return $status
}

function Test-GitHasChanges {
  return (@(Get-GitStatusLines).Count -gt 0)
}

function Write-GitStatusPreview {
  param([string[]]$StatusLines)

  if (-not $StatusLines -or $StatusLines.Count -eq 0) {
    return
  }

  Write-Host 'Changed files:' -ForegroundColor Yellow
  $StatusLines | Select-Object -First 12 | ForEach-Object { Write-Host "  $_" -ForegroundColor Yellow }
  if ($StatusLines.Count -gt 12) {
    Write-Host "  ... $($StatusLines.Count - 12) more" -ForegroundColor Yellow
  }
}

function Invoke-GitPullIfSafe {
  if ($SkipGitPull) {
    Write-Host 'Skipping git pull because -SkipGitPull was supplied.' -ForegroundColor Yellow
    return
  }

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
  if ($SkipAutoCommit) {
    Write-Host 'Skipping auto-commit because -SkipAutoCommit was supplied.' -ForegroundColor Yellow
    return
  }

  Set-Location $RepoRoot
  $statusLines = @(Get-GitStatusLines)
  if ($statusLines.Count -eq 0) {
    Write-Host 'No repository changes detected; no auto-commit created.' -ForegroundColor Yellow
    return
  }

  Write-GitStatusPreview -StatusLines $statusLines
  Invoke-Native 'git' @('add', '-A')

  $stagedFiles = @(& git diff --cached --name-only)
  if ($LASTEXITCODE -ne 0) {
    throw 'Unable to inspect staged changes before auto-commit.'
  }

  if ($stagedFiles.Count -eq 0) {
    Write-Host 'No staged changes detected after git add; no auto-commit created.' -ForegroundColor Yellow
    return
  }

  Invoke-Native 'git' @('commit', '-m', $CommitMessage)
  $script:AutoCommitSha = (& git rev-parse --short HEAD).Trim()
  if ($LASTEXITCODE -ne 0) {
    throw 'Auto-commit was created, but the commit SHA could not be read.'
  }
  Write-Host "Auto-commit created: $script:AutoCommitSha" -ForegroundColor Green
}

function Resolve-OutputDirectory {
  param(
    [Parameter(Mandatory = $true)][string]$Root,
    [string]$RequestedOutputDir
  )

  if ([string]::IsNullOrWhiteSpace($RequestedOutputDir)) {
    return (Join-Path $Root 'dist\windows')
  }

  if ([System.IO.Path]::IsPathRooted($RequestedOutputDir)) {
    return $RequestedOutputDir
  }

  return (Join-Path $Root $RequestedOutputDir)
}

function Show-NpmLockHelp {
  Write-Host ""
  Write-Host 'NPM dependency update failed.' -ForegroundColor Red
  Write-Host 'On Windows this is commonly caused by a locked native .node file in node_modules.' -ForegroundColor Yellow
  Write-Host 'Close any running frontend dev server, node.exe process, editor terminal, or antivirus scan that may be holding web\node_modules files open, then re-run the script.' -ForegroundColor Yellow
  Write-Host 'Normal updates use npm install instead of npm ci to avoid deleting locked native dependencies.' -ForegroundColor Yellow
  Write-Host 'Use .\update.ps1 -CleanWebDependencies only when you intentionally want a clean npm ci install.' -ForegroundColor Yellow
  Write-Host 'Use .\update.ps1 -SkipWebInstall when dependencies are already installed and locked files are expected.' -ForegroundColor Yellow
}

try {
  $RepoRoot = (Resolve-Path $RepoRoot).Path
  $WebRoot = Join-Path $RepoRoot 'web'
  $BuildOutputDir = Resolve-OutputDirectory -Root $RepoRoot -RequestedOutputDir $OutputDir
  $BackendBinary = Join-Path $BuildOutputDir 'dune-admin.exe'

  Set-Location $RepoRoot
  Write-Host "Repo folder:    $RepoRoot"
  Write-Host "Output folder:  $BuildOutputDir"
  Write-Host "Build version:  $Version"

  if (-not (Test-Path '.git')) {
    throw "Not a Git repository: $RepoRoot"
  }

  Assert-CommandAvailable -Name 'git' -InstallHint 'Install Git for Windows.'
  Assert-CommandAvailable -Name 'go' -InstallHint 'Install Go and reopen PowerShell so PATH is refreshed.'

  Invoke-GitPullIfSafe

  if (-not $SkipGoTests) {
    Invoke-Step 'Go tests' { Invoke-Native 'go' @('test', '-v', './...') }
  } else {
    Write-Host 'Skipping Go tests because -SkipGoTests was supplied.' -ForegroundColor Yellow
  }

  Invoke-Step 'Go backend build' {
    New-Item -ItemType Directory -Force -Path $BuildOutputDir | Out-Null
    $ldflags = "-s -w -X main.version=$Version"
    Invoke-Native 'go' @('build', '-trimpath', '-ldflags', $ldflags, '-o', $BackendBinary, '.')
  }

  if (-not (Test-Path $BackendBinary)) {
    throw "Backend build completed, but expected binary was not found: $BackendBinary"
  }

  foreach ($asset in @('.env.example', 'README.md')) {
    $source = Join-Path $RepoRoot $asset
    if (Test-Path $source) {
      Copy-Item -Force $source (Join-Path $BuildOutputDir $asset)
    }
  }

  if (Test-Path $WebRoot) {
    Set-Location $WebRoot
    Write-Host "Web folder:     $(Get-Location)"

    if (Test-Path 'package.json') {
      Assert-CommandAvailable -Name 'node' -InstallHint 'Install Node.js 22+.'
      Assert-CommandAvailable -Name 'npm' -InstallHint 'Install npm or repair the Node.js installation.'

      Invoke-Step 'Node version' { Invoke-Native 'node' @('--version') }
      Invoke-Step 'NPM version' { Invoke-Native 'npm' @('--version') }

      if (-not $SkipWebInstall) {
        if ($CleanWebDependencies) {
          Invoke-Step 'NPM clean install' { Invoke-Native 'npm' @('ci') }
        } else {
          Invoke-Step 'NPM install' { Invoke-Native 'npm' @('install') }
        }
      } else {
        Write-Host 'Skipping npm install because -SkipWebInstall was supplied.' -ForegroundColor Yellow
      }

      if (-not $SkipWebAudit) {
        Invoke-Step 'NPM audit' { Invoke-Native 'npm' @('audit', '--audit-level=high') }
      } else {
        Write-Host 'Skipping npm audit because -SkipWebAudit was supplied.' -ForegroundColor Yellow
      }

      if (-not $SkipWebTypecheck) {
        Invoke-Step 'Web typecheck' { Invoke-Native 'npm' @('run', 'typecheck') }
      } else {
        Write-Host 'Skipping web typecheck because -SkipWebTypecheck was supplied.' -ForegroundColor Yellow
      }

      if (-not $SkipWebLint) {
        Invoke-Step 'Web lint' { Invoke-Native 'npm' @('run', 'lint') }
      } else {
        Write-Host 'Skipping web lint because -SkipWebLint was supplied.' -ForegroundColor Yellow
      }

      if (-not $SkipWebBuild) {
        Invoke-Step 'Web build' { Invoke-Native 'npm' @('run', 'build') }
      } else {
        Write-Host 'Skipping web build because -SkipWebBuild was supplied.' -ForegroundColor Yellow
      }
    } else {
      Write-Host 'package.json not found; skipping web build' -ForegroundColor Yellow
    }
  } else {
    Write-Host "Web folder not found; skipping web build: $WebRoot" -ForegroundColor Yellow
  }

  Invoke-Step 'Git auto-commit successful changes' { Invoke-AutoCommitIfNeeded }

  $UpdateSucceeded = $true
}
catch {
  $ExitCode = 1
  Write-Host ""
  Write-Host 'Update failed.' -ForegroundColor Red
  Write-Host $_.Exception.Message -ForegroundColor Red

  $message = $_.Exception.Message
  if ($message -match 'NPM|npm|EPERM|unlink|node_modules') {
    Show-NpmLockHelp
  }
}
finally {
  if ((Test-Path variable:RepoRoot) -and (Test-Path $RepoRoot)) {
    Set-Location $RepoRoot
  } else {
    Set-Location $InitialLocation
  }

  Write-Host ""
  Write-Host "Current folder: $(Get-Location)"
}

if ($UpdateSucceeded) {
  Write-Host ""
  Write-Host 'Update complete.' -ForegroundColor Green
  Write-Host "Backend binary: $BackendBinary"
  if (-not $SkipWebBuild) {
    Write-Host "Frontend build:  $WebRoot\dist"
  }
  Write-Host "Copied assets:   $BuildOutputDir\.env.example, $BuildOutputDir\README.md"
  if (-not [string]::IsNullOrWhiteSpace($AutoCommitSha)) {
    Write-Host "Auto-commit:     $AutoCommitSha"
  }
} else {
  exit $ExitCode
}
