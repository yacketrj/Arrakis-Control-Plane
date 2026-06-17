# Runtime evidence audit

## Date

2026-06-16

## Scope

Runtime detection audit and correction.

## Problem

Runtime auto-detection could select Kubernetes first when both runtime tools were present, even when the active Dune workload was running in Docker.

## Change

- Auto runtime discovery now checks Dune Docker workload evidence before Kubernetes discovery.
- Kubernetes auto discovery now also requires Kubernetes workload evidence before choosing that path.
- Status payload generation refreshes runtime from active workload evidence when an SSH connection is available.

## Audit result

- Database discovery was the main problematic decision point.
- Battlegroup status commands are selected from the chosen runtime.
- Log discovery and log streaming are selected from the chosen runtime.
- Connectivity diagnostics reports remote tool availability, but database discovery uses the runtime discovery path.

## Validation

Run:

```bash
./update.sh
```
