package gostack

import (
	"html/template"
	"io/fs"
	"net/http"
	"path/filepath"
	"strings"
	"sync"
)

// RenderEngine renders html/template views with layouts.
type RenderEngine struct {
	fs      fs.FS
	mu      sync.RWMutex
	tmpl    map[string]*template.Template
	funcMap template.FuncMap
}

// RenderOptions configures the engine (reserved for future layout roots).
type RenderOptions struct {
	LayoutsDir    string
	PagesDir      string
	ComponentsDir string
}

// NewRenderEngine loads templates from filesystem (e.g. os.DirFS("views")).
func NewRenderEngine(fsys fs.FS, opts RenderOptions) *RenderEngine {
	_ = opts
	funcs := template.FuncMap{
		"csrf_token": func() string { return "" },
		"url":        func(name string, args ...any) string { return "#" },
		"asset":      func(p string) string { return "/public/" + strings.TrimLeft(p, "/") },
		"flash":   func(kind string) string { return "" },
		"partial": func(name string, data any) string { return "" },
	}
	return &RenderEngine{
		fs:      fsys,
		tmpl:    make(map[string]*template.Template),
		funcMap: funcs,
	}
}

// SetHelpers merges template func map entries (e.g. url, csrf).
func (e *RenderEngine) SetHelpers(f template.FuncMap) {
	e.mu.Lock()
	defer e.mu.Unlock()
	for k, v := range f {
		e.funcMap[k] = v
	}
}

// Render executes a page template; uses layouts/app.html when present.
func (e *RenderEngine) Render(w http.ResponseWriter, name string, data Data) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	path := filepath.ToSlash(name)
	if !strings.HasSuffix(path, ".html") {
		path += ".html"
	}
	t, err := e.parsePageLocked(path)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	return t.Execute(w, data)
}

// RenderPartial renders a component/partial only.
func (e *RenderEngine) RenderPartial(w http.ResponseWriter, name string, data Data) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	path := filepath.ToSlash(name)
	if !strings.HasSuffix(path, ".html") {
		path += ".html"
	}
	b, err := fs.ReadFile(e.fs, path)
	if err != nil {
		return err
	}
	t, err := template.New(path).Funcs(e.funcMap).Parse(string(b))
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	return t.Execute(w, data)
}

func (e *RenderEngine) parsePageLocked(pagePath string) (*template.Template, error) {
	if t, ok := e.tmpl[pagePath]; ok {
		return t, nil
	}
	pageBytes, err := fs.ReadFile(e.fs, pagePath)
	if err != nil {
		return nil, err
	}
	layoutPath := "layouts/app.html"
	layoutBytes, lerr := fs.ReadFile(e.fs, layoutPath)
	if lerr != nil {
		t, err := template.New(pagePath).Funcs(e.funcMap).Parse(string(pageBytes))
		if err != nil {
			return nil, err
		}
		e.tmpl[pagePath] = t
		return t, nil
	}
	t, err := template.New("layout").Funcs(e.funcMap).Parse(string(layoutBytes))
	if err != nil {
		return nil, err
	}
	t, err = t.Parse(string(pageBytes))
	if err != nil {
		return nil, err
	}
	e.tmpl[pagePath] = t
	return t, nil
}
