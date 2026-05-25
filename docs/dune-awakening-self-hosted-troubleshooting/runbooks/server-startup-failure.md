# Runbook: Server Startup Failure

Use this runbook when a Dune: Awakening server, map server, control-plane service, or supporting service fails to start or immediately exits.

Start only after the hosting platform and runtime/orchestration layer are identified.

## 1. Capture the Startup Symptom

Record:

```text
Service or process being started:
Who started it:
Start method used:
Expected behavior:
Actual behavior:
Exact error message:
UTC start time:
Recent configuration or file change:
```

Do not restart repeatedly until the first failure output is captured.

---

## 2. Check Service or Process State

Run on: Linux host or Linux VM shell

```bash
systemctl list-units --type=service --all | grep -Ei 'dune|awakening|sandbox|amp|docker|rabbit|compose' || true
ps -ef | grep -Ei 'DuneSandbox|Awakening|Sandbox|AMP|Rabbit|docker|compose' | grep -v grep || true
```

Run on: Windows host PowerShell

```powershell
Get-Service | Where-Object { $_.Name -match 'dune|awakening|sandbox|amp|docker|rabbit' -or $_.DisplayName -match 'dune|awakening|sandbox|amp|docker|rabbit' } | Select-Object Name, DisplayName, Status
Get-Process | Where-Object { $_.ProcessName -match 'Dune|Awakening|Sandbox|AMP|Rabbit|Docker' } | Select-Object ProcessName, Id, Path
```

Run in: control panel UI

```text
Open the affected instance.
Record status, latest console output, last start time, last stop time, and any visible task/job error.
```

Run on: Docker host shell, only if Docker is confirmed

```bash
docker ps -a --format 'table {{.Names}}\t{{.Image}}\t{{.Status}}\t{{.Ports}}'
```

---

## 3. Capture Startup Logs

Run on: Linux systemd host, only if systemd is used

```bash
journalctl -u "$SERVICE_NAME" --since "30 minutes ago" --no-pager | tee startup-journal.log
systemctl status "$SERVICE_NAME" --no-pager | tee startup-status.txt
```

Run on: Windows host PowerShell, only if Windows service or host process is used

```powershell
Get-WinEvent -LogName Application -MaxEvents 300 | Where-Object { $_.Message -match 'Dune|Awakening|Sandbox|AMP|Rabbit|Docker' } | Format-List TimeCreated, ProviderName, Id, LevelDisplayName, Message | Tee-Object -FilePath startup-eventlog.txt
```

Run on: Docker host shell, only if Docker is confirmed

```bash
docker logs --since 30m "$CONTAINER_NAME" 2>&1 | tee startup-container.log
```

Run in: control panel UI

```text
Export or copy console and task output covering the startup attempt.
```

---

## 4. Check Configuration and File Access

Run on: Linux host or Linux VM shell

```bash
ls -ld "$INSTANCE_PATH" "$SAVED_PATH" "$LOG_PATH" 2>/dev/null || true
find "$INSTANCE_PATH" -maxdepth 4 -type f \( -name '*.ini' -o -name '*.env' -o -name '*.json' -o -name '*.log' \) 2>/dev/null | head -100
```

Run on: Windows host PowerShell

```powershell
Get-Item $env:INSTANCE_PATH, $env:SAVED_PATH, $env:LOG_PATH -ErrorAction SilentlyContinue | Format-List FullName, Attributes, LastWriteTime
Get-ChildItem -Path $env:INSTANCE_PATH -Recurse -ErrorAction SilentlyContinue -Include *.ini,*.env,*.json,*.log | Select-Object -First 100 FullName
```

If the error contains permission language, continue to [Permission and ownership errors](./permission-and-ownership-errors.md).

---

## 5. Check Port Conflicts

Run on: Linux host or Linux VM shell

```bash
sudo ss -tulpen | tee startup-listeners.txt
```

Run on: Windows host PowerShell

```powershell
Get-NetUDPEndpoint | Sort-Object LocalPort | Tee-Object -FilePath startup-udp-listeners.txt
Get-NetTCPConnection | Sort-Object LocalPort | Tee-Object -FilePath startup-tcp-listeners.txt
```

Interpretation:

```text
Expected port already in use:
  Identify the owning process before restarting anything.

Expected port not listening after startup:
  Check service logs and launch arguments.
```

---

## 6. Check Runtime-Specific Startup Data

If AMP is used:

```text
Open AMP console/task output.
Record the exact command or module action that failed.
Capture visible file permission or container errors.
```

If Docker is used:

```bash
docker inspect "$CONTAINER_NAME" --format 'Image={{.Config.Image}} Status={{.State.Status}} ExitCode={{.State.ExitCode}} Error={{.State.Error}}'
docker inspect "$CONTAINER_NAME" --format '{{range .Mounts}}{{println .Source "->" .Destination}}{{end}}'
```

If systemd is used:

```bash
systemctl cat "$SERVICE_NAME"
systemctl show "$SERVICE_NAME" -p ExecStart -p User -p Group -p WorkingDirectory -p Environment --no-pager
```

If Windows service is used:

```powershell
Get-CimInstance Win32_Service | Where-Object { $_.Name -eq $env:SERVICE_NAME } | Select-Object Name, State, StartName, PathName
```

---

## 7. Interpretation

```text
Immediate exit with permission error:
  Use the permission and ownership runbook.

Immediate exit with missing file/config error:
  Verify INSTANCE_PATH, SAVED_PATH, config file path, and control-panel path mapping.

Immediate exit with port bind error:
  Use the port and listener validation runbook.

Container exits but host service remains running:
  Troubleshoot the container runtime and mounts.

Control panel says running but no process or listener exists:
  Capture panel logs and verify the real launch method.

No logs are produced:
  Verify the service actually starts, the log path is correct, and the process has write access.
```
