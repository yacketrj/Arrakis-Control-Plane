# WSL-native Dune self-hosting runbook

## Summary

Documented the known-good Windows + WSL 2 mirrored networking path for running a Dune: Awakening self-hosted stack from the Ubuntu WSL Docker Engine.

## Why

Docker Desktop networking did not reliably pass the public game UDP path to the world server on the tested host. The validated path moved the daemon into Ubuntu WSL, enabled WSL mirrored networking, and preserved game UDP ports while remapping colliding backend RMQ host ports.

## Added

- `docs/wsl-native-dune-self-hosting.md` as the operator runbook for the working WSL-native topology.
- Known-good runtime evidence, listener checks, router/firewall requirements, packet-capture checks, and cleanup follow-up items.

## Key validated configuration

```text
WSL mirrored networking: enabled
WSL LAN IP: 192.168.68.21
Docker daemon: Ubuntu WSL Docker Engine
Game UDP: 7777-7921/udp
RMQ Game: 31992/tcp
RMQ Game HTTP: 31983/tcp
RMQ Admin: 32673/tcp localhost-only
```

## Follow-up

- Centralize hardcoded RMQ host ports into variables.
- Remove the temporary Docker TCP API listener on `127.0.0.1:2375`.
- Resolve the host DNS behavior separately.
- Prevent Docker Desktop WSL integration from taking ownership of `/var/run/docker.sock` when WSL-native Docker is the selected runtime.
