# Platform Guide: OCI

Use this when the server is hosted on Oracle Cloud Infrastructure.

## 1. Confirm the OCI Instance

Run in: OCI Console

```text
Compute > Instances > select the suspected Dune server instance.
Record compartment, region, instance name, lifecycle state, shape, public IP, private IP, VCN, subnet, NSGs, security lists, route table, and attached boot/block volumes.
```

Run on: OCI CLI workstation

```bash
oci compute instance list --compartment-id <COMPARTMENT_ID> --all
oci network vnic list --compartment-id <COMPARTMENT_ID> --all
oci network security-list list --compartment-id <COMPARTMENT_ID> --all
oci network nsg list --compartment-id <COMPARTMENT_ID> --all
```

## 2. Validate Network Path

Run in: OCI Console

```text
Confirm ingress and egress rules allow the discovered game, server-to-server, and messaging ports.
Confirm route table path to the internet gateway or NAT path as appropriate.
```

## 3. Continue on the Guest OS

After identifying the OCI VM, SSH or log in to the guest OS and continue with:

- Linux guest: [Linux local or Linux VM](./linux-local-or-vm.md)
- Windows guest: [Windows / Hyper-V](./windows-hyper-v.md)

Do not conclude OCI firewall is the cause until host listener and packet evidence are collected.
