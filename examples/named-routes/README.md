# named-routes

`GETNamed` plus `App.RouteURL` for link generation.

```bash
cd examples/named-routes
PORT=3000 go run .
curl -s localhost:3000/_routes/resolve
curl -s localhost:3000/users/42
```
