# Platform Guide: Azure

Use this when the server is hosted on Microsoft Azure.

## 1. Confirm the Azure VM

Run in: Azure Portal

```text
Virtual machines > select the suspected Dune server VM.
Record subscription, resource group, region, VM name, VM size, power state, public IP, private IP, virtual network, subnet, NSG, route table, and whether a load balancer or NAT gateway is involved.
```

Run on: cloud CLI workstation with Azure CLI configured

```bash
az vm list -d -o table
az network nsg list -o table
az network public-ip list -o table
az network nic list -o table
```

## 2. Validate Azure Network Path

Run in: Azure Portal

```text
Check VM networking.
Check NSG inbound and outbound rules.
Check subnet NSG rules.
Check route tables.
Check public IP, load balancer, or NAT gateway path if used.
```

## 3. Continue on the Guest OS

After identifying the VM, log in to the guest OS and continue with:

- Linux guest: [Linux local or Linux VM](./linux-local-or-vm.md)
- Windows guest: [Windows / Hyper-V](./windows-hyper-v.md)

Azure network settings show the cloud edge only. Pair cloud checks with host listener and packet evidence.
