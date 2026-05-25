# Troubleshooting Decision Tree

Use this after completing intake and environment discovery.

This document helps entry-level support decide which guide to use next. Do not assume the root cause from the symptom alone.

## 1. Start With the User-Defined Symptom

```text
Is the issue a startup failure?
  -> Use Server Startup Failure runbook.

Is the issue a player login failure?
  -> Use Login and Authentication Failure runbook.

Is the issue a travel, map transfer, destination hang, dungeon/story instance failure, or instanced map disconnect?
  -> Use Map Travel and Instancing Failure runbook.
  -> If the issue affects dynamic/instanced destinations specifically, also use Instanced Travel Dynamic Spawn Validation.

Is the issue a port, firewall, NAT, or connectivity concern?
  -> Use Port and Network Listener Validation runbook.

Is the issue a permission denied, file write, save, config, or ownership error?
  -> Use Permission and Ownership Errors runbook.

Is the issue a RabbitMQ, queue, AMQP, or messaging error?
  -> Use RabbitMQ or Messaging Checks runbook.

Is the issue unclear?
  -> Return to Intake and Environment Discovery.
```

## 2. Confirm the Hosting Platform First

```text
If hosted through AMP:
  -> Use AMP-controlled hosting platform guide and AMP runtime guide.

If hosted on Linux local/bare metal/VM:
  -> Use Linux local or Linux VM platform guide.

If hosted on Windows local/VM:
  -> Use Windows host or Windows VM platform guide.

If hosted on Hyper-V:
  -> Use Hyper-V platform guide, then continue inside the guest VM.

If hosted on Proxmox:
  -> Use Proxmox platform guide, then continue inside the guest VM or LXC.

If hosted in OCI/AWS/Azure/GCP:
  -> Use the matching cloud platform guide, then continue inside the guest OS.
```

## 3. Confirm the Runtime or Orchestration Layer

```text
If AMP manages the server:
  -> Use AMP control panel runtime guide.

If Docker or Docker Compose is discovered:
  -> Use Docker or Docker Compose runtime guide.

If Linux systemd service is discovered:
  -> Use Linux systemd runtime guide.

If Windows service is discovered:
  -> Use Windows service runtime guide.

If a shell script, batch file, scheduled task, or manual command starts the server:
  -> Use Manual or Custom Script runtime guide.
```

## 4. Choose Evidence Capture Based on Failure Type

```text
Startup failure:
  Required evidence:
    - Service/container/process status
    - Startup logs
    - Launch command or service config
    - Permission and port checks

Login failure:
  Required evidence:
    - User-defined login symptom
    - Control-plane logs
    - Auth-related logs
    - Client-visible error
    - Network/listener state

Travel or destination hang:
  Required evidence:
    - One controlled reproduction
    - Source and destination logs
    - Control-plane/director logs
    - Listener before/during/after
    - Packet capture if possible

Dynamic or instanced destination failure:
  Required evidence:
    - One known-working destination and one known-failing destination
    - Director/control-plane queue response
    - Destination map lifecycle logs
    - Dynamic game/client port listener
    - Dynamic server-to-server/gateway/IGW-style listener, if applicable
    - Runtime launch arguments with secrets redacted
    - Packet capture if possible

Network/port issue:
  Required evidence:
    - Expected ports
    - Actual listeners
    - Firewall/security group/NAT path
    - Packet capture

Permission issue:
  Required evidence:
    - Exact path and error
    - File owner/ACL
    - Service/container user
    - Startup logs
```

## 5. Stop Conditions

Stop and escalate if:

```text
The same controlled test has been captured with logs, listeners, and packet output.
A vendor-owned binary or closed-source component appears to be failing internally.
A required secret or token is missing and support does not have authority to replace it.
The environment owner cannot provide access to the platform required for the next step.
Restarting or changing configuration would risk data loss without approval.
```

## 6. RCA Rule

Do not write a final RCA until the evidence proves:

```text
Confirmed symptom
Confirmed hosting platform
Confirmed runtime/orchestration layer
Confirmed working path
Confirmed failing path
Evidence showing where the request stops or fails
Competing explanations that were ruled out
```
