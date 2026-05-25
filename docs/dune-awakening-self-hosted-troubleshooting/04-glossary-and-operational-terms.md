# Glossary and Operational Terms

Use this glossary when support staff need plain-language definitions while reading logs, commands, or escalation notes.

## Hosting and Platform Terms

| Term | Plain-language meaning |
|---|---|
| Host | The physical or virtual machine where services run. |
| Guest VM | A virtual machine running inside a hypervisor such as Hyper-V or Proxmox. |
| Hypervisor | Software that runs virtual machines, such as Hyper-V or Proxmox. |
| Cloud VM | A virtual machine hosted by OCI, AWS, Azure, GCP, or another cloud provider. |
| Control panel | A web UI used to manage servers, such as AMP. |
| Runtime | The layer that actually starts and runs the server process, such as Docker, systemd, Windows Service, AMP, or a custom script. |
| Orchestration | The management layer that starts, stops, restarts, and configures services or containers. |
| Container | A packaged runtime environment, commonly managed by Docker or another container runtime. |
| Bind mount | A host directory mapped into a container. The same files are visible from both the host and container. |
| Volume | Storage managed by a container runtime. It may contain persistent server or database data. |

## Network Terms

| Term | Plain-language meaning |
|---|---|
| Public IP | Address clients use from the internet. |
| Private IP | Internal address used inside the host, VM, VPC, VNet, or local network. |
| Bind address | Address the process listens on. A process bound to `127.0.0.1` is usually local-only. |
| Advertised address | Address the server tells clients or services to connect to. |
| Listener | A process actively waiting for network connections or packets on a port. |
| UDP | Network protocol commonly used for real-time game traffic. |
| TCP | Network protocol commonly used for management, APIs, databases, or message brokers. |
| NAT | Network address translation. It maps traffic from one address/port to another. |
| Firewall rule | Rule that allows or blocks traffic. It may exist on the host, cloud provider, router, hypervisor, or control panel. |
| Security group / NSG / security list | Cloud-provider firewall object. Different providers use different names. |

## Dune Server Terms

| Term | Plain-language meaning |
|---|---|
| Game server process | The process that runs the Dune game world or map. |
| Control plane | Service layer that coordinates login, travel, map assignment, or server state. |
| Director | A control-plane role that may assign players to maps or destinations. |
| Destination | The map, zone, partition, dungeon, or server the player is trying to enter. |
| Source | The map or service the player starts from before travel. |
| Partition | A numbered game-world destination or map instance. |
| Instanced map | A destination that may start, stop, or scale separately from the main world. |
| Starting map | The first map or world location players enter after login. |
| Travel | Movement from one map, zone, partition, or server instance to another. |
| Handoff | The server-side process of moving a player from source to destination. |

## Messaging and Data Terms

| Term | Plain-language meaning |
|---|---|
| RabbitMQ | Messaging system that can pass events between services. |
| Queue | A named message holding area. Services read messages from queues. |
| Consumer | A service connected to a queue and reading messages. |
| messages_ready | Messages waiting to be processed. |
| messages_unacknowledged | Messages sent to a consumer but not confirmed as processed. |
| Database | Persistent data store for player, world, or service state. |
| Persistence | Saving and loading server, player, or world data. |

## Troubleshooting Terms

| Term | Plain-language meaning |
|---|---|
| Reproduction | A controlled attempt to make the issue happen again. |
| Evidence window | The exact time range where support captures logs and system state. |
| Known working path | A login, travel, or server action that succeeds. |
| Known failing path | A login, travel, or server action that fails. |
| RCA | Root cause analysis. It should only be written after evidence supports it. |
| Hypothesis | A possible explanation that still needs evidence. |
| Redaction | Removing sensitive values before sharing logs or reports. |
| Escalation package | The evidence bundle sent to an engineer, vendor, or higher-level support team. |
