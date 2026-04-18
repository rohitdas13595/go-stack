package routes

import (
	"os"
	"path/filepath"

	"github.com/rohitdas13595/go-stack"
	"github.com/rohitdas13595/go-stack/examples/ssr/handlers"
	"github.com/rohitdas13595/go-stack/middleware"
)

// Setup builds the SSR demo app (templates under ./views).
func Setup() *gostack.App {
	app := gostack.New()
	app.Use(middleware.Recover(), middleware.RequestID(), middleware.Logger())

	wd, _ := os.Getwd()
	views := filepath.Join(wd, "views")
	app.SetRenderer(gostack.NewRenderEngine(os.DirFS(views), gostack.RenderOptions{}))

	app.GET("/", handlers.Home)
	app.GET("/health", func(c *gostack.Context) error {
		return c.JSON(200, map[string]string{"status": "ok"})
	})
	return app
}
