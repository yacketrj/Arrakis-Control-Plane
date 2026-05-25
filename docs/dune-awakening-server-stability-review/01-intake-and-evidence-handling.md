# Intake and Evidence Handling

Use this before running technical commands.

## Capture the User-Defined Issue

Ask the user or environment owner:

```text
What is failing?
When did it start?
Who is affected?
What still works?
What changed recently?
How many times has it reproduced?
What does the player/client see?
What does the server/operator see?
Where is the server hosted? If unknown, say unknown.
How is the server managed? Examples: AMP, Docker, Windows service, Linux service, cloud VM, Hyper-V, Proxmox, unknown.
```

## Case Notes Template

```text
Reported issue:
Affected users/players:
Affected map, partition, service, or workflow, if known:
Known working map, partition, service, or workflow, if known:
Observed client-side behavior:
Observed server-side behavior:
First known failure time in UTC:
Recent change before issue started:
Reproduction status:
Impact to gameplay or operations:
Known or suspected hosting platform:
Known or suspected management layer:
```

## Evidence Confidence

```text
[ ] User-reported only
[ ] Reproduced by support
[ ] Supported by logs
[ ] Supported by logs and process/listener evidence
[ ] Supported by logs, process/listener evidence, and packet capture
```

## Redaction Rules

Usually redact before broad sharing:

```text
Personal names
Chat or Discord display names
Player names, unless approved
Raw player/account IDs, unless vendor-required and approved
Public or private IPs, if sharing broadly
Service tokens
Authentication secrets
RabbitMQ secrets
Database passwords
Cloud resource IDs when sharing broadly
```

Do not automatically redact operational values needed for internal troubleshooting, such as service names, container names, local paths, destination map names, discovered instance paths, and cloud instance IDs.

## Escalation Package

Collect:

```text
1. User-defined issue statement.
2. Environment discovery output.
3. Known working path and known failing path.
4. Hosting platform and management layer evidence.
5. One controlled reproduction capture archive.
6. Control-plane log capture.
7. Destination/server log capture.
8. Listener before/during/after files.
9. Packet capture output.
10. Service/container state before and after, if available.
11. Messaging queue snapshot, if applicable.
12. Exact UTC test times.
13. Client-side error or behavior.
```
