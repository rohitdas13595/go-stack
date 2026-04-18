# ssr

Server-rendered HTML with `html/template`.

Run from **this directory** (`examples/ssr`) so `views/` resolves correctly:

```bash
cd examples/ssr
PORT=3000 go run ./cmd/server
```

Open http://localhost:3000/ — use the button to call `/health` via HTMX.
