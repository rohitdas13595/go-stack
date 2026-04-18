# GoStack

Opinionated full-stack Go framework (see [PRD.md](PRD.md)). Module: `github.com/rohitdas13595/go-stack`.

## About

GoStack is a batteries-included Go toolkit for building web applications and APIs on a single codebase. It gives you a small **application kernel** (`App`, `Context`), a **radix-style router**, a **middleware** pipeline, and first-class support for **server-rendered HTML** (`html/template`) alongside **JSON**—so you can ship HTMX- or Alpine-friendly pages without giving up clean HTTP handlers.

The project also bundles practical infrastructure pieces you normally wire by hand: **named DB pools**, **SQL migrations**, a **lightweight ORM-style query layer**, **JWT and session helpers**, **SSE and WebSocket** utilities, and hooks for **cache, jobs, mail, object storage, and observability**. The included **`gostack` CLI** scaffolds apps and runs common tasks (migrate, routes, serve). Everything is designed to stay readable and Go-idiomatic rather than framework-magic heavy.

GoStack is **Apache-2.0 licensed**. Use it for prototypes, internal tools, or full-stack products where you want Go from the database to the browser—without adopting a heavyweight enterprise framework.

## Quick start (framework dev)

```bash
go build -o gostack ./cmd/gostack
./gostack new myapp
cd myapp
# Ensure go.mod replace points at your local framework checkout.
mkdir -p storage
export DATABASE_URL="file:./storage/app.db"
../gostack db migrate
PORT=3000 go run ./cmd/server
```

## Examples

Runnable demos live under [examples/](examples/). Each is a small Go module that replaces the framework with the local checkout (`replace` in `go.mod`). The repo root [go.work](go.work) lists these modules so `go run` and `go build` work without publishing `go-stack`.

## Packages

- Core: `App`, `Context`, router, middleware, SSR (`html/template`), config, DB manager, SQL migrations, minimal ORM, JWT auth, SSE/WebSocket helpers, ISR cache stub.
- Integrations: jobs (Redis list), cache (memory/Redis), storage (disk/S3), mail (SMTP), AI (OpenAI-compatible HTTP), observability (Prometheus, OTLP tracing), OpenAPI stub, OAuth/WebAuthn stubs.

## Community

- [Contributing](CONTRIBUTING.md) — issues, pull requests, and dev workflow
- [Security](SECURITY.md) — **report vulnerabilities privately** (not via public issues)
- [Code of conduct](CODE_OF_CONDUCT.md)
- [Developer notes](docs/DEVELOPMENT.md) — repo layout and tooling
- [Documentation index](docs/README.md)

## License

Licensed under the [Apache License 2.0](LICENSE). See [NOTICE](NOTICE) for attribution.
