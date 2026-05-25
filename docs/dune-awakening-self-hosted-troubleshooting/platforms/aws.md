# Platform Guide: AWS

Use this when the server is hosted on Amazon Web Services.

## 1. Confirm the EC2 Instance

Run in: AWS Console

```text
EC2 > Instances > select the suspected Dune server instance.
Record instance ID, name, state, region, availability zone, instance type, public IP, private IP, VPC, subnet, security groups, and whether a load balancer or NAT is involved.
```

Run on: cloud CLI workstation with AWS CLI configured

```bash
aws ec2 describe-instances
aws ec2 describe-security-groups
```

## 2. Continue on the Guest OS

After identifying the EC2 instance, log in to the guest OS and continue with:

- Linux guest: [Linux local or Linux VM](./linux-local-or-vm.md)
- Windows guest: [Windows / Hyper-V](./windows-hyper-v.md)

Cloud network settings show the edge path only. Pair cloud checks with host listener and packet evidence.
