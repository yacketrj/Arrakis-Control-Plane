Set-StrictMode -Version Latest
$ErrorActionPreference = 'Stop'
$ProgressPreference = 'SilentlyContinue'

$RepoRoot = 'Z:\dune-admin-fork'
$WebRoot = Join-Path $RepoRoot 'web'
$InitialLocation = Get-Location
$UpdateSucceeded = $false
$ExitCode = 0

function Run-Step {
  param(
    [Parameter(Mandatory = $true)][string]$Name,
    [Parameter(Mandatory = $true)][scriptblock]$Command
  )

  Write-Host ""
  Write-Host "=== $Name ===" -ForegroundColor Cyan
  & $Command
  $stepExitCode = $LASTEXITCODE
  if ($null -ne $stepExitCode -and $stepExitCode -ne 0) {
    throw "$Name failed with exit code $stepExitCode"
  }
}

try {
  if (-not (Test-Path $RepoRoot)) {
    throw "Repository root not found: $RepoRoot"
  }

  Set-Location $RepoRoot
  Write-Host "Repo folder: $(Get-Location)"

  if (-not (Test-Path '.git')) {
    throw "Not a Git repository: $RepoRoot"
  }

  Run-Step 'Git pull' { git pull }
  Run-Step 'Go tests' { go test -v ./... }
  Run-Step 'Go backend build' { go build -o dune-admin.exe . }

  if (-not (Test-Path (Join-Path $RepoRoot 'dune-admin.exe'))) {
    throw 'Backend build completed, but dune-admin.exe was not found.'
  }

  if (Test-Path $WebRoot) {
    Set-Location $WebRoot
    Write-Host "Web folder: $(Get-Location)"

    if (Test-Path 'package.json') {
      if (Test-Path 'package-lock.json') {
        Run-Step 'NPM install clean' { npm ci }
      } else {
        Run-Step 'NPM install' { npm install }
      }
      Run-Step 'Web build' { npm run build }
    } else {
      Write-Host 'package.json not found; skipping web build' -ForegroundColor Yellow
    }
  } else {
    Write-Host "Web folder not found; skipping web build: $WebRoot" -ForegroundColor Yellow
  }

  $UpdateSucceeded = $true
}
catch {
  $ExitCode = 1
  Write-Host ""
  Write-Host 'Update failed.' -ForegroundColor Red
  Write-Host $_.Exception.Message -ForegroundColor Red
}
finally {
  if (Test-Path $RepoRoot) {
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
  Write-Host "Backend binary: $RepoRoot\dune-admin.exe"
  Write-Host "Frontend build:  $WebRoot\dist"
} else {
  exit $ExitCode
}
