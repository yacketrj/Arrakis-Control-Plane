# Start Here: Troubleshooting Workflow

Use this document as the entry point for any Dune: Awakening self-hosted server issue.

The workflow is evidence-led. Confirm the environment first, collect the minimum required evidence, then select the appropriate runbook. Do not restart services, edit configuration, or change file permissions until the current state has been recorded.

## Standard Workflow

1. Record the issue in the user's own words.
2. Confirm the hosting platform and access location.
3. Confirm the runtime or orchestration layer.
4. Select the relevant platform guide.
5. Select the relevant runtime or orchestration guide.
6. Run the symptom-specific evidence-capture runbook.
7. Package the collected evidence for escalation or RCA.
8. Document the result, remaining unknowns, and next owner.

## Select the Platform Guide

Use [Environment Discovery](./02-environment-discovery.md) before selecting a platform guide.

- [AMP-controlled hosting](./platforms/amp-control-panel.md)
- [Linux local or Linux VM](./platforms/linux-local-or-vm.md)
- [Windows host or Windows VM](./platforms/windows-hyper-v.md)
- [Hyper-V](./platforms/hyper-v.md)
- [Proxmox](./platforms/proxmox.md)
- [OCI](./platforms/oci.md)
- [AWS](./platforms/aws.md)
- [Azure](./platforms/azure.md)
- [GCP](./platforms/gcp.md)

## Select the Runtime or Orchestration Guide

Choose a runtime guide only after the management layer is confirmed.

- [AMP control panel](./runtimes/amp-control-panel.md)
- [Docker or Docker Compose](./runtimes/docker-or-compose.md)
- [Linux systemd](./runtimes/linux-systemd.md)
- [Windows service](./runtimes/windows-service.md)
- [Manual or custom script](./runtimes/manual-or-custom-script.md)

## Select the Focused Runbook

- [Client-reported symptoms](./runbooks/client-reported-symptoms.md)
- [Server visibility and listing](./runbooks/server-visibility-and-listing.md)
- [Failed travel capture](./runbooks/failed-travel-capture.md)
- [Map travel and instancing failure](./runbooks/map-travel-and-instancing-failure.md)
- [Dynamic instancing and handoff validation](./runbooks/dynamic-instancing-and-handoff-validation.md)
- [Login and authentication failure](./runbooks/login-and-authentication-failure.md)
- [Server startup failure](./runbooks/server-startup-failure.md)
- [Crash, hang, and process exit analysis](./runbooks/crash-hang-and-process-exit-analysis.md)
- [Port and network listener validation](./runbooks/port-and-network-listener-validation.md)
- [Firewall, NAT, and cloud networking](./runbooks/firewall-nat-and-cloud-networking.md)
- [Resource and performance checks](./runbooks/resource-and-performance-checks.md)
- [Time sync and timestamp correlation](./runbooks/time-sync-and-timestamp-correlation.md)
- [Update, patch, and version validation](./runbooks/update-patch-and-version-validation.md)
- [Configuration and launch argument review](./runbooks/configuration-and-launch-argument-review.md)
- [Log collection and redaction](./runbooks/log-collection-and-redaction.md)
- [RabbitMQ or messaging checks](./runbooks/rabbitmq-or-messaging-checks.md)
- [Database and persistence checks](./runbooks/database-and-persistence-checks.md)
- [Permission and ownership errors](./runbooks/permission-and-ownership-errors.md)
- [Backup, restore, and change safety](./runbooks/backup-restore-and-change-safety.md)

## Escalation Standard

Escalate when the issue is reproducible and the evidence package shows:

```text
Confirmed symptom
Confirmed hosting platform
Confirmed runtime or orchestration layer
Known working path
Known failing path
Relevant logs
Process, service, or container state
Listener or network state where applicable
UTC timestamps
Outstanding unknowns
```

If the root cause is not proven, label it as a working hypothesis rather than a conclusion.