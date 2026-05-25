# Runbook: Instanced Travel Dynamic Spawn Validation

Use this runbook when players can reach some persistent or primary destinations, but travel to dynamically spawned or instanced destinations hangs, queues indefinitely, disconnects, or fails after partial login/travel progress.

This is a reusable troubleshooting runbook. Do not include client names, player names, Discord names, raw account IDs, passwords, tokens, public IPs, private IPs, or other PII in this document.

Start only after the hosting platform and runtime/orchestration layer are identified.

---

## 1. When to Use This Runbook

Use this runbook when the evidence pattern looks like one or more of the following:

```text
Primary or always-on maps work, but many instanced maps fail.
The director or control plane receives a travel request.
The destination is a ClassicalInstancing, dungeon, story, or dynamically spawned map.
The player hangs, disconnects, or returns instead of completing travel.
The destination sometimes reaches PreLogin, authorization, Join, or FinishSpawn and later disconnects.
The control plane shows queued/full/unavailable behavior even when the destination appears ready.
A destination process starts but expected dynamic listeners are missing.
```

Do not assume a network outage just because travel hangs. First prove where the travel lifecycle stops.

---

## 2. Capture the User-Defined Symptom

Ask the environment owner:

```text
Which destinations work?
Which destinations fail?
Is the failing destination always-on or dynamically spawned?
Does the player hang, disconnect, return, or receive an error?
Does the failure occur every time or intermittently?
Does the failure affect one player, some players, or all players?
Exact UTC time of the failed travel attempt:
```

Record:

```text
Known working source:
Known working destination:
Known failing source:
Known failing destination:
Observed client behavior:
Exact UTC test window:
```

---

## 3. Capture Three Evidence Streams at the Same Time

During one controlled failed travel attempt, capture all three streams together:

```text
1. Director or control-plane logs
2. Destination map/server logs
3. UDP listener and packet/traffic evidence
```

If these are captured separately, the timeline may not prove where the request stopped.

---

## 4. Director or Control-Plane Capture

Run on: Docker host shell, only if the director/control-plane is a container

```bash
docker logs -f <DIRECTOR_SERVICE> 2>&1 | \
  sed -E 's/(Token|Secret|Password|Key)=[^ ]+/\1=<redacted>/Ig' | \
  grep -Ei 'travel request|Travel response|Travel grant|travel completion|completion validation|ServerState|partition|serverId|queued|full|failed|error|warning' | \
  tee director-instanced-travel-capture.log
```

Run on: Linux host or Linux VM shell if logs are file-based

```bash
tail -F <DIRECTOR_LOG_FILE> 2>&1 | \
  sed -E 's/(Token|Secret|Password|Key)=[^ ]+/\1=<redacted>/Ig' | \
  grep -Ei 'travel request|Travel response|Travel grant|travel completion|completion validation|ServerState|partition|serverId|queued|full|failed|error|warning' | \
  tee director-instanced-travel-capture.log
```

Run in: control panel UI

```text
Open the director/control-plane console or log viewer.
Export logs covering the exact UTC test window.
Redact secrets before sharing broadly.
```

Record:

```text
Did the director receive the travel request?
Did the director create or queue the request?
Did the director return a destination partition/server ID?
Did the director return a login token or equivalent handoff token?
Did the director report queued, full, blocked, unavailable, failed, or timed out?
```

---

## 5. Destination Map Capture

Run on: Docker host shell, only if the destination map/server is a container

```bash
docker logs -f <DESTINATION_SERVICE> 2>&1 | \
  sed -E 's/(Token|Secret|Password|Key)=[^ ]+/\1=<redacted>/Ig' | \
  grep -Ei 'PreLogin|Welcome|VerifyFlsIdentity|VerifyFlsAuthorization|Completion|DatabaseLogin|CharacterDownload|Join|GameModeLogin|StartingNewPlayer|FlsLogin|LoadPlayerActors|FinishSpawn|Grace Period|Disconnected|Disconnect|TravelFromMap|client travel|failed|error|warning' | \
  tee destination-instanced-travel-capture.log
```

Run on: Linux host or Linux VM shell if logs are file-based

```bash
tail -F <DESTINATION_LOG_FILE> 2>&1 | \
  sed -E 's/(Token|Secret|Password|Key)=[^ ]+/\1=<redacted>/Ig' | \
  grep -Ei 'PreLogin|Welcome|VerifyFlsIdentity|VerifyFlsAuthorization|Completion|DatabaseLogin|CharacterDownload|Join|GameModeLogin|StartingNewPlayer|FlsLogin|LoadPlayerActors|FinishSpawn|Grace Period|Disconnected|Disconnect|TravelFromMap|client travel|failed|error|warning' | \
  tee destination-instanced-travel-capture.log
```

Run in: control panel UI

```text
Open the destination map/server console or log viewer.
Export logs covering the exact UTC test window.
```

Record:

```text
Did the destination process start?
Did the destination receive PreLogin or an equivalent first player connection stage?
Did the destination reach authorization, database login, character download, Join, or FinishSpawn?
Did the destination later disconnect the player through grace-period or cleanup handling?
```

---

## 6. Dynamic Port and Listener Validation

Instanced maps may require more than one port path. Validate both the player-facing game/client port and any server-to-server, gateway, or IGW-style port used by the runtime.

Run on: Linux host or Linux VM shell

```bash
sudo ss -uapn | tee udp-listeners-before.txt

while true; do
  echo "===== $(date -u --iso-8601=seconds) ====="
  sudo ss -uapn | grep -Ei 'Dune|777|778|779|780|781|788|789|790|791|792' || true
  sleep 2
done | tee udp-listeners-during-instanced-travel.log
```

Run on: Windows host PowerShell

```powershell
Get-NetUDPEndpoint | Sort-Object LocalPort | Tee-Object -FilePath udp-listeners-before.txt
while ($true) {
  "===== $(Get-Date -Format u) =====" | Tee-Object -FilePath udp-listeners-during-instanced-travel.log -Append
  Get-NetUDPEndpoint | Sort-Object LocalPort | Tee-Object -FilePath udp-listeners-during-instanced-travel.log -Append
  Start-Sleep -Seconds 2
}
```

Record:

```text
Expected player-facing game port:
Expected server-to-server / gateway / IGW-style port:
Actual player-facing listener:
Actual server-to-server / gateway / IGW-style listener:
Owning process/container:
Bind address:
```

Interpretation:

```text
Destination process starts but expected game listener is missing:
  Focus on destination startup arguments and bind address.

Destination process starts with game listener but expected server-to-server listener is missing:
  Focus on dynamic spawn argument passing and runtime/orchestration configuration.

Both listeners exist but no packets arrive:
  Focus on firewall, NAT, hypervisor networking, cloud rules, Docker publishing, or advertised address.

Packets arrive and replies leave, but travel still fails:
  Focus on game session state, travel completion, cleanup, and director handoff state.
```

---

## 7. Packet Capture

Run on: Linux host or Linux VM shell if the client IP is known and approved for capture

```bash
sudo tcpdump -ni any -vv "host ${CLIENT_IP} and (udp or tcp)" | tee instanced-travel-client-packets.log
```

Run on: Linux host or Linux VM shell if the client IP is unknown

```bash
sudo tcpdump -ni any -vv 'udp or tcp' | tee instanced-travel-general-packets.log
```

Run on: Windows host PowerShell as Administrator

```powershell
pktmon start --capture --pkt-size 0
```

Stop after the test:

```powershell
pktmon stop
pktmon format PktMon.etl -o instanced-travel-pktmon-capture.txt
```

---

## 8. Compare Working and Failing Destinations

Pick one known-working destination and one known-failing instanced destination. Compare them using the same evidence categories.

```text
Runtime launch command
Map or destination name
Partition or instance identifier
Game/client port
Server-to-server / gateway / IGW-style port
Bind address
Advertised address
Destination server ID
Queue response
Travel grant response
Login/handoff token presence
Destination PreLogin/authorization/Join/FinishSpawn stages
Disconnect or cleanup behavior
```

Run on: Docker host shell, only if Docker is confirmed

```bash
for c in $(docker ps --format '{{.Names}}'); do
  echo "===== $c ====="
  docker inspect "$c" --format 'Image={{.Config.Image}} NetworkMode={{.HostConfig.NetworkMode}} Args={{json .Args}}'
  docker inspect "$c" --format '{{range .Config.Env}}{{println .}}{{end}}' | sed -E 's/(TOKEN|SECRET|PASSWORD|KEY)=.*/\1=<redacted>/Ig'
done | tee instanced-travel-runtime-comparison.txt
```

---

## 9. Evidence Patterns and What They Mean

```text
Director never logs the request:
  The issue is before or outside the director path. Check source map, player action, queue path, and log source.

Director logs the request but does not create a travel grant:
  Focus on director queue logic, destination availability, server state, and capacity flags.

Director reports ready state but also queued/full/unavailable without a destination handoff:
  Focus on allocation state, capacity metadata, stale server state, or handoff-token issuance.

Destination never logs PreLogin or equivalent connection stage:
  The player likely never reached the destination process. Focus on director handoff, token issuance, listener, or routing.

Destination reaches PreLogin/authorization/Join/FinishSpawn and later disconnects:
  Initial handoff worked. Focus on post-login lifecycle, travel completion, persistence, cleanup, or return-handoff logic.

Destination process starts but expected listener is missing:
  Focus on runtime spawn arguments and dynamic port passing.

RabbitMQ authenticates and carries queue messages:
  Messaging path is at least partially working for that event. Continue correlating queue state with director and destination logs.
```

---

## 10. Escalation Evidence

Before escalating, include:

```text
User-defined symptom with PII removed
Known working destination
Known failing destination
Exact UTC test window
Director/control-plane logs
Destination map logs
UDP listener before/during/after
Packet capture, if available
Runtime launch arguments with secrets redacted
Queue state, if messaging is involved
Working-vs-failing comparison
Current hypothesis and confidence level
```

---

## 11. Privacy and Redaction Requirements

Before sharing outside the immediate support team, remove:

```text
Client names
User names
Discord names or handles
Player display names
Raw account IDs
Passwords
Tokens
Secrets
Public/private IPs when not required by the receiving team
```

Keep only generalized operational values needed to troubleshoot, such as placeholder service names, placeholder container names, placeholder destination names, and port roles.
