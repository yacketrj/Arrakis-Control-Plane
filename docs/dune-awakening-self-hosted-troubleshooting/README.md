# Dune: Awakening Self-Hosted Troubleshooting

This folder contains a discovery-first troubleshooting workflow for Dune: Awakening self-hosted servers.

The workflow is split by **hosting platform** and **runtime/orchestration layer** so support staff do not run Docker, AMP, Hyper-V, cloud, or Linux commands unless that environment has actually been discovered.

## Workflow

1. Start with [Start Here](./00-start-here-troubleshooting-guide.md).
2. Complete [Intake and Evidence Handling](./01-intake-and-evidence-handling.md).
3. Use [Environment Discovery](./02-environment-discovery.md) to identify the platform and runtime.
4. Pick the matching platform guide.
5. Pick the matching runtime/orchestration guide.
6. Use focused runbooks only after the environment is identified.
7. Package evidence for escalation.

## Platform Guides

- [AMP-controlled hosting](./platforms/amp-control-panel.md)
- [Linux local or Linux VM](./platforms/linux-local-or-vm.md)
- [Windows host or Windows VM](./platforms/windows-hyper-v.md)
- [Hyper-V](./platforms/hyper-v.md)
- [Proxmox](./platforms/proxmox.md)
- [OCI](./platforms/oci.md)
- [AWS](./platforms/aws.md)
- [Azure](./platforms/azure.md)
- [GCP](./platforms/gcp.md)

## Runtime / Orchestration Guides

- [AMP control panel](./runtimes/amp-control-panel.md)
- [Docker or Docker Compose](./runtimes/docker-or-compose.md)
- [Linux systemd](./runtimes/linux-systemd.md)
- [Windows service](./runtimes/windows-service.md)
- [Manual or custom script](./runtimes/manual-or-custom-script.md)

## Focused Runbooks

- [Failed travel capture](./runbooks/failed-travel-capture.md)
- [Server startup failure](./runbooks/server-startup-failure.md)
- [Port and network listener validation](./runbooks/port-and-network-listener-validation.md)
- [Log collection and redaction](./runbooks/log-collection-and-redaction.md)
- [RabbitMQ or messaging checks](./runbooks/rabbitmq-or-messaging-checks.md)
- [Permission and ownership errors](./runbooks/permission-and-ownership-errors.md)

## Rule for Support Staff

Do not assume the environment. If a value is unknown, write `unknown`, then use discovery steps to find it.
