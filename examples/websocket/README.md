# websocket

In-process pub/sub with `ws.Hub` and `gorilla/websocket`.

```bash
cd examples/websocket
PORT=3000 go run ./cmd/server
```

Connect to `ws://localhost:3000/ws?channel=demo`. Text messages are broadcast to every connection on the same channel.
