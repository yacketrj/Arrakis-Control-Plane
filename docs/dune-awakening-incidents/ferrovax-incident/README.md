# Ferrovax Incident

This folder contains the Ferrovax incident record.

## Status

```text
Incident status: open
Current phase: intake / initial diagnosis
Primary symptom: server appears to be running, but is not visible in the in-game server list
Resolution status: not resolved
Reusable-documentation promotion: none yet
```

## Evidence Boundary

```text
Use only Ferrovax-specific prompts, uploads, logs, commands, outputs, and notes provided after this incident label was created.
Do not import evidence, assumptions, client details, or conclusions from other incident labels.
Keep this incident separate from OIC reusable troubleshooting documentation.
Promote only generalized, sanitized lessons into OIC documentation after they are validated.
```

## Initial User-Reported Symptom

```text
The environment owner reports that the server status appears healthy or running, but the server does not appear in the in-game server list.
```

## Current Intake State

```text
Reported issue: server running/status appears good, but server is not visible in the server list
Affected workflow: server discovery / server listing / join path
Hosting platform: unknown
Runtime/orchestration layer: unknown
Control panel: unknown
Known working behavior: server reports running/healthy at the management layer
Known failing behavior: server does not appear in the in-game server list
First known failure time UTC: unknown
Recent changes: unknown
Evidence files received: none yet
Next evidence required: environment discovery, server registration/listing logs, active listener output, network path, and launch/config values
```

## Initial Working Assessment

The current symptom does not prove the root cause. A server can appear healthy in the local management layer while still failing to list publicly if registration, advertised address, region/listing configuration, port binding, firewall/NAT, or runtime launch arguments are incorrect.

## Next Evidence Required

Collect the following before changing configuration or restarting services:

```text
1. Hosting platform and runtime/orchestration discovery
2. Control panel or service status screenshot/export
3. Server startup logs covering the latest start
4. Server registration/listing-related logs
5. Active UDP/TCP listener output
6. Launch arguments and relevant config values with credentials redacted
7. Public/private IP and firewall/NAT path
8. Whether the server is visible to local/LAN players, remote players, or nobody
```

## Suggested First Runbooks

Use these OIC reusable runbooks after the platform and runtime are identified:

```text
1. Environment Discovery
2. Server Visibility and Listing
3. Port and Network Listener Validation
4. Firewall, NAT, and Cloud Networking
5. Configuration and Launch Argument Review
6. Log Collection and Redaction
```

## Privacy Handling

Do not include chat names, player names, raw account identifiers, server passwords, authentication tokens, private keys, or unrelated personal information in shared evidence packages or reusable OIC documentation.
