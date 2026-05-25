# Runtime Guide: Manual or Custom Script

Use this only when the server appears to be started by a shell script, batch file, scheduled task, manual command, or unknown custom launcher.

## 1. Identify the Running Process

Run on: Linux host or Linux VM shell

```bash
ps -ef | grep -Ei 'DuneSandbox|Awakening|Sandbox' | grep -v grep
```

Run on: Windows host PowerShell

```powershell
Get-Process | Where-Object { $_.ProcessName -match 'Dune|Awakening|Sandbox' } | Select-Object ProcessName, Id, Path
```

Record:

```text
SERVICE_NAME=
Process ID:
Executable path:
Launch user:
```

## 2. Locate Startup Scripts

Run on: Linux host or Linux VM shell

```bash
find /home /opt /srv -maxdepth 8 -type f \( -name '*.sh' -o -name '*.env' -o -name '*.ini' \) 2>/dev/null | grep -Ei 'dune|awakening|server|start|run' | head -100
```

Run on: Windows host PowerShell

```powershell
Get-ChildItem -Path C:\ -Recurse -ErrorAction SilentlyContinue -Include *.bat,*.cmd,*.ps1,*.ini,*.env | Where-Object { $_.FullName -match 'Dune|Awakening|server|start|run' } | Select-Object -First 100 FullName
```

## 3. Capture the Launch Command Safely

Run on: Linux host or Linux VM shell

```bash
ps -ef | grep DuneSandbox | grep -v grep | sed -E 's/(Token|Secret|Password|Key)=[^ ]+/\1=<redacted>/Ig'
```

Run on: Windows host PowerShell

```powershell
Get-CimInstance Win32_Process | Where-Object { $_.CommandLine -match 'Dune|Awakening|Sandbox' } | Select-Object ProcessId, CommandLine
```

Redact tokens, secrets, and passwords before sharing externally.

## 4. Continue to Focused Runbooks

- [Failed travel capture](../runbooks/failed-travel-capture.md)
- [Port and listener validation](../runbooks/port-and-network-listener-validation.md)
- [Log collection and redaction](../runbooks/log-collection-and-redaction.md)
