# Dune: Awakening Self-Hosted Server Troubleshooting Guide

**Project name:** Dune: Awakening Dedicated Server Stability Review  
**Document type:** Discovery-based troubleshooting guide  
**Audience:** Support, operations, server administrators, and technical escalation staff

---

## Purpose

This is a troubleshooting guide, not an incident report. It does not assume the customer environment, affected partitions, control panel, container runtime, file paths, port ranges, symptoms, or root cause.

The user or environment owner defines the issue first. Support staff then use this guide to discover the environment, collect evidence, and determine where the failure occurs.

Workflow:

```text
1. Capture the user-defined issue statement.
2. Discover the installation and runtime environment.
3. Identify the known-working path and known-failing path.
4. Capture logs, process state, listener state, and traffic during one controlled reproduction.
5. Compare the evidence against common Dune: Awakening server failure patterns.
6. Document findings, confidence level, next test, and escalation package.
```

Do not convert a symptom into a root cause until evidence supports it.

---

## Operational Variables

The following designations make the guide reusable across installations. They are installation-specific operational values, not automatically sensitive. Support staff should replace them with the actual values discovered in the environment and record those values in the case notes.

```text
INSTANCE_PATH          Active game/control-panel instance path on the host
LOG_PATH               Directory or file path where relevant logs are stored
SAVED_PATH             Active Saved directory for the game server
CONTAINER_NAME         Generic container name
DIRECTOR_CONTAINER     Container or service handling director/control-plane logs
RABBITMQ_CONTAINER     Container or service running RabbitMQ/messaging
DESTINATION_CONTAINER  Container or service for the destination map/server
DESTINATION_MAP        Map, partition, zone, or destination being tested
CLIENT_IP              Client public IP, when needed for packet capture
PUBLIC_IP              Server public IP or advertised external address
PRIVATE_IP             Server private, bind, or interface IP
PLAYER_ID              Player identifier, if needed for queue/log correlation
PLAYER_NAME            Player display name, if needed for operational notes
```

Case-note template:

```bash
INSTANCE_PATH=
LOG_PATH=
SAVED_PATH=
CONTAINER_NAME=
DIRECTOR_CONTAINER=
RABBITMQ_CONTAINER=
DESTINATION_CONTAINER=
DESTINATION_MAP=
CLIENT_IP=
PUBLIC_IP=
PRIVATE_IP=
```

Use actual values in local troubleshooting notes. Generalize or redact them only when the sharing audience requires it.

---

## Redaction Rules

Before sharing logs, screenshots, command output, or generated reports outside the support team, redact sensitive values.

Usually redact:

```text
Personal names
Chat or Discord display names
Player names, unless approved
Raw player/account IDs, unless needed by the vendor and approved
Public or private IPs, if sharing broadly
Service tokens
Authentication secrets
RabbitMQ secrets
Database passwords
Cloud resource IDs such as OCIDs
```

Do not automatically redact operational values that are required for troubleshooting, such as container names, service names, local paths, destination map names, and discovered instance paths, when sharing internally with trusted support staff.

---

## Intake: User-Defined Issue Statement

Capture the issue in the user's own words before adding assumptions or suspected causes.

Ask:

```text
What is failing?
When did it start?
Who is affected?
What still works?
What changed recently?
How many times has it reproduced?
What does the player/client see?
What does the server/operator see?
```

Record:

```text
Reported issue:
Affected users/players:
Affected map, partition, service, or workflow, if known:
Known working map, partition, service, or workflow, if known:
Observed client-side behavior:
Observed server-side behavior:
First known failure time in UTC:
Recent change before issue started:
Reproduction status:
Impact to gameplay or operations:
```

Do not fill in affected partitions, map names, host model, or suspected cause until the user provides them or evidence confirms them.

---

## Intake: Initial Classification Worksheet

Classify only after the user-defined issue statement is captured.

```text
Incident type:
  [ ] Player login failure
  [ ] Travel failure
  [ ] Instanced map failure
  [ ] Server startup failure
  [ ] Permission/config generation failure
  [ ] Messaging/RabbitMQ issue
  [ ] Network reachability issue
  [ ] Unknown / needs discovery

Current confidence level:
  [ ] User-reported only
  [ ] Reproduced by support
  [ ] Supported by logs
  [ ] Supported by logs and packet/process evidence

Known working path:
Known failing path:
Evidence collected so far:
Evidence still needed:
```

---

## Environment Discovery Workflow

Environment details must be discovered and verified before they are treated as confirmed.

### Step 1 - Identify the OS Visible to the Shell

```bash
hostnamectl
cat /etc/os-release
uname -a
systemd-detect-virt -v || true
```

Record:

```text
Hostname:
Operating system:
Kernel:
Virtualization result:
Architecture:
```

Interpretation:

```text
systemd-detect-virt returns none:
  The shell may be on bare metal, but confirm with hardware/provider evidence.

systemd-detect-virt returns a virtualization type:
  The shell is inside a virtualized guest.

/etc/os-release and uname identify the OS/kernel visible to the shell:
  They do not automatically identify a Docker container image or parent hypervisor.
```

### Step 2 - Determine Whether Docker Is Present and Usable

```bash
docker version 2>/dev/null || true
docker info 2>/dev/null | egrep -i 'Operating System|OSType|Architecture|Kernel Version|Docker Root Dir|Cgroup|Server Version' || true
docker ps --format 'table {{.Names}}\t{{.Image}}\t{{.Status}}\t{{.Ports}}' 2>/dev/null || true
```

Record:

```text
Docker CLI available: yes/no
Docker daemon reachable: yes/no
Docker root directory:
Active Dune-related containers:
Published ports or host network use:
```

### Step 3 - Identify the Control Panel or Service Manager

Look for control panels, services, or scripts used to launch the server.

```bash
ps -ef | egrep -i 'amp|cubecoders|DuneSandbox|docker|compose|systemd' | grep -v grep
systemctl list-units --type=service 2>/dev/null | egrep -i 'dune|amp|docker|compose' || true
find /home -maxdepth 4 -type d \( -iname '*amp*' -o -iname '*dune*' -o -iname '*awakening*' \) 2>/dev/null | head -100
```

Record:

```text
Control panel present:
Service manager:
Dune instance path:
User running the service:
Whether Docker is launched directly or through another tool:
```

### Step 4 - Map the Active Instance Path

```bash
find /home -maxdepth 8 -type f \( -name '*Engine.ini' -o -name 'UserEngine.ini' -o -name '*.log' \) 2>/dev/null | egrep -i 'dune|awakening|DuneSandbox' | head -100
```

If Docker is available:

```bash
docker ps --format '{{.Names}}' 2>/dev/null | while read c; do
  echo "===== $c ====="
  docker inspect "$c" --format '{{range .Mounts}}{{println .Source "->" .Destination}}{{end}}' 2>/dev/null
 done
```

Record actual values:

```text
INSTANCE_PATH=
LOG_PATH=
SAVED_PATH=
Container mount source:
Container mount destination:
```

Do not edit a path until it is proven to be used by the active running instance.

### Step 5 - Document the Discovered Environment

Only after discovery should support staff record a confirmed environment summary.

```text
Bare-metal or parent host, if known:
Virtualization platform, if known:
Guest OS:
Control panel/service manager:
Container runtime:
CPU allocation:
Memory allocation:
Network location/provider model:
Public IP model:
Active instance path:
Active log path:
```

---

## Runtime OS Discovery and Interpretation

Game logs may report an operating system that reflects what the game process can observe. That may differ from the shell OS, container image, VM OS, or physical host OS.

```bash
find "$INSTANCE_PATH" -type f -name '*.log' 2>/dev/null | head -100

grep -RniE 'LogInit: Log: OS:|ExecutableName|Binaries/Linux|Dreamworld platform|machine network name|user name is' "$INSTANCE_PATH" 2>/dev/null | head -100

cat /etc/os-release
uname -a

docker exec "$CONTAINER_NAME" sh -lc 'cat /etc/os-release; uname -a' 2>/dev/null || true
```

Interpretation:

```text
Game log OS differs from shell OS:
  The game log may be showing container userspace or another runtime layer.

Game log OS matches container exec output:
  The log likely reflects the container image/userspace.

Kernel version is the same inside and outside the container:
  This is expected for Linux containers because containers share the host kernel.

Game log reports a CPU model:
  This is the CPU visible to the runtime, not proof of bare metal by itself.
```

---

## Evidence Collection: Establish Working and Failing Paths

The user/environment owner must define the working and failing paths first. Support should then verify them with logs.

```text
Known working login path:
Known working map/partition path:
Known failing map/partition path:
Failure time in UTC:
User action that triggered failure:
Client-visible result:
Server-visible result:
```

---

## Common Evidence Patterns to Check

Use this section as a checklist, not as a conclusion.

### Director or Control-Plane Request Handling

```bash
grep -RniE 'travel queue|Travel request|LoginRequest|Travel grant|travel completion|ServerState' "$LOG_PATH" 2>/dev/null | head -200
```

### Destination Lifecycle Progress

```bash
grep -RniE 'PreLogin|VerifyFlsIdentity|VerifyFlsAuthorization|Completion|DatabaseLogin|CharacterDownload|Join|GameModeLogin|StartingNewPlayer|FlsLogin|LoadPlayerActors|FinishSpawn' "$LOG_PATH" 2>/dev/null | head -200
```

### Grace-Period or Delayed Disconnect

```bash
grep -RniE 'Grace Period|Disconnected from instanced map|Disconnect|Close|timeout|timed out' "$LOG_PATH" 2>/dev/null | head -200
```

### Port Topology and Runtime Arguments

Do not assume fixed port ranges. Discover them from active configuration, command lines, and listener output.

```bash
ps -ef | grep DuneSandbox | grep -v grep | sed -E 's/ServiceAuthToken=[^ ]+/ServiceAuthToken=<redacted>/g'
sudo ss -tulpen | egrep 'Dune|777|778|779|780|781|788|789|790|791|792|31982|31983' || true
```

Look for whether each process has:

```text
Client/game UDP port
Server-to-server or IGW UDP port, if required
Expected bind address
Expected advertised external/public address
```

---

## First-Response Triage Checklist

```text
1. Capture the user's problem statement.
2. Run environment discovery.
3. Identify working and failing paths.
4. Preserve evidence before changing configuration.
5. Capture one controlled reproduction attempt.
```

---

## Quick Health Check Commands

```bash
docker ps --format 'table {{.Names}}\t{{.Image}}\t{{.Status}}\t{{.Ports}}' 2>/dev/null || true
systemctl list-units --type=service 2>/dev/null | egrep -i 'dune|amp|docker|compose' || true

ps -ef | grep DuneSandbox | grep -v grep | sed -E 's/ServiceAuthToken=[^ ]+/ServiceAuthToken=<redacted>/g'

sudo ss -uapn | egrep 'Dune|777|778|779|780|781|788|789|790|791|792' || true

date -u
timedatectl
```

RabbitMQ, after `RABBITMQ_CONTAINER` is discovered:

```bash
docker exec "$RABBITMQ_CONTAINER" rabbitmqctl list_queues name messages messages_ready messages_unacknowledged consumers 2>/dev/null || true
```

---

## Runbook: Capture One Failed Attempt

Capture these at the same time:

```text
1. Director/control-plane logs
2. Destination/server logs
3. Active UDP listener and traffic output
```

Set actual values first:

```bash
DIRECTOR_CONTAINER=
RABBITMQ_CONTAINER=
DESTINATION_CONTAINER=
DESTINATION_MAP=
CLIENT_IP=
INSTANCE_PATH=
LOG_PATH=
```

Create a capture directory:

```bash
mkdir -p ~/dune-travel-capture-$(date -u +%Y%m%d-%H%M%S)
cd ~/dune-travel-capture-*
pwd
```

### Terminal 1 - Capture Current Containers/Services

```bash
docker ps --format 'table {{.Names}}\t{{.Image}}\t{{.Status}}\t{{.Ports}}' 2>/dev/null | tee docker-ps-before.txt
systemctl list-units --type=service 2>/dev/null | egrep -i 'dune|amp|docker|compose' | tee services-before.txt || true
```

### Terminal 2 - Capture Director Logs

```bash
docker logs -f "$DIRECTOR_CONTAINER" 2>&1 | tee director-travel-capture.log
```

If file-based:

```bash
tail -F "$DIRECTOR_LOG_FILE" 2>&1 | tee director-travel-capture.log
```

### Terminal 3 - Capture Destination Logs

```bash
docker logs -f "$DESTINATION_CONTAINER" 2>&1 | tee destination-map-capture.log
```

If file-based:

```bash
find "$INSTANCE_PATH" -type f -name '*.log' 2>/dev/null | tee possible-log-files.txt
tail -F "$DESTINATION_LOG_FILE" 2>&1 | tee destination-map-capture.log
```

### Terminal 4 - Capture UDP Listeners

```bash
sudo ss -uapn | tee udp-listeners-before.txt

while true; do
  echo "===== $(date -u --iso-8601=seconds) ====="
  sudo ss -uapn || true
  sleep 2
done | tee udp-listeners-during-travel.log
```

### Terminal 5 - Capture Traffic

If `CLIENT_IP` is known:

```bash
sudo tcpdump -ni any -vv "host ${CLIENT_IP}" | tee tcpdump-client-travel.log
```

If the client IP is not known:

```bash
sudo tcpdump -ni any -vv 'udp or tcp' | tee tcpdump-dune-ports-travel.log
```

Run one controlled reproduction attempt, then collect final state:

```bash
docker ps --format 'table {{.Names}}\t{{.Image}}\t{{.Status}}\t{{.Ports}}' 2>/dev/null | tee docker-ps-after.txt
sudo ss -uapn | tee udp-listeners-after.txt

docker exec "$RABBITMQ_CONTAINER" rabbitmqctl list_queues name messages messages_ready messages_unacknowledged consumers 2>/dev/null | tee rabbitmq-queues-after.txt || true
```

Package the capture after applying appropriate redaction for the sharing audience.

```bash
cd ..
tar -czf dune-travel-capture-$(date -u +%Y%m%d-%H%M%S).tar.gz dune-travel-capture-*/
ls -lh dune-travel-capture-*.tar.gz
```

---

## How to Interpret the Capture

```text
Director does not log the request:
  Check source action, log source, and queue path.

Director logs request but no grant/completion:
  Focus on director queue, destination availability, and server state.

Destination/server process never appears:
  Focus on spawn trigger, control panel, Docker/runtime, and service manager.

Destination/server process appears but expected listener does not:
  Focus on startup arguments and port binding.

Destination has game port but no server-to-server/IGW port when one is required:
  Focus on dynamic spawn command and runtime arguments.

Destination reaches early lifecycle stages:
  Initial handoff worked. Continue into session lifecycle, persistence, and cleanup.

Packets arrive from client but no server replies:
  Focus on destination process handling, auth/session state, or local firewall.

Server replies but client still hangs:
  Focus on client-side routing, session state, or higher-level travel handling.
```

---

## RabbitMQ / Messaging Checks

Run only after discovering `RABBITMQ_CONTAINER`.

```bash
docker exec "$RABBITMQ_CONTAINER" rabbitmqctl list_queues -p / name durable auto_delete arguments consumers messages messages_ready messages_unacknowledged state

docker exec "$RABBITMQ_CONTAINER" rabbitmqctl list_consumers -p /

docker exec "$RABBITMQ_CONTAINER" rabbitmqctl list_connections pid user peer_host peer_port state name
```

---

## Server State Ordering Checks

```bash
grep -RniE 'Received server state out of order|ServerState|server state' "$LOG_PATH" 2>/dev/null | head -200

docker ps -a --format 'table {{.Names}}\t{{.Status}}\t{{.Ports}}' 2>/dev/null
ps -ef | grep DuneSandbox | grep -v grep
sudo ss -uapn
date -u
timedatectl
```

Out-of-order state may indicate duplicate or recently stopped server instances, delayed messages, restart sequencing, multiple server IDs for the same role, or time ordering issues.

---

## UserEngine.ini / Permission Checks

Use this section when startup scripts fail to write settings files.

Discovery steps after setting `SAVED_PATH`:

```bash
ls -ld "$SAVED_PATH" "$SAVED_PATH/UserSettings" 2>/dev/null || true
ls -l "$SAVED_PATH/UserSettings" 2>/dev/null || true
stat "$SAVED_PATH/UserSettings/UserEngine.ini" 2>/dev/null || true
```

If Docker is available, identify the container user:

```bash
IMG="$(docker inspect "$CONTAINER_NAME" --format '{{.Config.Image}}' 2>/dev/null || true)"

if [ -n "$IMG" ]; then
  docker run --rm --entrypoint sh "$IMG" -lc 'id; id dune 2>/dev/null || true'
fi
```

Preferred repair pattern after mapping the real active path:

```bash
sudo apt-get update
sudo apt-get install -y acl

HOST_USER="${SUDO_USER:-$(id -un)}"
CONTAINER_UID=<CONTAINER_UID>
CONTAINER_GID=<CONTAINER_GID>

sudo mkdir -p "$SAVED_PATH/UserSettings"
sudo setfacl -R -m "u:${HOST_USER}:rwx,u:${CONTAINER_UID}:rwx,g:${CONTAINER_GID}:rwx" "$SAVED_PATH"
sudo setfacl -R -d -m "u:${HOST_USER}:rwx,u:${CONTAINER_UID}:rwx,g:${CONTAINER_GID}:rwx" "$SAVED_PATH"
```

---

## Escalation Package Checklist

Before escalating to a developer or vendor, collect:

```text
1. User-defined problem statement.
2. Environment discovery output.
3. Known working path and known failing path.
4. One controlled reproduction capture archive.
5. Director/control-plane log capture.
6. Destination/server log capture.
7. UDP listener before/during/after files.
8. tcpdump output.
9. Docker/service state before and after, if available.
10. RabbitMQ or messaging queue snapshot, if applicable.
11. Exact UTC test times.
12. Client-side error or behavior.
```

Apply appropriate redaction before external sharing.

---

## RCA Development Worksheet

Do not write a final RCA until evidence supports it.

```text
Confirmed symptom:
Confirmed working path:
Confirmed failing path:
Evidence source 1:
Evidence source 2:
Evidence source 3:
Most likely failure layer:
Competing explanations still possible:
Next test to confirm/refute leading theory:
Temporary workaround, if any:
Final RCA:
```
