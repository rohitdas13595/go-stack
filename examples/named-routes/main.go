// Package main demonstrates named routes and RouteURL resolution.
package main

import (
	"log"
	"net/http"
	"os"

	"github.com/rohitdas13595/go-stack"
	"github.com/rohitdas13595/go-stack/middleware"
)

func main() {
	app := gostack.New()
	app.Use(middleware.Recover(), middleware.Logger())

	app.GETNamed("user.profile", "/users/:id", func(c *gostack.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"user_id": c.Param("id"),
		})
	})

	app.GET("/_routes/resolve", func(c *gostack.Context) error {
		path, err := app.RouteURL("user.profile", map[string]string{"id": "42"})
		if err != nil {
			return c.BadRequest(err)
		}
		return c.JSON(http.StatusOK, map[string]string{
			"name": "user.profile",
			"path": path,
		})
	})

	addr := ":3000"
	if p := os.Getenv("PORT"); p != "" {
		addr = ":" + p
	}
	log.Printf("named-routes example listening on %s", addr)
	if err := app.ListenAndServe(addr); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
