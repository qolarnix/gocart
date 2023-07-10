package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/fiber/v2/middleware/redirect"
)

func SetupRoutes(app *fiber.App) {
	app.Use(favicon.New(favicon.Config{
		File: "./favicon.ico",
		URL: "/favicon.ico",
	}))
	
	app.Use(redirect.New(redirect.Config{
		Rules: map[string]string{
			"/": "/login",
		},
		StatusCode: 301,
	}))

	app.Get("/login", func(c *fiber.Ctx) error {
		return c.Render("login", fiber.Map{
			"title": "GoCart - Login",
		})
	})

	app.Get("/api", func(c *fiber.Ctx) error {
		return c.SendString("GoCart API")
	})
}