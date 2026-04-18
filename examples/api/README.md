# api

Versioned JSON API with CORS and `Bind` + struct validation.

```bash
cd examples/api
PORT=3000 go run ./cmd/server
```

```bash
curl -s localhost:3000/api/v1/items
curl -s localhost:3000/api/v1/items/42
curl -s -X POST localhost:3000/api/v1/items \
  -H 'Content-Type: application/json' \
  -d '{"name":"new item"}'
```
