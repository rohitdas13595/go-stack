package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var newAPI, newSPA bool

var newCmd = &cobra.Command{
	Use:   "new [name]",
	Short: "Scaffold a new application",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		if err := scaffoldProject(name, newAPI, newSPA); err != nil {
			exitErr(err)
		}
	},
}

func init() {
	newCmd.Flags().BoolVar(&newAPI, "api", false, "API-only skeleton")
	newCmd.Flags().BoolVar(&newSPA, "spa", false, "SPA mode stub")
}

func scaffoldProject(name string, apiOnly, spa bool) error {
	base := filepath.Join(".", name)
	if err := os.MkdirAll(base, 0o755); err != nil {
		return err
	}
	dirs := []string{
		"cmd/server", "config", "db/migrations", "internal/handlers",
		"views/layouts", "views/pages", "routes", "public", "storage",
	}
	for _, d := range dirs {
		if err := os.MkdirAll(filepath.Join(base, d), 0o755); err != nil {
			return err
		}
	}
	files := map[string]string{
		"go.mod": goModContents(name),
		"cmd/server/main.go": mainGo(name),
		"config/app.yaml":    appYAML(),
		"routes/web.go":      routesGo(name),
		"internal/handlers/home.go": homeHandler(),
		"views/layouts/app.html": layoutHTML(),
		"views/pages/home.html":  homeHTML(),
		".env.example":           dotEnvExample(),
	}
	if apiOnly {
		files["routes/web.go"] = routesAPIOnly(name)
	}
	if spa {
		files["public/index.html"] = spaIndexHTML()
	}
	for path, content := range files {
		full := filepath.Join(base, path)
		if err := os.WriteFile(full, []byte(content), 0o644); err != nil {
			return err
		}
	}
	// Initial migration
	mig := filepath.Join(base, "db/migrations/20260101000001_init.sql")
	return os.WriteFile(mig, []byte(initMigration()), 0o644)
}

func goModContents(name string) string {
	return fmt.Sprintf(`module %s

go 1.23

require github.com/rohitdas13595/go-stack v0.0.0

replace github.com/rohitdas13595/go-stack => ..
`, name)
}

func mainGo(mod string) string {
	return `package main

import (
	"log"
	"net/http"
	"os"

	"` + mod + `/routes"
)

func main() {
	app := routes.Setup()
	addr := ":3000"
	if p := os.Getenv("PORT"); p != "" {
		addr = ":" + p
	}
	log.Printf("listening %s", addr)
	if err := app.ListenAndServe(addr); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
`
}

func routesGo(mod string) string {
	return `package routes

import (
	"os"
	"path/filepath"

	"github.com/rohitdas13595/go-stack"
	"github.com/rohitdas13595/go-stack/db"
	"github.com/rohitdas13595/go-stack/middleware"
	"` + mod + `/internal/handlers"
)

func Setup() *gostack.App {
	app := gostack.New()
	app.Use(middleware.Recover(), middleware.RequestID(), middleware.Logger())

	wd, _ := os.Getwd()
	views := filepath.Join(wd, "views")
	app.SetRenderer(gostack.NewRenderEngine(os.DirFS(views), gostack.RenderOptions{}))

	mgr := db.NewManager()
	gostack.SetDBManager(mgr)

	app.GET("/", handlers.Home)
	app.GET("/health", func(c *gostack.Context) error {
		return c.JSON(200, map[string]string{"status": "ok"})
	})
	return app
}
`
}

func routesAPIOnly(mod string) string {
	return `package routes

import (
	"github.com/rohitdas13595/go-stack"
	"github.com/rohitdas13595/go-stack/middleware"
)

func Setup() *gostack.App {
	app := gostack.New()
	app.Use(middleware.Recover(), middleware.Logger())
	app.GET("/health", func(c *gostack.Context) error {
		return c.JSON(200, map[string]string{"status": "ok"})
	})
	return app
}
`
}

func homeHandler() string {
	return `package handlers

import "github.com/rohitdas13595/go-stack"

func Home(c *gostack.Context) error {
	return c.Render("pages/home", gostack.Data{"title": "Home"})
}
`
}

func layoutHTML() string {
	return `<!DOCTYPE html>
<html><head><title>{{.title}}</title>
<script src="https://unpkg.com/htmx.org@1.9.12"></script>
<script defer src="https://cdn.jsdelivr.net/npm/alpinejs@3.x.x/dist/cdn.min.js"></script>
</head><body>{{ block "content" . }}{{ end }}</body></html>`
}

func homeHTML() string {
	return `{{ define "content" }}<h1>GoStack</h1><p>Welcome.</p>{{ end }}`
}

func appYAML() string {
	return `app:
  name: demo
  env: development
server:
  port: 3000
`
}

func dotEnvExample() string {
	return "PORT=3000\nDATABASE_URL=file:./storage/app.db\nJWT_SECRET=change-me\n"
}

func initMigration() string {
	return `-- +gostack:up
CREATE TABLE IF NOT EXISTS users (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  email TEXT NOT NULL UNIQUE,
  name TEXT NOT NULL,
  password TEXT NOT NULL,
  role TEXT NOT NULL DEFAULT 'user',
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  deleted_at DATETIME
);

-- +gostack:down
DROP TABLE IF EXISTS users;
`
}

func spaIndexHTML() string {
	return `<!DOCTYPE html><html><body><div id="app"></div></body></html>`
}
