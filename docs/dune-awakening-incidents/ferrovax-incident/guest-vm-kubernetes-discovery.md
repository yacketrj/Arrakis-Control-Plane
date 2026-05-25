# Guest VM and Kubernetes Discovery

## Evidence

The latest command output was collected from inside the guest VM.

```text
Guest hostname prompt: duneawakening
Guest OS: Alpine Linux v3.23
Guest primary interface: eth0
Guest LAN IP: 192.168.1.125/24
Guest default gateway: 192.168.1.1
Container networking observed: flannel.1 and cni0
Pod network observed: 10.42.1.0/24 on cni0
Kubernetes-related firewall chains observed: KUBE-*, KUBE-ROUTER-*, KUBE-POD-FW-*
Docker CLI: not present or not in PATH
ss command: not present
RabbitMQ/Erlang processes: present
```

Sensitive or environment-specific namespace and pod identifiers from firewall output are not reproduced here.

## Current Interpretation

The server environment is not a simple standalone process on the VM. The guest VM is running an Alpine-based Kubernetes/CNI-style stack with flannel and Kubernetes iptables chains. RabbitMQ processes are visible on the VM, but Docker is not available from the shell that collected the evidence.

This changes the next investigation path. Router forwarding to `192.168.1.125` may still be required, but it is not sufficient by itself. The Dune-facing game/listing ports must also be exposed from the Kubernetes pod/service layer to the VM node IP.

## What This Proves

```text
The VM is Alpine Linux.
The VM is reachable on LAN IP 192.168.1.125.
The VM uses gateway 192.168.1.1.
Kubernetes-style networking is present.
Flannel/CNI networking is present.
RabbitMQ is running in the environment.
The Docker CLI is not available from the current VM shell.
The requested ss listener command cannot run because ss is not installed.
```

## What This Does Not Prove

```text
It does not prove the Dune game server ports are listening on the VM node IP.
It does not prove Kubernetes Services expose the game ports externally.
It does not prove router forwarding points to the correct node ports.
It does not prove server-list registration succeeded.
It does not prove the configured advertised public address is correct.
It does not prove players can reach the gateway or game server from outside the LAN.
```

## Current Working Assessment

The most important open question is whether the Dune game/listing ports are exposed from Kubernetes to the VM node address `192.168.1.125`.

If the Dune server pods are only reachable on pod-network addresses such as `10.42.1.x`, external players and the public server-list path may not reach them unless Kubernetes Service, hostPort, NodePort, LoadBalancer, or hostNetwork exposure is correctly configured.

## Next Evidence Required

```text
1. Confirm whether kubectl is available.
2. List Kubernetes nodes, namespaces, pods, services, endpoints, and ingresses.
3. Identify the Dune game/server/gateway services and their exposed ports.
4. Confirm whether services use NodePort, hostPort, LoadBalancer, or hostNetwork.
5. Capture listener state using netstat or /proc because ss is not installed.
6. Search server logs for registration/listing/advertised address messages.
7. Confirm router port forwards target 192.168.1.125 and the externally exposed node ports.
8. Confirm the public WAN IP and rule out double NAT or CGNAT.
```

## Recommended Next Commands

Run inside the guest VM.

```sh
command -v kubectl || true
kubectl get nodes -o wide 2>/dev/null || true
kubectl get namespaces 2>/dev/null || true
kubectl get pods -A -o wide 2>/dev/null || true
kubectl get svc -A -o wide 2>/dev/null || true
kubectl get endpoints -A 2>/dev/null || true
kubectl get ingress -A 2>/dev/null || true
```

If `kubectl` is not available:

```sh
find / -name kubectl -type f 2>/dev/null | head
find /etc /var -maxdepth 4 -type f \( -name '*kubeconfig*' -o -name 'config' \) 2>/dev/null | head
```

Listener fallback when `ss` is missing:

```sh
netstat -tulpen 2>/dev/null || netstat -tuln 2>/dev/null || true
cat /proc/net/udp
cat /proc/net/tcp
```

Process/runtime discovery:

```sh
ps -ef | grep -Ei 'DuneSandbox|Awakening|Sandbox|kube|containerd|rancher|k3s|rabbit|flannel' | grep -v grep || true
command -v crictl && crictl ps || true
command -v nerdctl && nerdctl ps || true
```

Log search:

```sh
grep -RniE 'register|registration|server list|listing|public|external|advertis|region|gateway|FLS|port|error|warning' /var /opt /home 2>/dev/null | head -300
```

## Decision Branch

```text
Kubernetes Service exposes NodePort/hostPort on 192.168.1.125:
  Check router forwarding, public IP, and registration logs.

Kubernetes Service is ClusterIP only:
  The service may be healthy internally but not reachable externally. Determine required exposure model.

No service/endpoints point to the Dune pods:
  Investigate deployment/service creation and operator status.

Listeners exist only on pod network 10.42.1.x:
  Investigate Kubernetes service exposure to the node IP.

Server/listing registration logs show external-address or auth errors:
  Focus on registration configuration, advertised address, tokens, or region/listing configuration.
```
