# Live Support Call Checklist

Use this checklist during a live support call or screen-share session. It is designed for operators who may not know the hosting platform, runtime, or correct command location yet.

Do not collect personal names, chat names, player names, raw account identifiers, passwords, authentication tokens, private keys, or unrelated customer details in reusable troubleshooting notes.

## 1. Start the Case

Record these first:

```text
Case title:
Support operator:
Environment owner or contact role:
Date/time UTC:
Reported issue in the user's own words:
What worked before:
What fails now:
When it last worked:
When it first failed:
Recent changes:
```

## 2. Confirm Where You Are Working

Ask the person sharing their screen which interface is open.

```text
Are you in a control panel web UI?
Are you on a Linux shell?
Are you on a Windows PowerShell window?
Are you on a hypervisor host such as Proxmox or Hyper-V?
Are you inside a VM?
Are you inside a container shell?
Are you in a cloud provider console?
```

Decision rule:

```text
If you do not know where the command will run, do not run it yet.
Use environment discovery first.
```

## 3. Identify the Platform

Run only the command that matches the current screen.

Run on: Linux host or Linux VM shell

```bash
hostnamectl
cat /etc/os-release
uname -a
systemd-detect-virt -v || true
```

Run on: Windows host PowerShell

```powershell
Get-ComputerInfo | Select-Object CsName, WindowsProductName, WindowsVersion, CsSystemType, HyperVisorPresent
```

Run in: control panel UI

```text
Record the product name, instance name, instance status, visible install path, log path, console output, and start/stop controls.
```

Run on: Proxmox host shell

```bash
pveversion
qm list
pct list
```

Run on: Hyper-V host PowerShell

```powershell
Get-VM | Select-Object Name, State, CPUUsage, MemoryAssigned, Uptime
Get-VMSwitch | Select-Object Name, SwitchType
```

## 4. Identify the Runtime or Orchestration Layer

Run on: Linux host or Linux VM shell

```bash
command -v docker && docker ps --format 'table {{.Names}}\t{{.Image}}\t{{.Status}}\t{{.Ports}}' || true
systemctl list-units --type=service --all | grep -Ei 'dune|awakening|sandbox|amp|docker|compose|rabbit|postgres' || true
ps -ef | grep -Ei 'DuneSandbox|Awakening|Sandbox|AMP|docker|rabbit|postgres' | grep -v grep || true
```

Run on: Windows host PowerShell

```powershell
Get-Service | Where-Object { $_.Name -match 'dune|awakening|sandbox|amp|docker|rabbit|postgres' -or $_.DisplayName -match 'dune|awakening|sandbox|amp|docker|rabbit|postgres' } | Select-Object Name, DisplayName, Status
Get-Process | Where-Object { $_.ProcessName -match 'Dune|Awakening|Sandbox|AMP|Docker|Rabbit|postgres' } | Select-Object ProcessName, Id, Path
```

Record:

```text
Hosting platform:
Runtime/orchestration:
Control panel, if any:
Guest OS:
Container runtime, if any:
```

## 5. Choose the Symptom Path

```text
Server does not start:
  Use Server Startup Failure.

Server not visible:
  Use Server Visibility and Listing.

Player cannot log in:
  Use Login and Authentication Failure.

Player can log in but travel hangs or fails:
  Use Map Travel and Instancing Failure.

Dynamically spawned destination fails:
  Use Dynamic Instancing and Handoff Validation.

Permission denied or file write failure:
  Use Permission and Ownership Errors.

Lag, stalls, or resource pressure:
  Use Resource and Performance Checks.

Crash, hang, or unexpected restart:
  Use Crash, Hang, and Process Exit Analysis.
```

## 6. Capture Minimum Evidence Before Changing Anything

Run on: Linux host or Linux VM shell

```bash
date -u
free -h
df -h
sudo ss -tulpen | tee listeners-before.txt
ps -ef | grep -Ei 'DuneSandbox|Awakening|Sandbox|AMP|docker|rabbit|postgres' | grep -v grep | tee processes-before.txt || true
```

Run on: Windows host PowerShell

```powershell
Get-Date -Format u
Get-Volume | Select-Object DriveLetter, FileSystemLabel, SizeRemaining, Size
Get-NetUDPEndpoint | Sort-Object LocalPort | Tee-Object -FilePath udp-listeners-before.txt
Get-NetTCPConnection | Sort-Object LocalPort | Tee-Object -FilePath tcp-listeners-before.txt
```

Run on: Docker host shell, only if Docker is confirmed

```bash
docker ps --format 'table {{.Names}}\t{{.Image}}\t{{.Status}}\t{{.Ports}}' | tee docker-ps-before.txt
```

## 7. Stop Conditions

Stop and escalate when:

```text
The operator cannot identify the platform or runtime.
The command location is unclear.
The next action risks deleting data or changing persistent files.
A service restart would interrupt active players without approval.
The evidence shows a closed-source or vendor-owned component failing internally.
The operator is asked for credentials, tokens, or secrets they are not authorized to use.
```

## 8. End the Call With a Clear Handoff

Record:

```text
What was confirmed:
What remains unknown:
Evidence files collected:
Exact UTC test window:
Next recommended runbook:
Next owner:
Risk before next action:
```
