function Invoke-WebValidationAndBuild {
  param(
    [Parameter(Mandatory = $true)][string]$WebRoot,
    [switch]$CleanWebDependencies,
    [switch]$SkipWebInstall,
    [switch]$SkipWebAudit,
    [switch]$SkipWebTypecheck,
    [switch]$SkipWebLint,
    [switch]$SkipWebBuild
  )

  if (-not (Test-Path $WebRoot)) {
    Write-Host "Web folder not found; skipping web build: $WebRoot" -ForegroundColor Yellow
    return
  }

  Set-Location $WebRoot
  Write-Host "Web folder:     $(Get-Location)"

  if (-not (Test-Path 'package.json')) {
    Write-Host 'package.json not found; skipping web build' -ForegroundColor Yellow
    return
  }

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
}
