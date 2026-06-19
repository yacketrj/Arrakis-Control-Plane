# WSL-native Dune self-hosting guide

## Purpose

This guide describes a reusable Windows + WSL 2 topology for running a Dune: Awakening self-hosted server when operators need public game access from a Windows host.

Recommended topology:

```text
Windows host
WSL 2 mirrored networking
Ubuntu WSL Docker Engine
Dune server containers launched from the WSL-native Docker daemon
Router forwards targeting the host LAN IP that WSL receives in mirrored mode
```

This guide is written for operators and upstream maintainers. Replace example IP addresses, public addresses, and port overrides with values appropriate for the target host.

## When to use this topology

Use the WSL-native topology when the operator is hosting from Windows and needs the public game client path to reach the world server ports reliably.

Prefer a native Linux host or dedicated Linux VM when available. Use this WSL-native guide when Windows + WSL is the required host model.

This topology is useful when Docker Desktop networking presents one or more of these symptoms:

| Symptom | Likely meaning |
|---|---|
| Server appears in the browser but character-to-world transition fails or hangs | Public discovery path works, but game/RMQ/world traffic is incomplete. |
| Windows receives forwarded UDP but reports port unreachable / no transport endpoint | Docker Desktop did not expose a usable Windows-side UDP endpoint. |
| Docker Desktop `-p .../udp` metadata exists but container captures see no UDP | Docker metadata and real packet delivery disagree. |
| Windows cannot route to Docker bridge container IPs | A Windows-side relay cannot target Docker bridge addresses directly. |
| Plain WSL cannot reach Docker Desktop bridge containers | A WSL-side relay cannot complete the Docker Desktop path. |

## Address and variable placeholders

Use placeholders instead of copying example addresses directly.

| Placeholder | Meaning | Example |
|---|---|---|
| `<LAN_IP>` | Windows/WSL mirrored LAN IP receiving router forwards | `192.168.68.21` |
| `<PUBLIC_IP>` | Public WAN IP or DNS name advertised to clients | `203.0.113.10` |
| `<WSL_DISTRO>` | Ubuntu WSL distro name | `Ubuntu-26.04` |
| `<RMQ_GAME_PORT>` | Public TCP AMQP port for game RMQ | default `31982`, example override `31992` |
| `<RMQ_GAME_HTTP_PORT>` | Public TCP HTTP/management path used by game RMQ where applicable | `31983` |
| `<RMQ_ADMIN_PORT>` | Localhost-only Admin RMQ AMQP port | default `32573`, example override `32673` |

## Required Windows WSL configuration

Create or update:

```text
%UserProfile%\.wslconfig
```

Recommended WSL 2 configuration:

```ini
[wsl2]
networkingMode=mirrored
firewall=true
localhostForwarding=true
```

Optional resource controls may be added for larger hosts:

```ini
memory=24GB
processors=4
swap=98GB
```

Restart WSL after changing this file:

```powershell
wsl --shutdown
```

Reopen the Ubuntu distro and verify that WSL has the LAN address:

```bash
ip addr
hostname -I
```

Expected shape:

```text
eth0 has <LAN_IP>/<prefix>
hostname -I includes <LAN_IP>
```

If WSL only shows a private NAT address such as `172.x.x.x`, mirrored networking is not active. Check `wsl --version`, run `wsl --update`, verify `.wslconfig`, then restart WSL again.

## Docker Engine requirement

Use Docker Engine inside Ubuntu WSL. Do not run the Dune server workload through Docker Desktop's daemon for this topology.

Verify the active daemon:

```bash
export DOCKER_HOST=unix:///var/run/docker.sock

docker info --format 'Name={{.Name}} OperatingSystem={{.OperatingSystem}} ServerVersion={{.ServerVersion}}'
```

Expected shape:

```text
Name=<windows-hostname> OperatingSystem=Ubuntu <version> ServerVersion=<docker-engine-version>
```

If the output says `Docker Desktop`, the shell is pointed at the wrong daemon or Docker Desktop WSL integration has taken over the socket.

Authoritative local daemon checks:

```bash
ps -fp "$(cat /var/run/docker.pid)"
readlink -f /proc/"$(cat /var/run/docker.pid)"/exe
```

Expected daemon path:

```text
/usr/bin/dockerd
```

### Docker Desktop UI warning

Docker Desktop may still show containers when WSL integration is enabled. Do not use the Desktop UI as proof of ownership. Trust `docker info`, `/var/run/docker.pid`, and the daemon executable path.

### Docker socket hygiene

The server stack does not require an unauthenticated Docker TCP API listener. If the local daemon was started with a listener such as:

```text
-H tcp://127.0.0.1:2375
```

remove it after the stack is stable. The desired listener is:

```text
-H unix:///var/run/docker.sock
```

## Port model

The public game ports should remain stable. Backend RMQ host ports may need to be overrideable because mirrored WSL shares more of the Windows localhost and LAN port surface.

| Purpose | Default | Recommended handling |
|---|---:|---|
| Game UDP range | `7777-7921/udp` | Keep default unless all game, gateway, router, and firewall settings are changed together. |
| RMQ Game AMQP | `31982/tcp` | Use default when free; otherwise set an override such as `31992`. |
| RMQ Game HTTP | `31983/tcp` | Keep default when free. |
| RMQ Admin AMQP | `32573/tcp` | Bind to localhost only; override if Windows already owns the port, for example `32673`. |
| Postgres | `15432/tcp` | Bind to localhost only. |
| Director | `11717/tcp` | Bind to localhost only. |

Upstream recommendation: expose RMQ host ports through centralized variables instead of hardcoded literals, for example:

```env
RMQ_GAME_HOST_PORT=31982
RMQ_GAME_HTTP_HOST_PORT=31983
RMQ_ADMIN_HOST_PORT=32573
```

This lets WSL operators avoid Windows port collisions without editing multiple runtime files.

## Detecting port collisions

Before starting the stack, check whether Windows or WSL already owns required ports.

From WSL:

```bash
ss -ltnup | grep -E ':31982|:31983|:32573|:7777|:7778|:7888|:7889' || true
```

From PowerShell:

```powershell
netstat -ano | findstr ":31982 :31983 :32573 :7777 :7778 :7888 :7889"
```

If Windows owns a backend RMQ port, prefer changing the RMQ host-port override. Do not kill generic Windows service host processes to free ports.

## Router forwarding

Router forwards must target `<LAN_IP>`, the address assigned to the Windows host and visible inside WSL mirrored networking.

Required forwards:

```text
UDP 7777-7921 -> <LAN_IP>
TCP <RMQ_GAME_PORT> -> <LAN_IP>
TCP <RMQ_GAME_HTTP_PORT> -> <LAN_IP>
```

Example with remapped RMQ Game:

```text
UDP 7777-7921 -> 192.168.68.21
TCP 31992      -> 192.168.68.21
TCP 31983      -> 192.168.68.21
```

Do not assume a TCP range that includes the default RMQ Game port still applies after remapping. For example, `31982-31983` does not include `31992`.

## Windows firewall and mirrored WSL firewall

Allow inbound traffic for the selected public ports.

PowerShell as Administrator:

```powershell
New-NetFirewallRule `
  -DisplayName "Dune WSL RMQ Game TCP" `
  -Direction Inbound `
  -Action Allow `
  -Protocol TCP `
  -LocalPort <RMQ_GAME_PORT>

New-NetFirewallRule `
  -DisplayName "Dune WSL RMQ Game HTTP TCP" `
  -Direction Inbound `
  -Action Allow `
  -Protocol TCP `
  -LocalPort <RMQ_GAME_HTTP_PORT>

New-NetFirewallRule `
  -DisplayName "Dune WSL UDP 7777-7921" `
  -Direction Inbound `
  -Action Allow `
  -Protocol UDP `
  -LocalPort 7777-7921
```

If normal Windows firewall rules are insufficient, inspect the WSL Hyper-V firewall policy. Treat Hyper-V firewall work as a troubleshooting step, not a universal requirement.

## Startup validation

After startup, verify listeners:

```bash
ss -ltnup | grep -E ':31982|:31983|:31992|:7777|:7778|:7888|:7889'
```

Expected shape:

```text
udp 0.0.0.0:7777
udp 0.0.0.0:7778
udp <LAN_IP>:7888
udp <LAN_IP>:7889
tcp 0.0.0.0:<RMQ_GAME_PORT>
tcp 0.0.0.0:<RMQ_GAME_HTTP_PORT>
```

Verify containers:

```bash
docker ps --format 'table {{.Names}}\t{{.Status}}\t{{.Ports}}\t{{.Networks}}'
```

Expected RMQ Game exposure:

```text
dune-rmq-game  0.0.0.0:<RMQ_GAME_PORT>->5672/tcp, 0.0.0.0:<RMQ_GAME_HTTP_PORT>->15672/tcp
```

Verify gateway-advertised ports:

```bash
docker inspect dune-server-gateway \
  --format '{{range .Config.Cmd}}{{println .}}{{end}}' \
  | grep -E 'RMQGameHostname|RMQGamePort|RMQGameHttpPort'
```

Expected shape:

```text
--RMQGameHostname=<PUBLIC_IP-or-public-DNS>
--RMQGamePort=<RMQ_GAME_PORT>
```

## Database bootstrap validation

After a clean initialization, verify the world partitions exist:

```bash
docker exec dune-postgres psql -U dune -d dune -c "
select partition_id, server_id, map, dimension_index, blocked, label
from dune.world_partition
order by partition_id;
"
```

Expected: rows exist for the configured world maps, including `Survival_1` and `Overmap` when those worlds are enabled.

If world servers log this error:

```text
LoadPartitionDefinition(... got 0 rows, expected exactly 1)
```

then the database bootstrap state is incomplete. Prefer rerunning the supported initialization/reset workflow rather than manually inserting partition rows.

## Packet-capture validation

Use packet captures to distinguish routing failure from application/session failure.

RMQ Game TCP:

```bash
sudo tcpdump -nnvvv -tttt -i any 'tcp and port <RMQ_GAME_PORT>'
```

Expected public-path evidence:

```text
<client-or-router-source> > <LAN_IP>.<RMQ_GAME_PORT>
<LAN_IP>.<RMQ_GAME_PORT> > <client-or-router-source>
```

World UDP:

```bash
sudo tcpdump -nnvvv -tttt -i any 'udp and (port 7777 or port 7778 or port 7888 or port 7889)'
```

Expected after successful handoff: client/public UDP reaches the world server port, especially `7778` for `Survival_1`.

If RMQ TCP arrives but world UDP never follows, continue debugging gateway/session/RMQ handoff rather than router TCP reachability.

## Docker Desktop image cache note

Some operators may already have required Funcom images cached in Docker Desktop while the WSL-native daemon cannot pull from the registry. If the registry host is not resolvable or not reachable, move images from Docker Desktop to WSL-native Docker:

```text
Docker Desktop image cache -> docker save -> tar file -> WSL-native docker load
```

Use this only for images the operator is authorized to run.

## Upstream recommendations

- Add a first-class WSL-native network profile that assumes WSL mirrored networking and Ubuntu WSL Docker Engine.
- Centralize RMQ host ports into overrideable configuration variables.
- Keep public game UDP ports stable by default.
- Document router/firewall requirements using `<LAN_IP>`, `<PUBLIC_IP>`, and selected RMQ ports rather than host-specific values.
- Add diagnostics that capture Docker daemon identity, WSL address mode, listener state, gateway-advertised ports, world partition rows, and short packet summaries.
- Warn operators when Docker Desktop owns `/var/run/docker.sock` but the selected profile expects WSL-native Docker.

## Troubleshooting matrix

| Observation | Next check |
|---|---|
| `docker info` reports Docker Desktop | Fix Docker context/socket ownership before starting the stack. |
| WSL has only `172.x.x.x` address | Mirrored networking is not active. |
| Server appears but connection hangs | Check RMQ Game TCP reachability and gateway-advertised RMQ port. |
| Character-to-world transition fails | Capture world UDP ports, especially `7778`. |
| `world_partition` is empty | Rerun supported init/reset workflow. |
| Backend RMQ port bind fails | Check Windows and WSL port ownership and choose an RMQ host-port override. |
