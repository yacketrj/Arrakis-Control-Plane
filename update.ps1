param(
  [switch]$CleanWebDependencies,
  [switch]$SkipWebInstall
)

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

function Show-NpmLockHelp {
  Write-Host ""
  Write-Host 'NPM dependency update failed.' -ForegroundColor Red
  Write-Host 'On Windows this is commonly caused by a locked native .node file in node_modules.' -ForegroundColor Yellow
  Write-Host 'Close any running frontend dev server, node.exe process, editor terminal, or antivirus scan that may be holding web\node_modules files open, then re-run the script.' -ForegroundColor Yellow
  Write-Host 'Normal updates now use npm install instead of npm ci to avoid deleting locked native dependencies.' -ForegroundColor Yellow
  Write-Host 'Use .\update.ps1 -CleanWebDependencies only when you intentionally want a clean npm ci install.' -ForegroundColor Yellow
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
      node --version
      npm --version

      if (-not $SkipWebInstall) {
        if ($CleanWebDependencies) {
          Run-Step 'NPM clean install' { npm ci }
        } else {
          Run-Step 'NPM install' { npm install }
        }
      } else {
        Write-Host 'Skipping npm install because -SkipWebInstall was supplied.' -ForegroundColor Yellow
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

  $message = $_.Exception.Message
  if ($message -match 'NPM|npm|EPERM|unlink|node_modules') {
    Show-NpmLockHelp
  }
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
