# Platform Guide: Hyper-V

Use this when the Dune: Awakening server is running inside a virtual machine managed by Microsoft Hyper-V, or when the environment owner says Hyper-V is involved.

Hyper-V is the virtualization layer. The actual Dune server usually runs inside a guest operating system, such as a Windows VM or Linux VM. Do not run Linux guest commands on the Hyper-V host unless the host itself is Linux, which is not typical for Hyper-V.

## 1. Confirm Whether You Are on the Hyper-V Host or Inside a Guest VM

Run on: Windows host PowerShell

```powershell
hostname
whoami
Get-ComputerInfo | Select-Object CsName, WindowsProductName, WindowsVersion, CsSystemType, HyperVisorPresent
Get-Service vmms
```

Interpretation:

```text
HyperVisorPresent is True and VMMS service exists/runs:
  This system may be the Hyper-V host.

The system looks like a normal Windows machine and has no VMMS service:
  You may be inside a guest VM or on a non-Hyper-V Windows host.
```

Record:

```text
Command location:
Hyper-V host confirmed: yes/no/unknown
Windows host name:
Logged-in user:
Hyper-V service state:
```

## 2. List Hyper-V VMs

Run on: Hyper-V host PowerShell

```powershell
Get-VM | Select-Object Name, State, CPUUsage, MemoryAssigned, Uptime, Generation
```

Record:

```text
Potential Dune VM name:
VM state:
CPU usage:
Memory assigned:
Uptime:
```

If no VMs are listed, confirm that you are actually on the Hyper-V host and that your account has permission to view Hyper-V.

## 3. Identify VM Networking

Run on: Hyper-V host PowerShell

```powershell
Get-VMSwitch | Select-Object Name, SwitchType, NetAdapterInterfaceDescription
Get-VMNetworkAdapter -VMName * | Select-Object VMName, Name, SwitchName, Status, MacAddress, IPAddresses
```

Record:

```text
VM name:
Virtual switch name:
Virtual switch type:
VM IP address:
MAC address:
```

Interpretation:

```text
External switch:
  VM may be directly reachable on the LAN or routed network.

Internal switch:
  VM usually needs host routing or NAT.

Private switch:
  VM is isolated from external network unless additional routing is configured.
```

## 4. Check Hyper-V NAT or Port Forwarding, If Used

Run on: Hyper-V host PowerShell

```powershell
Get-NetNat
Get-NetNatStaticMapping
Get-NetIPAddress | Sort-Object InterfaceAlias | Format-Table InterfaceAlias, IPAddress, AddressFamily
```

Record:

```text
NAT name:
Internal prefix:
Static mappings:
Host public/listening IP:
Guest private IP:
Forwarded ports:
```

Interpretation:

```text
No NAT or static mapping exists:
  Traffic must reach the VM through an external switch, routing, or another firewall/load-balancer path.

Static mapping exists:
  Confirm every required game, server-to-server, and messaging port is mapped to the correct guest IP.
```

## 5. Check Windows Defender Firewall on the Hyper-V Host

Run on: Hyper-V host PowerShell as Administrator

```powershell
Get-NetFirewallProfile | Select-Object Name, Enabled, DefaultInboundAction, DefaultOutboundAction
Get-NetFirewallRule -Enabled True | Where-Object { $_.DisplayName -match 'Dune|Awakening|Docker|AMP|Rabbit|Game' } | Select-Object DisplayName, Direction, Action, Enabled
```

If port-specific review is needed:

```powershell
Get-NetFirewallPortFilter | Where-Object { $_.Protocol -eq 'UDP' -or $_.Protocol -eq 'TCP' } | Select-Object Protocol, LocalPort
```

Record:

```text
Firewall enabled:
Relevant allow rules:
Relevant block rules:
Ports explicitly allowed:
```

## 6. Continue Inside the Guest VM

After identifying the VM that runs Dune, log into that guest OS and choose the matching platform guide.

```text
Linux guest VM -> Linux local or Linux VM guide
Windows guest VM -> Windows local or Windows VM guide
```

Do not stop at the Hyper-V host. Hyper-V only proves the virtualization/network layer. You still need guest OS process, listener, log, and runtime evidence.

## 7. Evidence to Capture for Escalation

```text
Hyper-V host name:
VM name:
VM state:
VM IP address:
Virtual switch name and type:
NAT/static mappings, if any:
Host firewall summary:
Guest OS guide used next:
Guest listener evidence collected: yes/no
Packet capture collected: yes/no
```

## 8. Common Hyper-V Findings

```text
VM is off or paused:
  Start/resume VM only after recording state.

VM has no IP address:
  Check guest network configuration and virtual switch.

VM uses an Internal or Private switch:
  Validate NAT/routing before assuming the game is unreachable.

Host firewall allows traffic but guest is not listening:
  Troubleshoot inside the guest OS/runtime.

Guest is listening but traffic never reaches guest:
  Check Hyper-V switch, NAT/static mapping, Windows firewall, upstream router, and cloud/provider firewall if applicable.
```
