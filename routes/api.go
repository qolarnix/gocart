package routes

import (
	"github.com/gofiber/fiber/v2"
)

func SetupApiRoutes(app fiber.Router) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("yayayayayay api workie")
	})
}