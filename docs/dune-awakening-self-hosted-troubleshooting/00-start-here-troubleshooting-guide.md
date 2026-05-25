# Start Here: Troubleshooting Workflow

This is the first document to use for any Dune: Awakening self-hosted server issue.

## Entry-Level Workflow

1. Do not assume the hosting platform or runtime.
2. Capture the user-defined issue statement.
3. Run environment discovery.
4. Pick one platform guide.
5. Pick one runtime/orchestration guide.
6. Run a focused evidence-capture runbook.
7. Package evidence and escalate.

## Choose a Platform Guide

Use [Environment Discovery](./02-environment-discovery.md), then choose one:

- [Linux local or Linux VM](./platforms/linux-local-or-vm.md)
- [Windows / Hyper-V](./platforms/windows-hyper-v.md)
- [Proxmox](./platforms/proxmox.md)
- [OCI](./platforms/oci.md)
- [AWS](./platforms/aws.md)
- [Azure](./platforms/azure.md)
- [GCP](./platforms/gcp.md)

## Choose a Runtime / Orchestration Guide

Choose only after the management layer is discovered:

- [AMP control panel](./runtimes/amp-control-panel.md)
- [Docker or Docker Compose](./runtimes/docker-or-compose.md)
- [Linux systemd](./runtimes/linux-systemd.md)
- [Windows service](./runtimes/windows-service.md)
- [Manual or custom script](./runtimes/manual-or-custom-script.md)

## Then Use Focused Runbooks

- [Failed travel capture](./runbooks/failed-travel-capture.md)
- [Port and listener validation](./runbooks/port-and-network-listener-validation.md)
- [Log collection and redaction](./runbooks/log-collection-and-redaction.md)
- [RabbitMQ or messaging checks](./runbooks/rabbitmq-or-messaging-checks.md)
- [Permission and ownership errors](./runbooks/permission-and-ownership-errors.md)
