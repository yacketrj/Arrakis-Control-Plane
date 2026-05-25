# Runbook: Firewall, NAT, and Cloud Networking

Use this runbook when players cannot reach the server, only some ports work, travel/login fails across network boundaries, or packet captures show traffic missing in one direction.

Start only after the hosting platform and runtime/orchestration layer are identified.

## 1. Record the Expected Network Path

Ask the environment owner to describe how traffic reaches the server.

```text
Client internet
  -> Cloud firewall / provider edge, if used
  -> Load balancer or NAT, if used
  -> Hypervisor bridge or virtual switch, if used
  -> Guest OS firewall
  -> Container port publishing or host network, if used
  -> Dune server process listener
```

Record:

```text
PUBLIC_IP=
PRIVATE_IP=
CLIENT_IP=
Cloud provider, if any:
Hypervisor, if any:
NAT or port forward, if any:
Load balancer, if any:
Host firewall enabled: yes/no/unknown
Container runtime: yes/no/unknown
```

## 2. Confirm the Server Is Listening Locally

Run on: Linux host or Linux VM shell

```bash
sudo ss -tulpen | tee firewall-listeners-linux.txt
```

Run on: Windows host PowerShell

```powershell
Get-NetUDPEndpoint | Sort-Object LocalPort | Tee-Object -FilePath firewall-udp-listeners-windows.txt
Get-NetTCPConnection | Sort-Object LocalPort | Tee-Object -FilePath firewall-tcp-listeners-windows.txt
```

If the server is not listening locally, stop network-edge troubleshooting and use the startup or port/listener runbook.

## 3. Check Linux Host Firewall

Run on: Linux host or Linux VM shell

```bash
sudo ufw status verbose 2>/dev/null || true
sudo iptables -S 2>/dev/null || true
sudo nft list ruleset 2>/dev/null | head -300 || true
```

Record:

```text
UFW enabled:
iptables rules found:
nftables rules found:
Relevant allow/block rules:
```

## 4. Check Windows Defender Firewall

Run on: Windows host PowerShell as Administrator

```powershell
Get-NetFirewallProfile | Select-Object Name, Enabled, DefaultInboundAction, DefaultOutboundAction
Get-NetFirewallRule -Enabled True | Where-Object { $_.DisplayName -match 'Dune|Awakening|Docker|AMP|Rabbit|Game' } | Select-Object DisplayName, Direction, Action, Enabled
Get-NetFirewallPortFilter | Where-Object { $_.Protocol -eq 'UDP' -or $_.Protocol -eq 'TCP' } | Select-Object Protocol, LocalPort
```

Record:

```text
Firewall profile enabled:
Relevant allow rules:
Relevant block rules:
Ports explicitly allowed:
```

## 5. Check Docker Port Publishing Only If Docker Is Confirmed

Run on: Docker host shell

```bash
docker ps --format 'table {{.Names}}\t{{.Ports}}'
docker ps --format '{{.Names}}' | while read c; do
  echo "===== $c ====="
  docker inspect "$c" --format 'NetworkMode={{.HostConfig.NetworkMode}} PortBindings={{json .HostConfig.PortBindings}}'
done | tee docker-networking.txt
```

Interpretation:

```text
Host network mode:
  The container uses host listeners directly.

Bridge mode with port bindings:
  Confirm each required port is published to the expected host IP and port.

Bridge mode without required port bindings:
  External clients may not reach that container port.
```

## 6. Check Hypervisor Network Only If Present

Run on: Hyper-V host PowerShell

```powershell
Get-VMSwitch | Select-Object Name, SwitchType, NetAdapterInterfaceDescription
Get-VMNetworkAdapter -VMName * | Select-Object VMName, SwitchName, Status, IPAddresses
Get-NetNat
Get-NetNatStaticMapping
```

Run on: Proxmox host shell

```bash
ip addr
bridge link
iptables -S 2>/dev/null || true
nft list ruleset 2>/dev/null | head -300 || true
```

## 7. Check Cloud Edge Only If Cloud-Hosted

Run in: cloud provider console

```text
Check inbound rules.
Check outbound rules.
Check subnet ACLs or NSGs.
Check route table.
Check NAT gateway, internet gateway, public IP, and load balancer path.
Confirm rules cover the actual discovered ports, not assumed ports.
```

Cloud CLI examples:

```bash
# OCI
oci network security-list list --compartment-id <COMPARTMENT_ID> --all
oci network nsg list --compartment-id <COMPARTMENT_ID> --all

# AWS
aws ec2 describe-security-groups
aws ec2 describe-network-acls

# Azure
az network nsg list -o table
az network route-table list -o table

# GCP
gcloud compute firewall-rules list
gcloud compute routes list
```

## 8. Capture Traffic

Run on: Linux host or Linux VM shell when CLIENT_IP is known

```bash
sudo tcpdump -ni any -vv "host ${CLIENT_IP}" | tee firewall-client-traffic.log
```

Run on: Linux host or Linux VM shell when CLIENT_IP is unknown

```bash
sudo tcpdump -ni any -vv 'udp or tcp' | tee firewall-all-traffic.log
```

Run on: Windows host PowerShell

```powershell
pktmon start --capture --pkt-size 0
```

Stop after test:

```powershell
pktmon stop
pktmon format PktMon.etl -o firewall-pktmon-capture.txt
```

## 9. Interpret Results

```text
Cloud allows traffic but host listener is missing:
  Host/runtime problem.

Host listener exists but no packets arrive:
  Cloud, firewall, NAT, hypervisor, route, or client source path problem.

Packets arrive but no replies leave:
  Application, local firewall, routing, or process handling problem.

Replies leave but client still fails:
  Check client-side routing, NAT return path, advertised public IP, and session/application state.
```
