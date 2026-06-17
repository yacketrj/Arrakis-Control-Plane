# Setup writes .env before SSH validation

## Date

2026-06-16

## Scope

Bug fix for interactive setup behavior.

## Problem

The setup wizard collected SSH, database, and admin-token configuration, but `.env` was only written after a successful SSH connection.

If SSH authentication failed, setup exited before writing `.env`, leaving the operator without a generated config file to inspect, edit, and retry.

## Change

`runSetup()` now writes a preliminary `.env` immediately after local configuration validation succeeds and before attempting the SSH dial.

The setup wizard still overwrites `.env` later after successful SSH/runtime/database discovery so detected runtime and database settings are preserved.

## Operator impact

If SSH fails during setup, `.env` should now exist and contain the entered values for:

- SSH host
- SSH user
- SSH key path
- known_hosts path
- database settings
- admin token
- listen address and allowed origins

The operator can edit `.env` and rerun setup instead of re-entering every prompt.

## Runtime impact

No runtime server behavior changed outside setup.

## Validation

Run:

```bash
./update.sh
```

On Windows:

```powershell
.\update.ps1 -SkipAutoPush
```
