# Runbook: Update, Patch, and Version Validation

Use this runbook when an issue begins after a game update, server image update, control panel update, OS patch, container image change, or rollback.

Goal: prove which version changed, when it changed, and whether all components are on compatible versions.

## 1. Capture the Change Story

Ask the user or environment owner:

```text
What was updated?
Who performed the update?
When was it updated in UTC?
Was the update automatic or manual?
Was there a rollback attempt?
Did the issue begin immediately after the update?
Did all services restart cleanly after the update?
```

Record:

```text
Reported change:
UTC change time:
Update method:
Rollback attempted: yes/no
Known good version, if known:
Current version, if known:
```

## 2. Capture Game Server Version or Build Evidence

Run on: Linux host or Linux VM shell

```bash
grep -RniE 'version|revision|build|changelist|Engine Version|LogInit|Project Version' "$LOG_PATH" 2>/dev/null | head -300 | tee version-log-search.txt
```

Run on: Windows host PowerShell

```powershell
Select-String -Path "$env:LOG_PATH\**\*.log" -Pattern 'version','revision','build','changelist','Engine Version','LogInit','Project Version' -ErrorAction SilentlyContinue | Select-Object -First 300 | Tee-Object -FilePath version-log-search.txt
```

Run in: control panel UI

```text
Open the instance details and console output.
Record any displayed game version, build ID, image tag, module version, update status, and last update time.
```

## 3. Capture Runtime Version Evidence

Run on: Docker host shell, only if Docker is confirmed

```bash
docker ps --format 'table {{.Names}}\t{{.Image}}\t{{.Status}}'
docker images | head -100
docker inspect "$CONTAINER_NAME" --format 'Image={{.Config.Image}} Created={{.Created}}'
```

Run on: Linux host or Linux VM shell

```bash
uname -a
cat /etc/os-release
apt list --upgradable 2>/dev/null | head -100 || true
```

Run on: Windows host PowerShell

```powershell
Get-ComputerInfo | Select-Object WindowsProductName, WindowsVersion, OsHardwareAbstractionLayer
Get-HotFix | Sort-Object InstalledOn -Descending | Select-Object -First 20 HotFixID, InstalledOn, Description
```

## 4. Capture Service Restart History

Run on: Linux systemd host

```bash
journalctl --since "24 hours ago" --no-pager | grep -Ei 'Started|Stopped|Restarted|Dune|Awakening|Sandbox|docker|amp|rabbit|postgres' | tail -300 | tee restart-history-linux.txt
```

Run on: Windows host PowerShell

```powershell
Get-WinEvent -LogName System -MaxEvents 500 | Where-Object { $_.Message -match 'service|started|stopped|restart|Dune|Awakening|Docker|AMP' } | Format-List TimeCreated, ProviderName, Id, LevelDisplayName, Message | Tee-Object -FilePath restart-history-windows.txt
```

Run in: control panel UI

```text
Export or screenshot the instance restart/update history if the panel provides it.
```

## 5. Interpret Results

```text
Only one component changed:
  Compare symptoms before and after that exact change.

Multiple components changed together:
  Do not assume which one caused the issue. Separate runtime, game, OS, control panel, and network changes.

Server and client versions appear mismatched:
  Escalate with version evidence and exact client/server behavior.

Container image changed but bind-mounted data did not:
  Check compatibility, migrations, permissions, and startup logs.

Rollback did not restore service:
  Capture rollback steps, current version evidence, and data/schema warnings.
```

## 6. Escalation Evidence

```text
User-reported update/change story
Version/build/revision logs
Container image tags or runtime versions
OS patch history
Control panel update history
Restart history
First failed UTC time after update
Known good version, if available
Current failing version
```
