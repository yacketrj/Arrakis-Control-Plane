# Environment Discovery

Use this guide to identify the hosting platform, management layer, and runtime before performing focused troubleshooting.

Do not assume Docker, AMP, Hyper-V, Proxmox, OCI, AWS, Azure, GCP, or any other platform until it is confirmed by the environment owner, control panel, shell output, provider console, or running process evidence.

## 1. Identify the Current Access Point

Run only the command set for the system currently open in the operator's session.

### Linux Host or Linux VM Shell

```bash
hostnamectl
cat /etc/os-release
uname -a
systemd-detect-virt -v || true
whoami
pwd
```

Record:

```text
Hostname:
Operating system:
Kernel:
Virtualization detected:
Current user:
Current directory:
```

### Windows Host PowerShell

```powershell
$env:COMPUTERNAME
whoami
Get-Location
Get-ComputerInfo | Select-Object CsName, WindowsProductName, WindowsVersion, OsHardwareAbstractionLayer, CsSystemType
systeminfo
```

Record:

```text
Computer name:
Windows version:
System type:
Current user:
Current directory:
```

### Control Panel UI

```text
Open the hosting panel.
Record the panel name, instance name, instance status, operating system shown by the panel, install path, log path, and service start/stop controls.
```

## 2. Identify the Hosting Platform and Management Layer

### Linux Host or Linux VM Shell

```bash
ps -ef | grep -Ei 'amp|cubecoders|DuneSandbox|docker|containerd|podman|compose|systemd|wine' | grep -v grep
systemctl list-units --type=service 2>/dev/null | grep -Ei 'dune|awakening|amp|docker|podman|compose' || true
find /home /opt /srv -maxdepth 5 -type d \( -iname '*amp*' -o -iname '*dune*' -o -iname '*awakening*' \) 2>/dev/null | head -100
```

### Windows Host PowerShell

```powershell
Get-Process | Where-Object { $_.ProcessName -match 'amp|docker|Dune|Awakening|Sandbox' } | Select-Object ProcessName, Id, Path
Get-Service | Where-Object { $_.Name -match 'amp|docker|dune|awakening' -or $_.DisplayName -match 'amp|docker|dune|awakening' } | Select-Object Name, DisplayName, Status
```

### Hyper-V Host PowerShell

```powershell
Get-VM | Select-Object Name, State, CPUUsage, MemoryAssigned, Uptime
Get-VMNetworkAdapter -VMName * | Select-Object VMName, SwitchName, IPAddresses
```

### Proxmox Host Shell

```bash
qm list
pct list
ip addr
bridge link
```

Record the confirmed management layer:

```text
[ ] AMP or another control panel
[ ] Linux systemd service
[ ] Windows service
[ ] Docker or another container runtime
[ ] Hyper-V VM
[ ] Proxmox VM or container
[ ] Cloud VM
[ ] Custom script or manual process
[ ] Unknown
Evidence source:
```

## 3. Identify the Cloud Provider, If Applicable

Skip this step only when the environment owner confirms that the server is not cloud-hosted.

### Cloud Provider Console

```text
Open the provider console.
Locate the VM or instance running the Dune server.
Record the provider, region, instance name, public IP, private IP, firewall or security group, subnet/VPC/VCN, and whether a load balancer or NAT gateway is involved.
```

### CLI Examples

Run only the command for the provider being used.

```bash
# OCI
oci compute instance list --compartment-id <COMPARTMENT_ID> --all

# AWS
aws ec2 describe-instances

# Azure
az vm list -d -o table

# GCP
gcloud compute instances list
```

Record:

```text
Cloud provider:
Region:
Instance or VM name:
Public IP:
Private IP:
Firewall/security object:
Subnet or network:
NAT or load balancer involved:
```

## 4. Check Containers Only When Indicated

Run this section only when the owner, control panel, process list, or service list indicates Docker, Podman, containers, or Compose may be involved.

### Docker Host Shell

```bash
docker version
docker ps --format 'table {{.Names}}\t{{.Image}}\t{{.Status}}\t{{.Ports}}'
```

Record:

```text
Container runtime:
Dune-related containers:
RabbitMQ or database containers:
Network mode or published ports, if visible:
```

## 5. Map Active Instance and Log Paths

The active instance path is the path used by the running server, not merely a path that contains similar files.

### Linux Host or Linux VM Shell

```bash
find /home /opt /srv -maxdepth 8 -type f \( -name '*Engine.ini' -o -name 'UserEngine.ini' -o -name '*.log' \) 2>/dev/null | grep -Ei 'dune|awakening|DuneSandbox' | head -100
```

### Windows Host PowerShell

```powershell
Get-ChildItem -Path C:\ -Recurse -ErrorAction SilentlyContinue -Include *Engine.ini,UserEngine.ini,*.log | Where-Object { $_.FullName -match 'Dune|Awakening|DuneSandbox' } | Select-Object -First 100 FullName
```

### Control Panel UI

```text
Open the instance file manager or settings page.
Record the install path, saved/config path, and log path shown by the panel.
```

Record:

```text
INSTANCE_PATH=
LOG_PATH=
SAVED_PATH=
Evidence proving these paths are active:
```

## 6. Document the Confirmed Environment

Complete this only after discovery evidence has been collected.

```text
Hosting model:
Control or management layer:
Runtime or orchestration layer:
Operating system observed from shell:
Operating system observed from game logs:
Container OS, if applicable:
Public IP:
Private IP:
Firewall or security group path:
Active instance path:
Active log path:
Evidence proving these values:
Unknowns that remain:
```

## 7. Common Discovery Errors

```text
Mistaking container OS for physical host OS:
  Confirm each layer separately.

Running Docker commands before Docker is confirmed:
  Check processes, services, or the control panel first.

Editing a file in an inactive path:
  Confirm the path is tied to the running service before changing it.

Treating the cloud edge as the root cause before checking listeners:
  Verify process and listener state before changing firewall rules.
```