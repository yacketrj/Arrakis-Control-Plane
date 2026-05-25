# Glossary and Operational Variables

Use this document when support staff need plain-language definitions for troubleshooting terms, platform terms, runtime terms, and reusable case variables.

## Operational Variables

These values are discovered during troubleshooting. They are not automatically sensitive. Use actual values in local case notes and redact only when the sharing audience requires it.

```text
CLOUD_PROVIDER          OCI, AWS, Azure, GCP, local, other, or unknown
CLOUD_INSTANCE_ID       Cloud VM/instance identifier, if hosted in cloud
INSTANCE_PATH           Active game/control-panel instance path on the host
LOG_PATH                Directory or file path where relevant logs are stored
SAVED_PATH              Active Saved directory for the game server
SERVICE_NAME            systemd service, Windows service, AMP instance, or other service name
CONTAINER_NAME          Generic container name, if containers are used
DIRECTOR_SERVICE        Service/container/process handling director or control-plane logs
RABBITMQ_SERVICE        Service/container/process running RabbitMQ or messaging, if present
DESTINATION_SERVICE     Service/container/process for the destination map/server
DESTINATION_MAP         Map, partition, zone, or destination being tested
CLIENT_IP               Client public IP, when needed for packet capture
PUBLIC_IP               Server public IP or advertised external address
PRIVATE_IP              Server private, bind, or interface IP
DATABASE_SERVICE        Database service/container/process, if present
DATABASE_HOST           Database hostname or IP, if known
DATABASE_PORT           Database port, if known
PLAYER_ID               Player identifier, if needed for queue/log correlation
PLAYER_NAME             Player display name, if needed for operational notes
```

## Hosting and Runtime Terms

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
