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
Where is the server hosted? If unknown, write unknown.
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

## Redaction Guidance

Redact tokens, passwords, secrets, raw account IDs, and personal information before broad sharing. Do not automatically redact operational values needed for troubleshooting, such as service names, container names, map names, local paths, and instance paths, when sharing internally with trusted support staff.
