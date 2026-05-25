# Runtime Guide: AMP Control Panel

Use this when the Dune: Awakening server is managed through AMP or another AMP-like control panel.

## 1. Confirm AMP Is the Management Layer

Run in: AMP web UI

```text
Open the Dune: Awakening instance.
Record instance name, status, module type, install path, configuration path, log path, startup command if visible, and whether Docker/container mode is enabled.
```

Run on: Linux host or Linux VM shell, if shell access is available

```bash
ps -ef | grep -Ei 'amp|cubecoders|DuneSandbox|docker' | grep -v grep
find /home -maxdepth 5 -type d -iname '*amp*' 2>/dev/null | head -100
```

Run on: Windows host PowerShell, if shell access is available

```powershell
Get-Process | Where-Object { $_.ProcessName -match 'amp|docker|Dune|Awakening|Sandbox' } | Select-Object ProcessName, Id, Path
Get-Service | Where-Object { $_.Name -match 'amp|docker|dune' -or $_.DisplayName -match 'amp|docker|dune' } | Select-Object Name, DisplayName, Status
```

## 2. Record AMP Operational Values

```text
SERVICE_NAME=
INSTANCE_PATH=
LOG_PATH=
SAVED_PATH=
Management mode shown by AMP:
Docker/container mode shown by AMP:
```

## 3. Use the Correct Next Runtime Layer

If AMP launches Docker containers, continue with [Docker or Docker Compose](./docker-or-compose.md).

If AMP launches a host process without Docker, continue with the matching platform guide and the focused runbooks.

## 4. Evidence to Capture From AMP

Run in: AMP web UI

```text
Export or copy console output covering the failure window.
Capture screenshots of instance status, recent restart history, visible ports, and configured startup options.
Record any AMP job/task failures or permission errors.
```
