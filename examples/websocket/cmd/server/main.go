package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
	"github.com/rohitdas13595/go-stack"
	"github.com/rohitdas13595/go-stack/middleware"
	"github.com/rohitdas13595/go-stack/ws"
)

func main() {
	hub := ws.NewHub()

	app := gostack.New()
	app.Use(middleware.Recover(), middleware.Logger())

	app.GET("/", func(c *gostack.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"hint": "connect a WebSocket client to /ws?channel=demo",
		})
	})

	app.GET("/ws", func(c *gostack.Context) error {
		channel := c.Query("channel")
		if channel == "" {
			channel = "demo"
		}

		conn, err := ws.Upgrade(c.ResponseWriter(), c.Request(), nil)
		if err != nil {
			return err
		}
		defer conn.Close()
		hub.Join(channel, conn)
		defer hub.Leave(conn)

		for {
			mt, msg, err := conn.ReadMessage()
			if err != nil {
				return nil
			}
			if mt == websocket.TextMessage {
				hub.Emit(channel, msg)
			}
		}
	})

	addr := ":3000"
	if p := os.Getenv("PORT"); p != "" {
		addr = ":" + p
	}
	log.Printf("websocket example listening on %s", addr)
	if err := app.ListenAndServe(addr); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
