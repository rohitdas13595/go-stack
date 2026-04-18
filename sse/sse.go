package sse

import (
	"fmt"
	"net/http"
	"strings"
)

// Sender sends named SSE events.
type Sender func(event string, data any)

// Handler writes SSE stream. Call send repeatedly; data is stringified with fmt.Sprint.
func Handler(w http.ResponseWriter, fn func(send Sender)) error {
	h := w.Header()
	h.Set("Content-Type", "text/event-stream")
	h.Set("Cache-Control", "no-cache")
	h.Set("Connection", "keep-alive")
	send := func(event string, data any) {
		if event != "" {
			fmt.Fprintf(w, "event: %s\n", event)
		}
		s := fmt.Sprint(data)
		for _, line := range strings.Split(s, "\n") {
			fmt.Fprintf(w, "data: %s\n", line)
		}
		fmt.Fprint(w, "\n")
		if f, ok := w.(http.Flusher); ok {
			f.Flush()
		}
	}
	fn(send)
	return nil
}
