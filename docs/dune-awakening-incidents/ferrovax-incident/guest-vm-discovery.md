# Guest VM Discovery

## Evidence Summary

Commands were run inside the guest VM.

```text
Guest OS: Alpine Linux v3.23
Guest IP: 192.168.1.125/24 on eth0
Default gateway: 192.168.1.1 via eth0
Pod network indicators: flannel.1 and cni0 present
Pod CIDR observed: 10.42.1.0/24
iptables default policy: INPUT ACCEPT, FORWARD ACCEPT, OUTPUT ACCEPT
nft ruleset: no output returned
Docker CLI: not present or not available in PATH
ss command: not present
RabbitMQ/Erlang processes: present
Direct DuneSandbox process from ps search: not shown in provided output
Kubernetes-style iptables chains: present
```

## Interpretation

The guest VM is Alpine Linux and appears to run a Kubernetes-style environment rather than a direct Docker workflow. The `flannel.1`, `cni0`, pod CIDR, and Kubernetes-style iptables chains indicate that the game services are likely managed through Kubernetes or a lightweight Kubernetes distribution.

The guest VM has the expected LAN address and default route. This supports the earlier Hyper-V evidence showing the VM at `192.168.1.125` on the external Hyper-V switch.

The output does not yet prove whether the Dune server ports are listening or exposed externally. The `ss` utility is missing, and the direct process search did not show the Dune game server process. The game services may be running inside pods/containers.

## What This Proves

```text
The guest VM is Alpine Linux.
The guest VM IP is 192.168.1.125.
The guest VM default gateway is 192.168.1.1.
The guest VM has Kubernetes-style pod networking.
RabbitMQ is running in the environment.
Docker is not the directly available runtime interface from this shell.
```

## What This Does Not Prove

```text
It does not prove the Dune game ports are listening.
It does not prove the Dune server pods are healthy.
It does not prove server-list registration succeeded.
It does not prove router port forwarding targets 192.168.1.125.
It does not prove external clients can reach the server.
```

## Next Evidence Required

Run these from inside the Alpine guest VM.

### Identify available Kubernetes/container tools

```sh
command -v kubectl || true
command -v k3s || true
command -v crictl || true
command -v nerdctl || true
command -v ctr || true
command -v containerd || true
```

### If kubectl is available

```sh
kubectl get nodes -o wide
kubectl get namespaces
kubectl get pods -A -o wide
kubectl get svc -A -o wide
kubectl get endpoints -A
kubectl get ingress -A 2>/dev/null || true
```

### If k3s is available but kubectl is not

```sh
k3s kubectl get nodes -o wide
k3s kubectl get namespaces
k3s kubectl get pods -A -o wide
k3s kubectl get svc -A -o wide
k3s kubectl get endpoints -A
k3s kubectl get ingress -A 2>/dev/null || true
```

### If crictl is available

```sh
crictl ps -a
crictl pods
```

### Listener alternatives because ss is missing

```sh
command -v netstat && netstat -tulpen || true
cat /proc/net/udp
cat /proc/net/tcp
```

## Current Branch

```text
If game services are not running or pods are not ready:
  Investigate runtime startup, pod events, image pulls, configuration, and resource constraints.

If game services are running but no service exposes required ports:
  Investigate Kubernetes Service, NodePort, hostNetwork, or ingress/service exposure configuration.

If ports are exposed on the VM but not reachable externally:
  Investigate router forwarding, double NAT, carrier-grade NAT, ISP firewall, and public IP mismatch.

If registration/listing logs show errors:
  Investigate service token, region, listing configuration, advertised address, and vendor-side registration path.
```
