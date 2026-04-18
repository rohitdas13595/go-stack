package handlers

import "github.com/rohitdas13595/go-stack"

// Home renders the main page.
func Home(c *gostack.Context) error {
	return c.Render("pages/home", gostack.Data{
		"title": "SSR example",
	})
}
