# Runtime Guide: Windows Service

Use this only after the Dune: Awakening server or related services are confirmed to be managed as Windows services.

## 1. Find Dune-Related Services

Run on: Windows host PowerShell

```powershell
Get-Service | Where-Object { $_.Name -match 'dune|awakening|sandbox|amp|docker|rabbit' -or $_.DisplayName -match 'dune|awakening|sandbox|amp|docker|rabbit' } | Select-Object Name, DisplayName, Status
```

Record:

```text
SERVICE_NAME=
DIRECTOR_SERVICE=
RABBITMQ_SERVICE=
DESTINATION_SERVICE=
```

## 2. Check Service Status

Run on: Windows host PowerShell

```powershell
Get-Service -Name $env:SERVICE_NAME
```

If multiple services exist, repeat for each service.

## 3. Identify Service Configuration

Run on: Windows host PowerShell

```powershell
Get-CimInstance Win32_Service | Where-Object { $_.Name -eq $env:SERVICE_NAME } | Select-Object Name, DisplayName, State, StartMode, StartName, PathName
```

Record:

```text
Service account:
Startup command:
Working directory, if known:
Executable path:
```

## 4. Capture Event Logs

Run on: Windows host PowerShell

```powershell
Get-WinEvent -LogName Application -MaxEvents 200 | Where-Object { $_.Message -match 'Dune|Awakening|Sandbox|Rabbit|AMP' } | Format-List TimeCreated, ProviderName, Id, LevelDisplayName, Message
```

## 5. Capture Process and Listener State

Run on: Windows host PowerShell

```powershell
Get-Process | Where-Object { $_.ProcessName -match 'Dune|Sandbox|AMP|Rabbit|Docker' } | Select-Object ProcessName, Id, Path
Get-NetUDPEndpoint | Sort-Object LocalPort | Format-Table LocalAddress, LocalPort, OwningProcess
Get-NetTCPConnection | Sort-Object LocalPort | Format-Table LocalAddress, LocalPort, State, OwningProcess
```

## 6. Stop Before Changing Anything

Do not restart, edit, or delete the service until service status, event logs, process state, listener state, and the user-defined failure window are captured.
