package testutil

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rohitdas13595/go-stack"
)

// App wraps a test HTTP server around *gostack.App.
type App struct {
	T   *testing.T
	App *gostack.App
	srv *httptest.Server
}

// NewApp starts a test server.
func NewApp(t *testing.T, app *gostack.App) *App {
	t.Helper()
	s := httptest.NewServer(app)
	t.Cleanup(s.Close)
	return &App{T: t, App: app, srv: s}
}

// URL returns base URL.
func (a *App) URL() string { return a.srv.URL }

// POST JSON to path.
func (a *App) POST(path string, body JSON) *Response {
	a.T.Helper()
	b, _ := json.Marshal(body)
	req, err := http.NewRequest(http.MethodPost, a.srv.URL+path, bytes.NewReader(b))
	if err != nil {
		a.T.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	return a.do(req)
}

// GET path.
func (a *App) GET(path string) *Response {
	a.T.Helper()
	req, err := http.NewRequest(http.MethodGet, a.srv.URL+path, nil)
	if err != nil {
		a.T.Fatal(err)
	}
	return a.do(req)
}

func (a *App) do(req *http.Request) *Response {
	a.T.Helper()
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		a.T.Fatal(err)
	}
	return &Response{T: a.T, Response: res}
}

// JSON is a JSON body map.
type JSON map[string]any

// Response wraps http.Response with assertions.
type Response struct {
	T *testing.T
	*http.Response
}

// AssertStatus checks status code.
func (r *Response) AssertStatus(code int) {
	r.T.Helper()
	if r.StatusCode != code {
		r.T.Fatalf("status: got %d want %d", r.StatusCode, code)
	}
}

// AssertJSONPath checks a top-level string field (simple dot path not supported).
func (r *Response) AssertJSON(field, want string) {
	r.T.Helper()
	defer r.Body.Close()
	var m map[string]any
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		r.T.Fatal(err)
	}
	got, _ := m[field].(string)
	if got != want {
		r.T.Fatalf("field %s: got %q want %q", field, got, want)
	}
}
