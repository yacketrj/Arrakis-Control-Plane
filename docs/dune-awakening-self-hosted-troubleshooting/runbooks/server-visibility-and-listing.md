# Runbook: Server Visibility and Listing

Use this runbook when the server does not appear in the browser/listing, appears with the wrong name, appears in the wrong region, appears offline, or appears but cannot be selected.

Start only after the hosting platform and runtime/orchestration layer are identified.

## 1. Capture the User-Visible Symptom

Ask the user or environment owner:

```text
Can players see the server in the list?
Is the server name correct?
Is the region correct?
Does the server show offline, full, locked, passworded, or unavailable?
Can players select the server?
Does the issue affect all players or only some players?
Exact UTC time checked:
```

Record:

```text
Expected server name:
Actual server name shown:
Expected region:
Actual region shown:
Visible to players: yes/no/intermittent
Selectable by players: yes/no/intermittent
Client-side message:
```

## 2. Capture Control-Plane, Gateway, or Listing Logs

Run on: Docker host shell, only if the gateway/control-plane is a container

```bash
docker logs --since 30m "$DIRECTOR_SERVICE" 2>&1 | tee server-listing-control-plane.log
```

Run on: Linux host or Linux VM shell if logs are file-based

```bash
grep -RniE 'Gateway|Declare|FarmStatus|Battlegroup|ServerName|DisplayName|Region|FLS|listing|revision|heartbeat|online|offline|failed|error|exception' "$LOG_PATH" 2>/dev/null | head -400 | tee server-listing-search.txt
```

Run on: Windows host PowerShell if logs are file-based

```powershell
Select-String -Path "$env:LOG_PATH\**\*.log" -Pattern 'Gateway','Declare','FarmStatus','Battlegroup','ServerName','DisplayName','Region','FLS','listing','revision','heartbeat','online','offline','failed','error','exception' -ErrorAction SilentlyContinue | Select-Object -First 400 | Tee-Object -FilePath server-listing-search.txt
```

Run in: control panel UI

```text
Open console/log output for gateway, director, or listing-related services.
Export logs covering the time the server should have appeared.
```

## 3. Confirm Runtime Configuration Values

Run on: Linux host or Linux VM shell

```bash
ps -ef | grep -Ei 'DuneSandbox|gateway|director|server' | grep -v grep | sed -E 's/(Token|Secret|Password|Key)=[^ ]+/\1=<redacted>/Ig' | tee server-listing-processes.txt
```

Run on: Windows host PowerShell

```powershell
Get-CimInstance Win32_Process | Where-Object { $_.CommandLine -match 'Dune|Awakening|Sandbox|gateway|director|server' } | Select-Object ProcessId, CommandLine | Tee-Object -FilePath server-listing-processes.txt
```

Run on: Docker host shell, only if Docker is confirmed

```bash
docker ps --format '{{.Names}}' | while read c; do
  echo "===== $c ====="
  docker inspect "$c" --format '{{range .Config.Env}}{{println .}}{{end}}' 2>/dev/null | sed -E 's/(TOKEN|SECRET|PASSWORD|KEY)=.*/\1=<redacted>/Ig'
done | tee server-listing-container-env.txt
```

Look for:

```text
Server name/title/display name
Region
Public IP / external address
Battlegroup or server ID
Revision/version
Authentication token presence, without exposing the token
```

## 4. Check Network Reachability for Listing and Game Paths

Run on: Linux host or Linux VM shell

```bash
sudo ss -tulpen | tee server-listing-listeners.txt
```

Run on: Windows host PowerShell

```powershell
Get-NetUDPEndpoint | Sort-Object LocalPort | Tee-Object -FilePath server-listing-udp-listeners.txt
Get-NetTCPConnection | Sort-Object LocalPort | Tee-Object -FilePath server-listing-tcp-listeners.txt
```

Run in: cloud provider console, only if cloud-hosted

```text
Check that security rules allow discovered listing, game, and management ports.
Check that the public IP shown by the cloud provider matches the public IP configured or advertised by the server.
```

## 5. Interpret Results

```text
Server name or region is wrong:
  Check control panel settings, launch arguments, and environment variables.

Server is not visible and no gateway/listing logs appear:
  Confirm the correct service/log source and whether the listing service started.

Gateway/listing logs show failed registration or heartbeat:
  Check service token, public IP, revision/version, and outbound network access.

Server is visible but cannot be selected:
  Check password/authentication, capacity/lock state, and game listener reachability.

Server is visible to some players but not others:
  Compare player region/filter settings, client cache, network path, and public listing state.
```

## 6. Evidence to Escalate

```text
Expected and actual server listing behavior
Exact UTC check time
Gateway/control-plane logs
Runtime config values with secrets redacted
Listener state
Cloud/security/firewall evidence if applicable
Screenshots of listing behavior, with personal information redacted if needed
```
