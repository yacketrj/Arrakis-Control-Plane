# Glossary and Operational Variables

Use this document when a support person is unsure what a term means or what value should be recorded in case notes.

## Operational Variables

These values are installation-specific. They are not automatically sensitive, but some may need redaction depending on who receives the final package.

```text
CLOUD_PROVIDER          OCI, AWS, Azure, GCP, local, hosted provider, or unknown
CLOUD_INSTANCE_ID       Cloud VM or instance identifier, if cloud-hosted
INSTANCE_PATH           Active game/control-panel instance path on the host
LOG_PATH                Directory or file path where relevant logs are stored
SAVED_PATH              Active Saved directory for the game server
SERVICE_NAME            systemd service, Windows service, AMP instance, or other service name
CONTAINER_NAME          Container name, only if containers are used
DIRECTOR_SERVICE        Service, container, or process handling control-plane/director behavior
RABBITMQ_SERVICE        Service, container, or process running RabbitMQ or messaging, if present
DESTINATION_SERVICE     Service, container, or process for the destination map/server
DESTINATION_MAP         Map, partition, zone, or destination being tested
CLIENT_IP               Client public IP, when needed for packet capture
PUBLIC_IP               Server public IP or advertised external address
PRIVATE_IP              Server private, bind, or interface IP
PLAYER_ID               Player identifier, if needed for queue/log correlation
PLAYER_NAME             Player display name, if needed for operational notes
DATABASE_SERVICE        Database service, container, or process name, if discovered
DATABASE_HOST           Database host or endpoint, if discovered
DATABASE_PORT           Database port, if discovered
```

## Case Notes Template

```text
CLOUD_PROVIDER=
CLOUD_INSTANCE_ID=
INSTANCE_PATH=
LOG_PATH=
SAVED_PATH=
SERVICE_NAME=
CONTAINER_NAME=
DIRECTOR_SERVICE=
RABBITMQ_SERVICE=
DESTINATION_SERVICE=
DESTINATION_MAP=
CLIENT_IP=
PUBLIC_IP=
PRIVATE_IP=
DATABASE_SERVICE=
DATABASE_HOST=
DATABASE_PORT=
```

## Glossary

| Term | Meaning |
|---|---|
| Control plane | The service or services that decide where a player should connect or travel. |
| Director | A control-plane service that handles login/travel routing decisions. |
| Destination | The map, partition, server, or instance the player is trying to reach. |
| Game/client port | The port the game client connects to for gameplay traffic. Usually UDP. |
| IGW/server-to-server port | A port used by game services to communicate internally. |
| Instance path | The active install or management directory for a server instance. |
| Listener | A process actively waiting for network traffic on a port. |
| NAT | Network address translation. It maps traffic from one address/port to another. |
| Orchestration layer | The tool that starts, stops, or manages services. Examples: AMP, Docker Compose, systemd, Windows service, scripts. |
| Platform | Where the server is hosted. Examples: Linux VM, Windows VM, Hyper-V, Proxmox, OCI, AWS, Azure, GCP. |
| RabbitMQ | A messaging service sometimes used for communication between server components. |
| Runtime | The actual execution model for the server. Examples: direct process, Docker container, Windows service, systemd service. |
| Saved path | The directory where runtime data, saves, user settings, or generated state may be written. |
| Service | A long-running managed process, such as a Windows service or Linux systemd service. |
| Travel | The gameplay action of moving from one map, partition, zone, or destination to another. |

## Sensitive vs Operational Values

Usually sensitive:

```text
Tokens
JWTs
Passwords
Database passwords
RabbitMQ secrets
Private keys
Personal names
Raw player/account IDs when not vendor-required
Cloud resource IDs when sharing broadly
```

Usually operational:

```text
Service names
Container names
Map names
Local paths
Instance paths
Runtime names
Platform names
Port numbers
```

Operational values may still be redacted if the sharing audience is public or untrusted.
