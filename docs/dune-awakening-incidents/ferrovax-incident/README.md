# Ferrovax Incident

This folder contains the Ferrovax incident record.

## Status

```text
Incident status: open
Current phase: initial diagnosis / network and listing validation
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

## Evidence Provided

### Local Status Screenshot

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

### Hyper-V Host Discovery Output

```text
VM name: dune-awakening
VM state: Running
CPU usage: 6
Memory assigned: 21,474,836,480 bytes, approximately 20 GiB
Virtual switch in use by VM: DuneAwakeningServerSwitch
Virtual switch type: External
Physical adapter backing switch: Realtek PCIe GBE Family Controller
VM IPv4 address: 192.168.1.125
VM IPv6 link-local address: present
Hyper-V NAT: no entries returned by Get-NetNat
Hyper-V static NAT mappings: no entries returned by Get-NetNatStaticMapping
```

### Additional Screenshot Confirmation

```text
Evidence type: screenshot of Hyper-V host PowerShell output
Assessment: screenshot confirms the same Hyper-V host discovery values already recorded above
Security note: unrelated token/credential material is visible in the background of the screenshot and is intentionally not reproduced in this incident record
Handling note: future screenshots should be cropped to the terminal output or taken with credential files closed
```

## Current Intake State

```text
Reported issue: server running/status appears good, but server is not visible in the server list
Affected workflow: server discovery / server listing / join path
Hosting platform: Hyper-V confirmed by host PowerShell output
Runtime/orchestration layer: unknown
Control panel: unknown
Guest OS: unknown
Known working behavior: local management/status screen reports healthy services
Known failing behavior: server does not appear in the in-game server list
First known failure time UTC: unknown
Recent changes: system was physically relocated near router and connected by wired Ethernet
Evidence files received: local-status screenshot; Hyper-V host networking output; Hyper-V host PowerShell screenshot
Next evidence required: guest OS/runtime discovery, server registration/listing logs, active listener output, advertised address, and router/firewall/NAT path
```

## Current Working Assessment

The Hyper-V layer now looks structurally reasonable for direct LAN connectivity: the VM is running, connected to an External virtual switch, and has a LAN IPv4 address of `192.168.1.125`. Hyper-V NAT does not appear to be in use because `Get-NetNat` and `Get-NetNatStaticMapping` returned no entries.

This shifts the next investigation step away from Hyper-V NAT and toward the guest VM, router/firewall path, server registration, advertised address, and listener state.

The current symptom remains consistent with one of the following unproven areas:

```text
1. Public listing or registration failure
2. Incorrect advertised/external address
3. Router/NAT/port-forwarding issue to 192.168.1.125
4. Host or guest firewall issue
5. Required UDP/TCP ports not listening inside the guest VM
6. Required ports listening locally but not reachable externally
7. Region/listing configuration mismatch
8. Runtime launch arguments not publishing the expected address/ports
```

## What the Evidence Proves

```text
The local status tool can see the battlegroup.
Database, gateway, and director are locally reported as healthy.
The service had been running for roughly one hour at the time of the screenshot.
The VM is running under Hyper-V.
The VM is attached to an External Hyper-V switch.
The VM has LAN IPv4 address 192.168.1.125.
Hyper-V NAT is not currently shown as the active NAT mechanism.
```

## What the Evidence Does Not Prove

```text
It does not prove the server registered successfully with the external server list.
It does not prove remote players can reach the game ports.
It does not prove the advertised public IP or region is correct.
It does not prove the router forwards required ports to 192.168.1.125.
It does not prove Windows firewall, guest firewall, or router firewall rules are correct.
It does not prove required UDP/TCP ports are listening inside the guest VM.
```

## Next Evidence Required

Collect the following before changing configuration or restarting services:

```text
1. Guest OS and runtime/orchestration discovery from inside the VM
2. Active UDP/TCP listener output from inside the VM
3. Server startup logs covering the latest start
4. Registration/listing-related log lines
5. Runtime launch arguments and relevant config values with credentials redacted
6. Public IP, private VM IP, and router/NAT/port-forward path
7. Whether the server is visible to LAN clients, remote clients, or neither
8. Windows host firewall profile and any allow/block rules related to the VM or game ports
9. Guest firewall status inside the VM
```

## Suggested Next Commands

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
Get-NetFirewallProfile | Select-Object Name, Enabled, DefaultInboundAction, DefaultOutboundAction
```

### Hyper-V host PowerShell - host firewall context

```powershell
Get-NetFirewallProfile | Select-Object Name, Enabled, DefaultInboundAction, DefaultOutboundAction
Get-NetFirewallRule -Enabled True | Where-Object { $_.DisplayName -match 'Dune|Awakening|AMP|Docker|Game|UDP|Server' } | Select-Object DisplayName, Direction, Action, Enabled
```

### Router / port-forward check

Use the router web UI.

```text
Confirm the WAN/public IP shown by the router.
Confirm the VM private IP is 192.168.1.125.
Confirm required game/listing ports are forwarded to 192.168.1.125.
Confirm forwarding protocol matches the server requirement, especially UDP where required.
Confirm no second router, ISP modem/router, CGNAT, or upstream firewall is between the internet and the Hyper-V host network.
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