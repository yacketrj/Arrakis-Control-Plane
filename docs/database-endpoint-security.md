# Database Endpoint Security Notes

## Purpose

This note tracks the AppSec review state for DA Manager database endpoints under `ASEA-004`.

The database tab is an administrative/high-risk surface because it can expose schema metadata, sampled game data, function metadata, and manual read-only SQL output. These endpoints must remain admin-only and must keep strict input validation, read-only SQL controls, result limits, and output redaction.

## Reviewed endpoints

| Method | Path | Handler | Current security posture |
|---|---|---|---|
| `GET` | `/api/v1/database/tables` | `handleDBTables` | Admin-only route; returns table names and estimated row counts. |
| `GET` | `/api/v1/database/describe` | `handleDBDescribe` | Requires bounded `table` query parameter; command uses parameterized information-schema lookup. |
| `GET` | `/api/v1/database/sample` | `handleDBSample` | Requires bounded `table` query parameter; clamps `limit` to `1..200`; command uses `pgx.Identifier` for schema/table quoting; response rows are redacted before return. |
| `GET` | `/api/v1/database/search` | `handleDBSearch` | Requires bounded `term` query parameter; command uses parameterized information-schema search; response rows are redacted before return. |
| `GET` | `/api/v1/database/functions` | `handleDBFunctions` | Bounds optional `term` and `category` query parameters. |
| `GET` | `/api/v1/database/functions/inspect` | `handleDBFunctionInspect` | Requires bounded numeric `oid` query parameter. |
| `POST` | `/api/v1/database/sql` | `handleDBSQL` | Request body is limited to `maxJSONBodyBytes`; SQL is trimmed; only single-statement read-only SQL passes `isReadOnlySQL`; result text is redacted before return. |

## Current guardrails

- Database endpoints remain protected by normal backend authentication.
- Database mutation risk is classified as high for `POST /api/v1/database/sql`.
- Query-string inputs are trimmed, bounded to 128 characters, and rejected when unsafe control characters are present.
- Function inspection now requires a numeric OID.
- Manual SQL still rejects semicolons and mutation-oriented keywords through `isReadOnlySQL`.
- Manual SQL output and sampled/search rows are passed through `RedactSensitiveText` before being returned.
- Sample row count remains clamped to a maximum of 200 rows.
- Manual SQL command output remains limited to 200 rows.

## Added regression tests

`handlers_database_test.go` covers:

- overlong database query parameters
- control-character query parameters
- trimmed query parameter behavior
- database row redaction helper behavior
- non-numeric function OID rejection
- unsafe SQL rejection before database access
- trimmed unsafe SQL rejection
- blank SQL rejection
- overlong search-term rejection before database access
- redacted SQL response payload shape

## Remaining ASEA-004 work

This is partial remediation. Before `ASEA-004` can be closed, finish:

- database handler-by-handler abuse-case review against a local/dev instance
- SQL timeout review for manual query execution
- expanded `isReadOnlySQL` bypass tests for comments, CTEs, nested statements, Unicode/control input, and long-running queries
- result redaction verification with representative real database rows
- manual verification that database errors do not leak credentials, tunnel details, or internal hostnames
- DAST/SAST/dependency evidence for this surface

## Validation

Required from the canonical local update path:

```bash
./update.sh
```
