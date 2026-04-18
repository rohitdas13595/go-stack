package routes

import (
	"net/http"
	"sync/atomic"

	"github.com/rohitdas13595/go-stack"
	"github.com/rohitdas13595/go-stack/middleware"
)

type createItemRequest struct {
	Name string `json:"name" validate:"required,min=1,max=200"`
}

// Setup demonstrates a versioned JSON API with CORS and request validation.
func Setup() *gostack.App {
	app := gostack.New()
	app.Use(
		middleware.Recover(),
		middleware.Logger(),
		middleware.CORS([]string{"*"}, []string{"GET", "POST", "OPTIONS"}),
	)

	var nextID atomic.Int64

	v1 := app.Group("/api/v1")
	v1.GET("/items", func(c *gostack.Context) error {
		return c.JSON(http.StatusOK, map[string]any{
			"items": []map[string]any{
				{"id": 1, "name": "alpha"},
				{"id": 2, "name": "beta"},
			},
		})
	})

	v1.GET("/items/:id", func(c *gostack.Context) error {
		id := c.Param("id")
		return c.JSON(http.StatusOK, map[string]any{
			"id":   id,
			"name": "item-" + id,
		})
	})

	v1.POST("/items", func(c *gostack.Context) error {
		var req createItemRequest
		if err := c.Bind(&req); err != nil {
			return c.BadRequest(err)
		}
		id := nextID.Add(1)
		return c.JSON(http.StatusCreated, map[string]any{
			"id":   id,
			"name": req.Name,
		})
	})

	app.GET("/health", func(c *gostack.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})

	return app
}
