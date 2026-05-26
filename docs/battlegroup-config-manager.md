# Battlegroup Config Manager

## Feature request

Area: Battlegroup

Feature: Ability to inspect and modify `UserEngine.ini` and `UserGame.ini`.

## Risk and impact

This is a high-risk configuration-management feature. These files can affect how the game is played, items gathered, progression balance, server behavior, and the in-game economy.

Risk categories:

- gameplay integrity
- player economy integrity
- server availability
- configuration drift
- accidental destructive edits
- rollback difficulty
- operator accountability

## Current status

```text
PLANNED / REQUIRES DISCOVERY
```

Do not implement direct editing until the discovery and safety model is complete.

## Unknowns to resolve

- Which container or runtime unit owns `UserEngine.ini` and `UserGame.ini`.
- Exact file paths for Kubernetes deployment.
- Exact file paths for Docker deployment.
- Whether AMP-managed Docker deployments mount these files into persistent host paths.
- Whether edits require container restart, battlegroup restart, or full server restart.
- Whether the files are overwritten during server updates.
- Whether changes should be applied per shard, per battlegroup, or globally.
- Whether runtime-generated config files differ from source/default config files.

## Required discovery behavior

The application must discover:

1. active runtime mode: Kubernetes, Docker, AMP-managed Docker, or unknown
2. battlegroup/container/pod identity
3. host-side path if bind-mounted
4. container-side path
5. file existence and readability
6. file ownership and permissions
7. whether write requires elevated privileges
8. restart requirement after write

Discovery must be read-only.

## Required safety controls

Before allowing edits, the feature must include:

- read-only inspect mode first
- explicit target file selection
- immutable before-change backup
- normalized diff preview
- validation of INI syntax before write
- allowlist or warning model for high-risk settings
- shared mutation confirmation
- required admin reason capture
- audit log entry
- post-write verification readback
- rollback to backup
- clear restart/apply instructions

## Required UI model

Recommended UI sections:

1. Runtime discovery
2. Config file inventory
3. Read-only file viewer
4. Diff editor
5. Validation results
6. Backup and restore panel
7. Apply/restart guidance
8. Audit trail

## Initial implementation phases

### Phase 1: Read-only discovery

- Locate candidate `UserEngine.ini` and `UserGame.ini` files.
- Display runtime, container/pod, and resolved path metadata.
- Read file contents safely.
- Do not write.

### Phase 2: Backup and diff preview

- Create backup artifact before proposed edit.
- Show before/after diff.
- Validate INI syntax.
- Still do not write by default.

### Phase 3: Confirmed write

- Write only after shared mutation confirmation and admin reason capture.
- Audit the change.
- Verify post-write contents.
- Provide restart/apply guidance.

### Phase 4: Rollback

- Allow restoring from a recorded backup.
- Require confirmation and reason.
- Audit rollback.

## Security review requirements

This feature must be reviewed under `docs/production-security-review.md`, specifically:

- Battlegroup command/config controls
- mutation safety and auditability
- logging and data leakage
- SSH tunnel/startup behavior
- secret handling if elevated credentials are needed

## Priority

Recommended priority: `P2` until production security gates and runtime discovery are stable.

Escalate to `P1` if server operators identify config editing as required for normal supported operations.