package gostack

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/rohitdas13595/go-stack/auth"
	"github.com/rohitdas13595/go-stack/router"
)

// Context wraps the HTTP request/response for handlers.
type Context struct {
	app          *App
	w            http.ResponseWriter
	r            *http.Request
	store        map[string]any
	routePattern string
	validator    *validator.Validate
}

func newContext(a *App, w http.ResponseWriter, r *http.Request) *Context {
	c := &Context{
		app:       a,
		w:         w,
		r:         r,
		store:     make(map[string]any),
		validator: a.validator.v,
	}
	if u := auth.UserFromContext(r.Context()); u != nil {
		c.store["user"] = u
	}
	return c
}

func (c *Context) cloneForRequest(w http.ResponseWriter, r *http.Request) *Context {
	nc := newContext(c.app, w, r)
	m2 := make(map[string]any, len(c.store))
	for k, v := range c.store {
		m2[k] = v
	}
	nc.store = m2
	return nc
}

// Request returns the underlying HTTP request.
func (c *Context) Request() *http.Request {
	return c.r
}

// ResponseWriter returns the response writer.
func (c *Context) ResponseWriter() http.ResponseWriter {
	return c.w
}

// Context returns the request context.
func (c *Context) Context() context.Context {
	return c.r.Context()
}

// Param returns a path parameter.
func (c *Context) Param(name string) string {
	m := router.ParamsFromContext(c.r.Context())
	if m == nil {
		return ""
	}
	return m[name]
}

// Query returns a query value.
func (c *Context) Query(key string) string {
	return c.r.URL.Query().Get(key)
}

// QueryInt returns query as int with default.
func (c *Context) QueryInt(key string, def int) int {
	s := c.r.URL.Query().Get(key)
	if s == "" {
		return def
	}
	n, err := strconv.Atoi(s)
	if err != nil {
		return def
	}
	return n
}

// Bind decodes JSON into dst and validates struct tags.
func (c *Context) Bind(dst any) error {
	ct := c.r.Header.Get("Content-Type")
	if !strings.HasPrefix(ct, "application/json") {
		return fmt.Errorf("gostack: Bind requires Content-Type application/json")
	}
	dec := json.NewDecoder(c.r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(dst); err != nil && !errors.Is(err, io.EOF) {
		return err
	}
	if c.validator != nil {
		return c.validator.Struct(dst)
	}
	return nil
}

// JSON writes JSON response.
func (c *Context) JSON(status int, v any) error {
	c.w.Header().Set("Content-Type", "application/json; charset=utf-8")
	c.w.WriteHeader(status)
	enc := json.NewEncoder(c.w)
	enc.SetEscapeHTML(true)
	return enc.Encode(v)
}

// Redirect sends redirect response.
func (c *Context) Redirect(status int, url string) error {
	http.Redirect(c.w, c.r, url, status)
	return nil
}

// BadRequest sends 400 with message.
func (c *Context) BadRequest(err error) error {
	return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
}

// Forbidden sends 403.
func (c *Context) Forbidden(msg string) error {
	return c.JSON(http.StatusForbidden, map[string]string{"error": msg})
}

// Error sends generic error JSON.
func (c *Context) Error(status int, msg string) error {
	return c.JSON(status, map[string]string{"error": msg})
}

// Set stores a value in request-local storage.
func (c *Context) Set(key string, val any) {
	c.store[key] = val
}

// Get retrieves from request-local storage.
func (c *Context) Get(key string) (any, bool) {
	v, ok := c.store[key]
	return v, ok
}

// User returns the authenticated user set by auth middleware.
func (c *Context) User() any {
	if u, ok := c.Get("user"); ok {
		return u
	}
	return nil
}

// Authorize checks a policy action on resource (Phase 4).
func (c *Context) Authorize(action string, resource any) error {
	return c.app.policyReg.authorize(c, action, resource)
}

// Render renders a template via the app's RenderEngine.
func (c *Context) Render(name string, data Data) error {
	if c.app.renderer == nil {
		return fmt.Errorf("gostack: no renderer configured")
	}
	return c.app.renderer.Render(c.w, name, data)
}

// RenderPartial renders a partial template (HTMX).
func (c *Context) RenderPartial(name string, data Data) error {
	if c.app.renderer == nil {
		return fmt.Errorf("gostack: no renderer configured")
	}
	return c.app.renderer.RenderPartial(c.w, name, data)
}

// Data is template data map.
type Data map[string]any
