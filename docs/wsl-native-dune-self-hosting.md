# WSL-native Dune self-hosting runbook

## Purpose

This runbook records the working Windows + WSL 2 topology for a Dune: Awakening self-hosted server after Docker Desktop networking could not reliably carry the public game UDP path.

Validated target:

```text
Windows host + WSL 2 mirrored networking
Ubuntu WSL Docker Engine
Dune server containers launched from the WSL-native Docker daemon
Router forwards targeting the Windows/WSL mirrored LAN IP
```

This is an operator runbook for the tested WSL path. It does not claim Docker Desktop is unusable for all development tasks; it documents the known-good public self-hosting path for this host class.

## Known-good evidence

Known-good diagnostic bundle:

```text
known-good-wsl-native-20260619T102922Z.tar.gz
```

Known-good runtime evidence:

```text
Docker Engine: 29.6.0
Operating system: Ubuntu 26.04 LTS
Kernel: 6.6.114.1-microsoft-standard-WSL2
Docker root: /var/lib/docker
WSL mirrored LAN IP: 192.168.68.21
```

Known-good container shape:

| Container | Network | Required exposure |
|---|---|---|
| `dune-server-survival-1` | `host` | UDP `7778`, UDP `7888` |
| `dune-server-overmap` | `host` | UDP `7777`, UDP `7889` |
| `dune-server-gateway` | `dune-net` | Advertises the public RMQ/game endpoints |
| `dune-rmq-game` | `dune-net` | TCP `31992`, TCP `31983` |
| `dune-rmq-admin` | `dune-net` | TCP `32673` bound to localhost |
| `dune-postgres` | `dune-net` | TCP `15432` bound to localhost |

Known-good DB state after a clean init included populated `dune.world_partition` rows. A prior empty `world_partition` state caused `LoadPartitionDefinition(... got 0 rows, expected exactly 1)` and was resolved by resetting to a clean default state.

## Why this path uses WSL-native Docker

Docker Desktop networking was tested and rejected for this host because the game client UDP path did not reliably reach the world server.

Observed failures:

| Attempt | Result |
|---|---|
| Docker Desktop host networking | Windows received UDP for `7778` but reported port unreachable / no transport endpoint. |
| Docker Desktop bridge with UDP `-p` publish | Docker metadata showed port bindings, but captures inside the container saw no UDP. |
| Explicit Docker Desktop HostIp publish to `192.168.68.21` | Docker metadata showed bindings, but packets still did not reach the container. |
| Windows relay direct to Docker bridge IP | Windows could not route to Docker bridge container addresses such as `172.18.x.x`. |
| WSL relay to Docker Desktop bridge IP | Plain WSL could not deliver packets into Docker Desktop's bridge network. |

The working solution moved the daemon into Ubuntu WSL and enabled WSL mirrored networking so the server sockets were reachable through the Windows/WSL LAN identity.

## Windows WSL configuration

Create or update:

```text
%UserProfile%\.wslconfig
```

Known-good shape:

```ini
[wsl2]
networkingMode=mirrored
firewall=true
localhostForwarding=true
memory=24GB
processors=4
swap=98GB
```

Restart WSL after changing this file:

```powershell
wsl --shutdown
```

Verify from Ubuntu:

```bash
ip addr
hostname -I
```

Expected evidence:

```text
eth0 has 192.168.68.21/22
hostname -I includes 192.168.68.21
```

If WSL only shows a `172.x.x.x` NAT address, mirrored networking is not active.

## Docker Engine requirements

Use Docker Engine inside Ubuntu WSL, not Docker Desktop's daemon.

Verify:

```bash
export DOCKER_HOST=unix:///var/run/docker.sock

docker info --format 'Name={{.Name}} OperatingSystem={{.OperatingSystem}} ServerVersion={{.ServerVersion}}'
```

Expected:

```text
Name=Tabr-Tau OperatingSystem=Ubuntu 26.04 LTS ServerVersion=29.6.0
```

If the output says `Docker Desktop`, the shell is pointed at the wrong daemon or Docker Desktop WSL integration has taken over the socket.

Authoritative daemon checks:

```bash
ps -fp "$(cat /var/run/docker.pid)"
readlink -f /proc/"$(cat /var/run/docker.pid)"/exe
```

Expected daemon path:

```text
/usr/bin/dockerd
```

### Docker Desktop UI warning

Docker Desktop may still display containers even when the WSL-native daemon owns them. Trust `docker info`, `/var/run/docker.pid`, and `/usr/bin/dockerd`; do not use the Desktop UI as the source of truth.

### Remove the temporary TCP Docker API listener

During troubleshooting, the daemon was observed with:

```text
-H tcp://127.0.0.1:2375
```

That listener is not required for this stack and should be removed once the server is stable. The target listener is:

```text
-H unix:///var/run/docker.sock
```

Do not leave an unauthenticated Docker API listener enabled longer than necessary.

## Port plan

Mirrored WSL caused Windows-local port collisions with the upstream default RMQ host ports. The working setup moved only the backend RMQ host ports and preserved game UDP ports.

| Purpose | Original | Known-good | Exposure |
|---|---:|---:|---|
| Game UDP range | `7777-7921/udp` | `7777-7921/udp` | Router + firewall to `192.168.68.21` |
| RMQ Game AMQP | `31982/tcp` | `31992/tcp` | Router + firewall to `192.168.68.21` |
| RMQ Game HTTP | `31983/tcp` | `31983/tcp` | Router + firewall to `192.168.68.21` |
| RMQ Admin AMQP | `32573/tcp` | `32673/tcp` | Localhost only |
| Postgres | `15432/tcp` | `15432/tcp` | Localhost only |
| Director | `11717/tcp` | `11717/tcp` | Localhost only |

Do not remap the public game UDP ports unless the game server, gateway, router, and firewall configuration are changed together.

## Router forwarding

Known-good target:

```text
192.168.68.21
```

Required forwards:

```text
UDP 7777-7921 -> 192.168.68.21
TCP 31992      -> 192.168.68.21
TCP 31983      -> 192.168.68.21
```

The old TCP rule was:

```text
TCP 31982-31983 -> 192.168.68.21
```

That is insufficient after moving RMQ Game to `31992`.

## Windows firewall / mirrored WSL firewall

Allow inbound traffic for:

```powershell
New-NetFirewallRule `
  -DisplayName "Dune WSL RMQ Game 31992 TCP" `
  -Direction Inbound `
  -Action Allow `
  -Protocol TCP `
  -LocalPort 31992

New-NetFirewallRule `
  -DisplayName "Dune WSL RMQ HTTP 31983 TCP" `
  -Direction Inbound `
  -Action Allow `
  -Protocol TCP `
  -LocalPort 31983

New-NetFirewallRule `
  -DisplayName "Dune WSL UDP 7777-7921" `
  -Direction Inbound `
  -Action Allow `
  -Protocol UDP `
  -LocalPort 7777-7921
```

If normal Windows firewall rules are insufficient, inspect WSL Hyper-V firewall policy. Keep this as a troubleshooting step; do not assume it is required on every host.

## Runtime verification

Verify listeners:

```bash
ss -ltnup | grep -E ':31983|:31992|:7777|:7778|:7888|:7889'
```

Known-good shape:

```text
udp 0.0.0.0:7777
udp 0.0.0.0:7778
udp 192.168.68.21:7888
udp 192.168.68.21:7889
tcp 0.0.0.0:31992
tcp 0.0.0.0:31983
```

Verify containers:

```bash
docker ps --format 'table {{.Names}}\t{{.Status}}\t{{.Ports}}\t{{.Networks}}'
```

Known-good RMQ Game exposure:

```text
dune-rmq-game  0.0.0.0:31992->5672/tcp, 0.0.0.0:31983->15672/tcp
```

Verify gateway advertised ports:

```bash
docker inspect dune-server-gateway \
  --format '{{range .Config.Cmd}}{{println .}}{{end}}' \
  | grep -E 'RMQGameHostname|RMQGamePort|RMQGameHttpPort'
```

Known-good shape:

```text
--RMQGameHostname=<public-ip>
--RMQGamePort=31992
```

## Packet-capture validation

For RMQ Game TCP:

```bash
sudo tcpdump -nnvvv -tttt -i any 'tcp and port 31992'
```

Known-good external evidence included inbound public-source traffic to:

```text
<public-source> > 192.168.68.21.31992
192.168.68.21.31992 > <public-source>
```

For world UDP:

```bash
sudo tcpdump -nnvvv -tttt -i any 'udp and (port 7777 or port 7778 or port 7888 or port 7889)'
```

Game entry should show client/public UDP reaching the world server port, especially `7778` for `Survival_1`.

## Registry image note

During WSL-native migration, `registry.funcom.com` returned public DNS `NXDOMAIN`. If the images already exist in Docker Desktop, export them from Docker Desktop and load them into the WSL-native engine rather than repeatedly retrying public DNS.

High-level flow:

```text
Docker Desktop image cache -> docker save -> tar file -> WSL-native docker load
```

## Known follow-up work

- Centralize hardcoded RMQ host ports into variables, for example `RMQ_GAME_HOST_PORT` and `RMQ_ADMIN_HOST_PORT`.
- Remove `tcp://127.0.0.1:2375` from `dockerd` startup.
- Resolve host DNS behavior separately.
- Avoid Docker Desktop WSL integration taking ownership of `/var/run/docker.sock` when WSL-native Docker is the target runtime.
- Preserve a known-good diagnostic bundle after every major runtime change.
