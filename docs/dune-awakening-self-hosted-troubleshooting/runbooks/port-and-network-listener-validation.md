# Runbook: Port and Network Listener Validation

Use this runbook after the hosting platform and runtime/orchestration layer are known.

Goal: prove whether the expected game, server-to-server, messaging, and management ports are actually listening on the host and reachable through the network path.

## 1. Record Expected Ports

Do not assume port ranges. Get expected ports from the user, control panel, configuration files, launch command, or vendor documentation.

Record:

```text
Game/client UDP ports:
Server-to-server or IGW UDP ports:
RabbitMQ/messaging ports, if used:
Management/API ports, if used:
Public IP:
Private/bind IP:
```

## 2. Check Active Listeners

Run on: Linux host or Linux VM shell

```bash
sudo ss -tulpen | tee listeners-all.txt
sudo ss -uapn | tee udp-listeners.txt
sudo ss -tapn | tee tcp-listeners.txt
```

Run on: Windows host PowerShell

```powershell
Get-NetUDPEndpoint | Sort-Object LocalPort | Tee-Object -FilePath udp-listeners.txt
Get-NetTCPConnection | Sort-Object LocalPort | Tee-Object -FilePath tcp-listeners.txt
```

Run on: Docker host shell, only if Docker is confirmed

```bash
docker ps --format 'table {{.Names}}\t{{.Ports}}' | tee docker-ports.txt
docker ps --format '{{.Names}}' | while read c; do
  echo "===== $c ====="
  docker inspect "$c" --format 'NetworkMode={{.HostConfig.NetworkMode}} PortBindings={{json .HostConfig.PortBindings}}'
done | tee docker-network-bindings.txt
```

## 3. Match Listener to Process

Run on: Linux host or Linux VM shell

```bash
ps -ef | grep -Ei 'DuneSandbox|Rabbit|AMP|docker' | grep -v grep | tee relevant-processes.txt
sudo ss -tulpen | grep -Ei 'Dune|rabbit|beam|docker|777|778|779|780|781|788|789|790|791|792|31982|31983' || true
```

Run on: Windows host PowerShell

```powershell
Get-Process | Where-Object { $_.ProcessName -match 'Dune|Sandbox|Rabbit|Docker|AMP' } | Select-Object ProcessName, Id, Path | Tee-Object -FilePath relevant-processes.txt
Get-NetUDPEndpoint | Sort-Object LocalPort | Format-Table LocalAddress, LocalPort, OwningProcess
```

## 4. Check Cloud or Hypervisor Edge Only If Present

Run in: cloud provider console

```text
Check inbound and outbound rules for the discovered ports.
Check whether public IP, private IP, NAT, load balancer, subnet, VPC/VCN, NSG/security group, and route table match the host listener path.
```

Run on: Hyper-V host PowerShell

```powershell
Get-VMNetworkAdapter -VMName * | Select-Object VMName, SwitchName, IPAddresses
Get-VMSwitch | Select-Object Name, SwitchType, NetAdapterInterfaceDescription
```

Run on: Proxmox host shell

```bash
ip addr
bridge link
iptables -S 2>/dev/null || true
nft list ruleset 2>/dev/null | head -200 || true
```

## 5. Interpret Results

```text
Expected port is not listening:
  Check process startup, launch arguments, bind address, and runtime configuration.

Port listens only on 127.0.0.1:
  External clients may not reach it. Check bind address and proxy/NAT design.

Port listens on private IP but not public IP:
  This may be normal in cloud/VM setups. Check NAT/security rules and advertised external address.

Cloud firewall allows traffic but host is not listening:
  The problem is on the host/runtime side, not the cloud edge.

Host is listening but no packets arrive:
  Check cloud/hypervisor firewall, network ACLs, routing, NAT, and client source path.

Packets arrive but no replies leave:
  Check application handling, auth/session state, local firewall, or process health.
```
