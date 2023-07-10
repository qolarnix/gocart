package routes

import (
	"github.com/gofiber/fiber/v2"
)

func ApiRoutes() *fiber.App {

	api := fiber.New()

	api.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("GoCart API")
	})


	return api
}