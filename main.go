package main

import (
	"database/sql"
	"flag"
	"io/ioutil"
	"log"
	"time"
	"os"
	"strings"

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

type Config struct {
    Key string `db:"key"`
    Value string `db:"value"`
}

// a map of valid environment variables 
var validEnvVars = map[string]bool{
	"HOST":         true,
	"PORT":         true,
	"REGISTRATION": true,
}

func updateConfigTableWithEnvVars(db *sqlx.DB) error {
	envVars := os.Environ()
	for _, envVar := range envVars {
		pair := strings.SplitN(envVar, "=", 2)
		key := pair[0]
		value := pair[1]

		// Check if the key is valid according to the schema
		if _, isValid := validEnvVars[key]; isValid {
			query := `INSERT OR REPLACE INTO config (key, value) VALUES (?, ?)`
			if _, err := db.Exec(query, key, value); err != nil {
				return err
			}
		}
	}

	return nil
}

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

	// Update the config table with environment variables
	if err := updateConfigTableWithEnvVars(conn); err != nil {
		log.Fatal(err)
	}

	// Create a new instance of the session middleware
	sessionMiddleware := session.New(session.Config{
		Expiration: 2 * time.Hour,
		// CookieSecure: true,
	})

	// Create a new instance of SessionManager with the database connection
	sessionManager := &auth.SessionManager{
		Conn: conn,
	}

	engine := html.New("./html", ".html")
	engine.Reload(true)
	engine.Debug(true)

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

	// Read all keys and values from the config table
	var configValues []Config
	err = conn.Select(&configValues, "SELECT key, value FROM config")
	if err != nil {
		log.Fatal("Failed to fetch configuration values from config table: ", err)
	}

	// Create a map to store the configuration values
	configMap := make(map[string]string)
	for _, cfg := range configValues {
		configMap[cfg.Key] = cfg.Value
	}

	registrationEnabled := configMap["REGISTRATION"] == "true"
	port := configMap["PORT"]

	// Middleware that checks for a valid session on every request
	app.Use(func(c *fiber.Ctx) error {
		sess := sessionMiddleware.Get(c)	
		fmt.Println(sess.Get("user"))
		return c.Next()
		
	})

	routes.SetupNoauthRoutes(app.Group("/"), sessionManager, sessionMiddleware, registrationEnabled)
	routes.SetupAuthRoutes(app.Group("/dash"), sessionMiddleware)
	routes.SetupApiRoutes(app.Group("/api"))

	log.Fatal(app.Listen(":" + port))
}
