package gostack

import (
	"net/http"
	"sync"
	"time"
)

// ISRPage caches full HTML responses per key with TTL (best-effort in-memory).
type ISRPage struct {
	mu   sync.RWMutex
	data map[string]isrEntry
}

type isrEntry struct {
	body []byte
	expiresAt time.Time
}

// NewISR returns a new ISR cache.
func NewISR() *ISRPage {
	return &ISRPage{data: make(map[string]isrEntry)}
}

// Render writes cached HTML if fresh; otherwise calls render and stores.
func (i *ISRPage) Render(w http.ResponseWriter, key string, ttl time.Duration, render func() ([]byte, error)) error {
	i.mu.RLock()
	e, ok := i.data[key]
	i.mu.RUnlock()
	if ok && time.Now().Before(e.expiresAt) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, err := w.Write(e.body)
		return err
	}
	b, err := render()
	if err != nil {
		return err
	}
	i.mu.Lock()
	i.data[key] = isrEntry{body: b, expiresAt: time.Now().Add(ttl)}
	i.mu.Unlock()
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, err = w.Write(b)
	return err
}
