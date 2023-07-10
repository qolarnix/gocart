package main

import (
	"database/sql"
	"flag"
	"io/ioutil"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/fiber/v2/middleware/redirect"
	"github.com/gofiber/template/html/v2"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"

	"github.com/gofiber/session/v2"
	"gocart/v2/auth"
	"gocart/v2/routes"

	"fmt"
)

var (
	port = flag.String("port", ":3000", "Port to listen on")
)

func main() {
	flag.Parse()

	// Initialize the database connection
	db, err := sql.Open("sqlite3", "./db/database.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Wrap the database connection with sqlx for additional features
	conn := sqlx.NewDb(db, "sqlite3")

	// Read the SQL schema file
	schemaSQL, err := ioutil.ReadFile("./db/schema.sql")
	if err != nil {
		log.Fatal(err)
	}

	// Execute the SQL statements to create the tables
	_, err = conn.Exec(string(schemaSQL))
	if err != nil {
		log.Fatal(err)
	}

	// Create a new instance of the session middleware
	sessionMiddleware := session.New(session.Config{
		Expiration: 2 * time.Hour,
	})

	// Create a new instance of SessionManager with the database connection
	sessionManager := &auth.SessionManager{
		Conn: conn,
	}

	engine := html.New("./html", ".html")

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Use(favicon.New(favicon.Config{
		File: "./favicon.ico",
		URL:  "/favicon.ico",
	}))

	app.Use(redirect.New(redirect.Config{
		Rules: map[string]string{
			"/": "/login",
		},
		StatusCode: 301,
	}))

	// Middleware that checks for a valid session on every request
	app.Use(func(c *fiber.Ctx) error {
		sess := sessionMiddleware.Get(c)	
		fmt.Println(sess.Get("user"))
		return c.Next()
		
	})

	routes.SetupNoauthRoutes(app.Group("/"), sessionManager, sessionMiddleware)
	routes.SetupAuthRoutes(app.Group("/dash"), sessionMiddleware)
	routes.SetupApiRoutes(app.Group("/api"))

	log.Fatal(app.Listen(*port))
}
