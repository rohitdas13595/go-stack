package gostack_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/rohitdas13595/go-stack"
	"github.com/rohitdas13595/go-stack/db"
	"github.com/rohitdas13595/go-stack/middleware"
	"github.com/rohitdas13595/go-stack/testutil"
)

func TestHealthJSON(t *testing.T) {
	app := gostack.New()
	app.Use(middleware.Recover())
	app.GET("/health", func(c *gostack.Context) error {
		return c.JSON(200, map[string]string{"status": "ok"})
	})
	ts := testutil.NewApp(t, app)
	res := ts.GET("/health")
	res.AssertStatus(200)
}

func TestScaffoldRoutesSQLite(t *testing.T) {
	tmp := t.TempDir()
	if err := os.MkdirAll(filepath.Join(tmp, "views/layouts"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(tmp, "views/pages"), 0o755); err != nil {
		t.Fatal(err)
	}
	layout := `<!DOCTYPE html><html><body>{{ block "content" . }}{{ end }}</body></html>`
	if err := os.WriteFile(filepath.Join(tmp, "views/layouts/app.html"), []byte(layout), 0o644); err != nil {
		t.Fatal(err)
	}
	page := `{{ define "content" }}<p>hi</p>{{ end }}`
	if err := os.WriteFile(filepath.Join(tmp, "views/pages/home.html"), []byte(page), 0o644); err != nil {
		t.Fatal(err)
	}
	app := gostack.New()
	app.SetRenderer(gostack.NewRenderEngine(os.DirFS(filepath.Join(tmp, "views")), gostack.RenderOptions{}))
	app.GET("/", func(c *gostack.Context) error {
		return c.Render("pages/home", gostack.Data{"title": "T"})
	})
	ts := testutil.NewApp(t, app)
	res := ts.GET("/")
	res.AssertStatus(200)
}

func TestDBManager(t *testing.T) {
	sdb, err := db.OpenSQLite("file:" + filepath.Join(t.TempDir(), "t.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer sdb.Close()
	m := db.NewManager()
	m.Register("default", sdb)
	gostack.SetDBManager(m)
	if gostack.DB() != sdb {
		t.Fatal("expected default db")
	}
}
