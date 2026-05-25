# Runbook: Crash, Hang, and Process Exit Analysis

Use this runbook when a server process crashes, exits unexpectedly, hangs, stops responding, or is killed by the operating system or container runtime.

Start only after the hosting platform and runtime/orchestration layer are identified.

## 1. Capture the Failure Type

Ask the user or environment owner:

```text
Did the service crash, hang, restart, or stop responding?
Did the control panel show crashed, stopped, restarting, or unknown?
Did the process disappear from the process list?
Did the process stay running but players could not use it?
Did this happen once or repeatedly?
Exact UTC time of the crash or hang:
```

Record:

```text
Failure type:
Affected service/process:
UTC time:
Last known good time:
Recent change:
```

## 2. Check Process State

Run on: Linux host or Linux VM shell

```bash
ps -ef | grep -Ei 'DuneSandbox|Awakening|Sandbox|AMP|Rabbit|docker|compose' | grep -v grep | tee crash-processes.txt || true
```

Run on: Windows host PowerShell

```powershell
Get-Process | Where-Object { $_.ProcessName -match 'Dune|Awakening|Sandbox|AMP|Rabbit|Docker' } | Select-Object ProcessName, Id, Path, StartTime | Tee-Object -FilePath crash-processes.txt
```

Run on: Docker host shell, only if Docker is confirmed

```bash
docker ps -a --format 'table {{.Names}}\t{{.Image}}\t{{.Status}}\t{{.Ports}}' | tee crash-docker-ps.txt
```

## 3. Check System Logs Around the Failure

Run on: Linux host or Linux VM shell

```bash
journalctl --since "60 minutes ago" --no-pager | grep -Ei 'Dune|Awakening|Sandbox|killed|oom|segfault|crash|fatal|core|docker|container|rabbit' | tee crash-system-journal.txt

dmesg -T | grep -Ei 'killed|oom|segfault|out of memory|Dune|Sandbox|docker|container' | tail -200 | tee crash-dmesg.txt
```

Run on: Windows host PowerShell

```powershell
Get-WinEvent -LogName Application -MaxEvents 500 | Where-Object { $_.Message -match 'Dune|Awakening|Sandbox|crash|fault|exception|Rabbit|Docker|AMP' } | Format-List TimeCreated, ProviderName, Id, LevelDisplayName, Message | Tee-Object -FilePath crash-application-events.txt

Get-WinEvent -LogName System -MaxEvents 500 | Where-Object { $_.Message -match 'service|crash|terminated|stopped|resource|memory|disk|network' } | Format-List TimeCreated, ProviderName, Id, LevelDisplayName, Message | Tee-Object -FilePath crash-system-events.txt
```

## 4. Check Runtime Logs

Run on: Docker host shell, only if Docker is confirmed

```bash
docker logs --since 60m "$CONTAINER_NAME" 2>&1 | tee crash-container.log
docker inspect "$CONTAINER_NAME" --format 'Status={{.State.Status}} ExitCode={{.State.ExitCode}} OOMKilled={{.State.OOMKilled}} Error={{.State.Error}} Started={{.State.StartedAt}} Finished={{.State.FinishedAt}}' | tee crash-container-state.txt
```

Run on: Linux systemd host, only if systemd is used

```bash
systemctl status "$SERVICE_NAME" --no-pager | tee crash-systemd-status.txt
journalctl -u "$SERVICE_NAME" --since "60 minutes ago" --no-pager | tee crash-systemd.log
```

Run in: control panel UI

```text
Export console output, task output, restart history, and crash status covering the failure window.
```

## 5. Check Resource Pressure

Run on: Linux host or Linux VM shell

```bash
free -h
df -h
uptime
```

Run on: Windows host PowerShell

```powershell
Get-Volume | Select-Object DriveLetter, FileSystemLabel, SizeRemaining, Size
Get-CimInstance Win32_OperatingSystem | Select-Object TotalVisibleMemorySize, FreePhysicalMemory, LastBootUpTime
```

If resource pressure appears, continue with [Performance, Resource, and Capacity Checks](./performance-resource-and-capacity-checks.md).

## 6. Interpret Results

```text
Process exited with non-zero exit code:
  Capture runtime logs and launch arguments.

Container shows OOMKilled=true:
  Investigate memory limits and host memory pressure.

OS logs show killed process or out-of-memory:
  Investigate host/VM memory pressure.

Process still exists but no longer responds:
  Capture listeners, packet flow, logs, and resource state before restarting.

Control panel shows running but no process exists:
  Investigate control-panel state mismatch and startup logs.

Crash happens after a specific player action:
  Capture one controlled reproduction and correlate logs by UTC time.
```

## 7. Evidence to Escalate

```text
Failure type
Exact UTC crash/hang time
Process state before/after
Runtime logs
System logs
Container/service exit state, if applicable
Resource state
Recent changes
Controlled reproduction notes, if available
```
