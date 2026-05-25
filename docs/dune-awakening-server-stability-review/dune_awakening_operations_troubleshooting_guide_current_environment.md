# Dune: Awakening Self-Hosted Server Troubleshooting Guide

**Project name:** Dune: Awakening Dedicated Server Stability Review  
**Document type:** Discovery-based troubleshooting guide  
**Audience:** Entry-level support, operations, server administrators, and escalation teams

---

## Purpose

This guide helps support staff troubleshoot a Dune: Awakening self-hosted server without assuming the customer environment, symptoms, affected partitions, control panel, container runtime, file paths, port ranges, or root cause.

The user or environment owner defines the issue first. Support staff then discover the environment and choose the correct command set for that environment.

Workflow:

```text
1. Capture the user-defined issue statement.
2. Discover the platform: local Windows, local Linux, Hyper-V, Proxmox, cloud VM, AMP, Docker, systemd, or another service model.
3. Identify what works and what fails.
4. Capture logs, process state, listener state, and traffic during one controlled reproduction.
5. Compare evidence against known Dune: Awakening self-hosted server failure patterns.
6. Document findings, confidence level, next test, and escalation package.
```

Do not convert a symptom into a root cause until evidence supports it.

---

## How to Read Command Steps

Every command in this guide includes where it should be run. Do not run every command. Pick the commands for the platform discovered in the environment.

Common labels:

```text
Run on: Linux host or Linux VM shell
  Use when connected by SSH or terminal to Linux, Ubuntu, Debian, Proxmox guest, or bare-metal Linux.

Run on: Windows host PowerShell
  Use when logged into Windows Server, Windows desktop hosting the server, or a Windows Hyper-V host.

Run on: Hyper-V host PowerShell
  Use only on the Windows machine running Hyper-V, not inside a Linux or Windows guest VM.

Run on: Proxmox host shell
  Use only on the Proxmox host itself, not inside the guest VM.

Run on: cloud CLI workstation
  Use when the admin has OCI, AWS, Azure, or GCP CLI access configured. This may be the admin laptop, a jump box, or cloud shell.

Run in: cloud provider console
  Use the web portal for OCI, AWS, Azure, or GCP when CLI access is not available.

Run in: control panel UI
  Use AMP or another hosting panel when the server is managed through a web interface.

Run on: Docker host shell
  Use only after Docker is confirmed. This may be Linux shell, Windows PowerShell, or a control-panel terminal.

Run inside: container shell
  Use only after entering a container with docker exec or an equivalent control-panel console.
```

If a step says Docker, AMP, Hyper-V, Proxmox, OCI, AWS, Azure, or GCP and that platform has not been discovered in the environment, skip that command and continue with the matching platform path.

---

## Operational Variables

These designations make the guide reusable across installations. They are installation-specific operational values, not automatically sensitive. Replace them with the actual values discovered in the environment and record those values in case notes.

```text
INSTANCE_PATH          Active game/control-panel instance path on the host
LOG_PATH               Directory or file path where relevant logs are stored
SAVED_PATH             Active Saved directory for the game server
SERVICE_NAME           systemd service, Windows service, AMP instance, or other service name
CONTAINER_NAME         Generic container name, if containers are used
DIRECTOR_SERVICE       Service/container/process handling director or control-plane logs
RABBITMQ_SERVICE       Service/container/process running RabbitMQ or messaging, if present
DESTINATION_SERVICE    Service/container/process for the destination map/server
DESTINATION_MAP        Map, partition, zone, or destination being tested
CLIENT_IP              Client public IP, when needed for packet capture
PUBLIC_IP              Server public IP or advertised external address
PRIVATE_IP             Server private, bind, or interface IP
PLAYER_ID              Player identifier, if needed for queue/log correlation
PLAYER_NAME            Player display name, if needed for operational notes
CLOUD_PROVIDER         OCI, AWS, Azure, GCP, other, or not applicable
CLOUD_INSTANCE_ID      Cloud VM/instance identifier, if hosted in cloud
```

Case-note template:

```text
CLOUD_PROVIDER=
CLOUD_INSTANCE_ID=
INSTANCE_PATH=
LOG_PATH=
SAVED_PATH=
SERVICE_NAME=
CONTAINER_NAME=
DIRECTOR_SERVICE=
RABBITMQ_SERVICE=
DESTINATION_SERVICE=
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
Cloud resource IDs when sharing broadly
Unredacted authentication material
```

Do not automatically redact operational values required for troubleshooting, such as service names, container names, local paths, destination map names, discovered instance paths, or cloud instance IDs when sharing internally with trusted support staff.

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
Where is the server hosted? If unknown, say unknown.
How is the server managed? Examples: AMP, Docker, Windows service, Linux service, cloud VM, Hyper-V, Proxmox, unknown.
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
Known or suspected hosting platform:
Known or suspected management layer:
```

Do not fill in affected partitions, map names, host model, or suspected cause until the user provides them or evidence confirms them.

---

## Environment Discovery Workflow

Environment details must be discovered and verified before they are treated as confirmed. Start broad, then choose the matching platform path.

### Step 1 — Identify Where You Are Logged In

Run the command set for the system you are actually using.

Run on: Linux host or Linux VM shell

```bash
hostnamectl
cat /etc/os-release
uname -a
systemd-detect-virt -v || true
whoami
pwd
```

Run on: Windows host PowerShell

```powershell
$env:COMPUTERNAME
whoami
Get-Location
Get-ComputerInfo | Select-Object CsName, WindowsProductName, WindowsVersion, OsHardwareAbstractionLayer, CsSystemType
systeminfo
```

Run in: control panel UI

```text
Open the hosting panel.
Record the panel name, server/instance name, operating system shown by the panel, install path, log path, and service start/stop controls.
```

Record:

```text
Where command/UI was run:
Hostname or system name:
Operating system shown:
Virtualization result, if shown:
Logged-in user:
Current path or panel instance path:
```

Interpretation:

```text
The OS shown here identifies only the system where the command or UI was run.
It does not automatically prove where the Dune server process is running.
A control panel may show a different path than the host shell.
A container may show a different OS than the VM or host.
```

---

### Step 2 — Identify the Hosting Platform and Management Layer

Do not assume Docker. First determine whether the server is managed by a control panel, Windows service, Linux service, cloud VM, hypervisor, container runtime, or custom script.

Run in: control panel UI

```text
Check whether the Dune server is managed by AMP or another panel.
Record:
- Panel name
- Instance name
- Install path
- Log path
- Startup command or launch method, if visible
- Whether the panel mentions Docker, container, service, or process mode
```

Run on: Linux host or Linux VM shell

```bash
ps -ef | grep -Ei 'amp|cubecoders|DuneSandbox|docker|containerd|podman|compose|systemd|wine' | grep -v grep
systemctl list-units --type=service 2>/dev/null | grep -Ei 'dune|awakening|amp|docker|podman|compose' || true
find /home /opt /srv -maxdepth 5 -type d \( -iname '*amp*' -o -iname '*dune*' -o -iname '*awakening*' \) 2>/dev/null | head -100
```

Run on: Windows host PowerShell

```powershell
Get-Process | Where-Object { $_.ProcessName -match 'amp|docker|Dune|Awakening|Sandbox' } | Select-Object ProcessName, Id, Path
Get-Service | Where-Object { $_.Name -match 'amp|docker|dune|awakening' -or $_.DisplayName -match 'amp|docker|dune|awakening' } | Select-Object Name, DisplayName, Status
Get-ChildItem -Path C:\ -Directory -ErrorAction SilentlyContinue | Where-Object { $_.Name -match 'Dune|Awakening|AMP|Docker' }
```

Run on: Hyper-V host PowerShell

```powershell
Get-VM | Select-Object Name, State, CPUUsage, MemoryAssigned, Uptime
Get-VMNetworkAdapter -VMName * | Select-Object VMName, SwitchName, IPAddresses
```

Run on: Proxmox host shell

```bash
qm list
pct list
ip addr
bridge link
```

Record:

```text
Management layer discovered:
  [ ] AMP or another control panel
  [ ] Linux systemd service
  [ ] Windows service
  [ ] Docker or another container runtime
  [ ] Hyper-V VM
  [ ] Proxmox VM or container
  [ ] Cloud VM
  [ ] Custom script/manual process
  [ ] Unknown

Evidence used to identify management layer:
```

---

### Step 3 — Identify Cloud Provider, If Any

Skip this step if the environment owner confirms the server is not cloud-hosted.

Run in: cloud provider console

```text
Open the provider console.
Find the VM or instance running the Dune server.
Record provider, region, instance name, public IP, private IP, attached firewall/security group, subnet/VPC/VCN, and whether a load balancer or NAT is involved.
```

Run on: cloud CLI workstation with OCI CLI configured

```bash
oci compute instance list --compartment-id <COMPARTMENT_ID> --all
oci network vnic list --compartment-id <COMPARTMENT_ID> --all
```

Run on: cloud CLI workstation with AWS CLI configured

```bash
aws ec2 describe-instances
aws ec2 describe-security-groups
```

Run on: cloud CLI workstation with Azure CLI configured

```bash
az vm list -d -o table
az network nsg list -o table
```

Run on: cloud CLI workstation with GCP CLI configured

```bash
gcloud compute instances list
gcloud compute firewall-rules list
```

Record:

```text
CLOUD_PROVIDER=
CLOUD_INSTANCE_ID=
Region/zone:
Public IP:
Private IP:
Firewall/security group/NSG/rule set:
Subnet/VPC/VCN:
Load balancer/NAT path:
```

Interpretation:

```text
Cloud metadata and firewall rules only describe the cloud edge.
They do not prove the game process is listening on the expected port.
Always pair cloud checks with host listener checks.
```

---

### Step 4 — Check Container Runtime Only If Containers Are Discovered or Suspected

Run this step only if the control panel, process list, service list, or environment owner indicates Docker, Podman, containers, or Compose may be involved.

Run on: Docker host shell, Linux or Windows PowerShell

```bash
docker version
docker ps
```

Run on: Linux Docker host shell

```bash
docker info 2>/dev/null | grep -Ei 'Operating System|OSType|Architecture|Kernel Version|Docker Root Dir|Cgroup|Server Version' || true
docker ps --format 'table {{.Names}}\t{{.Image}}\t{{.Status}}\t{{.Ports}}'
```

Run on: Windows Docker host PowerShell

```powershell
docker info | Select-String -Pattern 'Operating System|OSType|Architecture|Kernel Version|Docker Root Dir|Cgroup|Server Version'
docker ps --format "table {{.Names}}`t{{.Image}}`t{{.Status}}`t{{.Ports}}"
```

Record:

```text
Container runtime present:
Docker CLI available:
Docker daemon reachable:
Docker root directory:
Dune-related containers:
Published ports or host network use:
```

Interpretation:

```text
No Docker found:
  Continue with the discovered service/control panel/cloud/hypervisor path. Do not stop troubleshooting.

Docker found but no Dune containers:
  Dune may be managed by a service, panel, or another runtime.

Docker found with Dune containers:
  Use docker inspect to map mounts, ports, networks, and runtime arguments.
```

---

### Step 5 — Map the Active Instance Path and Log Path

Choose the command set matching the discovered platform.

Run on: Linux host or Linux VM shell

```bash
find /home /opt /srv -maxdepth 8 -type f \( -name '*Engine.ini' -o -name 'UserEngine.ini' -o -name '*.log' \) 2>/dev/null | grep -Ei 'dune|awakening|DuneSandbox' | head -100
```

Run on: Windows host PowerShell

```powershell
Get-ChildItem -Path C:\ -Recurse -ErrorAction SilentlyContinue -Include *Engine.ini,UserEngine.ini,*.log | Where-Object { $_.FullName -match 'Dune|Awakening|DuneSandbox' } | Select-Object -First 100 FullName
```

Run in: control panel UI

```text
Open the instance file manager or settings page.
Record the install path, saved/config path, and log path shown by the panel.
```

Run on: Linux Docker host shell, only if Docker is confirmed

```bash
docker ps --format '{{.Names}}' | while read c; do
  echo "===== $c ====="
  docker inspect "$c" --format '{{range .Mounts}}{{println .Source "->" .Destination}}{{end}}'
done
```

Run on: Windows Docker host PowerShell, only if Docker is confirmed

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
Container mount source, if applicable:
Container mount destination, if applicable:
```

Do not edit a path until it is proven to be used by the active running instance.

---

### Step 6 — Document the Discovered Environment

Only after discovery should support staff record the confirmed environment summary.

```text
Hosting model:
  [ ] Local Windows
  [ ] Local Linux
  [ ] Hyper-V
  [ ] Proxmox
  [ ] OCI
  [ ] AWS
  [ ] Azure
  [ ] GCP
  [ ] Other
  [ ] Unknown

Control/management layer:
  [ ] AMP/control panel
  [ ] Docker/container runtime
  [ ] Linux systemd service
  [ ] Windows service
  [ ] Custom script/manual process
  [ ] Unknown

Operating system observed from shell:
Operating system observed from game logs:
Container OS, if applicable:
Public IP:
Private IP:
Firewall/security group path:
Active instance path:
Active log path:
Evidence proving these values:
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

Run on: Docker host shell, only if containers are confirmed

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

## Quick Health Check Commands by Discovered Platform

Run only the section matching the discovered platform.

### Linux service or Linux VM path

Run on: Linux host or Linux VM shell

```bash
systemctl list-units --type=service 2>/dev/null | grep -Ei 'dune|amp|docker|compose' || true
ps -ef | grep DuneSandbox | grep -v grep | sed -E 's/ServiceAuthToken=[^ ]+/ServiceAuthToken=<redacted>/g'
sudo ss -uapn | grep -Ei 'Dune|777|778|779|780|781|788|789|790|791|792' || true
date -u
timedatectl
```

### Windows service or Hyper-V guest path

Run on: Windows host PowerShell

```powershell
Get-Service | Where-Object { $_.Name -match 'dune|amp|docker' -or $_.DisplayName -match 'dune|amp|docker' } | Select-Object Name, DisplayName, Status
Get-Process | Where-Object { $_.ProcessName -match 'Dune|Sandbox|Docker|AMP' } | Select-Object ProcessName, Id, Path
Get-NetUDPEndpoint | Sort-Object LocalPort | Format-Table LocalAddress, LocalPort, OwningProcess
Get-Date -Format u
w32tm /query /status
```

### Docker/container path

Run on: Docker host shell, only if Docker is confirmed

```bash
docker ps
docker exec "$RABBITMQ_SERVICE" rabbitmqctl list_queues name messages messages_ready messages_unacknowledged consumers 2>/dev/null || true
```

### Control panel path

Run in: control panel UI

```text
Check service status.
Check latest console output.
Open file manager or log viewer.
Record install path, log path, start command, visible ports, and restart history.
```

### Cloud path

Run in: cloud provider console

```text
Check instance state, public IP, private IP, firewall/security rules, subnet/VPC/VCN, and recent restart or maintenance events.
```

---

## Runbook: Capture One Failed Attempt

Capture these at the same time:

```text
1. Director/control-plane logs
2. Destination/server logs
3. Active listener and traffic output
```

Set actual values first in your notes or terminal.

Run on: Linux host or Linux VM shell

```bash
DIRECTOR_SERVICE=
RABBITMQ_SERVICE=
DESTINATION_SERVICE=
DESTINATION_MAP=
CLIENT_IP=
INSTANCE_PATH=
LOG_PATH=
```

Run on: Windows host PowerShell

```powershell
$env:DIRECTOR_SERVICE=""
$env:RABBITMQ_SERVICE=""
$env:DESTINATION_SERVICE=""
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

### Terminal 1 — Capture Current Services or Containers

Run on: Linux host or Linux VM shell

```bash
systemctl list-units --type=service 2>/dev/null | grep -Ei 'dune|amp|docker|compose' | tee services-before.txt || true
ps -ef | grep DuneSandbox | grep -v grep | tee processes-before.txt || true
```

Run on: Windows host PowerShell

```powershell
Get-Service | Where-Object { $_.Name -match 'dune|amp|docker' -or $_.DisplayName -match 'dune|amp|docker' } | Tee-Object -FilePath services-before.txt
Get-Process | Where-Object { $_.ProcessName -match 'Dune|Sandbox|Docker|AMP' } | Tee-Object -FilePath processes-before.txt
```

Run on: Docker host shell, only if Docker is confirmed

```bash
docker ps | tee docker-ps-before.txt
```

### Terminal 2 — Capture Director or Control-Plane Logs

Run on: Docker host shell, only if the director is a container

```bash
docker logs -f "$DIRECTOR_SERVICE" 2>&1 | tee director-travel-capture.log
```

Run on: Linux host or Linux VM shell if logs are file-based

```bash
tail -F "$DIRECTOR_LOG_FILE" 2>&1 | tee director-travel-capture.log
```

Run on: Windows host PowerShell if logs are file-based

```powershell
Get-Content -Path $env:DIRECTOR_LOG_FILE -Wait | Tee-Object -FilePath director-travel-capture.log
```

Run in: control panel UI

```text
Open the console or log viewer for the director/control-plane service.
Export or copy log output covering the test window.
```

### Terminal 3 — Capture Destination Logs

Run on: Docker host shell, only if the destination is a container

```bash
docker logs -f "$DESTINATION_SERVICE" 2>&1 | tee destination-map-capture.log
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

Run in: control panel UI

```text
Open the log viewer for the destination/server instance.
Export or copy log output covering the test window.
```

### Terminal 4 — Capture Active Listeners

Run on: Linux host or Linux VM shell

```bash
sudo ss -uapn | tee udp-listeners-before.txt
while true; do
  echo "===== $(date -u --iso-8601=seconds) ====="
  sudo ss -uapn || true
  sleep 2
done | tee udp-listeners-during-test.log
```

Run on: Windows host PowerShell

```powershell
Get-NetUDPEndpoint | Sort-Object LocalPort | Tee-Object -FilePath udp-listeners-before.txt
while ($true) {
  "===== $(Get-Date -Format u) =====" | Tee-Object -FilePath udp-listeners-during-test.log -Append
  Get-NetUDPEndpoint | Sort-Object LocalPort | Tee-Object -FilePath udp-listeners-during-test.log -Append
  Start-Sleep -Seconds 2
}
```

### Terminal 5 — Capture Traffic

Run on: Linux host or Linux VM shell if tcpdump is installed and CLIENT_IP is known

```bash
sudo tcpdump -ni any -vv "host ${CLIENT_IP}" | tee tcpdump-client-test.log
```

Run on: Linux host or Linux VM shell if CLIENT_IP is unknown

```bash
sudo tcpdump -ni any -vv 'udp or tcp' | tee tcpdump-server-test.log
```

Run on: Windows host PowerShell

```powershell
pktmon start --capture --pkt-size 0
```

Stop Windows packet capture after the test.

Run on: Windows host PowerShell

```powershell
pktmon stop
pktmon format PktMon.etl -o pktmon-test-capture.txt
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

Run on: Docker host shell, only if RabbitMQ runs in Docker

```bash
docker exec "$RABBITMQ_SERVICE" rabbitmqctl list_queues name messages messages_ready messages_unacknowledged consumers 2>/dev/null | tee rabbitmq-queues-after.txt || true
```

Package the capture.

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
Control plane does not log the request:
  Check source action, log source, and queue path.

Control plane logs request but no grant/completion:
  Focus on control-plane queue, destination availability, and server state.

Destination/server process never appears:
  Focus on spawn trigger, control panel, runtime, service manager, or cloud/host resources.

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
4. Hosting platform and management layer evidence.
5. One controlled reproduction capture archive.
6. Control-plane log capture.
7. Destination/server log capture.
8. Listener before/during/after files.
9. Packet capture output.
10. Service/container state before and after, if available.
11. Messaging queue snapshot, if applicable.
12. Exact UTC test times.
13. Client-side error or behavior.
```

Apply appropriate redaction before external sharing.

---

## RCA Development Worksheet

Do not write a final RCA until evidence supports it.

```text
Confirmed symptom:
Confirmed hosting platform:
Confirmed management layer:
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
