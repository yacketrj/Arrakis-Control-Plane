# Dune: Awakening Self-Hosted Server Troubleshooting Guide

**Project name:** Dune: Awakening Dedicated Server Stability Review  
**Document type:** Discovery-based troubleshooting guide  
**Audience:** Entry-level support, operations, server administrators, and escalation teams

---

## Purpose

This guide helps support staff troubleshoot a Dune: Awakening self-hosted server without assuming the customer environment, symptoms, affected partitions, control panel, container runtime, file paths, port ranges, or root cause.

The user or environment owner must define the issue first. Support staff then use this guide to discover the environment, collect evidence, and determine where the failure occurs.

Workflow:

```text
1. Capture the user-defined issue statement.
2. Discover the installation and runtime environment.
3. Identify what works and what fails.
4. Capture logs, process state, listener state, and traffic during one controlled reproduction.
5. Compare evidence against known Dune: Awakening self-hosted server failure patterns.
6. Document findings, confidence level, next test, and escalation package.
```

Do not convert a symptom into a root cause until evidence supports it.

---

## How to Read Command Steps

Every command in this guide states where it should be run.

Common platform labels:

```text
Run on: Linux host or Linux VM shell
  Use when you are SSH'd into a Linux server, Ubuntu VM, Proxmox guest, or bare-metal Linux host.

Run on: Windows host PowerShell
  Use when you are logged into a Windows server or Windows Hyper-V host.

Run on: Docker host shell
  Use on the machine where Docker is installed. This may be Linux shell, Windows PowerShell, or a control-panel terminal.

Run on: inside a container
  Use only after entering a container with docker exec or an equivalent control-panel console.

Run in: control panel UI
  Use when the value must be gathered from AMP or another hosting panel instead of a terminal command.
```

If you are unsure where you are, start with the environment discovery steps and record what you find.

---

## Operational Variables

These designations make the guide reusable across installations. They are installation-specific operational values, not automatically sensitive. Support staff should replace them with the actual values discovered in the environment and record those values in case notes.

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

```text
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
Cloud resource IDs
Unredacted authentication material
```

Do not automatically redact operational values that are required for troubleshooting, such as container names, service names, local paths, destination map names, and discovered instance paths, when sharing internally with trusted support staff.

---

## Intake: User-Defined Issue Statement

Capture the issue in the user's own words before adding assumptions or suspected causes.

Ask the user or environment owner:

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

### Step 1 — Identify the Operating System You Are Logged Into

Use the Linux command set if you are connected by SSH to a Linux server, Ubuntu VM, Proxmox guest, or Linux bare-metal host.

Run on: Linux host or Linux VM shell

```bash
hostnamectl
cat /etc/os-release
uname -a
systemd-detect-virt -v || true
```

Use the Windows command set if you are logged into a Windows server or Hyper-V host.

Run on: Windows host PowerShell

```powershell
$env:COMPUTERNAME
Get-ComputerInfo | Select-Object CsName, WindowsProductName, WindowsVersion, OsHardwareAbstractionLayer, CsSystemType
systeminfo
```

Record:

```text
Hostname:
Operating system:
Kernel or Windows version:
Virtualization result:
Architecture:
```

Interpretation:

```text
Linux systemd-detect-virt returns none:
  The shell may be on bare metal, but confirm with hardware/provider evidence.

Linux systemd-detect-virt returns a virtualization type:
  The shell is inside a virtualized guest.

Windows Get-ComputerInfo shows Hyper-V or virtual machine details:
  The shell may be on a Hyper-V host or a Windows VM. Record exactly what is shown.

OS discovery identifies only the system where the command ran:
  It does not automatically prove where the Dune game process is running.
```

---

### Step 2 — Determine Whether Docker Is Present and Usable

Run this on the machine that is expected to run the Dune server or containers.

Run on: Docker host shell, Linux or Windows PowerShell

```bash
docker version
docker ps
```

If `docker version` or `docker ps` fails, record the error exactly.

For more detail on Linux:

Run on: Linux Docker host shell

```bash
docker info 2>/dev/null | grep -Ei 'Operating System|OSType|Architecture|Kernel Version|Docker Root Dir|Cgroup|Server Version' || true
docker ps --format 'table {{.Names}}\t{{.Image}}\t{{.Status}}\t{{.Ports}}'
```

For more detail on Windows:

Run on: Windows Docker host PowerShell

```powershell
docker info | Select-String -Pattern 'Operating System|OSType|Architecture|Kernel Version|Docker Root Dir|Cgroup|Server Version'
docker ps --format "table {{.Names}}`t{{.Image}}`t{{.Status}}`t{{.Ports}}"
```

Record:

```text
Docker CLI available: yes/no
Docker daemon reachable: yes/no
Docker root directory:
Active Dune-related containers:
Published ports or host network use:
```

Interpretation:

```text
docker command not found:
  Docker may not be installed, may not be in PATH, or a control panel may abstract access.

docker command exists but permission is denied:
  The current user may need sudo, group membership, administrator rights, or control-panel-provided access.

docker ps shows Dune containers:
  Use docker inspect to map mounts, ports, and runtime arguments.

docker ps shows no Dune containers:
  Check whether the server is managed directly by a control panel, systemd, Windows Service, or another process supervisor.
```

---

### Step 3 — Identify the Control Panel or Service Manager

Look for control panels, services, or scripts used to launch the server.

Run on: Linux host or Linux VM shell

```bash
ps -ef | grep -Ei 'amp|cubecoders|DuneSandbox|docker|compose|systemd' | grep -v grep
systemctl list-units --type=service 2>/dev/null | grep -Ei 'dune|amp|docker|compose' || true
find /home -maxdepth 4 -type d \( -iname '*amp*' -o -iname '*dune*' -o -iname '*awakening*' \) 2>/dev/null | head -100
```

Run on: Windows host PowerShell

```powershell
Get-Process | Where-Object { $_.ProcessName -match 'amp|docker|Dune|Awakening' } | Select-Object ProcessName, Id, Path
Get-Service | Where-Object { $_.Name -match 'amp|docker|dune' -or $_.DisplayName -match 'amp|docker|dune' } | Select-Object Name, DisplayName, Status
Get-ChildItem -Path C:\ -Directory -ErrorAction SilentlyContinue | Where-Object { $_.Name -match 'Dune|Awakening|AMP' }
```

Run in: control panel UI

```text
Open the hosting control panel.
Find the Dune: Awakening instance.
Record the instance name, install path, log path, launch method, and any container/service name shown in the UI.
```

Record:

```text
Control panel present:
Service manager:
Dune instance path:
User running the service:
Whether Docker is launched directly or through another tool:
```

---

### Step 4 — Map the Active Instance Path

Use Linux commands when the server files are on Linux.

Run on: Linux host or Linux VM shell

```bash
find /home -maxdepth 8 -type f \( -name '*Engine.ini' -o -name 'UserEngine.ini' -o -name '*.log' \) 2>/dev/null | grep -Ei 'dune|awakening|DuneSandbox' | head -100
```

Use Windows commands when the server files are on Windows.

Run on: Windows host PowerShell

```powershell
Get-ChildItem -Path C:\ -Recurse -ErrorAction SilentlyContinue -Include *Engine.ini,UserEngine.ini,*.log | Where-Object { $_.FullName -match 'Dune|Awakening|DuneSandbox' } | Select-Object -First 100 FullName
```

If Docker is available, map bind mounts.

Run on: Linux Docker host shell

```bash
docker ps --format '{{.Names}}' | while read c; do
  echo "===== $c ====="
  docker inspect "$c" --format '{{range .Mounts}}{{println .Source "->" .Destination}}{{end}}'
done
```

Run on: Windows Docker host PowerShell

```powershell
docker ps --format '{{.Names}}' | ForEach-Object {
  Write-Host "===== $_ ====="
  docker inspect $_ --format '{{range .Mounts}}{{println .Source "->" .Destination}}{{end}}'
}
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

---

### Step 5 — Document the Discovered Environment

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

Run on: Linux host or Linux VM shell

```bash
find "$INSTANCE_PATH" -type f -name '*.log' 2>/dev/null | head -100
grep -RniE 'LogInit: Log: OS:|ExecutableName|Binaries/Linux|Dreamworld platform|machine network name|user name is' "$INSTANCE_PATH" 2>/dev/null | head -100
cat /etc/os-release
uname -a
```

Run on: Windows host PowerShell

```powershell
Get-ChildItem -Path $env:INSTANCE_PATH -Recurse -Filter *.log -ErrorAction SilentlyContinue | Select-Object -First 100 FullName
Select-String -Path "$env:INSTANCE_PATH\**\*.log" -Pattern 'LogInit: Log: OS:','ExecutableName','Binaries/Linux','Dreamworld platform','machine network name','user name is' -ErrorAction SilentlyContinue | Select-Object -First 100
Get-ComputerInfo | Select-Object CsName, WindowsProductName, WindowsVersion, CsSystemType
```

If Docker is in use, compare with container-visible OS.

Run on: Docker host shell

```bash
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

Record:

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

Run on: Linux host, Linux VM, or Linux container shell where logs are stored

```bash
grep -RniE 'travel queue|Travel request|LoginRequest|Travel grant|travel completion|ServerState' "$LOG_PATH" 2>/dev/null | head -200
```

Run on: Windows host PowerShell where logs are stored

```powershell
Select-String -Path "$env:LOG_PATH\**\*.log" -Pattern 'travel queue','Travel request','LoginRequest','Travel grant','travel completion','ServerState' -ErrorAction SilentlyContinue | Select-Object -First 200
```

### Destination Lifecycle Progress

Run on: Linux host, Linux VM, or Linux container shell where logs are stored

```bash
grep -RniE 'PreLogin|VerifyFlsIdentity|VerifyFlsAuthorization|Completion|DatabaseLogin|CharacterDownload|Join|GameModeLogin|StartingNewPlayer|FlsLogin|LoadPlayerActors|FinishSpawn' "$LOG_PATH" 2>/dev/null | head -200
```

Run on: Windows host PowerShell where logs are stored

```powershell
Select-String -Path "$env:LOG_PATH\**\*.log" -Pattern 'PreLogin','VerifyFlsIdentity','VerifyFlsAuthorization','Completion','DatabaseLogin','CharacterDownload','Join','GameModeLogin','StartingNewPlayer','FlsLogin','LoadPlayerActors','FinishSpawn' -ErrorAction SilentlyContinue | Select-Object -First 200
```

### Grace-Period or Delayed Disconnect

Run on: Linux host, Linux VM, or Linux container shell where logs are stored

```bash
grep -RniE 'Grace Period|Disconnected from instanced map|Disconnect|Close|timeout|timed out' "$LOG_PATH" 2>/dev/null | head -200
```

Run on: Windows host PowerShell where logs are stored

```powershell
Select-String -Path "$env:LOG_PATH\**\*.log" -Pattern 'Grace Period','Disconnected from instanced map','Disconnect','Close','timeout','timed out' -ErrorAction SilentlyContinue | Select-Object -First 200
```

### Port Topology and Runtime Arguments

Do not assume fixed port ranges. Discover them from active configuration, command lines, and listener output.

Run on: Linux host or Linux VM shell

```bash
ps -ef | grep DuneSandbox | grep -v grep | sed -E 's/ServiceAuthToken=[^ ]+/ServiceAuthToken=<redacted>/g'
sudo ss -tulpen | grep -Ei 'Dune|777|778|779|780|781|788|789|790|791|792|31982|31983' || true
```

Run on: Windows host PowerShell

```powershell
Get-Process | Where-Object { $_.ProcessName -match 'Dune|Sandbox' } | Select-Object ProcessName, Id, Path
Get-NetUDPEndpoint | Sort-Object LocalPort | Format-Table LocalAddress, LocalPort, OwningProcess
Get-NetTCPConnection | Sort-Object LocalPort | Format-Table LocalAddress, LocalPort, State, OwningProcess
```

Look for whether each process has:

```text
Client/game UDP port
Server-to-server or IGW UDP port, if required
Expected bind address
Expected advertised external/public address
```

---

## Quick Health Check Commands

Run on: Docker host shell, Linux or Windows PowerShell

```bash
docker ps
```

Run on: Linux host or Linux VM shell

```bash
systemctl list-units --type=service 2>/dev/null | grep -Ei 'dune|amp|docker|compose' || true
ps -ef | grep DuneSandbox | grep -v grep | sed -E 's/ServiceAuthToken=[^ ]+/ServiceAuthToken=<redacted>/g'
sudo ss -uapn | grep -Ei 'Dune|777|778|779|780|781|788|789|790|791|792' || true
date -u
timedatectl
```

Run on: Windows host PowerShell

```powershell
Get-Service | Where-Object { $_.Name -match 'dune|amp|docker' -or $_.DisplayName -match 'dune|amp|docker' } | Select-Object Name, DisplayName, Status
Get-Process | Where-Object { $_.ProcessName -match 'Dune|Sandbox|Docker|AMP' } | Select-Object ProcessName, Id, Path
Get-NetUDPEndpoint | Sort-Object LocalPort | Format-Table LocalAddress, LocalPort, OwningProcess
Get-Date -Format u
w32tm /query /status
```

Run on: Docker host shell after `RABBITMQ_CONTAINER` is discovered

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

Set actual values first in your notes or terminal.

Run on: Linux host or Linux VM shell

```bash
DIRECTOR_CONTAINER=
RABBITMQ_CONTAINER=
DESTINATION_CONTAINER=
DESTINATION_MAP=
CLIENT_IP=
INSTANCE_PATH=
LOG_PATH=
```

Run on: Windows host PowerShell

```powershell
$env:DIRECTOR_CONTAINER=""
$env:RABBITMQ_CONTAINER=""
$env:DESTINATION_CONTAINER=""
$env:DESTINATION_MAP=""
$env:CLIENT_IP=""
$env:INSTANCE_PATH=""
$env:LOG_PATH=""
```

Create a capture directory.

Run on: Linux host or Linux VM shell

```bash
mkdir -p ~/dune-travel-capture-$(date -u +%Y%m%d-%H%M%S)
cd ~/dune-travel-capture-*
pwd
```

Run on: Windows host PowerShell

```powershell
$CaptureDir = "$env:USERPROFILE\dune-travel-capture-$(Get-Date -Format yyyyMMdd-HHmmss)"
New-Item -ItemType Directory -Path $CaptureDir | Out-Null
Set-Location $CaptureDir
Get-Location
```

### Terminal 1 — Capture Current Containers or Services

Run on: Docker host shell, Linux or Windows PowerShell

```bash
docker ps
```

Run on: Linux host or Linux VM shell

```bash
systemctl list-units --type=service 2>/dev/null | grep -Ei 'dune|amp|docker|compose' || true
```

Run on: Windows host PowerShell

```powershell
Get-Service | Where-Object { $_.Name -match 'dune|amp|docker' -or $_.DisplayName -match 'dune|amp|docker' } | Select-Object Name, DisplayName, Status
```

### Terminal 2 — Capture Director Logs

Run on: Docker host shell after `DIRECTOR_CONTAINER` is discovered

```bash
docker logs -f "$DIRECTOR_CONTAINER" 2>&1 | tee director-travel-capture.log
```

Run on: Linux host or Linux VM shell if logs are file-based

```bash
tail -F "$DIRECTOR_LOG_FILE" 2>&1 | tee director-travel-capture.log
```

Run on: Windows host PowerShell if logs are file-based

```powershell
Get-Content -Path $env:DIRECTOR_LOG_FILE -Wait | Tee-Object -FilePath director-travel-capture.log
```

### Terminal 3 — Capture Destination Logs

Run on: Docker host shell after `DESTINATION_CONTAINER` is discovered

```bash
docker logs -f "$DESTINATION_CONTAINER" 2>&1 | tee destination-map-capture.log
```

Run on: Linux host or Linux VM shell if logs are file-based

```bash
find "$INSTANCE_PATH" -type f -name '*.log' 2>/dev/null | tee possible-log-files.txt
tail -F "$DESTINATION_LOG_FILE" 2>&1 | tee destination-map-capture.log
```

Run on: Windows host PowerShell if logs are file-based

```powershell
Get-ChildItem -Path $env:INSTANCE_PATH -Recurse -Filter *.log -ErrorAction SilentlyContinue | Select-Object FullName | Tee-Object -FilePath possible-log-files.txt
Get-Content -Path $env:DESTINATION_LOG_FILE -Wait | Tee-Object -FilePath destination-map-capture.log
```

### Terminal 4 — Capture UDP Listeners

Run on: Linux host or Linux VM shell

```bash
sudo ss -uapn | tee udp-listeners-before.txt
while true; do
  echo "===== $(date -u --iso-8601=seconds) ====="
  sudo ss -uapn || true
  sleep 2
done | tee udp-listeners-during-travel.log
```

Run on: Windows host PowerShell

```powershell
Get-NetUDPEndpoint | Sort-Object LocalPort | Tee-Object -FilePath udp-listeners-before.txt
while ($true) {
  "===== $(Get-Date -Format u) =====" | Tee-Object -FilePath udp-listeners-during-travel.log -Append
  Get-NetUDPEndpoint | Sort-Object LocalPort | Tee-Object -FilePath udp-listeners-during-travel.log -Append
  Start-Sleep -Seconds 2
}
```

### Terminal 5 — Capture Traffic

Run on: Linux host or Linux VM shell if `tcpdump` is installed and `CLIENT_IP` is known

```bash
sudo tcpdump -ni any -vv "host ${CLIENT_IP}" | tee tcpdump-client-travel.log
```

Run on: Linux host or Linux VM shell if client IP is not known

```bash
sudo tcpdump -ni any -vv 'udp or tcp' | tee tcpdump-dune-ports-travel.log
```

Run on: Windows host PowerShell

```powershell
pktmon start --capture --pkt-size 0
```

Stop Windows packet capture after the test.

Run on: Windows host PowerShell

```powershell
pktmon stop
pktmon format PktMon.etl -o pktmon-travel-capture.txt
```

Run one controlled reproduction attempt, then collect final state.

Run on: Linux host or Linux VM shell

```bash
sudo ss -uapn | tee udp-listeners-after.txt
```

Run on: Windows host PowerShell

```powershell
Get-NetUDPEndpoint | Sort-Object LocalPort | Tee-Object -FilePath udp-listeners-after.txt
```

Run on: Docker host shell after `RABBITMQ_CONTAINER` is discovered

```bash
docker exec "$RABBITMQ_CONTAINER" rabbitmqctl list_queues name messages messages_ready messages_unacknowledged consumers 2>/dev/null | tee rabbitmq-queues-after.txt || true
```

Package the capture after applying appropriate redaction for the sharing audience.

Run on: Linux host or Linux VM shell

```bash
cd ..
tar -czf dune-travel-capture-$(date -u +%Y%m%d-%H%M%S).tar.gz dune-travel-capture-*/
ls -lh dune-travel-capture-*.tar.gz
```

Run on: Windows host PowerShell

```powershell
Compress-Archive -Path $CaptureDir -DestinationPath "$CaptureDir.zip"
Get-Item "$CaptureDir.zip"
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
8. Packet capture output.
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
