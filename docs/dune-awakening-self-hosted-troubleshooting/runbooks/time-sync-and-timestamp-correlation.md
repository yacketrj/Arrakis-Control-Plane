# Runbook: Time Sync and Timestamp Correlation

Use this runbook when log timestamps do not line up, events appear out of order, client reports do not match server logs, or support needs to correlate events across cloud, host, VM, container, and game logs.

## 1. Record the User's Local Time and UTC Time

Ask the user or environment owner:

```text
What time did the issue happen in your local timezone?
What timezone are you in?
Can you reproduce the issue now while we record UTC time?
```

Record:

```text
User local time:
User timezone:
Converted UTC time:
Server UTC time:
```

## 2. Check Time on Linux

Run on: Linux host or Linux VM shell

```bash
date
date -u
timedatectl
```

Record:

```text
Local time:
UTC time:
Timezone:
NTP synchronized: yes/no/unknown
```

## 3. Check Time on Windows

Run on: Windows host PowerShell

```powershell
Get-Date
Get-Date -Format u
w32tm /query /status
```

Record:

```text
Local time:
UTC time:
Time source:
Windows time sync status:
```

## 4. Check Time in Docker Containers, Only If Docker Is Confirmed

Run on: Docker host shell

```bash
docker ps --format '{{.Names}}' | while read c; do
  echo "===== $c ====="
  docker exec "$c" date -u 2>/dev/null || true
done
```

Interpretation:

```text
Container UTC time differs significantly from host UTC time:
  Capture both values and investigate host/container time configuration.

Container UTC time matches host UTC time:
  Use UTC timestamps to correlate container and host logs.
```

## 5. Correlate Logs by UTC Window

Use a narrow window around the reproduction attempt.

```text
Test start UTC:
Test failure UTC:
Test end UTC:
Log sources included:
```

Run on: Linux host or Linux VM shell

```bash
grep -RniE 'error|warning|failed|exception|timeout|disconnect|LoginRequest|Travel|ServerState' "$LOG_PATH" 2>/dev/null | tee timestamp-correlation-search.txt
```

Run on: Windows host PowerShell

```powershell
Select-String -Path "$env:LOG_PATH\**\*.log" -Pattern 'error','warning','failed','exception','timeout','disconnect','LoginRequest','Travel','ServerState' -ErrorAction SilentlyContinue | Tee-Object -FilePath timestamp-correlation-search.txt
```

## 6. Interpret Results

```text
Server time is wrong:
  Fix time sync before relying on timestamp ordering.

Client time and server time are in different timezones:
  Convert both to UTC before drawing conclusions.

Cloud, hypervisor, guest, and container times differ:
  Record all values and use the most authoritative UTC source after time sync is corrected.

Events appear out of order only in one component:
  Check that component's log buffering, queue delay, restart time, and local clock.
```

## 7. Evidence to Escalate

```text
User local time and timezone
Server UTC time
Host time sync status
Container time, if applicable
Cloud/hypervisor time context, if applicable
UTC reproduction window
Log snippets from each component in the same UTC window
```
