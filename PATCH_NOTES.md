# Arrakis Control Panel Release Notes

## Current update: Runtime evidence audit

### Why this update was made

The web UI could still show Kubernetes status on a host where the active Dune workload was running in Docker.

### What changed

- Runtime auto-detection now checks active Dune Docker workload evidence before Kubernetes discovery.
- Kubernetes auto-detection now requires Kubernetes workload evidence before choosing that path.
- Status payload generation refreshes runtime from active workload evidence when SSH is connected.
- Durable detail is archived in:

```text
docs/changelog/unreleased/2026-06-16-runtime-evidence-audit.md
```

### Impact

- Runtime selection is now based on active Dune workloads instead of tool availability alone.
- Battlegroup status should report Docker when Dune Docker containers are present.

### Validation

Run the normal update script.

### Remaining final-v0.1.0 gate

- Post-release verification after tag/artifact install or launch.
