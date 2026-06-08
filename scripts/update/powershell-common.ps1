function Write-StepStatus {
  param(
    [Parameter(Mandatory = $true)][ValidateSet('RUN', 'PASS', 'FAIL', 'WARN')][string]$Status,
    [Parameter(Mandatory = $true)][string]$Message
  )

  $color = switch ($Status) {
    'RUN' { 'Cyan' }
    'PASS' { 'Green' }
    'FAIL' { 'Red' }
    'WARN' { 'Yellow' }
  }

  Write-Host ("{0}: {1}" -f $Status, $Message) -ForegroundColor $color
}

function Write-Section {
  param([Parameter(Mandatory = $true)][string]$Name)
  Write-Host ""
  Write-StepStatus -Status 'RUN' -Message $Name
}

function Invoke-Step {
  param([Parameter(Mandatory = $true)][string]$Name, [Parameter(Mandatory = $true)][scriptblock]$Command)
  Write-Section $Name
  try {
    & $Command
    Write-StepStatus -Status 'PASS' -Message $Name
  }
  catch {
    Write-StepStatus -Status 'FAIL' -Message $Name
    throw
  }
}

function Invoke-Native {
  param([Parameter(Mandatory = $true)][string]$FilePath, [string[]]$Arguments = @())
  Write-StepStatus -Status 'RUN' -Message ("{0} {1}" -f $FilePath, ($Arguments -join ' '))
  & $FilePath @Arguments
  $nativeExitCode = $LASTEXITCODE
  if ($null -ne $nativeExitCode -and $nativeExitCode -ne 0) { throw "$FilePath $($Arguments -join ' ') failed with exit code $nativeExitCode" }
}

function Update-ProcessPath {
  $machinePath = [Environment]::GetEnvironmentVariable('Path', 'Machine')
  $userPath = [Environment]::GetEnvironmentVariable('Path', 'User')
  $common = @(
    "$env:ProgramFiles\Git\cmd",
    "$env:ProgramFiles\Git\bin",
    "$env:ProgramFiles\Go\bin",
    "$env:ProgramFiles\nodejs",
    "$env:LOCALAPPDATA\Microsoft\WindowsApps"
  )
  $env:Path = (@($env:Path, $machinePath, $userPath) + $common | Where-Object { -not [string]::IsNullOrWhiteSpace($_) }) -join ';'
}

function Install-PrerequisiteForCommand {
  param([Parameter(Mandatory = $true)][string]$Name)

  if ($SkipPrereqInstall) { throw "$Name was not found on PATH and prerequisite auto-install is disabled by -SkipPrereqInstall." }
  if (-not (Get-Command winget -ErrorAction SilentlyContinue)) { throw "$Name was not found on PATH and winget is unavailable. Install the missing prerequisite manually." }

  $packageId = switch ($Name) {
    'git' { 'Git.Git' }
    'go' { 'GoLang.Go' }
    'node' { 'OpenJS.NodeJS.LTS' }
    'npm' { 'OpenJS.NodeJS.LTS' }
    'gh' { 'GitHub.cli' }
    default { throw "No auto-install mapping exists for missing command: $Name" }
  }

  Invoke-Step "Install prerequisite: $Name" { Invoke-Native 'winget' @('install', '--id', $packageId, '-e', '--accept-package-agreements', '--accept-source-agreements') }
  Update-ProcessPath
}

function Assert-CommandAvailable {
  param([Parameter(Mandatory = $true)][string]$Name, [string]$InstallHint = '')
  Update-ProcessPath
  if (Get-Command $Name -ErrorAction SilentlyContinue) { return }

  Write-StepStatus -Status 'WARN' -Message "$Name was not found on PATH. Attempting prerequisite auto-install."
  Install-PrerequisiteForCommand -Name $Name
  if (Get-Command $Name -ErrorAction SilentlyContinue) { return }

  $message = "$Name is still unavailable after auto-install."
  if (-not [string]::IsNullOrWhiteSpace($InstallHint)) { $message = "$message $InstallHint" }
  throw $message
}

function Resolve-OutputDirectory {
  param([Parameter(Mandatory = $true)][string]$Root, [string]$RequestedOutputDir)
  if ([string]::IsNullOrWhiteSpace($RequestedOutputDir)) { return (Join-Path $Root 'dist\windows') }
  if ([System.IO.Path]::IsPathRooted($RequestedOutputDir)) { return $RequestedOutputDir }
  return (Join-Path $Root $RequestedOutputDir)
}