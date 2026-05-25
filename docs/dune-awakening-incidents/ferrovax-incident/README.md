# Ferrovax Incident

This folder contains the Ferrovax incident record.

## Status

```text
Incident status: open
Current phase: intake / initial diagnosis
Primary symptom: server appears to be running locally, but is not visible in the in-game server list
Resolution status: not resolved
Reusable-documentation promotion: none yet
```

## Evidence Boundary

```text
Use only Ferrovax-specific prompts, uploads, logs, commands, outputs, screenshots, and notes provided after this incident label was created.
Do not import evidence, assumptions, client details, or conclusions from other incident labels.
Keep this incident separate from OIC reusable troubleshooting documentation.
Promote only generalized, sanitized lessons into OIC documentation after they are validated.
```

## User-Reported Symptom

```text
The environment owner reports that the server status appears healthy/running locally, but the server does not appear in the in-game server list.
```

## Latest Evidence Provided

```text
Evidence type: screenshot and operator note
Observed local status: Healthy
Database status: Ready
Gateway status: Healthy
Director status: Healthy
Observed uptime: approximately 58 minutes
Game server section: visible in local status screen; screenshot does not prove public server-list registration
Hosting note: currently running on Hyper-V
Network note: system was moved near the router and connected by wired Ethernet
Battlegroup identifier: visible in screenshot but not reproduced in this report
```

## Current Intake State

```text
Reported issue: server running/status appears good, but server is not visible in the server list
Affected workflow: server discovery / server listing / join path
Hosting platform: Hyper-V reported by environment owner
Runtime/orchestration layer: unknown
Control panel: unknown
Guest OS: unknown
Known working behavior: local management/status screen reports healthy services
Known failing behavior: server does not appear in the in-game server list
First known failure time UTC: unknown
Recent changes: system was physically relocated near router and connected by wired Ethernet
Evidence files received: one local-status screenshot
Next evidence required: Hyper-V networking details, guest OS/runtime discovery, server registration/listing logs, active listener output, advertised address, and router/firewall/NAT path
```

## Current Working Assessment

The server appears healthy from the local control/status view. This does not prove that the server is registered with the public server list or reachable from outside the local environment.

The current symptom is most consistent with one of the following unproven areas:

```text
1. Public listing or registration failure
2. Incorrect advertised/external address
3. Hyper-V virtual switch or NAT mismatch
4. Router/NAT/port-forwarding issue
5. Host or guest firewall issue
6. Required UDP/TCP ports not listening or not reachable externally
7. Region/listing configuration mismatch
8. Runtime launch arguments not publishing the expected address/ports
```

## What the Screenshot Proves

```text
The local status tool can see the battlegroup.
Database, gateway, and director are locally reported as healthy.
The service had been running for roughly one hour at the time of capture.
```

## What the Screenshot Does Not Prove

```text
It does not prove the server registered successfully with the external server list.
It does not prove remote players can reach the game ports.
It does not prove the advertised public IP or region is correct.
It does not prove Hyper-V networking is using an external switch instead of NAT/internal/private networking.
It does not prove router port forwarding or firewall rules are correct.
```

## Next Evidence Required

Collect the following before changing configuration or restarting services:

```text
1. Hyper-V VM name, state, assigned IP, and virtual switch type
2. Guest OS and runtime/orchestration discovery from inside the VM
3. Server startup logs covering the latest start
4. Registration/listing-related log lines
5. Active UDP/TCP listener output from inside the VM
6. Runtime launch arguments and relevant config values with credentials redacted
7. Public IP, private VM IP, and router/NAT/port-forward path
8. Whether the server is visible to LAN clients, remote clients, or neither
```

## Suggested First Commands

### Hyper-V host PowerShell

```powershell
Get-VM | Select-Object Name, State, CPUUsage, MemoryAssigned, Uptime
Get-VMSwitch | Select-Object Name, SwitchType, NetAdapterInterfaceDescription
Get-VMNetworkAdapter -VMName * | Select-Object VMName, SwitchName, Status, IPAddresses
Get-NetNat
Get-NetNatStaticMapping
```

### Inside the guest VM - Linux shell, if Linux

```bash
hostnamectl
cat /etc/os-release
ip addr
ip route
sudo ss -tulpen
ps -ef | grep -Ei 'DuneSandbox|Awakening|Sandbox|docker|amp|rabbit' | grep -v grep || true
command -v docker && docker ps --format 'table {{.Names}}\t{{.Image}}\t{{.Status}}\t{{.Ports}}' || true
```

### Inside the guest VM - Windows PowerShell, if Windows

```powershell
Get-ComputerInfo | Select-Object CsName, WindowsProductName, WindowsVersion, CsSystemType
Get-NetIPAddress | Sort-Object InterfaceAlias | Format-Table InterfaceAlias, IPAddress, AddressFamily
Get-NetRoute | Sort-Object DestinationPrefix | Format-Table DestinationPrefix, NextHop, InterfaceAlias
Get-NetUDPEndpoint | Sort-Object LocalPort
Get-NetTCPConnection | Sort-Object LocalPort
```

### Server listing / registration log search

Run where the server logs are stored.

```bash
grep -RniE 'register|registration|listing|server list|public|external|advertis|region|gateway|FLS|login|error|warning' <LOG_PATH> 2>/dev/null | head -300
```

PowerShell equivalent:

```powershell
Select-String -Path "<LOG_PATH>\**\*.log" -Pattern 'register','registration','listing','server list','public','external','advertis','region','gateway','FLS','login','error','warning' -ErrorAction SilentlyContinue | Select-Object -First 300
```

## Suggested OIC Runbooks

Use these reusable runbooks after platform and runtime are identified:

```text
1. Environment Discovery
2. Server Visibility and Listing
3. Port and Network Listener Validation
4. Firewall, NAT, and Cloud Networking
5. Hyper-V Platform Guide
6. Configuration and Launch Argument Review
7. Log Collection and Redaction
```

## Privacy Handling

Do not include chat names, player names, raw account identifiers, server passwords, authentication tokens, private keys, or unrelated personal information in shared evidence packages or reusable OIC documentation.
