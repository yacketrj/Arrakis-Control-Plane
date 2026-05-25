# Start Here: Dune: Awakening Self-Hosted Server Troubleshooting

This is the entry point for troubleshooting. Do not assume the hosting platform, orchestration tool, control panel, runtime, or root cause.

## Support Workflow

1. Capture the issue in the user's own words.
2. Identify the hosting platform.
3. Identify the management/runtime layer.
4. Choose the matching platform and runtime runbooks.
5. Capture one controlled reproduction attempt.
6. Package and redact evidence for escalation.

## Pick the Correct Path

Use the discovery guide first: [Environment Discovery](./02-environment-discovery.md).

Then choose one platform guide:

- [Linux local or Linux VM](./platforms/linux-local-or-vm.md)
- [Windows / Hyper-V](./platforms/windows-hyper-v.md)
- [Proxmox](./platforms/proxmox.md)
- [OCI](./platforms/oci.md)
- [AWS](./platforms/aws.md)
- [Azure](./platforms/azure.md)
- [GCP](./platforms/gcp.md)

Then choose one runtime/orchestration guide:

- [AMP control panel](./runtimes/amp-control-panel.md)
- [Docker or Docker Compose](./runtimes/docker-or-compose.md)
- [Linux systemd](./runtimes/linux-systemd.md)
- [Windows service](./runtimes/windows-service.md)
- [Manual or custom script](./runtimes/manual-or-custom-script.md)

Use focused runbooks only after the environment is identified:

- [Failed travel capture](./runbooks/failed-travel-capture.md)
- [Port and listener validation](./runbooks/port-and-network-listener-validation.md)
- [Log collection and redaction](./runbooks/log-collection-and-redaction.md)
- [RabbitMQ or messaging checks](./runbooks/rabbitmq-or-messaging-checks.md)
- [Permission and ownership errors](./runbooks/permission-and-ownership-errors.md)

## Operational Variables

Record discovered values in the case notes. These are not automatically sensitive; they are needed for troubleshooting.

```text
CLOUD_PROVIDER=
CLOUD_INSTANCE_ID=
INSTANCE_PATH=
LOG_PATH=
SAVED_PATH=
SERVICE_NAME=
CONTAINER_NAME=
DIRECTOR_SERVICE=
RABBITMQ_SERVICE=
DESTINATION_SERVICE=
DESTINATION_MAP=
CLIENT_IP=
PUBLIC_IP=
PRIVATE_IP=
```

## Do Not Guess

If a value is unknown, write `unknown` and use the discovery steps. Do not fill in values from a different customer, prior incident, or earlier troubleshooting thread.
