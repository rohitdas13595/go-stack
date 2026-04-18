package gostack

import (
	"github.com/gorilla/websocket"
	"github.com/rohitdas13595/go-stack/sse"
	"github.com/rohitdas13595/go-stack/ws"
)

// SSE streams server-sent events (see sse package).
func (c *Context) SSE(fn func(send sse.Sender)) error {
	return sse.Handler(c.w, fn)
}

// WebSocket upgrades the connection and runs fn with *websocket.Conn.
func (c *Context) WebSocket(fn func(conn *websocket.Conn)) error {
	conn, err := ws.Upgrade(c.w, c.r, nil)
	if err != nil {
		return err
	}
	defer conn.Close()
	fn(conn)
	return nil
}
