# Runbook: Dynamic Instancing and Handoff Validation

Use this runbook when players can reach some destinations but fail, hang, queue indefinitely, or disconnect when traveling to dynamically spawned or instanced destinations.

This runbook is intentionally generic. Do not include customer names, chat names, player names, raw account IDs, Discord identifiers, or other personal information in the troubleshooting notes.

Start only after the hosting platform and runtime/orchestration layer are identified.

## 1. Define the Working and Failing Paths

Ask the environment owner to identify one known-working destination and one known-failing destination.

```text
Known working source:
Known working destination:
Known failing source:
Known failing destination:
Failure type:
Exact UTC time of test:
```

Do not test multiple failing destinations at once. Pick one destination and capture one clean reproduction.

## 2. Identify Whether the Destination Is Static or Dynamically Spawned

Run on: Linux host or Linux VM shell where logs are stored

```bash
grep -RniE 'ClassicalInstancing|Dimension|ServerState|travel queue|spawn|partition|DestinationPartition|ServerFull|ServerLoginToken' "$LOG_PATH" 2>/dev/null | head -500 | tee instancing-mode-search.txt
```

Run on: Windows host PowerShell where logs are stored

```powershell
Select-String -Path "$env:LOG_PATH\**\*.log" -Pattern 'ClassicalInstancing','Dimension','ServerState','travel queue','spawn','partition','DestinationPartition','ServerFull','ServerLoginToken' -ErrorAction SilentlyContinue | Select-Object -First 500 | Tee-Object -FilePath instancing-mode-search.txt
```

Record:

```text
Destination map:
Destination partition, if known:
Instancing mode, if known:
Static or dynamic destination:
Server ID, if known:
```

## 3. Capture Director / Control-Plane Travel Response

Run on: Docker host shell, only if the director/control-plane is a container

```bash
docker logs --since 30m "$DIRECTOR_SERVICE" 2>&1 | \
  grep -Ei 'Travel request|Received travel request|Travel response|Travel grant|completion validation|ServerState|ServerFull|DestinationPartitionId|ServerLoginToken|partition|serverId|ClassicalInstancing|failed|error|warning' | \
  tee dynamic-instancing-director.log
```

Run on: Linux host or Linux VM shell if logs are file-based

```bash
grep -RniE 'Travel request|Received travel request|Travel response|Travel grant|completion validation|ServerState|ServerFull|DestinationPartitionId|ServerLoginToken|partition|serverId|ClassicalInstancing|failed|error|warning' "$LOG_PATH" 2>/dev/null | tee dynamic-instancing-director-search.txt
```

Run on: Windows host PowerShell if logs are file-based

```powershell
Select-String -Path "$env:LOG_PATH\**\*.log" -Pattern 'Travel request','Received travel request','Travel response','Travel grant','completion validation','ServerState','ServerFull','DestinationPartitionId','ServerLoginToken','partition','serverId','ClassicalInstancing','failed','error','warning' -ErrorAction SilentlyContinue | Tee-Object -FilePath dynamic-instancing-director-search.txt
```

## 4. Validate the Destination Handoff Values

A usable destination handoff normally needs a concrete destination and a way for the client/session to authenticate or continue into that destination.

Record from the travel response:

```text
Response code:
Server state:
ServerFull value:
DestinationPartitionId:
ServerLoginToken present: yes/no
Queue length:
Destination server ID:
Destination address/port, if present:
```

Interpretation:

```text
ServerState is Ready but DestinationPartitionId is empty:
  The destination may exist but was not allocated to the player.

ServerState is Ready but ServerFull is True while no active players are shown:
  Investigate stale capacity state, queue state, or instancing metadata.

ServerLoginToken is empty:
  The client may not have a usable destination login handoff.

Response remains Queued and never changes:
  Investigate director queue handling, destination availability, and queue consumers.
```

## 5. Capture Destination Lifecycle Logs

Run on: Docker host shell, only if the destination is a container

```bash
docker logs --since 30m "$DESTINATION_SERVICE" 2>&1 | \
  grep -Ei 'PreLogin|Welcome|VerifyFlsIdentity|VerifyFlsAuthorization|Completion|DatabaseLogin|CharacterDownload|Join|GameModeLogin|StartingNewPlayer|FlsLogin|LoadPlayerActors|FinishSpawn|Grace Period|Disconnected|Travel|failed|error|warning' | \
  tee dynamic-instancing-destination.log
```

Run on: Linux host or Linux VM shell if logs are file-based

```bash
grep -RniE 'PreLogin|Welcome|VerifyFlsIdentity|VerifyFlsAuthorization|Completion|DatabaseLogin|CharacterDownload|Join|GameModeLogin|StartingNewPlayer|FlsLogin|LoadPlayerActors|FinishSpawn|Grace Period|Disconnected|Travel|failed|error|warning' "$LOG_PATH" 2>/dev/null | tee dynamic-instancing-destination-search.txt
```

Interpretation:

```text
Destination shows no PreLogin or login lifecycle for the test:
  The handoff likely failed before the player reached the destination process.

Destination reaches PreLogin or authorization stages:
  The handoff reached the destination; continue into auth/session/persistence checks.

Destination reaches FinishSpawn then later disconnects:
  Initial handoff worked; investigate post-login session state, travel completion, cleanup, or return lifecycle.
```

## 6. Validate Dynamic Game and Server-to-Server Ports

Do not assume fixed ports. Discover the allocated ports from logs, launch arguments, and active listeners.

Run on: Linux host or Linux VM shell

```bash
ps -ef | grep -Ei 'DuneSandbox|Awakening|Sandbox' | grep -v grep | sed -E 's/(Token|Secret|Password|Key)=[^ ]+/\1=<redacted>/Ig' | tee dynamic-instancing-processes.txt
sudo ss -uapn | tee dynamic-instancing-udp-listeners.txt
```

Run on: Windows host PowerShell

```powershell
Get-CimInstance Win32_Process | Where-Object { $_.CommandLine -match 'Dune|Awakening|Sandbox' } | Select-Object ProcessId, CommandLine | Tee-Object -FilePath dynamic-instancing-processes.txt
Get-NetUDPEndpoint | Sort-Object LocalPort | Tee-Object -FilePath dynamic-instancing-udp-listeners.txt
```

Run on: Docker host shell, only if Docker is confirmed

```bash
docker ps --format 'table {{.Names}}\t{{.Image}}\t{{.Status}}\t{{.Ports}}' | tee dynamic-instancing-docker-ps.txt
```

Check whether the destination process has both expected roles:

```text
Game/client UDP port present: yes/no/unknown
Server-to-server or IGW UDP port present: yes/no/unknown
Bind address correct: yes/no/unknown
External or advertised address correct: yes/no/unknown
```

## 7. Compare Working and Failing Destinations

Compare one known-working destination with one known-failing destination.

```text
Map/destination name
Partition ID
Instancing mode
Server ID
Launch arguments
Game/client port
Server-to-server or IGW port
Bind address
External/advertised address
Travel response code
ServerFull value
DestinationPartitionId present
ServerLoginToken present
Destination lifecycle stage reached
```

Interpretation:

```text
Working destination returns a concrete partition and token but failing destination does not:
  Focus on director allocation, queue handling, capacity state, or instancing metadata.

Working destination has both game and server-to-server ports but failing destination does not:
  Focus on dynamic spawn argument passing and runtime port binding.

Both destinations receive handoff but only one disconnects after spawn:
  Focus on session lifecycle, persistence, completion, cleanup, or map-specific runtime errors.
```

## 8. Evidence to Escalate

```text
Known working destination evidence
Known failing destination evidence
Director/control-plane travel response
Destination lifecycle logs
Process launch arguments with secrets redacted
Active UDP listener output
RabbitMQ queue state, if applicable
Packet capture, if available
Exact UTC test time
```

## 9. Redaction Rules

Before sharing, remove:

```text
Personal names
Chat or Discord names
Raw player/account IDs unless vendor-required and approved
Server passwords
Tokens
Secrets
Private keys
Public IPs if sharing broadly
```

Keep operational values needed for troubleshooting when sharing with trusted support staff:

```text
Map names
Partition numbers
Service names
Container names
Port numbers
Local file paths
```
