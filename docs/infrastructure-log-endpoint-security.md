# Infrastructure and Log Endpoint Security Notes

## Purpose

This note tracks the AppSec review state for DA Manager infrastructure and log endpoints under `ASEA-005`.

The infrastructure/log surface is high risk because it can run server-control commands, discover runtime targets, issue log-stream tickets, stream remote logs, and expose diagnostic output. These endpoints must remain admin-only and must enforce strict command allowlists, target validation, replay-resistant log-stream tickets, TTLs, and output redaction.

## Reviewed endpoints

| Method | Path | Handler | Current security posture |
|---|---|---|---|
| `GET` | `/api/v1/battlegroup/status` | `handleBGStatus` | Admin-only. Validates Kubernetes namespace before command construction and redacts command output before returning it. |
| `GET` | `/api/v1/battlegroup/health` | `handleBGHealth` | Admin-only. Validates Kubernetes namespace before fixed health commands and redacts section output/errors. |
| `POST` | `/api/v1/battlegroup/exec` | `handleBGExec` | Admin mutation/high risk. Uses a normalized allowlist for battlegroup script commands only. Docker runtimes remain unsupported for script commands. |
| `GET` | `/api/v1/battlegroup/pods` | `handleBGPods` | Admin-only. Validates Kubernetes namespace before command construction and redacts returned pod/container list lines. |
| `GET` | `/api/v1/logs/pods` | `handleLogPods` | Admin-only. Validates runtime namespace/targets and redacts Docker display names. |
| `POST` | `/api/v1/logs/stream-ticket` | `handleIssueLogStreamTicket` | Admin mutation/high risk. Issues scoped, one-time, 60-second log-stream tickets after target validation. |
| `GET` | `/api/v1/logs/stream` | `handleLogStream` | WebSocket ticket/high risk. Rejects legacy `ws_token`; validates target; streams redacted log lines. |
| `GET` | `/api/v1/logs/cheats` | `handleGetCheatLog` | Admin-only. Redacts returned cheat-log fields before returning rows. |

## Current guardrails

- Infrastructure and log endpoints remain protected by normal backend authentication.
- Battlegroup exec commands are restricted to the static allowlist: `start`, `stop`, `restart`, `update`, `backup`, and `restore`.
- Battlegroup exec command input is trimmed, lowercased, control-character checked, and rejected when not explicitly allowed.
- Kubernetes runtime command surfaces validate `globalPodNS` before interpolating it into command strings.
- Docker log targets must match either container ID format or the strict Docker target name pattern.
- Kubernetes log targets use runtime target validation before command construction.
- Log-stream tickets are generated from 32 random bytes, base64url encoded, scoped to namespace and target, single-use, and limited to 60 seconds.
- Wrong-target ticket use consumes the ticket and fails.
- Expired tickets are pruned/rejected.
- Remote command output, health output, log lines, cheat-log rows, and relevant errors pass through `RedactSensitiveText` before returning to the browser.

## Added regression tests

`infrastructure_security_test.go` covers:

- battlegroup command normalization and strict allowlist behavior
- command control-character/metacharacter rejection
- Kubernetes namespace validation for runtime command construction
- Docker runtime namespace bypass behavior
- split-and-redact line handling
- Docker/Kubernetes log target rejection for unsafe targets
- log-stream ticket single-use behavior
- log-stream ticket wrong-target behavior
- expired-ticket rejection
- invalid-target ticket issuance rejection
- cheat-log field redaction

## Remaining ASEA-005 work

This is partial remediation. Before `ASEA-005` can be closed, finish:

- handler-level tests around `handleBGStatus`, `handleBGHealth`, `handleBGExec`, `handleBGPods`, `handleLogPods`, `handleLogStream`, and `handleGetCheatLog` where SSH/database dependencies can be stubbed or abstracted safely
- manual validation against Kubernetes and Docker runtime modes
- command timeout review for SSH exec/stream operations
- WebSocket origin behavior review with configured `ALLOWED_ORIGINS`
- replay/TTL verification through a browser/WebSocket client path
- representative real log/output redaction review
- SAST/DAST/dependency evidence for this surface

## Validation

Required from the canonical local update path:

```bash
./update.sh
```
