# Environment Discovery

Use this guide to identify the hosting platform and management/runtime layer before running focused troubleshooting steps.

Do not assume Docker, AMP, Hyper-V, Proxmox, OCI, AWS, Azure, GCP, or any other platform until it is discovered or reported by the user.

## Step 1 — Identify Where You Are Logged In

Run only the command set for the system you are using.

### Linux Host or Linux VM Shell

```bash
hostnamectl
cat /etc/os-release
uname -a
systemd-detect-virt -v || true
whoami
pwd
```

### Windows Host PowerShell

```powershell
$env:COMPUTERNAME
whoami
Get-Location
Get-ComputerInfo | Select-Object CsName, WindowsProductName, WindowsVersion, OsHardwareAbstractionLayer, CsSystemType
systeminfo
```

### Control Panel UI

```text
Open the hosting panel.
Record panel name, instance name, operating system shown by the panel, install path, log path, and service start/stop controls.
```

## Step 2 — Identify the Hosting Platform and Management Layer

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

Record the discovered management layer:

```text
[ ] AMP or another control panel
[ ] Linux systemd service
[ ] Windows service
[ ] Docker or another container runtime
[ ] Hyper-V VM
[ ] Proxmox VM or container
[ ] Cloud VM
[ ] Custom script/manual process
[ ] Unknown
```

## Step 3 — Identify Cloud Provider, If Any

Skip this step if the environment owner confirms the server is not cloud-hosted.

### Cloud Provider Console

```text
Open the provider console.
Find the VM or instance running the Dune server.
Record provider, region, instance name, public IP, private IP, firewall/security group, subnet/VPC/VCN, and whether a load balancer or NAT is involved.
```

### CLI Examples

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

## Step 4 — Check Containers Only If Discovered or Suspected

Run this only if the owner, control panel, process list, or service list indicates Docker, Podman, containers, or Compose may be involved.

### Docker Host Shell, Linux or Windows PowerShell

```bash
docker version
docker ps
```

## Step 5 — Map Active Instance and Log Paths

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
```

## Step 6 — Document Confirmed Environment

Only after discovery should support staff record the confirmed environment summary.

```text
Hosting model:
Control/management layer:
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
