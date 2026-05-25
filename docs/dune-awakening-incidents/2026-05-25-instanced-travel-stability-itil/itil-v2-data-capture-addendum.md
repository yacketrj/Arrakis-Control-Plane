# ITIL v2 Data-Capture Addendum

This addendum updates the case-specific ITIL incident/RCA record with the next required evidence capture. It is intentionally separate from the reusable troubleshooting guide.

## Status

```text
Incident status: Investigation / Diagnosis
Resolution status: Not proven
Current strongest RCA direction: Director / ClassicalInstancing allocation or handoff behavior
Required next action: synchronized evidence capture of one controlled failed travel attempt
```

## Environment Context Used for the Capture Plan

```text
Platform: Proxmox bare-metal host
Guest OS: Ubuntu 24.04 VM
Control layer: AMP
Runtime/orchestration: Docker behind AMP
Primary command location: Ubuntu VM shell where Docker is available
```

## Non-PII Handling

Do not include client names, user names, chat names, player names, raw account identifiers, passwords, tokens, public IP values, private keys, or other PII in the incident report or shared evidence package.

## Evidence Needed Next

Collect these at the same time during one controlled failed travel attempt:

```text
1. Director/control-plane logs
2. Source map logs
3. Destination map logs
4. RabbitMQ queue/consumer state
5. UDP listener output
6. Optional packet capture, only if approved
7. Final Docker/container state
```

## Key Question

```text
Does the destination map ever receive PreLogin / login / authorization for the failed travel attempt?
```

If the destination does not receive PreLogin/login/authorization, the investigation remains focused on director allocation, queue response, capacity state, token issuance, or pre-destination handoff.

If the destination does receive PreLogin/login/authorization, the investigation shifts to destination runtime, authentication/session state, persistence, or network response behavior.

## Primary Capture Workflow

Run all commands from the Ubuntu 24.04 VM shell unless noted otherwise.

### 1. Create capture directory

```bash
CAPTURE_ID="$(date -u +%Y%m%d-%H%M%S)"
mkdir -p "$HOME/dune-incident-capture-$CAPTURE_ID"
cd "$HOME/dune-incident-capture-$CAPTURE_ID"
```

### 2. Discover containers

```bash
docker ps --format 'table {{.Names}}\t{{.Image}}\t{{.Status}}\t{{.Ports}}' | tee docker-ps-before.txt
docker ps --format '{{.Names}}' | sort | tee all-container-names.txt
```

### 3. Set variables

```bash
DIRECTOR_SERVICE="<DIRECTOR_SERVICE>"
SOURCE_SERVICE="<SOURCE_SERVICE>"
DESTINATION_SERVICE="<DESTINATION_SERVICE>"
RABBITMQ_SERVICE="<RABBITMQ_SERVICE>"
DESTINATION_MAP="<DESTINATION_MAP>"
CLIENT_IP="<CLIENT_IP_IF_APPROVED>"
```

### 4. Start captures in separate terminals

```bash
docker logs -f "$DIRECTOR_SERVICE" 2>&1 | tee director-travel-capture.log
docker logs -f "$SOURCE_SERVICE" 2>&1 | tee source-map-travel-capture.log
docker logs -f "$DESTINATION_SERVICE" 2>&1 | tee destination-map-lifecycle-capture.log
```

```bash
while true; do
  echo "===== $(date -u --iso-8601=seconds) ====="
  docker exec "$RABBITMQ_SERVICE" rabbitmqctl list_queues -p / \
    name messages messages_ready messages_unacknowledged consumers state 2>/dev/null
  sleep 2
done | tee rabbitmq-during-travel.log
```

```bash
while true; do
  echo "===== $(date -u --iso-8601=seconds) ====="
  sudo ss -uapn | tee -a udp-listeners-during-travel.log
  sleep 2
done
```

Optional packet capture, only if approved:

```bash
sudo tcpdump -ni any -vv "host ${CLIENT_IP}" | tee packet-capture-client.log
```

General UDP capture without recording a client IP:

```bash
sudo tcpdump -ni any -vv 'udp' | tee packet-capture-udp-general.log
```

### 5. Run one controlled failed travel attempt

Run exactly one failed travel attempt after all captures are active. Do not retry repeatedly in the same capture window.

### 6. Capture final state

```bash
docker ps --format 'table {{.Names}}\t{{.Image}}\t{{.Status}}\t{{.Ports}}' | tee docker-ps-after.txt
sudo ss -uapn | tee udp-listeners-after.txt

docker inspect "$DESTINATION_SERVICE" --format \
'Image={{.Config.Image}} User={{.Config.User}} NetworkMode={{.HostConfig.NetworkMode}} Status={{.State.Status}} ExitCode={{.State.ExitCode}} OOMKilled={{.State.OOMKilled}}' \
| tee destination-container-state.txt
```

## Interpretation

```text
No destination PreLogin/login lifecycle:
  Focus on director allocation, queue response, capacity state, token issuance, or pre-destination handoff.

Destination receives PreLogin/auth but fails later:
  Focus on destination runtime, authentication/session, persistence, or local network response.

Destination reaches FinishSpawn then disconnects:
  Focus on post-login lifecycle, travel completion, cleanup, return path, or instance lifecycle.

UDP listener missing:
  Focus on dynamic spawn arguments, port allocation, bind address, or runtime startup.

Packets arrive but no replies leave:
  Focus on application handling, destination process health, local firewall, or auth/session state.
```

## Closure Reminder

Do not close the incident until at least one previously failing ClassicalInstancing destination is successfully reached and the director, destination, listener, and traffic evidence support stable travel completion.
