# Arrakis Control Panel Release Notes

## Current update: Setup writes .env before SSH validation

### Why this update was made

The setup wizard was collecting configuration but only wrote `.env` after SSH succeeded. If SSH authentication failed, the exe exited before creating `.env`.

### What changed

- `runSetup()` now writes a preliminary `.env` before the SSH dial.
- If SSH fails, the operator can edit `.env` and rerun setup instead of re-entering every prompt.
- Setup still rewrites `.env` after successful SSH/runtime/database discovery.
- Durable detail is archived in:

```text
docs/changelog/unreleased/2026-06-16-setup-env-before-ssh.md
```

### Impact

- Setup behavior changed.
- Runtime server behavior outside setup did not change.
- No endpoint or mutation behavior was added.

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
