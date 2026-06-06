Set-StrictMode -Version Latest

function Write-Section {
  param([Parameter(Mandatory = $true)][string]$Name)
  Write-Host ""
  Write-Host "=== $Name ===" -ForegroundColor Cyan
}

function Invoke-Step {
  param([Parameter(Mandatory = $true)][string]$Name, [Parameter(Mandatory = $true)][scriptblock]$Command)
  Write-Section $Name
  & $Command
}

function Invoke-Native {
  param([Parameter(Mandatory = $true)][string]$FilePath, [string[]]$Arguments = @())
  Write-Host ((">>> {0} {1}" -f $FilePath, ($Arguments -join ' '))) -ForegroundColor DarkGray
  & $FilePath @Arguments
  $nativeExitCode = $LASTEXITCODE
  if ($null -ne $nativeExitCode -and $nativeExitCode -ne 0) {
    throw "$FilePath $($Arguments -join ' ') failed with exit code $nativeExitCode"
  }
}

function Resolve-OutputDirectory {
  param([Parameter(Mandatory = $true)][string]$Root, [string]$RequestedOutputDir)
  if ([string]::IsNullOrWhiteSpace($RequestedOutputDir)) { return (Join-Path $Root 'dist\windows') }
  if ([System.IO.Path]::IsPathRooted($RequestedOutputDir)) { return $RequestedOutputDir }
  return (Join-Path $Root $RequestedOutputDir)
}
