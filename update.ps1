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
  [switch]$SkipNpmRepair,
  [switch]$SkipAutoCommit,
  [switch]$SkipAutoPush,
  [switch]$AllowDirtyWorktree,
  [switch]$SkipRootBinaryCopy,
  [switch]$SkipPrereqInstall
)

Set-StrictMode -Version Latest
$ErrorActionPreference = 'Stop'
$ProgressPreference = 'SilentlyContinue'

$CommonHelpers = Join-Path $PSScriptRoot 'scripts\update\powershell-common.ps1'
if (-not (Test-Path $CommonHelpers)) { throw "PowerShell update helper not found: $CommonHelpers" }
. $CommonHelpers

$GitHelpers = Join-Path $PSScriptRoot 'scripts\update\powershell-git.ps1'
if (-not (Test-Path $GitHelpers)) { throw "PowerShell update Git helper not found: $GitHelpers" }
. $GitHelpers

$NpmHelpers = Join-Path $PSScriptRoot 'scripts\update\powershell-npm.ps1'
if (-not (Test-Path $NpmHelpers)) { throw "PowerShell update npm helper not found: $NpmHelpers" }
. $NpmHelpers

$BackendHelpers = Join-Path $PSScriptRoot 'scripts\update\powershell-backend.ps1'
if (-not (Test-Path $BackendHelpers)) { throw "PowerShell update backend helper not found: $BackendHelpers" }
. $BackendHelpers

if ([string]::IsNullOrWhiteSpace($RepoRoot)) { $RepoRoot = $PSScriptRoot }
if ([string]::IsNullOrWhiteSpace($Version)) { $Version = if ($env:VERSION) { $env:VERSION } else { 'dev' } }
if ([string]::IsNullOrWhiteSpace($CommitMessage)) { $CommitMessage = if ($env:COMMIT_MESSAGE) { $env:COMMIT_MESSAGE } else { 'Automated successful update' } }

$InitialLocation = Get-Location
$UpdateSucceeded = $false
$ExitCode = 0
$AutoCommitSha = ''

try {
  Update-ProcessPath
  $RepoRoot = (Resolve-Path $RepoRoot).Path
  $WebRoot = Join-Path $RepoRoot 'web'
  $BuildOutputDir = Resolve-OutputDirectory -Root $RepoRoot -RequestedOutputDir $OutputDir
  $BackendBinary = Join-Path $BuildOutputDir 'arrakis-control-panel.exe'
  $RepoRootBinary = Join-Path $RepoRoot 'arrakis-control-panel.exe'

  Set-Location $RepoRoot
  Write-Host "Repo folder:    $RepoRoot"
  Write-Host "Output folder:  $BuildOutputDir"
  Write-Host "Build version:  $Version"

  if (-not (Test-Path '.git')) { throw "Not a Git repository: $RepoRoot" }

  Assert-CommandAvailable -Name 'git' -InstallHint 'Install Git for Windows.'
  Assert-CommandAvailable -Name 'go' -InstallHint 'Install Go and reopen PowerShell so PATH is refreshed.'

  Invoke-GitPullIfSafe

  if (-not $SkipGoTests) { Invoke-Step 'Go tests' { Invoke-BackendGoTests } }
  else { Write-Host 'Skipping Go tests because -SkipGoTests was supplied.' -ForegroundColor Yellow }

  Invoke-BackendBuild -BuildOutputDir $BuildOutputDir -BackendBinary $BackendBinary -Version $Version
  Copy-BackendBinaryAndAssets -RepoRoot $RepoRoot -BuildOutputDir $BuildOutputDir -BackendBinary $BackendBinary -RepoRootBinary $RepoRootBinary -SkipRootBinaryCopy:$SkipRootBinaryCopy

  if (Test-Path $WebRoot) {
    Set-Location $WebRoot
    Write-Host "Web folder:     $(Get-Location)"
    if (Test-Path 'package.json') {
      Assert-CommandAvailable -Name 'node' -InstallHint 'Install Node.js 22+.'
      Assert-CommandAvailable -Name 'npm' -InstallHint 'Install npm or repair the Node.js installation.'
      Invoke-Step 'Node version' { Invoke-Native 'node' @('--version') }
      Invoke-Step 'NPM version' { Invoke-Native 'npm' @('--version') }
      if (-not $SkipWebInstall) { Invoke-NpmInstallWithRepair -Clean:$CleanWebDependencies }
      else { Write-Host 'Skipping npm install because -SkipWebInstall was supplied.' -ForegroundColor Yellow }
      Assert-WebPackageToolchain
      if (-not $SkipWebAudit) { Invoke-Step 'NPM audit' { Invoke-Native 'npm' @('audit', '--audit-level=high') } }
      else { Write-Host 'Skipping npm audit because -SkipWebAudit was supplied.' -ForegroundColor Yellow }
      if (-not $SkipWebTypecheck) { Invoke-Step 'Web typecheck' { Invoke-Native 'npm' @('run', 'typecheck') } }
      else { Write-Host 'Skipping web typecheck because -SkipWebTypecheck was supplied.' -ForegroundColor Yellow }
      if (-not $SkipWebLint) { Invoke-Step 'Web lint' { Invoke-Native 'npm' @('run', 'lint') } }
      else { Write-Host 'Skipping web lint because -SkipWebLint was supplied.' -ForegroundColor Yellow }
      if (-not $SkipWebBuild) { Invoke-Step 'Web build' { Invoke-Native 'npm' @('run', 'build') } }
      else { Write-Host 'Skipping web build because -SkipWebBuild was supplied.' -ForegroundColor Yellow }
    } else { Write-Host 'package.json not found; skipping web build' -ForegroundColor Yellow }
  } else { Write-Host "Web folder not found; skipping web build: $WebRoot" -ForegroundColor Yellow }

  Invoke-Step 'Git auto-commit successful changes' { Invoke-AutoCommitIfNeeded }
  Invoke-Step 'Git auto-push committed changes' { Invoke-AutoPushIfNeeded }
  $UpdateSucceeded = $true
}
catch {
  $ExitCode = 1
  Write-Host ""
  Write-Host 'Update failed.' -ForegroundColor Red
  Write-Host $_.Exception.Message -ForegroundColor Red
  $message = $_.Exception.Message
  if ($message -match 'npm install|npm ci|EPERM|EBUSY|ENOTEMPTY|unlink|rmdir|node_modules|tsc|eslint|vite') { Show-NpmLockHelp }
}
finally {
  if ((Test-Path variable:RepoRoot) -and (Test-Path $RepoRoot)) { Set-Location $RepoRoot }
  else { Set-Location $InitialLocation }
  Write-Host ""
  Write-Host "Current folder: $(Get-Location)"
}

if ($UpdateSucceeded) {
  Write-Host ""
  Write-Host 'Update complete.' -ForegroundColor Green
  Write-Host "Backend binary: $BackendBinary"
  Write-Host "Repo root exe:   $RepoRootBinary"
  if (-not $SkipWebBuild) { Write-Host "Frontend build:  $WebRoot\dist" }
  Write-Host "Copied assets:   $BuildOutputDir\.env.example, $BuildOutputDir\README.md"
  if (-not [string]::IsNullOrWhiteSpace($AutoCommitSha)) { Write-Host "Auto-commit:     $AutoCommitSha" }
} else { exit $ExitCode }
