// Package main is the smallest GoStack demo: JSON-only routes and default middleware.
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

	app.GET("/", func(c *gostack.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"message": "hello from GoStack"})
	})
	app.GET("/health", func(c *gostack.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})

	addr := ":3000"
	if p := os.Getenv("PORT"); p != "" {
		addr = ":" + p
	}
	log.Printf("hello example listening on %s (GET / and /health)", addr)
	if err := app.ListenAndServe(addr); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
