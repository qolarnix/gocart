package main

import (
	"flag"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/fiber/v2/middleware/redirect"

	"gocart/v2/routes"
)

var (
	port = flag.String("port", ":3000", "Port to listen on")
)

func main() {
	flag.Parse()

	engine := html.New("./html", ".html")

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	// favicon, duuh
	app.Use(favicon.New(favicon.Config{
		File: "./favicon.ico",
		URL: "/favicon.ico",
	}))

	// redirect to login if going to root 
	app.Use(redirect.New(redirect.Config{
		Rules: map[string]string{
			"/": "/login",
		},
		StatusCode: 301,
	}))

	app.Mount("/", routes.NoauthRoutes())
	app.Mount("/dash", routes.AuthRoutes())
	app.Mount("/api", routes.ApiRoutes())
	


	log.Fatal(app.Listen(*port))
}