Set-StrictMode -Version Latest
$ErrorActionPreference = 'Stop'
$ProgressPreference = 'SilentlyContinue'

$RepoRoot = 'Z:\dune-admin-fork'
$WebRoot = Join-Path $RepoRoot 'web'

function Run-Step {
  param(
    [string]$Name,
    [scriptblock]$Command
  )
  Write-Host ""
  Write-Host "=== $Name ===" -ForegroundColor Cyan
  & $Command
  if ($LASTEXITCODE -ne 0) {
    throw "$Name failed with exit code $LASTEXITCODE"
  }
}

if (-not (Test-Path $RepoRoot)) {
  throw "Repository root not found: $RepoRoot"
}

Set-Location $RepoRoot
Write-Host "Repo folder: $(Get-Location)"

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

Set-Location $RepoRoot
Write-Host ""
Write-Host 'Update complete.' -ForegroundColor Green
Write-Host "Backend binary: $RepoRoot\dune-admin.exe"
Write-Host "Frontend build:  $WebRoot\dist"
