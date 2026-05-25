# Platform Guide: GCP

Use this when the server is hosted on Google Cloud Platform.

## 1. Confirm the Compute Engine Instance

Run in: GCP Console

```text
Compute Engine > VM instances > select the suspected Dune server VM.
Record project, zone, instance name, machine type, status, public IP, private IP, VPC network, subnet, firewall rules, tags, service account, and whether a load balancer or NAT is involved.
```

Run on: cloud CLI workstation with gcloud configured

```bash
gcloud compute instances list
gcloud compute firewall-rules list
gcloud compute networks list
gcloud compute routes list
```

## 2. Continue on the Guest OS

After identifying the VM, log in to the guest OS and continue with:

- Linux guest: [Linux local or Linux VM](./linux-local-or-vm.md)
- Windows guest: [Windows / Hyper-V](./windows-hyper-v.md)

Cloud network settings show the edge path only. Pair cloud checks with host listener and packet evidence.
