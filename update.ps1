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

$WebHelpers = Join-Path $PSScriptRoot 'scripts\update\powershell-web.ps1'
if (-not (Test-Path $WebHelpers)) { throw "PowerShell update web helper not found: $WebHelpers" }
. $WebHelpers

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

  Invoke-WebValidationAndBuild `
    -WebRoot $WebRoot `
    -CleanWebDependencies:$CleanWebDependencies `
    -SkipWebInstall:$SkipWebInstall `
    -SkipWebAudit:$SkipWebAudit `
    -SkipWebTypecheck:$SkipWebTypecheck `
    -SkipWebLint:$SkipWebLint `
    -SkipWebBuild:$SkipWebBuild

  Set-Location $RepoRoot
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
