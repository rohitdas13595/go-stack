# GoStack

Opinionated full-stack Go framework (see [PRD.md](PRD.md)). Module: `github.com/rohitdas13595/go-stack`.

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

## Packages

- Core: `App`, `Context`, router, middleware, SSR (`html/template`), config, DB manager, SQL migrations, minimal ORM, JWT auth, SSE/WebSocket helpers, ISR cache stub.
- Integrations: jobs (Redis list), cache (memory/Redis), storage (disk/S3), mail (SMTP), AI (OpenAI-compatible HTTP), observability (Prometheus, OTLP tracing), OpenAPI stub, OAuth/WebAuthn stubs.

## License

See repository license (if any).
