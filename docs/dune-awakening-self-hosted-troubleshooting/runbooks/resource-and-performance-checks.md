# Runbook: Resource and Performance Checks

Use this runbook when the server is slow, players rubber-band, travel hangs under load, startup takes too long, processes are killed, or the environment owner reports high CPU, RAM, disk, or network usage.

Start only after the hosting platform and runtime/orchestration layer are identified.

## 1. Capture the Performance Symptom

Ask the user or environment owner:

```text
What feels slow or broken?
Does it affect all players or only some players?
Does it happen all the time or only during peak activity?
Does it happen during login, travel, combat, building, harvesting, storms, or shutdown/startup?
When did it start?
Was there a recent player-count increase, map change, update, backup, restore, or host change?
Exact UTC time of the issue:
```

Record:

```text
Reported performance issue:
Affected workflow:
Affected users:
Known working time period:
Known failing time period:
Recent change:
Exact UTC time:
```

## 2. Check CPU, Memory, and Load

Run on: Linux host or Linux VM shell

```bash
uptime
free -h
ps -eo pid,ppid,user,%cpu,%mem,etime,cmd --sort=-%cpu | head -30
ps -eo pid,ppid,user,%cpu,%mem,etime,cmd --sort=-%mem | head -30
```

Run on: Windows host PowerShell

```powershell
Get-Counter '\Processor(_Total)\% Processor Time','\Memory\Available MBytes'
Get-Process | Sort-Object CPU -Descending | Select-Object -First 30 ProcessName, Id, CPU, WorkingSet, Path
Get-Process | Sort-Object WorkingSet -Descending | Select-Object -First 30 ProcessName, Id, CPU, WorkingSet, Path
```

Run on: Docker host shell, only if Docker is confirmed

```bash
docker stats --no-stream
```

## 3. Check Disk Space and I/O Risk

Run on: Linux host or Linux VM shell

```bash
df -h
findmnt
lsblk
```

If available:

```bash
iostat -xz 1 5 2>/dev/null || true
```

Run on: Windows host PowerShell

```powershell
Get-Volume | Select-Object DriveLetter, FileSystemLabel, SizeRemaining, Size
Get-PSDrive -PSProvider FileSystem
Get-Counter '\PhysicalDisk(_Total)\% Disk Time','\PhysicalDisk(_Total)\Avg. Disk Queue Length'
```

Interpretation:

```text
Disk full or nearly full:
  Logs, database writes, save files, and backups may fail.

High disk queue or high disk time:
  Startup, save, database, and travel operations may stall.
```

## 4. Check Network Interface Health

Run on: Linux host or Linux VM shell

```bash
ip -s link
ss -s
```

Run on: Windows host PowerShell

```powershell
Get-NetAdapterStatistics
Get-NetAdapter | Select-Object Name, Status, LinkSpeed, MacAddress
```

## 5. Check Process Restarts or Kills

Run on: Linux host or Linux VM shell

```bash
dmesg -T | grep -Ei 'killed process|out of memory|oom|segfault|blocked for more than|I/O error' | tail -100
journalctl --since "2 hours ago" --no-pager | grep -Ei 'killed process|out of memory|oom|segfault|Dune|docker|rabbit|postgres|failed|error' | tail -300
```

Run on: Windows host PowerShell

```powershell
Get-WinEvent -LogName System -MaxEvents 300 | Where-Object { $_.Message -match 'resource|memory|disk|network|service|terminated|failed|error' } | Format-List TimeCreated, ProviderName, Id, LevelDisplayName, Message
Get-WinEvent -LogName Application -MaxEvents 300 | Where-Object { $_.Message -match 'Dune|Awakening|Sandbox|Docker|AMP|Rabbit|Postgres|failed|error' } | Format-List TimeCreated, ProviderName, Id, LevelDisplayName, Message
```

## 6. Correlate With Gameplay Logs

Run on: Linux host or Linux VM shell

```bash
grep -RniE 'timeout|timed out|lag|stall|blocked|slow|failed|error|disconnect|ServerState|Travel|Login' "$LOG_PATH" 2>/dev/null | head -500 | tee performance-log-search.txt
```

Run on: Windows host PowerShell

```powershell
Select-String -Path "$env:LOG_PATH\**\*.log" -Pattern 'timeout','timed out','lag','stall','blocked','slow','failed','error','disconnect','ServerState','Travel','Login' -ErrorAction SilentlyContinue | Select-Object -First 500 | Tee-Object -FilePath performance-log-search.txt
```

## 7. Interpret Results

```text
CPU saturated:
  Identify the top process and correlate with player activity or map startup.

Memory low or OOM events present:
  Reduce workload, check container/service memory limits, and capture logs before restarting.

Disk full or I/O saturated:
  Free space or move heavy backups/logs only after confirming what is safe to remove.

Network errors or drops increasing:
  Check host NIC, hypervisor switch/bridge, cloud network, and upstream provider path.

No resource pressure found:
  Continue with the specific symptom runbook, such as login, travel, database, or messaging.
```

## 8. Evidence to Escalate

```text
Exact UTC issue window
CPU/memory output
Disk space output
Network interface output
Service/container stats, if applicable
Gameplay logs from same time window
Host/system event logs from same time window
```
