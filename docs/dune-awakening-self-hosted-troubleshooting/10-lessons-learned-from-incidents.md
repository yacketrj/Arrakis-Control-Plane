# Lessons Learned From Incidents

This document captures generalized lessons learned from Dune: Awakening self-hosted server incidents. It must not include client names, user names, Discord names, raw player/account IDs, passwords, tokens, IP addresses, or other PII.

Case-specific evidence belongs in the incident records area. Only reusable troubleshooting lessons should be promoted here.

---

## 1. Keep Incident Records Separate From Troubleshooting Guides

Lesson:

```text
Incident records preserve facts for one case.
Troubleshooting guides provide reusable process steps.
```

Operational guidance:

```text
Do not paste case-specific logs, names, identifiers, screenshots, or private environment details into the general troubleshooting guide.
Promote only generalized patterns, commands, interpretations, and decision logic.
Keep ITIL incident reports in a separate incident/case folder.
```

---

## 2. Do Not Treat Process-Visible OS Logs as Physical Host Proof

Lesson:

```text
A game process may report the operating system visible from its runtime context. In containerized or virtualized deployments, that may represent container userspace, guest VM context, or shared kernel details rather than the physical host or parent hypervisor.
```

Troubleshooting step:

```text
Validate host, guest, container, and hypervisor context separately.
```

Use:

```text
Linux shell: hostnamectl, /etc/os-release, uname -a, systemd-detect-virt
Docker runtime: docker exec <container> cat /etc/os-release; docker info
Hypervisor: Hyper-V or Proxmox platform checks
Cloud: provider console and instance metadata
```

---

## 3. AMP-Managed Docker Requires Path Discovery Before File Edits

Lesson:

```text
When a control panel launches Docker containers, the apparent game path, container path, and active host bind-mount path may differ.
```

Troubleshooting step:

```text
Before editing files, discover the active instance path, Saved path, log path, Docker mounts, and control-panel file paths.
```

Interpretation:

```text
Editing a copied, inactive, or non-mounted path may appear to succeed but will not affect the running server.
```

---

## 4. Instanced Travel Requires Synchronized Evidence

Lesson:

```text
Instanced travel failures cannot be diagnosed from one log stream alone. The key question is where the request stops: source, director/control plane, destination spawn, listener binding, network path, login lifecycle, or cleanup.
```

Troubleshooting step:

```text
Capture director/control-plane logs, destination logs, and listener/traffic evidence during the same failed travel attempt.
```

Use the dedicated runbook:

```text
runbooks/instanced-travel-dynamic-spawn-validation.md
```

---

## 5. Compare One Working Destination Against One Failing Destination

Lesson:

```text
Comparing many destinations at once creates noise. A controlled comparison between one known-working destination and one known-failing destination is more useful.
```

Compare:

```text
Runtime launch command
Destination/map mode
Partition or instance identifier
Game/client port
Server-to-server/gateway/IGW-style port
Bind address
Advertised address
Queue response
Travel grant response
Login/handoff token presence
Destination lifecycle stage reached
Disconnect or cleanup behavior
```

---

## 6. Validate Both Game and Server-to-Server Port Paths

Lesson:

```text
A dynamically spawned destination may require both a player-facing game/client port and a server-to-server, gateway, or IGW-style port. A process may appear partially healthy if one path is present while the other is missing.
```

Troubleshooting step:

```text
Confirm that the runtime allocates both port roles, passes both into the game process, and that the process listens on both expected ports.
```

Interpretation:

```text
Destination process starts but expected game listener is missing:
  Focus on destination startup arguments and bind address.

Destination process starts with game listener but expected server-to-server listener is missing:
  Focus on dynamic spawn argument passing and runtime/orchestration configuration.

Both listeners exist but no packets arrive:
  Focus on firewall, NAT, hypervisor networking, cloud rules, Docker publishing, or advertised address.
```

---

## 7. Queue Authentication Success Does Not Prove Travel Completion

Lesson:

```text
A messaging service may authenticate a client or service successfully while the travel flow still fails later in allocation, handoff, destination login, session state, or cleanup.
```

Troubleshooting step:

```text
Correlate queue authentication, queue state, director response, destination lifecycle logs, and client-visible result in the same UTC window.
```

---

## 8. Grace-Period Disconnects Point Later Than Initial Ingress

Lesson:

```text
If a destination reaches late login or spawn stages and later disconnects through grace-period or cleanup handling, initial network ingress was likely not the only issue.
```

Troubleshooting step:

```text
Focus next on session lifecycle, travel completion, persistence, cleanup, return-handoff behavior, and server-state ordering.
```

---

## 9. Server-State Ordering Errors Are Symptoms Until Correlated

Lesson:

```text
Out-of-order server-state reports may indicate stale messages, delayed queue events, duplicate processes, restart overlap, time-sync issues, or replacement server instances. They should be treated as symptoms until correlated with failed travel windows.
```

Troubleshooting step:

```text
Compare server IDs, partition IDs, process start times, queue timestamps, and host/container time sync before treating state ordering as root cause.
```

---

## 10. UID/GID Mismatch on Bind Mounts Should Be Solved Deliberately

Lesson:

```text
A host automation user and a container runtime user may both need write access to the same bind-mounted Saved or configuration path.
```

Troubleshooting step:

```text
Identify the host user, container/service user, file owner, group, and ACLs before changing ownership.
```

Preferred approach:

```text
Use shared ACLs where appropriate.
Avoid repeated ownership flipping.
Avoid changing the container runtime user unless the requirement is proven.
```

---

## 11. Privacy Requirement for Promoted Lessons

Before moving an incident lesson into this general guide, remove:

```text
Client names
User names
Discord names or handles
Player display names
Raw account IDs
Passwords
Tokens
Secrets
Public/private IPs unless required and approved
Exact private filesystem paths that identify a person or organization, unless replaced with placeholders
```

Use placeholders such as:

```text
<INSTANCE_PATH>
<LOG_PATH>
<SAVED_PATH>
<DIRECTOR_SERVICE>
<DESTINATION_SERVICE>
<RABBITMQ_SERVICE>
<DESTINATION_MAP>
<CLIENT_IP>
<PUBLIC_IP>
<PRIVATE_IP>
```
