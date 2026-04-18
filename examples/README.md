# GoStack examples

Small apps that show common patterns. Each folder is its own Go module with a `replace` directive pointing at the framework root (`../..`).

From the **repository root**, the workspace (`go.work`) lists these modules so builds resolve without publishing the framework.

## Run

```bash
# From repo root (uses go.work)
cd examples/hello && PORT=3000 go run .

cd examples/named-routes && PORT=3001 go run .
cd examples/ssr && PORT=3002 go run ./cmd/server
cd examples/api && PORT=3003 go run ./cmd/server
cd examples/sse && PORT=3004 go run ./cmd/server
cd examples/websocket && PORT=3005 go run ./cmd/server
```

Or with workspace disabled from an example directory:

```bash
cd examples/hello && GOWORK=off go run .
```

## Index

| Example | What it shows |
|--------|----------------|
| [hello](hello/) | Minimal app: JSON routes, `Recover` + `Logger` middleware. |
| [named-routes](named-routes/) | `GETNamed` and `RouteURL` for stable path building. |
| [ssr](ssr/) | Server-rendered HTML with `html/template`, layouts, and HTMX/Alpine CDNs. |
| [api](api/) | Versioned `/api/v1` route group, CORS, JSON `Bind` + validation. |
| [sse](sse/) | Server-Sent Events stream via `sse.Handler`. |
| [websocket](websocket/) | `ws.Hub` broadcast channel and `ws.Upgrade`. |
