# Runbook: Performance, Resource, and Capacity Checks

Use this runbook when players report lag, rubber-banding, slow travel, delayed login, high memory use, CPU pressure, disk pressure, or server instability under load.

Start only after the hosting platform and runtime/orchestration layer are identified.

## 1. Capture the Performance Symptom

Ask the user or environment owner:

```text
What feels slow or unstable?
When did it begin?
Does it affect all players or only some players?
How many players were online?
Was the server recently restarted, updated, backed up, or moved?
Does it happen during login, travel, combat, building, harvesting, storms, or general movement?
Exact UTC time of the symptom:
```

Record:

```text
Reported symptom:
Affected action:
Player count:
UTC time window:
Recent change:
```

## 2. Check CPU, Memory, Disk, and Load

Run on: Linux host or Linux VM shell

```bash
uptime
free -h
df -h
top -b -n 1 | head -60
ps -eo pid,ppid,user,%cpu,%mem,cmd --sort=-%cpu | head -30
ps -eo pid,ppid,user,%cpu,%mem,cmd --sort=-%mem | head -30
```

Run on: Windows host PowerShell

```powershell
Get-CimInstance Win32_OperatingSystem | Select-Object TotalVisibleMemorySize, FreePhysicalMemory
Get-Volume | Select-Object DriveLetter, FileSystemLabel, SizeRemaining, Size
Get-Process | Sort-Object CPU -Descending | Select-Object -First 30 ProcessName, Id, CPU, WorkingSet
Get-Counter '\Processor(_Total)\% Processor Time','\Memory\Available MBytes','\LogicalDisk(_Total)\% Free Space'
```

Run on: Docker host shell, only if Docker is confirmed

```bash
docker stats --no-stream
docker ps --format 'table {{.Names}}\t{{.Status}}\t{{.Ports}}'
```

## 3. Check Host or VM Resource Allocation

Run on: Hyper-V host PowerShell, only if Hyper-V is used

```powershell
Get-VM | Select-Object Name, State, CPUUsage, MemoryAssigned, MemoryDemand, Uptime
```

Run on: Proxmox host shell, only if Proxmox is used

```bash
qm list
pct list
pvesh get /nodes/$(hostname)/status
```

Run in: cloud provider console, only if cloud-hosted

```text
Check VM size/shape, CPU allocation, memory allocation, disk type, disk throughput, network bandwidth, and recent maintenance events.
```

## 4. Check Logs for Resource Symptoms

Run on: Linux host or Linux VM shell

```bash
grep -RniE 'out of memory|oom|killed process|timeout|timed out|slow|stalled|hitch|lag|disk|space|failed to allocate|GC|garbage|crash|fatal' "$LOG_PATH" 2>/dev/null | head -400 | tee resource-log-search.txt
```

Run on: Windows host PowerShell

```powershell
Select-String -Path "$env:LOG_PATH\**\*.log" -Pattern 'out of memory','oom','timeout','timed out','slow','stalled','hitch','lag','disk','space','failed to allocate','GC','garbage','crash','fatal' -ErrorAction SilentlyContinue | Select-Object -First 400 | Tee-Object -FilePath resource-log-search.txt
```

## 5. Interpret Results

```text
CPU saturated:
  Capture top processes, player count, and exact time window before changing allocation.

Memory nearly full with swap pressure or OOM logs:
  Check game process memory, container limits, VM memory allocation, and recent restarts.

Disk nearly full:
  Logs, database writes, saves, and backups may fail.

Docker memory limit reached:
  Check container memory settings and host capacity.

Cloud VM shape too small or burst-limited:
  Capture provider metrics and compare against player count and workload.

No resource pressure visible during failure:
  Continue with login, travel, network, or persistence runbooks.
```

## 6. Evidence to Escalate

```text
Exact symptom and UTC time window
Player count
CPU/memory/disk output
Docker stats, if applicable
Hypervisor/cloud resource allocation, if applicable
Relevant logs
Recent restart/update/backup events
```
