# Runbook: Crash, Hang, and Resource Pressure

Use this runbook when the server process crashes, hangs, becomes unresponsive, consumes high CPU or memory, or the host appears resource constrained.

Start only after the hosting platform and runtime/orchestration layer are identified.

## 1. Capture the Symptom

Ask the user or environment owner:

```text
Did the server crash, hang, restart, freeze, or slow down?
Did it recover by itself?
Did a player action trigger it?
Does the issue affect one map/server or the whole deployment?
Exact UTC time of the crash or hang:
Was the server manually restarted?
Was there a recent update, config change, backup, or restore?
```

Record:

```text
Observed symptom:
Affected service/process:
Exact UTC time:
Manual restart performed: yes/no
Recent change:
```

## 2. Capture Host Resource State

Run on: Linux host or Linux VM shell

```bash
uptime
free -h
df -h
ps -eo pid,ppid,user,%cpu,%mem,etime,cmd --sort=-%cpu | head -30
ps -eo pid,ppid,user,%cpu,%mem,etime,cmd --sort=-%mem | head -30
dmesg -T | grep -Ei 'killed process|oom|out of memory|segfault|blocked for more than|I/O error' | tail -100
```

Run on: Windows host PowerShell

```powershell
Get-Date -Format u
Get-ComputerInfo | Select-Object CsName, OsTotalVisibleMemorySize, OsFreePhysicalMemory
Get-Process | Sort-Object CPU -Descending | Select-Object -First 30 ProcessName, Id, CPU, WorkingSet, Path
Get-Volume | Select-Object DriveLetter, FileSystemLabel, SizeRemaining, Size
Get-WinEvent -LogName System -MaxEvents 300 | Where-Object { $_.Message -match 'memory|resource|disk|crash|terminated|failed|timeout' } | Format-List TimeCreated, ProviderName, Id, LevelDisplayName, Message
```

## 3. Capture Service, Container, or Control Panel State

Run on: Linux systemd host, only if systemd is used

```bash
systemctl status "$SERVICE_NAME" --no-pager
journalctl -u "$SERVICE_NAME" --since "2 hours ago" --no-pager | tee crash-hang-service.log
```

Run on: Windows host PowerShell, only if Windows service is used

```powershell
Get-Service -Name $env:SERVICE_NAME
Get-WinEvent -LogName Application -MaxEvents 300 | Where-Object { $_.Message -match 'Dune|Awakening|Sandbox|crash|fault|terminated|exception' } | Format-List TimeCreated, ProviderName, Id, LevelDisplayName, Message
```

Run on: Docker host shell, only if Docker is confirmed

```bash
docker ps -a --format 'table {{.Names}}\t{{.Status}}\t{{.Image}}'
docker stats --no-stream
docker logs --since 2h "$CONTAINER_NAME" 2>&1 | tee crash-hang-container.log
docker inspect "$CONTAINER_NAME" --format 'Status={{.State.Status}} ExitCode={{.State.ExitCode}} OOMKilled={{.State.OOMKilled}} Error={{.State.Error}}'
```

Run in: control panel UI

```text
Open the instance status, console, task/job log, and restart history.
Record whether the panel shows crash, restart, stopped, killed, memory limit, or timeout messages.
```

## 4. Search Game Logs

Run on: Linux host or Linux VM shell

```bash
grep -RniE 'crash|fatal|assert|ensure|segfault|out of memory|OOM|timeout|hang|stalled|blocked|exception|error' "$LOG_PATH" 2>/dev/null | head -500 | tee crash-hang-log-search.txt
```

Run on: Windows host PowerShell

```powershell
Select-String -Path "$env:LOG_PATH\**\*.log" -Pattern 'crash','fatal','assert','ensure','segfault','out of memory','OOM','timeout','hang','stalled','blocked','exception','error' -ErrorAction SilentlyContinue | Select-Object -First 500 | Tee-Object -FilePath crash-hang-log-search.txt
```

## 5. Interpret Results

```text
OOMKilled is true or Linux dmesg shows killed process:
  Memory pressure or container memory limit likely contributed. Capture limits and workload before changing memory.

CPU is saturated and server hangs:
  Capture process list and affected map/server logs before restarting.

Disk is full:
  Logs, saves, database writes, and backups may fail. Free space only after backup review.

Control panel shows repeated restarts:
  Capture restart history and startup logs. Avoid restart loops.

Game logs show fatal/assert/exception:
  Package logs and escalate with exact UTC time and reproduction steps.
```

## 6. Evidence to Escalate

```text
User-defined crash/hang symptom
Exact UTC failure time
Host resource output
Service/container/control-panel state
Game logs around failure
System or event logs around failure
Recent changes
Whether restart recovered service
```
