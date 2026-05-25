# Platform Guide: Proxmox

Use this when the environment owner says the server runs inside a Proxmox VM or LXC, or when Proxmox is suspected.

## 1. Confirm Whether You Are on the Proxmox Host or Inside a Guest

Run on: Proxmox host shell

```bash
hostnamectl
pveversion
qm list
pct list
ip addr
bridge link
```

Run on: Linux guest VM shell

```bash
hostnamectl
cat /etc/os-release
uname -a
systemd-detect-virt -v || true
ip addr
```

Record:

```text
Are you on Proxmox host or guest VM:
Proxmox version, if on host:
Guest VM ID/name:
Guest IP:
Bridge/switch:
```

## 2. Identify the Guest Running the Server

Run on: Proxmox host shell

```bash
qm list
qm config <VMID>
qm guest cmd <VMID> network-get-interfaces 2>/dev/null || true
```

If the server runs inside an LXC container:

```bash
pct list
pct config <CTID>
pct exec <CTID> -- hostnamectl
```

## 3. Continue Inside the Guest

After finding the guest VM or LXC, run the platform guide that matches the guest OS:

- Linux guest: [Linux local or Linux VM](./linux-local-or-vm.md)
- Windows guest: [Windows / Hyper-V](./windows-hyper-v.md)

## 4. Check Proxmox-Level Network Only After Guest Listener Checks

Run on: Proxmox host shell

```bash
ip addr
bridge link
iptables -S 2>/dev/null || true
nft list ruleset 2>/dev/null | head -200 || true
```

Do not conclude Proxmox is the root cause until guest listener and packet evidence are collected.
