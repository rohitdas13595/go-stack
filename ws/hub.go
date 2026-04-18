package ws

import (
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// Hub tracks connections per channel (in-process). Use Redis bridge in production.
type Hub struct {
	mu    sync.RWMutex
	chans map[string]map[*websocket.Conn]struct{}
}

// NewHub creates a hub.
func NewHub() *Hub {
	return &Hub{chans: make(map[string]map[*websocket.Conn]struct{})}
}

// Upgrade upgrades HTTP to WebSocket.
func Upgrade(w http.ResponseWriter, r *http.Request, responseHeader http.Header) (*websocket.Conn, error) {
	return upgrader.Upgrade(w, r, responseHeader)
}

// Join adds connection to a channel.
func (h *Hub) Join(channel string, c *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.chans[channel] == nil {
		h.chans[channel] = make(map[*websocket.Conn]struct{})
	}
	h.chans[channel][c] = struct{}{}
}

// Emit sends a text message to all connections in channel.
func (h *Hub) Emit(channel string, msg []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	for c := range h.chans[channel] {
_ = c.WriteMessage(websocket.TextMessage, msg)
	}
}

// Leave removes connection from all channels.
func (h *Hub) Leave(c *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	for ch, m := range h.chans {
		delete(m, c)
		if len(m) == 0 {
			delete(h.chans, ch)
		}
	}
}
