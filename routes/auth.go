package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/session/v2"
	"fmt"
)

func AuthRequired(sessionMiddleware *session.Session, c *fiber.Ctx) error {
	sess := sessionMiddleware.Get(c)

	user := sess.Get("user")
	if user == nil {
		return c.Redirect("/login")
	}
	return c.Next()
}

func SetupAuthRoutes(app fiber.Router, sessionMiddleware *session.Session) {
	app.Use(func(c *fiber.Ctx) error {
		return AuthRequired(sessionMiddleware, c)
	})

	app.Get("/", func(c *fiber.Ctx) error {
		session_key := sessionMiddleware.Get(c)

		return c.SendString(fmt.Sprintf("You're in.. Good job.. \n\n%s", session_key))
	})
}