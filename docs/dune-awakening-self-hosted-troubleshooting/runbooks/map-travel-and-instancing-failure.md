# Runbook: Map Travel and Instancing Failure

Use this runbook when players can log in but fail when moving between maps, zones, partitions, dungeons, deep desert, overmap, or any instanced destination.

Start only after the hosting platform and runtime/orchestration layer are identified.

Do not record client names, player names, chat names, Discord names, raw account identifiers, or other personal information in the troubleshooting notes. Use generic labels such as `PLAYER_1`, `CLIENT_IP`, `SOURCE_MAP`, and `DESTINATION_MAP`.

## 1. Capture the User-Visible Travel Symptom

Ask the player or environment owner:

```text
Where did the player start?
Where was the player trying to go?
What action triggered travel?
Did the player hang, disconnect, return, or receive an error?
Does the issue happen for one player or many players?
Does the issue happen every time or intermittently?
Which destinations still work?
Which destinations fail?
Exact UTC time of the test:
```

Record:

```text
Known working source:
Known working destination:
Known failing source:
Known failing destination:
Expected behavior:
Actual behavior:
```

## 2. Capture Control-Plane or Director Logs

Run on: Docker host shell, only if the control-plane/director is a container

```bash
docker logs --since 30m "$DIRECTOR_SERVICE" 2>&1 | tee travel-control-plane.log
```

Run on: Linux host or Linux VM shell if logs are file-based

```bash
grep -RniE 'Travel request|Received travel request|travel queue|Travel grant|travel completion|completion validation|LoginRequest|ServerState|partition|serverId|instancing|ServerFull|DestinationPartitionId|ServerLoginToken|failed|error|exception|timeout' "$LOG_PATH" 2>/dev/null | head -400 | tee travel-control-plane-search.txt
```

Run on: Windows host PowerShell if logs are file-based

```powershell
Select-String -Path "$env:LOG_PATH\**\*.log" -Pattern 'Travel request','Received travel request','travel queue','Travel grant','travel completion','completion validation','LoginRequest','ServerState','partition','serverId','instancing','ServerFull','DestinationPartitionId','ServerLoginToken','failed','error','exception','timeout' -ErrorAction SilentlyContinue | Select-Object -First 400 | Tee-Object -FilePath travel-control-plane-search.txt
```

Run in: control panel UI

```text
Open the control-plane/director log viewer.
Export or copy logs covering the failed travel test window.
```

## 3. Capture Source and Destination Server Logs

Run on: Docker host shell, only if source or destination are containers

```bash
docker logs --since 30m "$DESTINATION_SERVICE" 2>&1 | tee travel-destination.log
```

Run on: Linux host or Linux VM shell if logs are file-based

```bash
grep -RniE 'TravelEvent|PreLogin|Welcome|VerifyFlsIdentity|VerifyFlsAuthorization|Completion|DatabaseLogin|CharacterDownload|Join|GameModeLogin|StartingNewPlayer|FlsLogin|LoadPlayerActors|FinishSpawn|Grace Period|Disconnected|TravelFromMap|client travel|failed|error|exception|timeout' "$LOG_PATH" 2>/dev/null | head -500 | tee travel-destination-search.txt
```

Run on: Windows host PowerShell if logs are file-based

```powershell
Select-String -Path "$env:LOG_PATH\**\*.log" -Pattern 'TravelEvent','PreLogin','Welcome','VerifyFlsIdentity','VerifyFlsAuthorization','Completion','DatabaseLogin','CharacterDownload','Join','GameModeLogin','StartingNewPlayer','FlsLogin','LoadPlayerActors','FinishSpawn','Grace Period','Disconnected','TravelFromMap','client travel','failed','error','exception','timeout' -ErrorAction SilentlyContinue | Select-Object -First 500 | Tee-Object -FilePath travel-destination-search.txt
```

## 4. Capture Listener State During One Controlled Test

Run on: Linux host or Linux VM shell

```bash
sudo ss -uapn | tee travel-udp-listeners-before.txt
while true; do
  echo "===== $(date -u --iso-8601=seconds) ====="
  sudo ss -uapn || true
  sleep 2
done | tee travel-udp-listeners-during.txt
```

Run on: Windows host PowerShell

```powershell
Get-NetUDPEndpoint | Sort-Object LocalPort | Tee-Object -FilePath travel-udp-listeners-before.txt
while ($true) {
  "===== $(Get-Date -Format u) =====" | Tee-Object -FilePath travel-udp-listeners-during.txt -Append
  Get-NetUDPEndpoint | Sort-Object LocalPort | Tee-Object -FilePath travel-udp-listeners-during.txt -Append
  Start-Sleep -Seconds 2
}
```

## 5. Capture Traffic During One Controlled Test

Run on: Linux host or Linux VM shell if `CLIENT_IP` is known

```bash
sudo tcpdump -ni any -vv "host ${CLIENT_IP}" | tee travel-client-packets.log
```

Run on: Linux host or Linux VM shell if `CLIENT_IP` is unknown

```bash
sudo tcpdump -ni any -vv 'udp or tcp' | tee travel-server-packets.log
```

Run on: Windows host PowerShell

```powershell
pktmon start --capture --pkt-size 0
```

Stop Windows capture after the test:

```powershell
pktmon stop
pktmon format PktMon.etl -o travel-pktmon-capture.txt
```

## 6. Validate Dynamic Handoff Values

For dynamically spawned or instanced destinations, inspect the control-plane response and destination lifecycle together.

Record:

```text
Travel response code:
Server state:
ServerFull value:
DestinationPartitionId present: yes/no
ServerLoginToken present: yes/no
Destination server ID:
Destination active players, if shown:
Queue length, if shown:
```

Interpretation:

```text
ServerState=Ready but DestinationPartitionId is empty:
  The destination may exist, but the player was not given a usable destination assignment.

ServerFull=True while destination state shows no players:
  Investigate stale capacity state, queue state, or instancing metadata.

ServerLoginToken is empty:
  The source/client may not have a valid destination login handoff.

The same request later appears with missing or generic origin context:
  Treat this as a possible orphaned, retry, or cleanup-path symptom. Correlate with queue state and destination logs.
```

## 7. Interpret Results

```text
Control plane never logs the travel request:
  Check source server, client action, log source, and messaging path.

Control plane logs travel request but no destination appears:
  Check instance spawn, service manager, runtime/orchestration, and resource limits.

Control plane reports destination ready but no usable partition/token is returned:
  Check director allocation, capacity state, queue state, and instancing metadata.

Destination appears but no expected listener appears:
  Check launch arguments, bind address, runtime config, and port allocation.

Destination listener exists but no packets arrive:
  Check cloud firewall, host firewall, NAT, Hyper-V/Proxmox networking, Docker publishing, and advertised address.

Packets arrive but no replies leave:
  Check destination process health, auth/session state, and application errors.

Destination logs PreLogin or lifecycle stages then disconnects:
  Handoff reached the destination. Continue investigating session lifecycle, persistence, cleanup, and return/instance completion.

Failure happens only for some destinations:
  Compare working and failing destination launch arguments, listeners, logs, map mode, capacity state, and runtime path.
```

## 8. Evidence to Escalate

```text
Known working source/destination
Known failing source/destination
Exact UTC test time
Control-plane/director logs
Source and destination server logs
Travel response values
Listener before/during/after
Packet capture
Runtime/orchestration state before and after
Cloud/hypervisor/firewall evidence if applicable
```

## 9. Privacy and Redaction Requirements

Before adding evidence to notes or reports, remove:

```text
Client names
User names
Discord names
Player display names
Raw account IDs unless vendor-required and approved
Server passwords
Tokens
Secrets
Private keys
```

Keep generalized operational values when needed for troubleshooting:

```text
Map names
Partition numbers
Service names
Container names
Port numbers
Local file paths
```
