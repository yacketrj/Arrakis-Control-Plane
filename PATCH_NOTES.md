# Arrakis Control Panel Release Notes

## Current update: Runtime-aware status handling

### Why this update was made

A Docker deployment reached the web UI, but the Battlegroup status view was still parsing and labeling status as Kubernetes pods.

### What changed

- Docker database discovery now prefers the published Docker port mapping.
- The Battlegroup tab now reads the backend runtime from status responses.
- Docker status output renders as containers.
- Kubernetes status output renders as pods.
- Docker runtime keeps battlegroup script controls disabled until that command path is Docker-safe.
- Durable detail is archived in:

```text
docs/changelog/unreleased/2026-06-16-runtime-aware-status.md
```

### Impact

- Runtime display behavior is clearer for both Docker and Kubernetes deployments.
- Runtime server behavior outside status and discovery did not change.

### Validation

Run:

```bash
./update.sh
```

On Windows:

```powershell
.\update.ps1 -SkipAutoPush
```

### Remaining final-`v0.1.0` gate

- Post-release verification after tag/artifact install or launch.
