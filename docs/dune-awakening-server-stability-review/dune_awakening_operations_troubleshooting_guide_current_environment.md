# Dune: Awakening Self-Hosted Server Operations Troubleshooting Guide

**Project name:** Dune: Awakening Dedicated Server Stability Review

**Document type:** Operations troubleshooting guide for support and escalation teams

**Audience:** Front-line support, operations, server administrators, and technical escalation staff

---

## Documentation Evidence Scope

This guide is scoped to the current client environment and should not mix in findings from unrelated environments.

Allowed evidence sources for this guide:

```text
1. logs.zip
2. User-provided commands, outputs, client conversation excerpts, and troubleshooting notes shared after logs.zip was uploaded
3. The established operations troubleshooting-guide structure, used only as a documentation format reference
```

Do not import conclusions from unrelated deployments unless that information is explicitly reintroduced for this environment.

---

## Redaction Rules

Before sharing logs, screenshots, command output, or generated reports outside the support team, redact sensitive values.

Use these placeholders:

```text
<ENVIRONMENT_OWNER>
<PLAYER_ID>
<PLAYER_NAME>
<CLIENT_IP>
<PUBLIC_IP>
<PRIVATE_IP>
<TOKEN>
<SECRET>
<ACCOUNT_ID>
<OCID>
<INSTANCE_PATH>
<CONTAINER_NAME>
<DESTINATION_MAP>
<DIRECTOR_CONTAINER>
<RABBITMQ_CONTAINER>
<DESTINATION_CONTAINER>
```

Do not include personal names, chat display names, player names, raw player IDs, raw account IDs, raw IPs, tokens, JWTs, RabbitMQ secrets, cloud resource IDs, or unredacted authentication material.

---

## Plain-Language Summary

Players can reach the normal starting areas, including partition `1` and partition `8`, but many other travel destinations hang or disconnect.

The current evidence does not show a simple “the server is offline” or “all ports are blocked” issue. The logs show that some instanced destinations receive travel activity and progress through several travel/login steps before later failing.

The most likely area to investigate is the dynamic instanced travel path. These destinations differ from always-running maps because they may need to be spawned on demand, assigned ports, receive the player, complete travel, and clean up correctly afterward.

Highest-priority technical check:

```text
When an instanced map is spawned, does it receive both:
1. its dynamic game/client UDP port, and
2. its dynamic IGW/server-to-server UDP port?
```

If both ports are allocated but only one is passed into the game process, players may reach part of the travel flow and then hang or disconnect.

---

## Confirmed Environment Summary

The environment owner provided the following hosting model:

```text
Bare-metal hypervisor: Proxmox
Guest OS: Ubuntu 24.04 VM
Control panel: AMP
Container runtime: Docker behind AMP
CPU allocation: 16 cores from an AMD Ryzen 9 9950X3D
VM memory allocation: 96 GB
Host memory: 192 GB
Network: data center hosted, 5 Gbps fiber connection, 5 static public IPs
AMP instance path example: <INSTANCE_PATH>
```

Support impact:

```text
1. Treat AMP as the control panel and service manager.
2. Docker is available, but paths may not match a normal standalone Docker install.
3. Do not assume standard Docker Compose directories.
4. Always map the active AMP instance path before editing files.
5. Network traffic may pass through data center networking, Proxmox, the Ubuntu VM, AMP, Docker, and the game process.
6. If a command fails because of AMP pathing or syntax, collect the equivalent AMP-managed path/output instead of assuming the component is absent.
```

---

## Important OS Clarification

UE5 server logs may show a line similar to:

```text
LogInit: Log: OS: Debian GNU/Linux 13 (trixie) (...), CPU: AMD Ryzen 9 9950X 16-Core Processor, GPU: UnknownVendor
```

This line shows what the Dune game process can see. It does not automatically prove the physical host OS or parent hypervisor.

```text
If the process is inside Docker:
  Debian may be the container userspace.
  The kernel is usually the Linux kernel visible to the container.

If the process is inside a Linux VM:
  The line reflects the Linux runtime context, not the Proxmox host OS.

Confirmed environment model:
  AMP-managed Docker inside an Ubuntu 24.04 VM on Proxmox bare metal.
```

---

## Key Terms for Support Staff

| Term | Meaning |
|---|---|
| Partition | A numbered game-world destination or map instance. |
| Partition 1 | `Survival_1`, a key starting/normal gameplay area. |
| Partition 8 | `DeepDesert_1`, another reachable destination in the current issue. |
| ClassicalInstancing | A map mode where destination servers may be spawned dynamically for travel. |
| Director | Service that receives travel/login requests and decides where players should go. |
| RabbitMQ | Messaging service used by the stack for game service communication and per-player queues. |
| IGW port | Server-to-server / internal game-world communication port used by map servers. |
| Game/client port | UDP port the game client connects to for a map server. |
| AMP | Control panel managing this server environment and Docker runtime layout. |
| Grace-period disconnect | Log pattern where a player is eventually disconnected after travel or instanced-map lifecycle handling. |

---

## Current Problem Statement

```text
Players can reach partition 1 and partition 8.
Players hang indefinitely or disconnect when traveling to many other partitions.
```

```text
Incident type: Gameplay service degradation
Affected function: Instanced travel
Most likely failure area: ClassicalInstancing dynamic map lifecycle
Not currently proven as: general internet outage, all-port block, or total Docker failure
```

---

## What We Know So Far

### 1. The issue is not limited to the director receiving nothing

The logs show the director processing ClassicalInstancing travel queues, including patterns like:

```text
Processing travel queue for ClassicalInstancing group <DESTINATION_MAP> (servers: [], num: 0)
```

Support interpretation:

```text
The director knows about the travel queues.
The destination may not have an active spawned server at that moment.
Dynamic spawn behavior must be checked.
```

### 2. Some instanced maps progress through travel stages

Observed travel/login stages include:

```text
PreLogin
VerifyFlsIdentity
VerifyFlsAuthorization
Completion
DatabaseLogin
CharacterDownload
Join
GameModeLogin
StartingNewPlayer
FlsLogin
LoadPlayerActors
FinishSpawn
```

If these stages appear, the player reached more than the first step. This points away from a simple ingress-only problem and toward later travel/session lifecycle behavior.

### 3. Some failures end as grace-period disconnects

Observed failure pattern:

```text
Grace Period:Disconnected from instanced map: <DESTINATION_MAP>
```

The destination may receive the player and then fail later during instanced map lifecycle, cleanup, persistence, or return-handoff handling.

### 4. Dynamic port handling is the top suspect

Dynamic maps should use:

```text
Game/client UDP ports: 7779-7810
IGW/server-to-server UDP ports: 7890-7921
```

Expected arguments include:

```text
-Port=<dynamic game port>
-ini:engine:[URL]:Port=<dynamic game port>
-IGWPort=<dynamic igw port>
-ini:engine:[URL]:IGWPort=<dynamic igw port>
-IGWBindAddress=<expected bind IP>
-ExternalAddress=<public IP>
```

---

## First-Response Triage Checklist

### Step 1 - Confirm the symptom

Record:

```text
Which destination map failed?
Did partition 1 still work?
Did partition 8 still work?
Did the player hang, disconnect, return, or receive an error?
What UTC time did the failed travel occur?
Was this fresh login travel or travel from inside the game?
```

### Step 2 - Confirm the environment

```bash
hostnamectl
cat /etc/os-release
uname -a
systemd-detect-virt -v || true
```

Expected environment classification:

```text
AMP-managed Docker deployment inside Ubuntu 24.04 VM on Proxmox bare metal.
```

### Step 3 - Confirm Docker visibility under AMP

```bash
docker ps --format 'table {{.Names}}\t{{.Image}}\t{{.Status}}\t{{.Ports}}' 2>/dev/null || true
```

If Docker does not work directly, note that AMP may require a different command path or permissions.

### Step 4 - Map the AMP instance path

```bash
find <INSTANCE_PATH> -maxdepth 3 -type d 2>/dev/null | head -100
```

### Step 5 - Do not change configuration yet

Before editing files or restarting services, collect enough evidence to identify where the travel attempt stops:

```text
Director request/grant/completion
Destination map spawn and lifecycle
UDP listener creation
Client packet path
RabbitMQ queue/consumer state
```

---

## Quick Health Check Commands

### Container and service view

```bash
docker ps --format 'table {{.Names}}\t{{.Image}}\t{{.Status}}\t{{.Ports}}' 2>/dev/null || true
```

### Dune server process view

```bash
ps -ef | grep DuneSandbox | grep -v grep | \
  sed -E 's/ServiceAuthToken=[^ ]+/ServiceAuthToken=<redacted>/g'
```

### Active UDP listener view

```bash
sudo ss -uapn | egrep ':7777|:7778|:7779|:778[0-9]|:779[0-9]|:780[0-9]|:7810|:7888|:7889|:789[0-9]|:790[0-9]|:791[0-9]|:792[0-1]' || true
```

### RabbitMQ queue view

```bash
docker exec <RABBITMQ_CONTAINER> rabbitmqctl list_queues \
  name messages messages_ready messages_unacknowledged consumers 2>/dev/null | \
  egrep -i 'name|login|travel|survival|overmap|deep|completion|serverState' || true
```

### Time sync view

```bash
date -u
timedatectl
```

---

## Runbook: Capture One Failed Travel Attempt

Capture these three things at the same time:

```text
1. Director logs
2. Destination map logs
3. Active UDP listener and traffic output
```

### Preparation

Pick one failing destination for the test. Do not test several maps at once.

Record:

```text
UTC start time:
Starting location:
Destination map:
Observed result:
Approximate failure time:
Client-side error or behavior:
```

Create a capture directory:

```bash
mkdir -p ~/dune-travel-capture-$(date -u +%Y%m%d-%H%M%S)
cd ~/dune-travel-capture-*
pwd
```

### Terminal 1 - Capture current containers

```bash
docker ps --format 'table {{.Names}}\t{{.Image}}\t{{.Status}}\t{{.Ports}}' 2>/dev/null | tee docker-ps-before.txt
```

### Terminal 2 - Capture director logs

```bash
docker logs -f <DIRECTOR_CONTAINER> 2>&1 | \
  sed -E 's/ServiceAuthToken=[^ ]+/ServiceAuthToken=<redacted>/g' | \
  egrep -i 'LoginRequest|Received travel request|Travel response|Travel grant|travel completion|completion validation|ServerState|OnServerStateReportReceived|out of order|partition|serverId|failed|error|warning' | \
  tee director-travel-capture.log
```

Look for `Received travel request`, `Travel grant`, `Travel completion`, `ServerState`, out-of-order state errors, and failed/error/warning lines.

### Terminal 3 - Find and capture destination map logs

Try Docker container names first:

```bash
docker ps --format '{{.Names}}' 2>/dev/null | \
  egrep -i 'arrakeen|harko|proces|dungeon|ecolab|hephaestus|waterfat|pit|artofkanly|overland|server' | \
  tee possible-destination-containers.txt
```

If the destination container exists:

```bash
docker logs -f <DESTINATION_CONTAINER> 2>&1 | \
  sed -E 's/ServiceAuthToken=[^ ]+/ServiceAuthToken=<redacted>/g' | \
  egrep -i 'Login request|TravelEvent|PreLogin|Welcome|CleaningUpOldConnections|VerifyFlsIdentity|VerifyFlsAuthorization|Completion|DatabaseLogin|CharacterDownload|Join|GameModeLogin|StartingNewPlayer|FlsLogin|LoadPlayerActors|FinishSpawn|Grace Period|Disconnected from instanced map|Disconnect|Close|TravelFromMap|Telling.*client travel|failed|error|warning' | \
  tee destination-map-capture.log
```

If AMP writes logs to files instead:

```bash
find <INSTANCE_PATH> -type f -name '*.log' 2>/dev/null | \
  egrep -i 'Arrakeen|Harko|Proces|Dungeon|Overland|Ecolab|Hephaestus|WaterFat|Pit|ArtOfKanly' | \
  tee possible-destination-log-files.txt
```

Then tail the selected log:

```bash
tail -F <DESTINATION_LOG_FILE> 2>&1 | \
  sed -E 's/ServiceAuthToken=[^ ]+/ServiceAuthToken=<redacted>/g' | \
  egrep -i 'Login request|TravelEvent|PreLogin|Welcome|CleaningUpOldConnections|VerifyFlsIdentity|VerifyFlsAuthorization|Completion|DatabaseLogin|CharacterDownload|Join|GameModeLogin|StartingNewPlayer|FlsLogin|LoadPlayerActors|FinishSpawn|Grace Period|Disconnected from instanced map|Disconnect|Close|TravelFromMap|Telling.*client travel|failed|error|warning' | \
  tee destination-map-capture.log
```

### Terminal 4 - Capture UDP listeners

Before travel:

```bash
sudo ss -uapn | egrep ':7777|:7778|:7779|:778[0-9]|:779[0-9]|:780[0-9]|:7810|:7888|:7889|:789[0-9]|:790[0-9]|:791[0-9]|:792[0-1]' | tee udp-listeners-before.txt
```

During travel:

```bash
while true; do
  echo "===== $(date -u --iso-8601=seconds) ====="
  sudo ss -uapn | egrep ':7777|:7778|:7779|:778[0-9]|:779[0-9]|:780[0-9]|:7810|:7888|:7889|:789[0-9]|:790[0-9]|:791[0-9]|:792[0-1]' || true
  sleep 2
done | tee udp-listeners-during-travel.log
```

Look for a dynamic game/client listener in `7779-7810` and a matching IGW/server-to-server listener in `7890-7921`.

### Terminal 5 - Capture UDP/TCP traffic

If the client IP is known:

```bash
sudo tcpdump -ni any -vv "host <CLIENT_IP> and (udp portrange 7777-7810 or udp portrange 7888-7921 or tcp port 31982 or tcp port 31983)" | tee tcpdump-client-travel.log
```

If the client IP is not known:

```bash
sudo tcpdump -ni any -vv '(udp portrange 7777-7810 or udp portrange 7888-7921 or tcp port 31982 or tcp port 31983)' | tee tcpdump-dune-ports-travel.log
```

Look for whether packets arrive, whether the server replies, and whether packets go to the port the destination server is listening on.

### Run the test

```text
1. Have the player start from a known working location.
2. Have the player travel to the selected failing destination.
3. Do not retry repeatedly during the same capture.
4. Let the attempt succeed, hang, disconnect, or return.
5. Record exact UTC start and failure time.
```

### Collect final state

```bash
docker ps --format 'table {{.Names}}\t{{.Image}}\t{{.Status}}\t{{.Ports}}' 2>/dev/null | tee docker-ps-after.txt

sudo ss -uapn | egrep ':7777|:7778|:7779|:778[0-9]|:779[0-9]|:780[0-9]|:7810|:7888|:7889|:789[0-9]|:790[0-9]|:791[0-9]|:792[0-1]' | tee udp-listeners-after.txt

docker exec <RABBITMQ_CONTAINER> rabbitmqctl list_queues name messages messages_ready messages_unacknowledged consumers 2>/dev/null | \
  egrep -i 'name|login|travel|survival|overmap|deep|completion|serverState' | tee rabbitmq-queues-after.txt || true
```

### Redact and package

```bash
for f in *.log *.txt; do
  [ -f "$f" ] || continue
  sed -i \
    -e 's/ServiceAuthToken=[^ ]*/ServiceAuthToken=<redacted>/g' \
    -e 's/eyJ[A-Za-z0-9._-]*/<redacted-jwt>/g' \
    -e 's/ocid1\.[A-Za-z0-9._-]*/<redacted-ocid>/g' \
    "$f"
done

cd ..
tar -czf dune-travel-capture-$(date -u +%Y%m%d-%H%M%S).tar.gz dune-travel-capture-*/
ls -lh dune-travel-capture-*.tar.gz
```

---

## How to Interpret the Capture

```text
Director does not log a travel request
  -> The request is not reaching director. Check source map, player action, and queue path.

Director logs travel request but no grant/completion
  -> Focus on director queue, destination availability, and server state for that map.

Destination container/log never appears
  -> Focus on ClassicalInstancing spawn trigger and AMP/Docker spawn path.

Destination appears but no dynamic UDP listener appears
  -> Focus on UE5 startup arguments and port binding.

Destination has game port but no IGW/server-to-server port
  -> Focus on dynamic spawn command passing IGWPort / URL IGWPort.

Destination receives player and reaches PreLogin/VerifyFls/Completion
  -> Initial handoff worked. Continue into travel lifecycle and persistence.

Destination reaches FinishSpawn then Grace Period disconnects
  -> Focus on post-login session state, travel completion cleanup, and return/instance lifecycle.

UDP packets arrive from client but no server replies
  -> Focus on destination game process handling, auth/session state, or local firewall.

Server replies but client still hangs
  -> Focus on client-side routing, session state, or higher-level game travel handling.
```

---

## Active Issue: RabbitMQ Per-Player Queue Delete Fails While Queue Is In Use

Observed pattern:

```text
accepting AMQP connection (<CLIENT_IP>:<port> -> <container-ip>:5672)
user '<PLAYER_ID>' authenticated and granted access to vhost '/'
operation queue.delete caused a channel exception precondition_failed: queue '<PLAYER_ID>_queue' in vhost '/' in use
client unexpectedly closed TCP connection
```

Meaning:

```text
RabbitMQ accepted the connection and authentication, but a queue delete operation failed because the per-player queue was still in use.
```

Likely causes:

```text
1. Overlapping AMQP connections for the same player.
2. Stale consumer on the per-player queue.
3. Reconnect behavior during a failed travel attempt.
4. Cleanup attempted while another channel still had the queue open.
```

Follow-up commands:

```bash
# List matching per-player queues
docker exec <RABBITMQ_CONTAINER> rabbitmqctl list_queues -p / \
  name durable auto_delete arguments consumers messages messages_ready messages_unacknowledged state | \
  grep '<PLAYER_ID>'

# List consumers
docker exec <RABBITMQ_CONTAINER> rabbitmqctl list_consumers -p / | \
  grep '<PLAYER_ID>'

# List active connections
docker exec <RABBITMQ_CONTAINER> rabbitmqctl list_connections \
  pid user peer_host peer_port state name | \
  grep '<PLAYER_ID>'
```

Do not delete a queue while it has consumers or while the player is actively connecting.

---

## Active Issue: Server State Arrives Out of Order

Observed director error:

```text
OnServerStateReportReceived: Received server state out of order
```

Meaning:

```text
The director received an older server-state report after already processing a newer report for the same server or partition.
```

Possible causes:

```text
1. Duplicate or recently stopped server instance.
2. Delayed server-state messages.
3. Director restarted while servers continued publishing state.
4. Multiple server IDs for the same partition during repeated restarts.
5. Time synchronization or timestamp ordering issue.
```

Follow-up commands:

```bash
docker ps -a --format 'table {{.Names}}\t{{.Status}}\t{{.Ports}}' | \
  egrep 'survival|overmap|deepdesert|dune-server'

ps -ef | grep DuneSandboxServ | grep -v grep

sudo ss -uapn | egrep ':7777|:7778|:7779|:7888|:7889|:7890'

docker exec <RABBITMQ_CONTAINER> rabbitmqctl list_queues \
  name messages messages_ready messages_unacknowledged consumers | \
  egrep -i 'name|login|travel|survival|overmap|deep|completion|serverState'

date -u
timedatectl
```

Support interpretation:

```text
Treat this as a state-ordering symptom unless it directly correlates with the failed travel attempt, queue buildup, or changing server IDs for the same partition.
```

---

## Active Issue: UserEngine.ini Permission Error

Observed error:

```text
PermissionError: [Errno 13] Permission denied: '<path>/Saved/UserSettings/UserEngine.ini'
```

Working explanation:

```text
The host-side script writes UserEngine.ini before startup. If the bind-mounted Saved/UserSettings directory was previously written by the container runtime user, the host automation user may not have write permission.
```

Preferred fix:

```text
Use ACL-based shared write access for both the host automation user and the container's dune:dune UID/GID.
Do not repeatedly chown the files back and forth.
Do not change the container runtime user unless confirmed necessary.
```

Confirm container UID/GID:

```bash
IMG="$(docker inspect <DUNE_SERVER_CONTAINER> --format '{{.Config.Image}}' 2>/dev/null || true)"

if [ -n "$IMG" ]; then
  docker run --rm --entrypoint sh "$IMG" -lc 'id dune || id'
fi
```

Apply ACL repair after mapping the real AMP `Saved` path:

```bash
sudo apt-get update
sudo apt-get install -y acl

HOST_USER="${SUDO_USER:-$(id -un)}"
CONTAINER_UID=1000
CONTAINER_GID=1000

for d in \
  <INSTANCE_PATH>/*/Saved \
  <INSTANCE_PATH>/*/*/Saved
do
  [ -d "$d" ] || continue
  sudo mkdir -p "$d/UserSettings"
  sudo setfacl -R -m "u:${HOST_USER}:rwx,u:${CONTAINER_UID}:rwx,g:${CONTAINER_GID}:rwx" "$d"
  sudo setfacl -R -d -m "u:${HOST_USER}:rwx,u:${CONTAINER_UID}:rwx,g:${CONTAINER_GID}:rwx" "$d"
done
```

---

## Decision Tree

```text
Players can travel to partition 1 and 8 but hang on other maps
  -> Focus on ClassicalInstancing dynamic spawn and travel lifecycle.

Director has no active server for a ClassicalInstancing group
  -> Validate autoscaling/spawn trigger and destination server availability.

Director issues travel grant but destination map never sees the client
  -> Check dynamic UDP game port exposure and tcpdump.

Destination map receives client and reaches early travel/login stages
  -> Network handoff is not the first failure. Inspect UE5 lifecycle and grace-period disconnect reason.

Destination reaches FinishSpawn then grace-period disconnects
  -> Focus on post-login session persistence, completion, cleanup, or return-handoff path.

Spawn tooling allocates game + IGW ports, but UE5 process lacks IGWPort
  -> Treat as a high-priority dynamic-spawn topology defect.

RabbitMQ queue.delete fails because player queue is in use
  -> Check stale or overlapping AMQP connections and consumers.

ServerState out-of-order appears
  -> Check duplicate processes, stale queues, server ID alignment, and time sync.
```

---

## Escalation Package Checklist

Before escalating to a developer or vendor, collect:

```text
1. The failed-travel capture directory archive.
2. Director log capture.
3. Destination map log capture.
4. UDP listener before/during/after files.
5. tcpdump output.
6. Docker ps before/after output, if available.
7. RabbitMQ queue snapshot after failure.
8. Exact UTC test times.
9. Destination map tested.
10. Redacted client-side error or behavior.
```

All files should be redacted before external sharing.

---

## Current RCA Statement

The current evidence supports a working RCA direction that the issue is tied to ClassicalInstancing travel and dynamic map lifecycle behavior. Partition `1` and partition `8` follow a simpler or more stable runtime path, while many other destinations rely on dynamic spawning and handoff.

The leading technical suspect is incomplete or mismatched dynamic spawn topology, especially whether allocated IGW/server-to-server ports are passed into each dynamically spawned UE5 process.

RabbitMQ per-player queue lifecycle errors and server-state ordering warnings should remain tracked as supporting symptoms or correlated issues, but they do not currently supersede the dynamic ClassicalInstancing handoff hypothesis.
