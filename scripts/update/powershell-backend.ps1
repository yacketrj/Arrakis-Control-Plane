function Get-BackendGoPackages {
  $packages = @(& go list ./...)
  if ($LASTEXITCODE -ne 0) { throw 'Unable to list Go packages.' }
  return @($packages | Where-Object { $_ -notmatch '/web/node_modules/' -and $_ -notmatch '/web/dist/' })
}

function Invoke-BackendGoTests {
  $packages = @(Get-BackendGoPackages)
  if ($packages.Count -eq 0) { throw 'No backend Go packages found to test.' }
  Invoke-Native 'go' (@('test', '-v') + $packages)
}

function Invoke-BackendBuild {
  param(
    [Parameter(Mandatory = $true)][string]$BuildOutputDir,
    [Parameter(Mandatory = $true)][string]$BackendBinary,
    [Parameter(Mandatory = $true)][string]$Version
  )

  Invoke-Step 'Go backend build' {
    New-Item -ItemType Directory -Force -Path $BuildOutputDir | Out-Null
    $ldflags = "-s -w -X main.version=$Version"
    Invoke-Native 'go' @('build', '-trimpath', '-ldflags', $ldflags, '-o', $BackendBinary, '.')
  }

  if (-not (Test-Path $BackendBinary)) { throw "Backend build completed, but expected binary was not found: $BackendBinary" }
  $binaryInfo = Get-Item -LiteralPath $BackendBinary
  Write-Host ("Backend build output: {0} ({1} bytes, modified {2})" -f $binaryInfo.FullName, $binaryInfo.Length, $binaryInfo.LastWriteTime) -ForegroundColor Green
}

function Copy-BackendBinaryAndAssets {
  param(
    [Parameter(Mandatory = $true)][string]$RepoRoot,
    [Parameter(Mandatory = $true)][string]$BuildOutputDir,
    [Parameter(Mandatory = $true)][string]$BackendBinary,
    [Parameter(Mandatory = $true)][string]$RepoRootBinary,
    [switch]$SkipRootBinaryCopy
  )

  if (-not $SkipRootBinaryCopy) {
    Invoke-Step 'Copy backend binary to repo root' {
      Copy-Item -Force -LiteralPath $BackendBinary -Destination $RepoRootBinary
      $rootBinaryInfo = Get-Item -LiteralPath $RepoRootBinary
      Write-Host ("Repo root binary:    {0} ({1} bytes, modified {2})" -f $rootBinaryInfo.FullName, $rootBinaryInfo.Length, $rootBinaryInfo.LastWriteTime) -ForegroundColor Green
    }
  } else {
    Write-Host 'Skipping repo-root binary copy because -SkipRootBinaryCopy was supplied.' -ForegroundColor Yellow
  }

  foreach ($asset in @('.env.example', 'README.md')) {
    $source = Join-Path $RepoRoot $asset
    if (Test-Path $source) { Copy-Item -Force $source (Join-Path $BuildOutputDir $asset) }
  }
}
