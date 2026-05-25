# Platform Guide: Windows and Hyper-V

Use this when the server runs on Windows directly, inside a Windows VM, or inside a guest VM managed by Hyper-V.

## 1. Confirm Windows Context

Run on: Windows host PowerShell

```powershell
$env:COMPUTERNAME
whoami
Get-Location
Get-ComputerInfo | Select-Object CsName, WindowsProductName, WindowsVersion, OsHardwareAbstractionLayer, CsSystemType
systeminfo
```

Record:

```text
Hostname:
Windows version:
Logged-in user:
Is this a host or guest VM:
```

## 2. Check Hyper-V Only If This Is a Hyper-V Host

Run on: Hyper-V host PowerShell

```powershell
Get-VM | Select-Object Name, State, CPUUsage, MemoryAssigned, Uptime
Get-VMNetworkAdapter -VMName * | Select-Object VMName, SwitchName, IPAddresses
Get-VMSwitch | Select-Object Name, SwitchType, NetAdapterInterfaceDescription
```

Record:

```text
VM name:
VM state:
VM IP address:
Virtual switch:
```

## 3. Find the Server Management Layer

Run on: Windows host PowerShell

```powershell
Get-Process | Where-Object { $_.ProcessName -match 'amp|docker|Dune|Awakening|Sandbox' } | Select-Object ProcessName, Id, Path
Get-Service | Where-Object { $_.Name -match 'amp|docker|dune|awakening' -or $_.DisplayName -match 'amp|docker|dune|awakening' } | Select-Object Name, DisplayName, Status
```

Interpretation:

```text
Docker visible:
  Use the Docker runtime guide.

AMP visible:
  Use the AMP runtime guide.

Windows service visible:
  Use the Windows service runtime guide.

Dune process visible without a known service:
  Use the manual/custom script runtime guide.
```

## 4. Find Install and Log Paths

Run on: Windows host PowerShell

```powershell
Get-ChildItem -Path C:\ -Recurse -ErrorAction SilentlyContinue -Include *Engine.ini,UserEngine.ini,*.log | Where-Object { $_.FullName -match 'Dune|Awakening|DuneSandbox' } | Select-Object -First 100 FullName
```

Record:

```text
INSTANCE_PATH=
LOG_PATH=
SAVED_PATH=
```

## 5. Capture Listener State

Run on: Windows host PowerShell

```powershell
Get-NetUDPEndpoint | Sort-Object LocalPort | Format-Table LocalAddress, LocalPort, OwningProcess
Get-NetTCPConnection | Sort-Object LocalPort | Format-Table LocalAddress, LocalPort, State, OwningProcess
```

## 6. Next Runbooks

- Runtime layer found as AMP: [AMP control panel](../runtimes/amp-control-panel.md)
- Runtime layer found as Docker: [Docker or Docker Compose](../runtimes/docker-or-compose.md)
- Runtime layer found as Windows service: [Windows service](../runtimes/windows-service.md)
- Unknown/custom launch: [Manual or custom script](../runtimes/manual-or-custom-script.md)
