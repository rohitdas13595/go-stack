package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/rohitdas13595/go-stack"
	"github.com/rohitdas13595/go-stack/middleware"
	"github.com/rohitdas13595/go-stack/sse"
)

func main() {
	app := gostack.New()
	app.Use(middleware.Recover(), middleware.Logger())

	app.GET("/", func(c *gostack.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"hint": "open GET /stream in a browser or curl -N",
		})
	})

	app.GET("/stream", func(c *gostack.Context) error {
		return sse.Handler(c.ResponseWriter(), func(send sse.Sender) {
			for i := 0; i < 5; i++ {
				send("tick", time.Now().Format(time.RFC3339Nano))
				time.Sleep(500 * time.Millisecond)
			}
			send("done", "stream finished")
		})
	})

	addr := ":3000"
	if p := os.Getenv("PORT"); p != "" {
		addr = ":" + p
	}
	log.Printf("sse example listening on %s (GET /stream)", addr)
	if err := app.ListenAndServe(addr); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
