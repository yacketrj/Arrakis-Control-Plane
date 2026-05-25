# Host Firewall Discovery

## Evidence

PowerShell was run on the Hyper-V host.

```text
Get-NetFirewallProfile result:
- Domain profile: disabled
- Private profile: disabled
- Public profile: disabled

Filtered firewall rules returned Dune-related entries, including:
- Dune: Awakening Public Test Client
- DuneSandboxWin64Shipping
```

## Interpretation

The Hyper-V host firewall is unlikely to be the immediate blocker because all host firewall profiles are disabled. This does not validate the guest VM firewall, router port forwarding, advertised address, server registration, or listener state inside the VM.

## Next Step

Continue inside the guest VM and collect:

```text
1. Guest OS
2. Active UDP/TCP listeners
3. Guest firewall status
4. Runtime/process/container state
5. Server registration/listing logs
6. Router forwarding to VM IP 192.168.1.125
```
