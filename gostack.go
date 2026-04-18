package gostack

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/rohitdas13595/go-stack/router"
)

// Middleware is HTTP middleware.
type Middleware func(http.Handler) http.Handler

// HandlerFunc is a GoStack HTTP handler returning error.
type HandlerFunc func(*Context) error

// App is the HTTP application kernel.
type App struct {
	mux          *router.Router
	chain        []Middleware
	routeNames   map[string]string // name -> "METHOD path"
	routes       []router.Route
	mu           sync.RWMutex
	shutdown     []func(context.Context) error
	renderer     *RenderEngine
	validator    *validateWrapper
	policyReg    *policyRegistry
}

// New creates an application with defaults.
func New() *App {
	return &App{
		mux:        router.New(),
		routeNames: make(map[string]string),
		validator:  newValidateWrapper(),
		policyReg:  newPolicyRegistry(),
	}
}

// Use appends global middleware.
func (a *App) Use(mw ...Middleware) {
	a.chain = append(a.chain, mw...)
}

// SetRenderer sets the template engine (optional until first Render).
func (a *App) SetRenderer(r *RenderEngine) {
	a.renderer = r
}

// Renderer returns the render engine if set.
func (a *App) Renderer() *RenderEngine {
	return a.renderer
}

func (a *App) addShutdown(fn func(context.Context) error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.shutdown = append(a.shutdown, fn)
}

// Shutdown runs registered cleanup hooks.
func (a *App) Shutdown(ctx context.Context) error {
	a.mu.RLock()
	defer a.mu.RUnlock()
	var errs []error
	for _, fn := range a.shutdown {
		if err := fn(ctx); err != nil {
			errs = append(errs, err)
		}
	}
	return joinErrors(errs)
}

func joinErrors(errs []error) error {
	if len(errs) == 0 {
		return nil
	}
	var b strings.Builder
	for i, e := range errs {
		if i > 0 {
			b.WriteString("; ")
		}
		b.WriteString(e.Error())
	}
	return fmt.Errorf("%s", b.String())
}

// GET registers a GET route.
func (a *App) GET(path string, h HandlerFunc) {
	a.handle(http.MethodGet, path, h)
}

// POST registers a POST route.
func (a *App) POST(path string, h HandlerFunc) {
	a.handle(http.MethodPost, path, h)
}

// PUT registers a PUT route.
func (a *App) PUT(path string, h HandlerFunc) {
	a.handle(http.MethodPut, path, h)
}

// PATCH registers a PATCH route.
func (a *App) PATCH(path string, h HandlerFunc) {
	a.handle(http.MethodPatch, path, h)
}

// DELETE registers a DELETE route.
func (a *App) DELETE(path string, h HandlerFunc) {
	a.handle(http.MethodDelete, path, h)
}

// Handle registers any method.
func (a *App) Handle(method, path string, h HandlerFunc) {
	a.handle(method, path, h)
}

// Name names the last registered route pattern (call chained after registration).
// For simplicity we store names on the next handle via RouteName option pattern —
// here we use a different API: GET(...).Name("x") would need wrapper.
// PRD: app.GET("/profile", h).Name("user.profile")
// We implement: a.GETNamed("user.profile", "/profile", h)

// GETNamed registers GET with a route name.
func (a *App) GETNamed(name, path string, h HandlerFunc) {
	a.routeNames[name] = http.MethodGet + " " + path
	a.GET(path, h)
}

func (a *App) handle(method, path string, h HandlerFunc) {
	fullPattern := path
	a.mux.Handle(method, path, func(w http.ResponseWriter, r *http.Request) {
		ctx := newContext(a, w, r)
		ctx.routePattern = method + " " + fullPattern
		if err := h(ctx); err != nil {
			_ = ctx.Error(http.StatusInternalServerError, err.Error())
		}
	})
	a.routes = append(a.routes, router.Route{Method: method, Pattern: path})
}

// Group returns a route group with prefix and middleware.
func (a *App) Group(prefix string, mws ...Middleware) *RouteGroup {
	return &RouteGroup{app: a, prefix: strings.TrimSuffix(prefix, "/"), mws: mws}
}

// RouteGroup is a prefixed route group.
type RouteGroup struct {
	app    *App
	prefix string
	mws    []Middleware
}

func (g *RouteGroup) join(path string) string {
	if g.prefix == "" {
		return path
	}
	if path == "" || path == "/" {
		return g.prefix
	}
	return g.prefix + "/" + strings.TrimPrefix(path, "/")
}

func (g *RouteGroup) GET(path string, h HandlerFunc) {
	g.app.handleGroup(g.join(path), http.MethodGet, h, g.mws)
}

func (g *RouteGroup) POST(path string, h HandlerFunc) {
	g.app.handleGroup(g.join(path), http.MethodPost, h, g.mws)
}

func (g *RouteGroup) PUT(path string, h HandlerFunc) {
	g.app.handleGroup(g.join(path), http.MethodPut, h, g.mws)
}

func (g *RouteGroup) DELETE(path string, h HandlerFunc) {
	g.app.handleGroup(g.join(path), http.MethodDelete, h, g.mws)
}

func (g *RouteGroup) PATCH(path string, h HandlerFunc) {
	g.app.handleGroup(g.join(path), http.MethodPatch, h, g.mws)
}

func (g *RouteGroup) Group(prefix string, mws ...Middleware) *RouteGroup {
	p := g.join(prefix)
	return &RouteGroup{app: g.app, prefix: p, mws: append(append([]Middleware{}, g.mws...), mws...)}
}

func (a *App) handleGroup(path, method string, h HandlerFunc, groupMW []Middleware) {
	a.mux.Handle(method, path, func(w http.ResponseWriter, r *http.Request) {
		handler := http.Handler(http.HandlerFunc(func(w2 http.ResponseWriter, r2 *http.Request) {
			ctx := newContext(a, w2, r2)
			ctx.routePattern = method + " " + path
			if err := h(ctx); err != nil {
				_ = ctx.Error(http.StatusInternalServerError, err.Error())
			}
		}))
		for i := len(groupMW) - 1; i >= 0; i-- {
			handler = groupMW[i](handler)
		}
		handler.ServeHTTP(w, r)
	})
	a.routes = append(a.routes, router.Route{Method: method, Pattern: path})
}

// Resource registers REST routes for a resource prefix using a ResourceHandler.
func (a *App) Resource(prefix string, rh ResourceHandler) {
	p := strings.TrimSuffix(prefix, "/")
	a.GET(p, rh.Index)
	a.GET(p+"/new", rh.New)
	a.POST(p, rh.Create)
	a.GET(p+"/:id", rh.Show)
	a.GET(p+"/:id/edit", rh.Edit)
	a.PUT(p+"/:id", rh.Update)
	a.PATCH(p+"/:id", rh.Patch)
	a.DELETE(p+"/:id", rh.Destroy)
}

// ResourceHandler supplies CRUD handlers; optional methods can be no-op returning 404.
type ResourceHandler struct {
	Index   HandlerFunc
	New     HandlerFunc
	Create  HandlerFunc
	Show    HandlerFunc
	Edit    HandlerFunc
	Update  HandlerFunc
	Patch   HandlerFunc
	Destroy HandlerFunc
}

// RouteURL resolves a named route to a path (no query). Params are substituted for :param segments.
func (a *App) RouteURL(name string, params map[string]string) (string, error) {
	a.mu.RLock()
	key, ok := a.routeNames[name]
	a.mu.RUnlock()
	if !ok {
		return "", fmt.Errorf("gostack: unknown route %q", name)
	}
	parts := strings.SplitN(key, " ", 2)
	if len(parts) != 2 {
		return "", fmt.Errorf("gostack: invalid route key")
	}
	path := parts[1]
	for k, v := range params {
		path = strings.ReplaceAll(path, ":"+k, v)
	}
	if strings.Contains(path, ":") {
		return "", fmt.Errorf("gostack: missing params for route %q", name)
	}
	return path, nil
}

// Routes returns registered routes (method + pattern).
func (a *App) Routes() []router.Route {
	out := make([]router.Route, len(a.routes))
	copy(out, a.routes)
	return out
}

// ServeHTTP implements http.Handler.
func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var h http.Handler = a.mux
	for i := len(a.chain) - 1; i >= 0; i-- {
		h = a.chain[i](h)
	}
	h.ServeHTTP(w, r)
}

// ListenAndServe runs the app on addr with optional std http.Server overrides.
func (a *App) ListenAndServe(addr string) error {
	srv := &http.Server{
		Addr:              addr,
		Handler:           a,
		ReadHeaderTimeout: 10 * time.Second,
	}
	return srv.ListenAndServe()
}

// ListenAndServeContext runs until ctx is cancelled, then shuts down gracefully.
func (a *App) ListenAndServeContext(ctx context.Context, addr string) error {
	srv := &http.Server{
		Addr:              addr,
		Handler:           a,
		ReadHeaderTimeout: 10 * time.Second,
	}
	errCh := make(chan error, 1)
	go func() {
		errCh <- srv.ListenAndServe()
	}()
	select {
	case <-ctx.Done():
		shCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		_ = srv.Shutdown(shCtx)
_ = a.Shutdown(shCtx)
		return ctx.Err()
	case err := <-errCh:
		if err == http.ErrServerClosed {
			return nil
		}
		return err
	}
}

// Env returns an environment variable with optional default.
func Env(key string, def ...string) string {
	v := os.Getenv(key)
	if v == "" && len(def) > 0 {
		return def[0]
	}
	return v
}
