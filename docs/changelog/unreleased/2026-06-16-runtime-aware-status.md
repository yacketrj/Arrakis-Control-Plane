# Runtime-aware status handling

## Date

2026-06-16

## Scope

Bug fix for Docker and Kubernetes runtime handling.

## Change

- Docker database discovery now prefers the published Docker port mapping.
- The Battlegroup UI now reads the runtime returned by the backend.
- Docker status output is rendered as containers.
- Kubernetes status output is rendered as pods.
- Docker runtime keeps battlegroup script controls disabled because that command path is not Docker-safe yet.

## Validation

Run:

```bash
./update.sh
```

On Windows:

```powershell
.\update.ps1 -SkipAutoPush
```
