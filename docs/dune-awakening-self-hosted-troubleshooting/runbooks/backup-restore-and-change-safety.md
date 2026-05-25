# Runbook: Backup, Restore, and Change Safety

Use this runbook before making changes that could affect saved game data, configuration, database state, or service startup behavior.

Goal: prevent avoidable data loss while troubleshooting.

## 1. Decide Whether a Backup Is Required

Create or verify a backup before:

```text
Editing configuration files
Changing ownership or permissions
Changing Docker bind mounts
Changing AMP instance files
Changing systemd or Windows service definitions
Changing database files or database containers
Restarting repeatedly after crashes
Applying updates or image changes
Deleting queues, containers, volumes, or old instances
```

## 2. Identify What Must Be Backed Up

Record actual values:

```text
INSTANCE_PATH=
SAVED_PATH=
LOG_PATH=
Database path or container/service:
Configuration path:
Control panel export path, if available:
```

## 3. Linux File Backup

Run on: Linux host or Linux VM shell

```bash
BACKUP_DIR="$HOME/dune-backup-$(date -u +%Y%m%d-%H%M%S)"
mkdir -p "$BACKUP_DIR"

[ -n "$INSTANCE_PATH" ] && sudo tar -czf "$BACKUP_DIR/instance-path.tar.gz" "$INSTANCE_PATH" 2>/dev/null || true
[ -n "$SAVED_PATH" ] && sudo tar -czf "$BACKUP_DIR/saved-path.tar.gz" "$SAVED_PATH" 2>/dev/null || true
[ -n "$LOG_PATH" ] && sudo tar -czf "$BACKUP_DIR/log-path.tar.gz" "$LOG_PATH" 2>/dev/null || true

ls -lh "$BACKUP_DIR"
```

## 4. Windows File Backup

Run on: Windows host PowerShell

```powershell
$BackupDir = "$env:USERPROFILE\dune-backup-$(Get-Date -Format yyyyMMdd-HHmmss)"
New-Item -ItemType Directory -Path $BackupDir | Out-Null

if ($env:INSTANCE_PATH) { Compress-Archive -Path $env:INSTANCE_PATH -DestinationPath "$BackupDir\instance-path.zip" -Force }
if ($env:SAVED_PATH) { Compress-Archive -Path $env:SAVED_PATH -DestinationPath "$BackupDir\saved-path.zip" -Force }
if ($env:LOG_PATH) { Compress-Archive -Path $env:LOG_PATH -DestinationPath "$BackupDir\log-path.zip" -Force }

Get-ChildItem $BackupDir
```

## 5. Docker Volume and Mount Safety

Run on: Docker host shell, only if Docker is confirmed

```bash
docker ps --format 'table {{.Names}}\t{{.Image}}\t{{.Status}}'
docker volume ls
docker ps --format '{{.Names}}' | while read c; do
  echo "===== $c ====="
  docker inspect "$c" --format '{{range .Mounts}}{{println .Type .Name .Source "->" .Destination}}{{end}}'
done | tee docker-mounts-before-change.txt
```

Before deleting or recreating a container, confirm whether important data lives in:

```text
Bind mounts on the host
Docker named volumes
The container writable layer
A database container or service
A control panel instance directory
```

## 6. Database Backup Reminder

If a database is discovered, use the database-specific backup method. Do not assume file copy alone is safe for a running database.

Record:

```text
Database engine:
Database service/container:
Backup method used:
Backup file path:
UTC backup time:
```

## 7. Change Record Template

Before making a change, record:

```text
UTC time:
Person making change:
Reason for change:
Files/services affected:
Backup created: yes/no
Rollback method:
Expected result:
Actual result:
```

## 8. Rollback Rule

If a change makes the issue worse or creates a new failure, stop and roll back before making another change.

Do not stack multiple untracked changes. That makes RCA unreliable and increases data-loss risk.
