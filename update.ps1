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

if ([string]::IsNullOrWhiteSpace($RepoRoot)) { $RepoRoot = $PSScriptRoot }
if ([string]::IsNullOrWhiteSpace($Version)) { $Version = if ($env:VERSION) { $env:VERSION } else { 'dev' } }
if ([string]::IsNullOrWhiteSpace($CommitMessage)) { $CommitMessage = if ($env:COMMIT_MESSAGE) { $env:COMMIT_MESSAGE } else { 'Automated successful update' } }

$InitialLocation = Get-Location
$UpdateSucceeded = $false
$ExitCode = 0
$AutoCommitSha = ''

function Get-NodeProcessSummary {
  $processes = @(Get-Process -Name 'node' -ErrorAction SilentlyContinue)
  if ($processes.Count -eq 0) { Write-Host 'No running node.exe processes were detected.' -ForegroundColor Yellow; return }
  Write-Host 'Running node.exe processes detected:' -ForegroundColor Yellow
  foreach ($process in $processes | Select-Object -First 12) {
    $path = 'unknown path'
    try { if ($process.Path) { $path = $process.Path } } catch { $path = 'access denied reading process path' }
    Write-Host ("  PID {0}: {1}" -f $process.Id, $path) -ForegroundColor Yellow
  }
  if ($processes.Count -gt 12) { Write-Host "  ... $($processes.Count - 12) more node.exe processes" -ForegroundColor Yellow }
}

function Remove-NodeModulesForRepair {
  param([Parameter(Mandatory = $true)][string]$WebRootPath)
  $nodeModules = Join-Path $WebRootPath 'node_modules'
  if (-not (Test-Path $nodeModules)) { Write-Host 'node_modules is not present; repair will perform a fresh install.' -ForegroundColor Yellow; return }
  Write-Host 'Removing web node_modules for automatic dependency repair...' -ForegroundColor Yellow
  try { Remove-Item -LiteralPath $nodeModules -Recurse -Force -ErrorAction Stop }
  catch {
    Get-NodeProcessSummary
    throw "Unable to remove web node_modules automatically. A process or security scanner is probably locking files under node_modules. Original error: $($_.Exception.Message)"
  }
}

function Invoke-NpmInstallWithRepair {
  param([switch]$Clean)
  $installLabel = if ($Clean) { 'NPM clean install' } else { 'NPM install' }
  $installArgs = if ($Clean) { @('ci') } else { @('install') }
  try { Invoke-Step $installLabel { Invoke-Native 'npm' $installArgs }; return }
  catch { Write-Host "Initial $installLabel failed: $($_.Exception.Message)" -ForegroundColor Yellow }
  if ($SkipNpmRepair) { throw "$installLabel failed and automatic npm repair is disabled by -SkipNpmRepair." }
  Write-Host 'Attempting npm recovery: cache verify, retry, then dependency repair if needed.' -ForegroundColor Yellow
  try {
    Invoke-Step 'NPM cache verify' { Invoke-Native 'npm' @('cache', 'verify') }
    Invoke-Step "$installLabel retry" { Invoke-Native 'npm' $installArgs }
    return
  } catch { Write-Host "NPM retry failed: $($_.Exception.Message)" -ForegroundColor Yellow }
  Get-NodeProcessSummary
  Remove-NodeModulesForRepair -WebRootPath $WebRoot
  Invoke-Step "$installLabel after dependency repair" { Invoke-Native 'npm' $installArgs }
}

function Test-WebPackageBinary {
  param([Parameter(Mandatory = $true)][string]$Name)
  $binDir = Join-Path $WebRoot 'node_modules\.bin'
  return (
    (Test-Path (Join-Path $binDir $Name)) -or
    (Test-Path (Join-Path $binDir "$Name.cmd")) -or
    (Test-Path (Join-Path $binDir "$Name.ps1"))
  )
}

function Assert-WebPackageToolchain {
  $required = @('tsc', 'eslint', 'vite')
  $missing = @($required | Where-Object { -not (Test-WebPackageBinary -Name $_) })

  if ($missing.Count -eq 0) {
    Write-StepStatus -Status 'PASS' -Message "Web package toolchain present: $($required -join ', ')"
    return
  }

  Write-StepStatus -Status 'WARN' -Message "Missing web package toolchain: $($missing -join ', ')"
  Write-StepStatus -Status 'RUN' -Message 'Running npm install to restore local package binaries.'
  Invoke-NpmInstallWithRepair

  $missing = @($required | Where-Object { -not (Test-WebPackageBinary -Name $_) })
  if ($missing.Count -gt 0) {
    throw "Missing web package toolchain after npm install: $($missing -join ', '). Verify web/package.json devDependencies and package-lock.json, then rerun update.ps1."
  }

  Write-StepStatus -Status 'PASS' -Message "Web package toolchain restored: $($required -join ', ')"
}

function Show-NpmLockHelp {
  Write-Host ""
  Write-Host 'NPM dependency update failed after automatic recovery attempts.' -ForegroundColor Red
  Write-Host 'The script already retried npm and attempted node_modules repair unless -SkipNpmRepair was supplied.' -ForegroundColor Yellow
  Write-Host 'If removal failed, a running process or security scanner is still locking web\node_modules files.' -ForegroundColor Yellow
  Write-Host 'Check the node.exe process list printed above, then close the listed process or pause the scanner and rerun the script.' -ForegroundColor Yellow
  Write-Host 'You can bypass dependency installation only when dependencies are already valid by running: .\update.ps1 -SkipWebInstall' -ForegroundColor Yellow
}

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

  if (-not $SkipGoTests) { Invoke-Step 'Go tests' { Invoke-Native 'go' @('test', '-v', './...') } }
  else { Write-Host 'Skipping Go tests because -SkipGoTests was supplied.' -ForegroundColor Yellow }

  Invoke-Step 'Go backend build' {
    New-Item -ItemType Directory -Force -Path $BuildOutputDir | Out-Null
    $ldflags = "-s -w -X main.version=$Version"
    Invoke-Native 'go' @('build', '-trimpath', '-ldflags', $ldflags, '-o', $BackendBinary, '.')
  }

  if (-not (Test-Path $BackendBinary)) { throw "Backend build completed, but expected binary was not found: $BackendBinary" }
  $binaryInfo = Get-Item -LiteralPath $BackendBinary
  Write-Host ("Backend build output: {0} ({1} bytes, modified {2})" -f $binaryInfo.FullName, $binaryInfo.Length, $binaryInfo.LastWriteTime) -ForegroundColor Green

  if (-not $SkipRootBinaryCopy) {
    Invoke-Step 'Copy backend binary to repo root' {
      Copy-Item -Force -LiteralPath $BackendBinary -Destination $RepoRootBinary
      $rootBinaryInfo = Get-Item -LiteralPath $RepoRootBinary
      Write-Host ("Repo root binary:    {0} ({1} bytes, modified {2})" -f $rootBinaryInfo.FullName, $rootBinaryInfo.Length, $rootBinaryInfo.LastWriteTime) -ForegroundColor Green
    }
  } else { Write-Host 'Skipping repo-root binary copy because -SkipRootBinaryCopy was supplied.' -ForegroundColor Yellow }

  foreach ($asset in @('.env.example', 'README.md')) {
    $source = Join-Path $RepoRoot $asset
    if (Test-Path $source) { Copy-Item -Force $source (Join-Path $BuildOutputDir $asset) }
  }

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
