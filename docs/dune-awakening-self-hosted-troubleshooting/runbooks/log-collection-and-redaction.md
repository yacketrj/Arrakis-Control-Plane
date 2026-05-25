# Runbook: Log Collection and Redaction

Use this runbook after the hosting platform and runtime/orchestration layer are known.

Goal: collect useful logs while protecting sensitive values before broad sharing.

## 1. Identify Log Sources

Run in: control panel UI

```text
Open the instance console or log viewer.
Record the log file name, time range, and export method.
```

Run on: Linux host or Linux VM shell

```bash
find "$INSTANCE_PATH" -type f -name '*.log' 2>/dev/null | sort | tee discovered-log-files.txt
```

Run on: Windows host PowerShell

```powershell
Get-ChildItem -Path $env:INSTANCE_PATH -Recurse -Filter *.log -ErrorAction SilentlyContinue | Select-Object FullName | Tee-Object -FilePath discovered-log-files.txt
```

Run on: Docker host shell, only if Docker is confirmed

```bash
docker ps --format '{{.Names}}' | while read c; do
  echo "===== $c ====="
  docker logs --since 30m "$c" 2>&1 | tail -200
 done | tee docker-recent-logs.txt
```

## 2. Collect Logs for a Known Time Window

Run on: Linux host or Linux VM shell

```bash
mkdir -p logs-collected
find "$LOG_PATH" -type f -name '*.log' -mmin -180 -print -exec cp --parents {} logs-collected/ \; 2>/dev/null
```

Run on: Windows host PowerShell

```powershell
New-Item -ItemType Directory -Path logs-collected -Force | Out-Null
Get-ChildItem -Path $env:LOG_PATH -Recurse -Filter *.log -ErrorAction SilentlyContinue |
  Where-Object { $_.LastWriteTime -gt (Get-Date).AddHours(-3) } |
  Copy-Item -Destination logs-collected -Force
```

## 3. Redact Sensitive Values

Use redaction appropriate for the sharing audience. Do not blindly redact operational values needed for troubleshooting, such as service names, container names, map names, or local paths when sharing internally.

Run on: Linux host or Linux VM shell

```bash
for f in $(find logs-collected -type f); do
  sed -i \
    -e 's/ServiceAuthToken=[^ ]*/ServiceAuthToken=<redacted>/g' \
    -e 's/eyJ[A-Za-z0-9._-]*/<redacted-jwt>/g' \
    -e 's/ocid1\.[A-Za-z0-9._-]*/<redacted-ocid>/g' \
    -e 's/Password=[^ ]*/Password=<redacted>/Ig' \
    -e 's/Secret=[^ ]*/Secret=<redacted>/Ig' \
    -e 's/Token=[^ ]*/Token=<redacted>/Ig' \
    "$f"
done
```

Run on: Windows host PowerShell

```powershell
Get-ChildItem -Path logs-collected -Recurse -File | ForEach-Object {
  $content = Get-Content $_.FullName -Raw
  $content = $content -replace 'ServiceAuthToken=[^ ]+', 'ServiceAuthToken=<redacted>'
  $content = $content -replace 'eyJ[A-Za-z0-9._-]+', '<redacted-jwt>'
  $content = $content -replace 'ocid1\.[A-Za-z0-9._-]+', '<redacted-ocid>'
  $content = $content -replace '(?i)Password=[^ ]+', 'Password=<redacted>'
  $content = $content -replace '(?i)Secret=[^ ]+', 'Secret=<redacted>'
  $content = $content -replace '(?i)Token=[^ ]+', 'Token=<redacted>'
  Set-Content -Path $_.FullName -Value $content
}
```

## 4. Package Logs

Run on: Linux host or Linux VM shell

```bash
tar -czf dune-logs-$(date -u +%Y%m%d-%H%M%S).tar.gz logs-collected/
ls -lh dune-logs-*.tar.gz
```

Run on: Windows host PowerShell

```powershell
Compress-Archive -Path logs-collected -DestinationPath "dune-logs-$(Get-Date -Format yyyyMMdd-HHmmss).zip"
Get-ChildItem dune-logs-*.zip
```
