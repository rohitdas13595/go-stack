# Development notes

Internal pointers for people hacking on the GoStack repository.

## Layout

| Path | Role |
|------|------|
| `router/` | HTTP router (method + path, params, wildcards) |
| `middleware/` | Cross-cutting HTTP middleware |
| `gostack.go`, `context.go` | `App`, `Context`, route groups |
| `cmd/gostack`, `internal/cli/` | CLI (`new`, `serve`, `db`, …) |
| `examples/` | Standalone modules; each replaces the framework with `../..` |
| `PRD.md` | Product requirements / design reference |

## Commands

```bash
go test ./...
go vet ./...
go fmt ./...
```

Vulnerability scan (also run in CI):

```bash
go install golang.org/x/vuln/cmd/govulncheck@latest
govulncheck ./...
```

## Workspace

The root `go.work` includes the main module and all `examples/*` modules so
local `replace` directives resolve without installing the module from the network.

## Releases

Versioning and tagging are project-specific. When cutting a release, update
downstream docs and ensure `go test ./...` and example builds pass on the tagged
commit.
