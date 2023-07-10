package main

import (
	"flag"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"

	"gocart/v2/router"
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

	routes.SetupRoutes(app)

	log.Fatal(app.Listen(*port))
}