# Environment Discovery

Use this guide to identify the hosting platform and management/runtime layer before running focused troubleshooting steps.

Do not assume Docker, AMP, Hyper-V, Proxmox, or any cloud provider until discovered or reported by the user.

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

Record:

```text
Where command/UI was run:
Hostname or system name:
Operating system shown:
Virtualization result, if shown:
Logged-in user:
Current path or panel instance path:
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
Get-ChildItem -Path C:\ -Directory -ErrorAction SilentlyContinue | Where-Object { $_.Name -match 'Dune|Awakening|AMP|Docker' }
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

## Step 3 — Identify Cloud Provider, If Any

Skip this step if the environment owner confirms the server is not cloud-hosted.

### Cloud Provider Console

```text
Open the provider console.
Find the VM or instance running the Dune server.
Record provider, region, instance name, public IP, private IP, attached firewall/security group, subnet/VPC/VCN, and whether a load balancer or NAT is involved.
```

### OCI CLI

```bash
oci compute instance list --compartment-id <COMPARTMENT_ID> --all
oci network vnic list --compartment-id <COMPARTMENT_ID> --all
```

### AWS CLI

```bash
aws ec2 describe-instances
aws ec2 describe-security-groups
```

### Azure CLI

```bash
az vm list -d -o table
az network nsg list -o table
```

### GCP CLI

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

## Step 4 — Check Container Runtime Only If Discovered or Suspected

### Docker Host Shell, Linux or Windows PowerShell

```bash
docker version
docker ps
```

### Linux Docker Host Shell

```bash
docker info 2>/dev/null | grep -Ei 'Operating System|OSType|Architecture|Kernel Version|Docker Root Dir|Cgroup|Server Version' || true
docker ps --format 'table {{.Names}}\t{{.Image}}\t{{.Status}}\t{{.Ports}}'
```

### Windows Docker Host PowerShell

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

### Linux Docker Host Shell, Only If Docker Is Confirmed

```bash
docker ps --format '{{.Names}}' | while read c; do
  echo "===== $c ====="
  docker inspect "$c" --format '{{range .Mounts}}{{println .Source "->" .Destination}}{{end}}'
done
```

### Windows Docker Host PowerShell, Only If Docker Is Confirmed

```powershell
docker ps --format '{{.Names}}' | ForEach-Object {
  Write-Host "===== $_ ====="
  docker inspect $_ --format '{{range .Mounts}}{{println .Source "->" .Destination}}{{end}}'
}
```

Record:

```text
INSTANCE_PATH=
LOG_PATH=
SAVED_PATH=
Container mount source, if applicable:
Container mount destination, if applicable:
```

## Step 6 — Document Confirmed Environment

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
