# Ferrovax Incident - Current Status and Operator Discovery Plan v3

Prepared for Dune: Awakening self-hosted server support.

Date: 2026-05-25

Purpose: provide a clear, non-technical evidence collection plan for the current server-listing issue. This document focuses on what has been confirmed, what remains unknown, and exactly where each discovery step must be performed.

> Do not paste passwords, tokens, private keys, or account identifiers into shared notes. Crop screenshots to the terminal or control panel only.

## 1. Current Incident Status

| Field | Current value |
|---|---|
| Incident status | Open |
| Primary symptom | Server appears healthy/running locally but does not appear in the in-game server list. |
| Current phase | Initial diagnosis: network, Kubernetes exposure, and server-list registration validation. |
| Resolution | Not resolved; final RCA not proven. |
| Most important open question | Are the Dune-facing services exposed from Kubernetes to the VM node IP `192.168.1.125` and mapped correctly through the router? |

## 2. Confirmed Evidence So Far

| Layer | Confirmed evidence |
|---|---|
| Hyper-V host | Windows 10 Pro host named `DUNE`. Host IP on the external Hyper-V switch is `192.168.1.192`. |
| Hyper-V VM | VM named `dune-awakening` is running and attached to External switch `DuneAwakeningServerSwitch`. |
| VM network | Guest VM IP is `192.168.1.125/24`. Default gateway is `192.168.1.1`. |
| Host firewall | Windows host firewall profiles are disabled. This reduces likelihood of host firewall blocking but does not prove external reachability. |
| Guest OS | Alpine Linux v3.23. |
| Guest orchestration | Kubernetes-style networking is present: `flannel.1`, `cni0`, pod network `10.42.1.0/24`, and `KUBE-*` iptables chains. |
| Runtime clue | RabbitMQ/Erlang processes are visible. Docker CLI did not return container output from the current shell. |
| Missing tool | `ss` is not installed in the Alpine guest; use `netstat` or `/proc` fallback commands. |

## 3. Current Working Assessment

The Hyper-V and Windows host layers no longer appear to be the primary suspect. The VM is on the expected LAN and connected through an External Hyper-V switch. The next diagnostic layer is inside the Alpine VM and the Kubernetes exposure path.

A local status screen can show database, gateway, and director as healthy even when the public server list cannot reach the correct externally exposed game/listing ports. The Kubernetes pod network may be healthy internally while still not publishing the required ports to the VM node IP or router.

## 4. What We Still Need to Prove

- Which Kubernetes services represent the Dune gateway, director, game server, RabbitMQ, and related components.
- Whether those services are ClusterIP-only, NodePort, LoadBalancer, hostPort, or hostNetwork.
- Which ports are actually listening on the VM node address `192.168.1.125`.
- Whether router port forwards target `192.168.1.125`, not the Windows host IP `192.168.1.192`.
- Whether server registration/listing logs show success, rejection, timeout, bad region, bad advertised address, or auth/token errors.
- Whether a client on the same LAN can see the server and whether a remote client can see it.

## 5. Have Kubernetes and Server Logs Already Been Captured?

Short answer for this incident: not yet.

What has been captured so far:

- Evidence that Kubernetes-style networking exists on the VM.
- Evidence that `flannel.1`, `cni0`, pod-network addressing, and Kubernetes iptables chains are present.
- Evidence that RabbitMQ/Erlang processes are running.
- Evidence that the `ss` command is missing.
- Evidence that Docker CLI output was not available from the current shell.

What has not been captured yet:

- `kubectl get pods -A -o wide` output.
- `kubectl get svc -A -o wide` output.
- `kubectl get endpoints -A` output.
- `kubectl get ingress -A` output.
- `crictl ps` or equivalent container runtime output.
- Server-list registration logs.
- Dune gateway/director/game-server logs from the active pod or service.

Operator rule:

- If Kubernetes service output has not been collected, run Step 2, Step 3, or Step 4 below.
- If server registration logs have not been collected, run Step 8 below.
- If those outputs were already collected later, do not repeat them. Reference the saved files from `/tmp/ferrovax-capture` and continue with Step 9 or Step 10.

## 6. Recommended Command Set - Operator Instructions

Use the prompt or screen title to decide where you are. Each step below states the required location. Do not run Linux commands in Windows PowerShell. Do not run Windows PowerShell commands inside the Alpine VM.

| Where you are | How to recognize it |
|---|---|
| Windows Hyper-V host | Prompt starts with `PS C:\WINDOWS\system32>` or another Windows PowerShell path. This is the physical Windows machine running Hyper-V. |
| Alpine guest VM | Prompt looks like `duneawakening:~#` or another Linux shell prompt. This is the VM where the Dune environment is running. |
| Router web UI | A browser page for the home/business router, usually reached through `192.168.1.1` or the router vendor app. |
| External test device | A phone hotspot, remote PC, or machine not connected to the same LAN/Wi-Fi as the server. |

### Step 1 - Confirm the VM-side network state

Run this inside the Alpine guest VM. This confirms the VM IP, default route, and network interfaces. It should show `eth0` with `192.168.1.125` and default route via `192.168.1.1`.

```sh
# Location: Alpine guest VM shell
cat /etc/os-release
ip addr
ip route
```

Send back: the full output. Expected: `eth0` has `192.168.1.125/24` and default route uses `192.168.1.1`.

Status: already collected for this incident. Re-run only if the VM was rebooted, moved, or reconfigured.

### Step 2 - Check whether Kubernetes tools are available

Run this inside the Alpine guest VM. This tells us whether `kubectl`, `k3s`, `crictl`, or containerd tools can inspect the running services.

```sh
# Location: Alpine guest VM shell
command -v kubectl || true
command -v k3s || true
command -v crictl || true
command -v ctr || true
command -v nerdctl || true
```

Send back: the output. If the command prints a path, that tool exists. If it prints nothing for a tool, that tool may not be installed or not in `PATH`.

Status: not yet captured in the evidence provided for this incident.

### Step 3 - List Kubernetes objects if kubectl works

Run this inside the Alpine guest VM only if Step 2 showed `kubectl` is available. This identifies pods, services, service types, ports, and endpoints.

```sh
# Location: Alpine guest VM shell
OUT=/tmp/ferrovax-capture
mkdir -p "$OUT"
kubectl get nodes -o wide | tee "$OUT/kube-nodes.txt"
kubectl get namespaces | tee "$OUT/kube-namespaces.txt"
kubectl get pods -A -o wide | tee "$OUT/kube-pods.txt"
kubectl get svc -A -o wide | tee "$OUT/kube-services.txt"
kubectl get endpoints -A | tee "$OUT/kube-endpoints.txt"
kubectl get ingress -A | tee "$OUT/kube-ingress.txt"
```

Send back: the files created under `/tmp/ferrovax-capture` or paste the output. Most important file: `kube-services.txt`.

If this was already collected, reference these files instead of re-running:

- `/tmp/ferrovax-capture/kube-pods.txt`
- `/tmp/ferrovax-capture/kube-services.txt`
- `/tmp/ferrovax-capture/kube-endpoints.txt`
- `/tmp/ferrovax-capture/kube-ingress.txt`

Status: not yet captured in the evidence provided for this incident.

### Step 4 - Use k3s kubectl if kubectl alone does not work

Run this inside the Alpine guest VM only if `kubectl` is not available but `k3s` exists. Some lightweight Kubernetes installs use `k3s kubectl` instead of `kubectl`.

```sh
# Location: Alpine guest VM shell
OUT=/tmp/ferrovax-capture
mkdir -p "$OUT"
k3s kubectl get nodes -o wide | tee "$OUT/kube-nodes.txt"
k3s kubectl get namespaces | tee "$OUT/kube-namespaces.txt"
k3s kubectl get pods -A -o wide | tee "$OUT/kube-pods.txt"
k3s kubectl get svc -A -o wide | tee "$OUT/kube-services.txt"
k3s kubectl get endpoints -A | tee "$OUT/kube-endpoints.txt"
k3s kubectl get ingress -A | tee "$OUT/kube-ingress.txt"
```

Send back: the files created under `/tmp/ferrovax-capture` or paste the output. If these commands fail, send the exact error message.

If this was already collected, reference the same files listed under Step 3.

Status: not yet captured in the evidence provided for this incident.

### Step 5 - Capture listener state without ss

Run this inside the Alpine guest VM. The `ss` command is missing, so use `netstat` first. If `netstat` is also missing, collect `/proc` listener data.

```sh
# Location: Alpine guest VM shell
OUT=/tmp/ferrovax-capture
mkdir -p "$OUT"
netstat -tulpen 2>/dev/null | tee "$OUT/netstat-listeners.txt" || \
netstat -tuln 2>/dev/null | tee "$OUT/netstat-listeners.txt" || true

cat /proc/net/udp  | tee "$OUT/proc-net-udp.txt"
cat /proc/net/tcp  | tee "$OUT/proc-net-tcp.txt"
cat /proc/net/udp6 | tee "$OUT/proc-net-udp6.txt"
cat /proc/net/tcp6 | tee "$OUT/proc-net-tcp6.txt"
```

Send back: `netstat-listeners.txt` if it exists, plus the `/proc` files if `netstat` did not show useful listener details.

Status: not yet captured in the evidence provided for this incident.

### Step 6 - Identify runtime processes

Run this inside the Alpine guest VM. This identifies whether Dune, Kubernetes, containerd, k3s, RabbitMQ, or other runtime components are active.

```sh
# Location: Alpine guest VM shell
OUT=/tmp/ferrovax-capture
mkdir -p "$OUT"
PAT='DuneSandbox|Awakening|Sandbox|kube|containerd|k3s|rabbit|flannel'
ps -ef | grep -Ei "$PAT" | grep -v grep | tee "$OUT/runtime-processes.txt"
```

Send back: `runtime-processes.txt`. If no Dune process appears here, that may be normal if the game server runs inside pods.

Status: partially captured. RabbitMQ/Erlang process output was provided, but Kubernetes/container runtime process output with this broader pattern is still needed.

### Step 7 - Capture pod/container runtime state if crictl exists

Run this inside the Alpine guest VM only if Step 2 showed `crictl` is available. This is useful when Docker is not installed but Kubernetes uses containerd underneath.

```sh
# Location: Alpine guest VM shell
OUT=/tmp/ferrovax-capture
mkdir -p "$OUT"
crictl ps -a | tee "$OUT/crictl-ps.txt"
crictl pods | tee "$OUT/crictl-pods.txt"
```

Send back: the two files. They help identify running game, gateway, director, and message broker containers.

If this was already collected, reference:

- `/tmp/ferrovax-capture/crictl-ps.txt`
- `/tmp/ferrovax-capture/crictl-pods.txt`

Status: not yet captured in the evidence provided for this incident.

### Step 8 - Search logs for server-list registration messages

Run this inside the Alpine guest VM. This searches common log locations for listing, registration, external address, region, gateway, FLS, and warning/error messages.

```sh
# Location: Alpine guest VM shell
OUT=/tmp/ferrovax-capture
mkdir -p "$OUT"
PAT='register|registration|server list|listing|public|external|advertis'
PAT="$PAT|region|gateway|FLS|port|error|warning"
grep -RniE "$PAT" /var /opt /home 2>/dev/null | head -500 | \
  tee "$OUT/registration-log-search.txt"
```

Send back: `registration-log-search.txt`. If the file is empty, say it returned no results. Replace any secrets with `REDACTED`.

If this was already collected, reference:

- `/tmp/ferrovax-capture/registration-log-search.txt`

Status: not yet captured in the evidence provided for this incident.

### Step 9 - Router web UI check

Run this in the router web UI, not in a terminal. This confirms whether internet traffic is forwarded to the VM, not the Windows host.

| Router item | What to record |
|---|---|
| WAN/public IP | Record the WAN IP shown by the router. Do not post it publicly; send privately to the technical reviewer if needed. |
| Forward target | Confirm all Dune/game/listing forwards target `192.168.1.125`, not `192.168.1.192`. |
| Protocol | Confirm UDP is forwarded where UDP is required. Do not rely on TCP-only forwarding for game traffic. |
| Port list | Record every forwarded external port and internal port. |
| Double NAT/CGNAT | Check if the router WAN IP matches a public IP lookup. If it does not, there may be ISP router, double NAT, or CGNAT. |

Status: not yet captured in the evidence provided for this incident.

### Step 10 - External reachability test

Run this from outside the LAN, such as a phone hotspot or remote PC. Do not use the same Wi-Fi/LAN as the server. UDP testing is difficult, so the best test is to watch the VM while a remote client tries to find or connect to the server.

Before this step, confirm whether Step 3 or Step 4 and Step 8 were completed:

- If Kubernetes service output has not been collected, complete Step 3 or Step 4 first.
- If registration logs have not been collected, complete Step 8 first.
- If both have been collected, reference `kube-services.txt`, `kube-endpoints.txt`, and `registration-log-search.txt` while interpreting packet results.

```sh
# Location: Alpine guest VM shell
# Run while a remote client searches/connects
OUT=/tmp/ferrovax-capture
mkdir -p "$OUT"
command -v tcpdump && tcpdump -ni eth0 'udp or tcp' | \
  tee "$OUT/external-attempt-packets.log"
```

If `tcpdump` is not installed:

1. Record that `tcpdump` is missing.
2. Do not stop the investigation.
3. Use the already collected Kubernetes files from Step 3 or Step 4.
4. Use the server registration log file from Step 8.
5. Use the router web UI results from Step 9.

If packets arrive from the external client but the server still does not list or connect, the issue is likely above the router layer. Compare packet timing with Kubernetes services and server registration logs. If no packets arrive, focus on router forwarding, ISP gateway, double NAT, CGNAT, or wrong forwarded ports.

### Step 11 - Package the capture files

Run this inside the Alpine guest VM after the commands above. It creates one archive that can be sent to the technical reviewer.

```sh
# Location: Alpine guest VM shell
cd /tmp
tar -czf ferrovax-capture-$(date +%Y%m%d-%H%M%S).tar.gz ferrovax-capture
ls -lh ferrovax-capture-*.tar.gz
```

Send back the `.tar.gz` file. Review the files first and replace secrets with `REDACTED` if any appear.

## 7. How to Interpret the Next Results

| Result | Likely next focus |
|---|---|
| Kubernetes services are ClusterIP only | The stack may be healthy internally but not exposed to LAN/WAN. Determine required NodePort/hostPort/LoadBalancer/hostNetwork model. |
| NodePort or hostPort exists | Router forwards must target `192.168.1.125` and the exposed node/host ports. |
| Listeners exist only on `10.42.1.x` | Traffic may be confined to pod network. Validate service exposure to node IP. |
| Listeners exist on `0.0.0.0` or `192.168.1.125` | Check router forwarding, public IP, external reachability, and server-list registration logs. |
| Registration logs show external address, region, auth, or FLS errors | Focus on server-list configuration, advertised address, token/auth, or region configuration. |
| No external packets reach `eth0` during remote test | Focus on router forwarding, ISP gateway, double NAT, CGNAT, or wrong public IP. |

## 8. Immediate Discovery Questions

- Is the server visible to a client on the same LAN as `192.168.1.125`?
- Is the server visible to a remote client outside the LAN?
- Which ports does the server documentation or generated config expect to expose publicly?
- Does the router forward those exact ports to `192.168.1.125`?
- Does Kubernetes expose those ports as NodePort, hostPort, LoadBalancer, or hostNetwork?
- Do server logs say registration succeeded, failed, timed out, or used a specific external address?
- Does the router WAN IP match a public IP lookup, or is there double NAT/CGNAT?

## 9. Current Non-Resolution Statement

At this point, the incident is narrowed but not resolved. Hyper-V host networking and host firewall are less likely based on available evidence. The next likely failure layer is Kubernetes service exposure, router forwarding, public/advertised address, or server-list registration. A final RCA should not be written until the Kubernetes service exposure, listener state, router forwards, and registration logs are reviewed together.
