// Package router provides a radix-style HTTP router with method + path matching,
// named parameters, and catch-all segments.
package router

import (
	"context"
	"net/http"
	"strings"
)

type ctxKey int

const paramsKey ctxKey = iota

// ParamsFromContext returns path parameters set by the router.
func ParamsFromContext(ctx context.Context) map[string]string {
	if ctx == nil {
		return nil
	}
	v := ctx.Value(paramsKey)
	if v == nil {
		return nil
	}
	return v.(map[string]string)
}

// WithParams attaches path parameters to the context.
func WithParams(ctx context.Context, params map[string]string) context.Context {
	return context.WithValue(ctx, paramsKey, params)
}

// Handler handles a matched request.
type Handler func(w http.ResponseWriter, r *http.Request)

// Route describes a registered route (for introspection).
type Route struct {
	Method  string
	Pattern string
	Name    string
}

// Router matches HTTP requests to handlers.
type Router struct {
	root *node
}

// New creates an empty router.
func New() *Router {
	return &Router{root: &node{static: make(map[string]*node), handlers: make(map[string]Handler)}}
}

// Handle registers method + pattern → handler.
func (rt *Router) Handle(method, pattern string, h Handler) {
	method = strings.ToUpper(strings.TrimSpace(method))
	segs := splitPattern(pattern)
	rt.root.insert(method, segs, 0, h)
}

// ServeHTTP dispatches the request.
func (rt *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	method := strings.ToUpper(r.Method)
	segs := splitPath(r.URL.Path)
	params := make(map[string]string)
	n := rt.root.lookup(method, segs, 0, params)
	if n == nil {
		http.NotFound(w, r)
		return
	}
	h := n.handlers[method]
	if h == nil {
		http.NotFound(w, r)
		return
	}
	h(w, r.WithContext(WithParams(r.Context(), params)))
}

func splitPattern(p string) []string {
	p = strings.TrimSpace(p)
	if p == "" || p == "/" {
		return nil
	}
	p = strings.TrimPrefix(p, "/")
	if p == "" {
		return nil
	}
	return strings.Split(p, "/")
}

func splitPath(p string) []string {
	p = strings.TrimSpace(p)
	if p == "" || p == "/" {
		return nil
	}
	p = strings.TrimPrefix(p, "/")
	if p == "" {
		return nil
	}
	return strings.Split(p, "/")
}

type node struct {
	static map[string]*node
	paramChild *node
	paramName  string
	wildChild  *node
	wildName   string
	handlers   map[string]Handler
}

func (n *node) insert(method string, segs []string, i int, h Handler) {
	if n.handlers == nil {
		n.handlers = make(map[string]Handler)
	}
	if i >= len(segs) {
		n.handlers[method] = h
		return
	}
	s := segs[i]
	if strings.HasPrefix(s, "*") {
		if n.wildChild == nil {
			n.wildChild = &node{static: make(map[string]*node), handlers: make(map[string]Handler)}
		}
		n.wildChild.wildName = strings.TrimPrefix(s, "*")
		if n.wildChild.wildName == "" {
			n.wildChild.wildName = "path"
		}
		n.wildChild.handlers[method] = h
		return
	}
	if strings.HasPrefix(s, ":") {
		if n.paramChild == nil {
			n.paramChild = &node{static: make(map[string]*node), handlers: make(map[string]Handler)}
		}
		n.paramChild.paramName = strings.TrimPrefix(s, ":")
		n.paramChild.insert(method, segs, i+1, h)
		return
	}
	child := n.static[s]
	if child == nil {
		child = &node{static: make(map[string]*node), handlers: make(map[string]Handler)}
		n.static[s] = child
	}
	child.insert(method, segs, i+1, h)
}

func (n *node) lookup(method string, segs []string, i int, params map[string]string) *node {
	if n == nil {
		return nil
	}
	if i >= len(segs) {
		if n.handlers[method] != nil {
			return n
		}
		return nil
	}
	if child := n.static[segs[i]]; child != nil {
		if found := child.lookup(method, segs, i+1, params); found != nil {
			return found
		}
	}
	if n.paramChild != nil {
		params[n.paramChild.paramName] = segs[i]
		if found := n.paramChild.lookup(method, segs, i+1, params); found != nil {
			return found
		}
		delete(params, n.paramChild.paramName)
	}
	if n.wildChild != nil && n.wildChild.handlers[method] != nil {
		params[n.wildChild.wildName] = strings.Join(segs[i:], "/")
		return n.wildChild
	}
	return nil
}
