package gostack

import (
	"sync"

	"github.com/gorilla/websocket"
	"github.com/rohitdas13595/go-stack/ws"
)

var (
	hubMu sync.RWMutex
	hub   *ws.Hub
)

// SetWebSocketHub sets the global hub used by Broadcast helpers.
func SetWebSocketHub(h *ws.Hub) {
	hubMu.Lock()
	defer hubMu.Unlock()
	hub = h
}

func getHub() *ws.Hub {
	hubMu.RLock()
	defer hubMu.RUnlock()
	if hub == nil {
		return nil
	}
	return hub
}

// BroadcastTo emits a text message on an in-process channel.
func BroadcastTo(channel string, msg []byte) {
	if h := getHub(); h != nil {
		h.Emit(channel, msg)
	}
}

// JoinWebSocket registers conn to channel (typically after upgrade elsewhere).
func JoinWebSocket(channel string, conn *websocket.Conn) {
	if h := getHub(); h != nil {
		h.Join(channel, conn)
	}
}
