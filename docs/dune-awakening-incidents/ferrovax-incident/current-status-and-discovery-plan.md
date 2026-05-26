# Ferrovax Incident - Current Status and Discovery Plan

This is the current operator-facing discovery plan. Versioned files may remain in the backend for revision history, but the active user-facing document should not expose version labels.

## Current Status

```text
Incident status: Open
Primary symptom: Server appears healthy/running locally but does not appear in the in-game server list.
Current phase: Kubernetes service exposure, router forwarding, and server-list registration validation.
Resolution: Not resolved; final RCA not proven.
Primary question: Are the Dune-facing services exposed from Kubernetes to VM node IP 192.168.1.125 and mapped correctly through the router?
```

## Confirmed Evidence

```text
Hyper-V host: Windows 10 Pro host named DUNE.
Hyper-V host LAN IP on external switch: 192.168.1.192.
Hyper-V VM: dune-awakening, running on External switch DuneAwakeningServerSwitch.
Guest VM IP: 192.168.1.125/24.
Guest VM default gateway: 192.168.1.1.
Guest OS: Alpine Linux v3.23.
Guest orchestration: Kubernetes-style networking present, including flannel.1, cni0, 10.42.1.0/24 pod network, and KUBE-* iptables chains.
Runtime clue: RabbitMQ/Erlang processes are visible.
Missing tools: ss is not installed; tcpdump availability has not yet been confirmed.
```

## Evidence Not Yet Captured

```text
Kubernetes services: not yet captured. Use Step 3 or Step 4.
Kubernetes pods/endpoints: not yet captured. Use Step 3 or Step 4.
Server registration/listing logs: not yet captured. Use Step 8.
Listener state: partially blocked because ss is missing. Use Step 5.
Packet capture: not yet captured. Use Step 10 after services/logs are collected, or while actively testing external traffic.
```

## Operator Location Guide

```text
Windows Hyper-V host:
  Prompt starts with PS C:\WINDOWS\system32> or another Windows path.
  Use for Hyper-V VM state, virtual switch, and Windows host firewall.

Alpine guest VM:
  Prompt looks like duneawakening:~# or another Linux shell prompt.
  Use for Kubernetes service exposure, server logs, and listeners.

Router web UI:
  Browser page for router, often 192.168.1.1 or vendor app.
  Use for WAN IP, port forwards, and double NAT/CGNAT checks.

External test device:
  Phone hotspot, remote PC, or machine outside the LAN/Wi-Fi.
  Use for remote visibility or packet-arrival testing.
```

## Recommended Command Set

### Step 1 - Confirm VM-side network state

Run on: Alpine guest VM.

```sh
cat /etc/os-release
ip addr
ip route
```

Send back the full output. This was captured once; rerun if the VM was rebooted or networking changed.

### Step 2 - Check whether Kubernetes tools are available

Run on: Alpine guest VM.

```sh
command -v kubectl || true
command -v k3s || true
command -v crictl || true
command -v ctr || true
command -v nerdctl || true
```

Send back the output. This has not been captured yet.

### Step 3 - Capture Kubernetes services and pods if kubectl works

Run on: Alpine guest VM, only if `kubectl` exists.

```sh
OUT=/tmp/ferrovax-capture
mkdir -p "$OUT"
kubectl get nodes -o wide | tee "$OUT/kube-nodes.txt"
kubectl get namespaces | tee "$OUT/kube-namespaces.txt"
kubectl get pods -A -o wide | tee "$OUT/kube-pods.txt"
kubectl get svc -A -o wide | tee "$OUT/kube-services.txt"
kubectl get endpoints -A | tee "$OUT/kube-endpoints.txt"
kubectl get ingress -A | tee "$OUT/kube-ingress.txt"
```

Send back all files under `/tmp/ferrovax-capture`, especially `kube-services.txt` and `kube-endpoints.txt`.

### Step 4 - Capture Kubernetes services and pods if only k3s works

Run on: Alpine guest VM, only if `kubectl` is unavailable but `k3s` exists.

```sh
OUT=/tmp/ferrovax-capture
mkdir -p "$OUT"
k3s kubectl get nodes -o wide | tee "$OUT/kube-nodes.txt"
k3s kubectl get namespaces | tee "$OUT/kube-namespaces.txt"
k3s kubectl get pods -A -o wide | tee "$OUT/kube-pods.txt"
k3s kubectl get svc -A -o wide | tee "$OUT/kube-services.txt"
k3s kubectl get endpoints -A | tee "$OUT/kube-endpoints.txt"
k3s kubectl get ingress -A | tee "$OUT/kube-ingress.txt"
```

Send back all files under `/tmp/ferrovax-capture`, especially `kube-services.txt`. If commands fail, send the exact error message.

### Step 5 - Capture listener state using netstat

Run on: Alpine guest VM.

```sh
OUT=/tmp/ferrovax-capture
mkdir -p "$OUT"
netstat -tulpen 2>/dev/null | tee "$OUT/netstat-listeners.txt" || \
netstat -tuln 2>/dev/null | tee "$OUT/netstat-listeners.txt" || true
cat /proc/net/udp  | tee "$OUT/proc-net-udp.txt"
cat /proc/net/tcp  | tee "$OUT/proc-net-tcp.txt"
cat /proc/net/udp6 | tee "$OUT/proc-net-udp6.txt"
cat /proc/net/tcp6 | tee "$OUT/proc-net-tcp6.txt"
```

Send back `netstat-listeners.txt` if present, plus the `/proc` files.

### Step 6 - Identify runtime processes

Run on: Alpine guest VM.

```sh
OUT=/tmp/ferrovax-capture
mkdir -p "$OUT"
PAT="DuneSandbox|Awakening|Sandbox|kube|containerd|k3s|rabbit|flannel"
ps -ef | grep -Ei "$PAT" | grep -v grep | tee "$OUT/runtime-processes.txt"
```

Send back `runtime-processes.txt`.

### Step 7 - Capture pod/container runtime state if crictl exists

Run on: Alpine guest VM, only if `crictl` exists.

```sh
OUT=/tmp/ferrovax-capture
mkdir -p "$OUT"
crictl ps -a | tee "$OUT/crictl-ps.txt"
crictl pods | tee "$OUT/crictl-pods.txt"
```

Send back `crictl-ps.txt` and `crictl-pods.txt`.

### Step 8 - Capture server-list registration and server logs

Run on: Alpine guest VM.

```sh
OUT=/tmp/ferrovax-capture
mkdir -p "$OUT"
PAT="register|registration|server list|listing|public|external|advertis"
PAT="$PAT|region|gateway|FLS|port|error|warning"
grep -RniE "$PAT" /var /opt /home 2>/dev/null | head -500 | \
  tee "$OUT/registration-log-search.txt"
```

Send back `registration-log-search.txt`. If empty, state that it returned no results. Replace secrets with `REDACTED`.

### Step 9 - Router web UI check

Run in: Router web UI, not in terminal.

```text
WAN/public IP:
  Record privately for technical review; do not post publicly.

Forward target:
  Confirm all Dune/game/listing forwards target 192.168.1.125, not 192.168.1.192.

Protocol:
  Confirm UDP is forwarded where UDP is required. TCP-only is not sufficient for game traffic.

Port list:
  Record every external port and internal port.

Double NAT/CGNAT:
  Check whether router WAN IP matches public IP lookup. If not, suspect ISP router, double NAT, or CGNAT.
```

### Step 10 - External reachability and tcpdump fallback

Run from: an external test device for the client attempt, and Alpine guest VM for packet observation.

Before this step, confirm Step 3 or Step 4 and Step 8 are complete when possible. Kubernetes service output and server logs tell us which services and ports matter.

```sh
OUT=/tmp/ferrovax-capture
mkdir -p "$OUT"
command -v tcpdump && tcpdump -ni eth0 "udp or tcp" | \
  tee "$OUT/external-attempt-packets.log"
```

If `tcpdump` is missing, record `tcpdump missing` in the notes. Then rely on:

```text
Step 3 or Step 4: Kubernetes service and endpoint output.
Step 8: Server registration/listing log output.
Step 9: Router WAN IP and forwarding configuration.
```

If Step 3/4 or Step 8 have not been collected, collect them before guessing.

### Step 11 - Package capture files

Run on: Alpine guest VM after all applicable steps.

```sh
cd /tmp
tar -czf ferrovax-capture-$(date +%Y%m%d-%H%M%S).tar.gz ferrovax-capture
ls -lh ferrovax-capture-*.tar.gz
```

Send back the `.tar.gz` file. Review files first and replace secrets with `REDACTED` if any appear.

## Decision Table

```text
Kubernetes services are ClusterIP only:
  Stack may be healthy internally but not exposed to LAN/WAN.

NodePort or hostPort exists:
  Router forwards must target 192.168.1.125 and the exposed node/host ports.

Listeners exist only on 10.42.1.x:
  Traffic may be confined to pod network. Validate service exposure to node IP.

Listeners exist on 0.0.0.0 or 192.168.1.125:
  Check router forwarding, public IP, external reachability, and registration logs.

Registration logs show address, region, auth, or FLS errors:
  Focus on listing configuration, advertised address, token/auth, or region configuration.

No external packets reach eth0 during remote test:
  Focus on router forwarding, ISP gateway, double NAT, CGNAT, or wrong public IP.
```

## Current Non-Resolution Statement

The incident is narrowed but not resolved. Hyper-V host networking and host firewall are less likely based on available evidence. The next likely failure layer is Kubernetes service exposure, router forwarding, public/advertised address, or server-list registration. A final RCA should not be written until Kubernetes service exposure, listener state, router forwards, packet behavior, and registration logs are reviewed together.
