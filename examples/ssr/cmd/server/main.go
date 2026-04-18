package main

import (
	"log"
	"net/http"
	"os"

	"github.com/rohitdas13595/go-stack/examples/ssr/routes"
)

func main() {
	app := routes.Setup()
	addr := ":3000"
	if p := os.Getenv("PORT"); p != "" {
		addr = ":" + p
	}
	log.Printf("ssr example listening on %s", addr)
	if err := app.ListenAndServe(addr); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
