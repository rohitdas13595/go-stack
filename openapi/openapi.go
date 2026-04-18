package openapi

import (
	"encoding/json"
	"net/http"
	"strings"
)

// Spec is a minimal OpenAPI 3.0 document.
type Spec struct {
	OpenAPI string         `json:"openapi"`
	Info    map[string]any `json:"info"`
	Paths   map[string]any `json:"paths"`
}

// Handler serves a static spec as JSON.
func Handler(spec *Spec) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(spec)
	})
}

// FromRoutes builds a stub spec from method+path list.
func FromRoutes(routes [][2]string) *Spec {
	paths := map[string]any{}
	for _, rp := range routes {
		method, path := rp[0], rp[1]
		entry, _ := paths[path].(map[string]any)
		if entry == nil {
			entry = map[string]any{}
		}
		entry[strings.ToLower(method)] = map[string]any{"summary": method + " " + path}
		paths[path] = entry
	}
	return &Spec{
		OpenAPI: "3.0.3",
		Info: map[string]any{
			"title":   "GoStack API",
			"version": "0.1.0",
		},
		Paths: paths,
	}
}
