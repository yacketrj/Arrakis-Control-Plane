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
