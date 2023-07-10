package routes

import (
	"github.com/gofiber/fiber/v2"
)

func NoauthRoutes() *fiber.App {
	noauth := fiber.New()

	noauth.Get("/login", func(c *fiber.Ctx) error {
		return c.Render("login", fiber.Map{
			"title": "GoCart - Login",
		})
	})
	
	return noauth
}
