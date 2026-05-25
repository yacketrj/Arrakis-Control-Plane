# Dune: Awakening Self-Hosted Troubleshooting

This documentation provides a discovery-first troubleshooting workflow for Dune: Awakening self-hosted server operations.

The guide is organized by hosting platform, runtime/orchestration layer, and symptom-specific runbook. Operators should identify the environment before running platform-specific commands. This prevents inaccurate assumptions, preserves evidence quality, and reduces the risk of unnecessary service disruption.

Case-specific incident evidence must remain outside this reusable guide. Only generalized procedures, decision logic, and sanitized lessons learned should be promoted here.

## Operating Model

1. Start with [Start Here](./00-start-here-troubleshooting-guide.md).
2. Complete [Intake and Evidence Handling](./01-intake-and-evidence-handling.md).
3. Use [Environment Discovery](./02-environment-discovery.md) to identify the platform and runtime.
4. Use the [Troubleshooting Decision Tree](./03-troubleshooting-decision-tree.md) to select the correct runbook.
5. Use the [Glossary and Operational Variables](./04-glossary-and-operational-variables.md) when terminology or required values are unclear.
6. Follow the applicable platform guide.
7. Follow the applicable runtime or orchestration guide.
8. Run focused evidence capture only after the environment is confirmed.
9. Package evidence with the [Escalation Package Template](./05-escalation-package-template.md) and [Standard Evidence Bundle](./06-standard-evidence-bundle.md).
10. Prepare final analysis with the [RCA Report Template](./07-rca-report-template.md) only after evidence supports the conclusion.
11. Use [Project Continuity Notes](./08-project-continuity-notes.md) when handing off or continuing work in a later session.
12. Review the [Documentation Maintenance and QA Checklist](./09-documentation-maintenance-and-qa-checklist.md) before publishing or handing off updates.
13. Promote sanitized, reusable findings through [Lessons Learned From Incidents](./10-lessons-learned-from-incidents.md).
14. Use the [Live Support Call Checklist](./11-live-support-call-checklist.md) during screen-share or real-time support sessions.

## Core Workflow Documents

- [Start Here](./00-start-here-troubleshooting-guide.md)
- [Intake and Evidence Handling](./01-intake-and-evidence-handling.md)
- [Environment Discovery](./02-environment-discovery.md)
- [Troubleshooting Decision Tree](./03-troubleshooting-decision-tree.md)
- [Glossary and Operational Variables](./04-glossary-and-operational-variables.md)
- [Escalation Package Template](./05-escalation-package-template.md)
- [Standard Evidence Bundle](./06-standard-evidence-bundle.md)
- [RCA Report Template](./07-rca-report-template.md)
- [Project Continuity Notes](./08-project-continuity-notes.md)
- [Documentation Maintenance and QA Checklist](./09-documentation-maintenance-and-qa-checklist.md)
- [Lessons Learned From Incidents](./10-lessons-learned-from-incidents.md)
- [Live Support Call Checklist](./11-live-support-call-checklist.md)

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

## Runtime and Orchestration Guides

- [AMP control panel](./runtimes/amp-control-panel.md)
- [Docker or Docker Compose](./runtimes/docker-or-compose.md)
- [Linux systemd](./runtimes/linux-systemd.md)
- [Windows service](./runtimes/windows-service.md)
- [Manual or custom script](./runtimes/manual-or-custom-script.md)

## Focused Runbooks

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

## Evidence and Privacy Standard

Use `unknown` when a value has not been confirmed. Do not infer the hosting platform, runtime layer, or root cause from partial evidence.

Do not add case-specific personal data, credentials, unrelated environment details, or incident-only conclusions to this reusable troubleshooting guide. Preserve operational values when they are required for troubleshooting, such as service names, container names, map names, partition numbers, local paths, and port numbers.