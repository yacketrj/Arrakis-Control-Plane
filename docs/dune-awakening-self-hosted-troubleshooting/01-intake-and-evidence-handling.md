# Intake and Evidence Handling

Use this document before running technical commands or making changes to the environment.

Effective troubleshooting starts with a clear issue statement, a defined time window, and evidence that can be reviewed later. Capture the facts first. Avoid changing services, files, permissions, firewall rules, or container settings until the current state is recorded.

## 1. Record the Issue Statement

Ask the environment owner for the issue in plain language.

```text
What is failing?
When did it start?
What still works?
What changed recently?
How often does the issue reproduce?
What does the player or client see?
What does the server operator see?
Where is the server hosted? If unknown, write unknown.
How is the server managed? Examples: AMP, Docker, Windows service, Linux service, cloud VM, Hyper-V, Proxmox, unknown.
```

## 2. Case Notes Template

```text
Reported issue:
Affected workflow:
Affected map, partition, service, or destination, if known:
Known working map, partition, service, or workflow, if known:
Observed client-side behavior:
Observed server-side behavior:
First known failure time in UTC:
Recent change before issue started:
Reproduction status:
Impact to gameplay or operations:
Known or suspected hosting platform:
Known or suspected management layer:
Known or suspected runtime layer:
```

## 3. Evidence Confidence

Use this scale when documenting how strongly the evidence supports the current working theory.

```text
Level 1 - User-reported only
Level 2 - Reproduced by support
Level 3 - Supported by logs
Level 4 - Supported by logs and process/listener evidence
Level 5 - Supported by logs, process/listener evidence, and packet capture
```

Do not present a root cause as final unless the evidence supports it. If the evidence is incomplete, document the current statement as a working hypothesis.

## 4. Evidence Handling Standard

Record evidence using UTC timestamps whenever possible.

```text
Evidence source:
Collection time UTC:
Time window covered:
Command or UI action used:
Operator who collected it:
File name or screenshot name:
What the evidence proves:
What the evidence does not prove:
```

## 5. Redaction Standard

Remove credentials and personal identifiers before sharing evidence outside the trusted support group.

Generally redact:

```text
Passwords
Tokens
Private keys
Credential material
Raw account identifiers unless vendor-required and approved
Personal names
Chat handles
Player display names
```

Do not automatically remove operational values needed for troubleshooting when sharing internally with trusted support staff. These may be required to understand the system path.

Operational values that are often needed:

```text
Service names
Container names
Map names
Partition numbers
Port numbers
Local file paths
Instance paths
```

## 6. Change Control Reminder

Before applying any fix or workaround, record:

```text
Proposed change:
Reason for change:
Evidence supporting the change:
Risk:
Backup or rollback plan:
Approver:
Expected validation result:
```

Apply one change at a time. Validate the result before making additional changes.