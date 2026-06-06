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

  if ($script:SkipPrereqInstall) { throw "$Name was not found on PATH and prerequisite auto-install is disabled by -SkipPrereqInstall." }
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

  Write-Host "$Name was not found on PATH. Attempting prerequisite auto-install." -ForegroundColor Yellow
  Install-PrerequisiteForCommand -Name $Name
  if (Get-Command $Name -ErrorAction SilentlyContinue) { return }

  $message = "$Name is still unavailable after auto-install."
  if (-not [string]::IsNullOrWhiteSpace($InstallHint)) { $message = "$message $InstallHint" }
  throw $message
}
