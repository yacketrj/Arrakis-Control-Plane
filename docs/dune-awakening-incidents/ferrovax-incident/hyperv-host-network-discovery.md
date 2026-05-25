# Hyper-V Host Network Discovery

## Evidence

PowerShell output was provided from the Windows host.

```text
Host name: DUNE
Host OS: Windows 10 Pro
Hyper-V host vEthernet IP for DuneAwakeningServerSwitch: 192.168.1.192
Default gateway from host: 192.168.1.1
Default Switch network: 172.28.224.0/20
DuneAwakeningServerSwitch network: 192.168.1.0/24
Host firewall profiles: disabled
```

Previously collected Hyper-V output showed the VM `dune-awakening` attached to `DuneAwakeningServerSwitch` with VM IPv4 `192.168.1.125`.

## Interpretation

This output is from the Hyper-V host, not from inside the guest VM. The presence of `vEthernet` interfaces indicates the Windows host virtual switch layer. It does not show whether the Dune server is listening inside the VM.

The host has normal outbound connectivity through `192.168.1.1`, and the Hyper-V external switch is on the same LAN as the VM. Host firewall is disabled, so host firewall is unlikely to be the immediate blocker.

## Current Assessment

The investigation should not spend more time on Hyper-V NAT or host firewall unless new evidence contradicts this. The next required evidence must come from inside the guest VM and from the router.

## Required Next Evidence

```text
1. Log into the guest VM named dune-awakening.
2. Confirm the guest OS.
3. Confirm the guest IP is 192.168.1.125 or identify the current guest IP.
4. Capture UDP/TCP listeners inside the guest VM.
5. Capture guest firewall status.
6. Confirm the server process/runtime inside the guest VM.
7. Confirm router port forwarding points to the guest VM IP, not the Hyper-V host IP.
8. Check server registration/listing logs from the active server log path.
```

## Key Branch

```text
If server ports are not listening inside the guest VM:
  Investigate server runtime, launch arguments, and startup logs.

If server ports are listening inside the guest VM:
  Investigate router forwarding, public IP, CGNAT/double NAT, advertised address, and public listing registration.
```
