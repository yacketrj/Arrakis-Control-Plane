# Go Project Layout

Target layout:

```text
cmd/dune-admin/       application entrypoint
internal/config/      configuration loading and setup output
internal/sshx/        SSH client and remote command helpers
internal/tunnel/      managed SSH tunnel lifecycle and status
internal/db/          database connection and SQL helpers
internal/httpapi/     HTTP server, auth middleware, and route registration
internal/game/        Dune-specific player, inventory, battlegroup, and runtime logic
web/                  React frontend
docs/                 operator and developer documentation
```

## Migration order

1. Move current root Go files into `cmd/dune-admin` without changing package names.
2. Update build commands from `go build .` to `go build ./cmd/dune-admin`.
3. Extract low-risk support code into `internal/config`, `internal/sshx`, and `internal/tunnel`.
4. Extract database helpers into `internal/db`.
5. Extract HTTP route registration and middleware into `internal/httpapi`.
6. Extract Dune-specific player, inventory, battlegroup, and runtime logic into `internal/game`.

## Validation after each slice

```powershell
gofmt -w .
go test ./...
go build -o dune-admin.exe ./cmd/dune-admin
cd web
npm run lint
npm run build
```

The migration should be done in small compile-clean slices. Avoid creating separate Windows and Linux source trees. Use Go build tags only for files that truly require OS-specific implementations.
