# WSL-native Dune self-hosting documentation summary

Documented the known-good Windows + WSL 2 mirrored networking path for running a Dune: Awakening self-hosted stack from the Ubuntu WSL Docker Engine.

Added operator guidance for:

- WSL mirrored networking requirements.
- Ubuntu WSL Docker Engine verification.
- Known-good container layout.
- Required router forwards.
- Required Windows firewall allowances.
- RMQ host-port remaps used to avoid mirrored-networking port collisions.
- Runtime listener and packet-capture validation.
- Follow-up work to centralize RMQ port configuration and clean up Docker daemon startup.

Known-good validated ports:

```text
UDP 7777-7921 -> 192.168.68.21
TCP 31992      -> 192.168.68.21
TCP 31983      -> 192.168.68.21
TCP 32673      -> localhost only for Admin RMQ
```

The runbook is stored in:

```text
docs/wsl-native-dune-self-hosting.md
```
