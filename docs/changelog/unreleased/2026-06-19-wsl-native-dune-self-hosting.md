# WSL-native Dune self-hosting guide

## Summary

Added an operator-facing guide for running a Dune: Awakening self-hosted stack on Windows with WSL 2 mirrored networking and the Ubuntu WSL Docker Engine.

## Why

Windows + WSL operators need a documented path that distinguishes Docker Desktop development convenience from the WSL-native runtime required for reliable public game connectivity on affected hosts. The guide converts troubleshooting evidence into general setup guidance and upstream recommendations.

## Added

- `docs/wsl-native-dune-self-hosting.md` as a reusable guide for WSL-native operation.
- Placeholder-based configuration guidance using `<LAN_IP>`, `<PUBLIC_IP>`, and overrideable RMQ port names.
- Router, firewall, listener, database-bootstrap, and packet-capture validation steps.
- Upstream recommendations for a first-class WSL-native profile and centralized RMQ host-port variables.

## Key guidance

```text
WSL mirrored networking: required for this topology
Docker daemon: Ubuntu WSL Docker Engine
Public game UDP: keep 7777-7921/udp by default
RMQ Game TCP: default 31982, override when Windows owns the port
RMQ Game HTTP TCP: 31983
RMQ Admin TCP: localhost-only, override when Windows owns the port
```

## Follow-up

- Centralize RMQ host ports into variables.
- Remove any temporary local Docker TCP API listener.
- Resolve host DNS behavior separately.
- Prevent Docker Desktop WSL integration from taking ownership of `/var/run/docker.sock` when WSL-native Docker is selected.
