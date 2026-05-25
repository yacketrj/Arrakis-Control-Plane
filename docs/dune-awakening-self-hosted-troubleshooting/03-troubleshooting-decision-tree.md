# Troubleshooting Decision Tree

Use this document after intake and environment discovery are complete.

The decision tree directs the operator to the correct evidence path. It does not determine root cause by itself. A symptom may have several possible causes, so each branch requires supporting logs, runtime state, and network or listener evidence where applicable.

## 1. Classify the Primary Symptom

```text
Startup failure:
  Use Server Startup Failure.

Login or authentication failure:
  Use Login and Authentication Failure.

Travel, map transfer, dungeon/story instance, destination hang, or instanced map disconnect:
  Use Map Travel and Instancing Failure.
  If the failing target is dynamically spawned or instanced, also use Dynamic Instancing and Handoff Validation.

Port, firewall, NAT, routing, or connectivity concern:
  Use Port and Network Listener Validation.
  Use Firewall, NAT, and Cloud Networking when the issue involves cloud, hypervisor, router, or edge firewall paths.

Permission denied, file write, save, configuration, or ownership error:
  Use Permission and Ownership Errors.

RabbitMQ, queue, AMQP, or messaging error:
  Use RabbitMQ or Messaging Checks.

Crash, hang, unexpected restart, or process exit:
  Use Crash, Hang, and Process Exit Analysis.

Resource pressure, lag, stalls, or capacity concern:
  Use Resource and Performance Checks.

Unclear symptom:
  Return to Intake and Evidence Handling, then Environment Discovery.
```

## 2. Confirm the Hosting Platform

Select one platform path only after it is confirmed by the operator, the control panel, shell output, or provider console.

```text
AMP or another control panel:
  Use AMP-Controlled Hosting.

Linux local, bare metal, or Linux VM:
  Use Linux Local or Linux VM.

Windows host or Windows VM:
  Use Windows Host or Windows VM.

Hyper-V:
  Use Hyper-V first, then continue inside the guest VM.

Proxmox:
  Use Proxmox first, then continue inside the guest VM or container.

OCI, AWS, Azure, or GCP:
  Use the matching cloud provider guide, then continue inside the guest OS.
```

## 3. Confirm the Runtime or Orchestration Layer

Choose the runtime guide after determining how the server process is started and managed.

```text
AMP manages the instance:
  Use AMP Control Panel.

Docker or Docker Compose is present:
  Use Docker or Docker Compose.

Linux systemd service is present:
  Use Linux systemd.

Windows service is present:
  Use Windows Service.

Custom script, scheduled task, batch file, shell script, or manual command:
  Use Manual or Custom Script.
```

## 4. Evidence Required by Failure Type

### Startup Failure

```text
Required evidence:
- Service, container, or process status
- Startup logs
- Launch command or service configuration
- File permissions for active instance paths
- Port/listener state before and after startup
```

### Login Failure

```text
Required evidence:
- User-defined login symptom
- Control-plane or login service logs
- Destination or starting-map logs
- Client-visible error or behavior
- Network and listener state
```

### Travel or Destination Hang

```text
Required evidence:
- One controlled reproduction
- Known working source and destination
- Known failing source and destination
- Source map logs
- Destination map logs
- Control-plane or director logs
- Listener state before, during, and after the attempt
- Packet capture when available
```

### Dynamic or Instanced Destination Failure

```text
Required evidence:
- One known-working destination and one known-failing destination
- Director/control-plane travel response
- Destination map lifecycle logs
- Dynamic game/client port listener
- Dynamic server-to-server, gateway, or IGW-style listener, if applicable
- Runtime launch arguments with credentials redacted
- Packet capture when available
```

### Network or Port Issue

```text
Required evidence:
- Expected ports from configuration or launch arguments
- Actual listeners on the host or VM
- Cloud, host, hypervisor, router, or firewall path
- NAT or port-forwarding configuration, if used
- Packet capture when available
```

### Permission or Ownership Issue

```text
Required evidence:
- Exact error message
- Exact failed path
- File owner and permissions or ACL
- Service, process, or container user
- Startup or write-attempt logs
```

## 5. Stop and Escalate When Required

Stop the troubleshooting session and escalate when any of the following conditions apply:

```text
The same controlled test has been captured with logs, listeners, and packet output.
The next step would change persistent data without a backup or approval.
A restart would disrupt active players without authorization.
A required credential is missing and support is not authorized to replace it.
The environment owner cannot provide access to the required platform or runtime.
The evidence points to closed-source or vendor-owned behavior that cannot be validated locally.
```

## 6. RCA Standard

Do not document a final root cause until the evidence establishes:

```text
Confirmed symptom
Confirmed hosting platform
Confirmed runtime or orchestration layer
Known working path
Known failing path
Evidence showing the first failed or missing step
Competing explanations considered and ruled out
Remaining unknowns, if any
```

If the evidence is incomplete, document the finding as a working hypothesis.