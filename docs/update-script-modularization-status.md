# Update-Script Modularization Status

## Status

Complete for final `v0.1.0` readiness, subject to the standard release validation path.

## Scope reviewed

### Bash update path

The Bash entrypoint is now a thin coordinator and sources helper modules from `scripts/update/`:

- `defaults.sh`
- `ui.sh`
- `prereqs.sh`
- `git.sh`
- `npm.sh`
- `backend.sh`
- `web.sh`

The Bash entrypoint keeps argument parsing and orchestration in `update.sh`; reusable behavior lives in helper modules.

### PowerShell update path

The PowerShell entrypoint is now split across helper modules from `scripts/update/`:

- `powershell-common.ps1`
- `powershell-git.ps1`
- `powershell-npm.ps1`
- `powershell-backend.ps1`
- `powershell-web.ps1`

The remaining logic in `update.ps1` is limited to parameter handling, helper loading, high-level orchestration, error handling, and final status output.

## Completed final-release corrections

- Extracted PowerShell backend Go package discovery and test execution.
- Extracted PowerShell backend build execution.
- Extracted PowerShell backend binary and asset copy behavior.
- Extracted PowerShell web folder/package checks.
- Extracted PowerShell Node/npm prerequisite checks and version probes.
- Extracted PowerShell npm install/repair, toolchain validation, audit, typecheck, lint, and web build execution.

## Validation evidence

The final backend and web helper extractions were validated from both update paths:

```bash
./update.sh
```

```powershell
.\update.ps1 -SkipAutoPush
```

## Release decision

The update-script modularization gate is closed for final `v0.1.0` readiness. No deferral is required for this gate.

## Remaining release gates outside this scope

- Go code-quality/refactor review or explicit deferral.
- Full documentation review beyond primary release/security docs, or explicit deferral.
- Post-release verification after tag/artifact install or launch.
