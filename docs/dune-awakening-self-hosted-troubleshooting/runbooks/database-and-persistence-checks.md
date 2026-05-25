# Runbook: Database and Persistence Checks

Use this runbook when logs or symptoms suggest database, persistence, character loading, saved state, or world-state issues.

Examples:

```text
Character download fails
DatabaseLogin fails
World state does not save
Server starts but cannot load saved data
Player progress rolls back
Database connection errors appear
```

Start only after the hosting platform and runtime/orchestration layer are identified.

## 1. Capture the Persistence Symptom

Ask the user or environment owner:

```text
What data failed to load or save?
Does the issue affect one player, many players, or the whole server?
Did the issue begin after a restart, update, migration, restore, crash, or disk event?
Does the server still start?
Can new characters join?
Can existing characters join?
Exact UTC time of the failure:
```

Record:

```text
Reported persistence symptom:
Affected data type:
Affected player or world scope:
Recent restart/update/restore:
Exact UTC time:
```

## 2. Search Logs for Persistence and Database Errors

Run on: Linux host or Linux VM shell

```bash
grep -RniE 'DatabaseLogin|Database|Postgres|postgres|SQL|persistence|Persistent|CharacterDownload|LoadPlayer|Save|rollback|transaction|connection refused|timeout|failed|error|exception' "$LOG_PATH" 2>/dev/null | head -500 | tee database-persistence-search.txt
```

Run on: Windows host PowerShell

```powershell
Select-String -Path "$env:LOG_PATH\**\*.log" -Pattern 'DatabaseLogin','Database','Postgres','postgres','SQL','persistence','Persistent','CharacterDownload','LoadPlayer','Save','rollback','transaction','connection refused','timeout','failed','error','exception' -ErrorAction SilentlyContinue | Select-Object -First 500 | Tee-Object -FilePath database-persistence-search.txt
```

Run in: control panel UI

```text
Open server, database, and startup logs.
Export logs covering the failure window.
```

## 3. Identify the Database Service

Run on: Linux host or Linux VM shell

```bash
ps -ef | grep -Ei 'postgres|mysql|mariadb|sqlite|database' | grep -v grep
systemctl list-units --type=service 2>/dev/null | grep -Ei 'postgres|mysql|mariadb|database' || true
```

Run on: Windows host PowerShell

```powershell
Get-Process | Where-Object { $_.ProcessName -match 'postgres|mysql|mariadb|database' } | Select-Object ProcessName, Id, Path
Get-Service | Where-Object { $_.Name -match 'postgres|mysql|mariadb|database' -or $_.DisplayName -match 'postgres|mysql|mariadb|database' } | Select-Object Name, DisplayName, Status
```

Run on: Docker host shell, only if Docker is confirmed

```bash
docker ps --format 'table {{.Names}}\t{{.Image}}\t{{.Status}}\t{{.Ports}}' | grep -Ei 'postgres|mysql|mariadb|database' || true
```

Record:

```text
Database service/container name:
Database host:
Database port:
Database name, if known:
Database user, if known:
```

Do not record passwords in notes unless absolutely required and approved.

## 4. Check Database Service Health

Run on: Linux systemd host, if database is a systemd service

```bash
systemctl status "$DATABASE_SERVICE" --no-pager
journalctl -u "$DATABASE_SERVICE" --since "30 minutes ago" --no-pager | tee database-service.log
```

Run on: Docker host shell, if database is a container

```bash
docker logs --since 30m "$DATABASE_SERVICE" 2>&1 | tee database-container.log
docker inspect "$DATABASE_SERVICE" --format 'Status={{.State.Status}} ExitCode={{.State.ExitCode}} Error={{.State.Error}}'
```

Run on: Windows host PowerShell, if database is a Windows service

```powershell
Get-Service -Name $env:DATABASE_SERVICE
Get-WinEvent -LogName Application -MaxEvents 300 | Where-Object { $_.Message -match 'postgres|mysql|mariadb|database|sql' } | Format-List TimeCreated, ProviderName, Id, LevelDisplayName, Message
```

## 5. Check Disk Space

Run on: Linux host or Linux VM shell

```bash
df -h
findmnt
```

Run on: Windows host PowerShell

```powershell
Get-Volume | Select-Object DriveLetter, FileSystemLabel, SizeRemaining, Size
Get-PSDrive -PSProvider FileSystem
```

Interpretation:

```text
Disk full or nearly full:
  Database writes, save files, logs, and backups may fail.

Database service unhealthy:
  Fix service health before deeper gameplay troubleshooting.

Game logs show database login or character download failure:
  Correlate with database logs at the same UTC time.
```

## 6. Escalation Evidence

```text
Exact persistence symptom
Database-related game logs
Database service/container status
Database logs
Disk space output
Recent backup/restore/update history
Exact UTC failure time
```
