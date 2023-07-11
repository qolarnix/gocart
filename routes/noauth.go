package routes

import (
	"fmt"
	//"time"

	"github.com/gofiber/fiber/v2"
	"gocart/v2/auth"
	"github.com/gofiber/session/v2"
	// "github.com/gofiber/storage/memory"
	
)

func SetupNoauthRoutes(app fiber.Router, sm *auth.SessionManager, sessionMiddleware *session.Session) {
	app.Get("/login", func(c *fiber.Ctx) error {
		return c.Render("login", fiber.Map{
			"title": "GoCart - Login",
		})
	})

	app.Post("/login", func(c *fiber.Ctx) error {
		type LoginRequest struct {
			Email    string `form:"email"`
			Password string `form:"password"`
		}
	
		var loginReq LoginRequest
		if err := c.BodyParser(&loginReq); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Error on request",
			})
		}
	
		userSession, err := sm.Login(loginReq.Email, loginReq.Password)
    	if err != nil {
    	    return c.Render("login", fiber.Map{
    	        "title":          "GoCart - Login",
    	        "errorMessage":   "Invalid email or password",
    	        "emailValue":     loginReq.Email,
    	        "passwordValue":  loginReq.Password,
    	    })
    	}
	
		
		// Destroy old session
        oldSession := sessionMiddleware.Get(c)
        if err := oldSession.Destroy(); err != nil {
            panic(err)
        }

        sess := sessionMiddleware.Get(c)


        sess.Regenerate()

		fmt.Println("Session ID:" + userSession.ID)

		
		sess.Set("user", userSession.ID)
	 
		if err := sess.Save(); err != nil {
			fmt.Println("Getting an error with: ", err)
		}
		 
		// Go to the /dash authenticated
		return c.Redirect("/dash")
	})

	type RegisterRequest struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email"`
		Password  string `json:"password"`
	}

	app.Get("/register", func(c *fiber.Ctx) error {
		return c.Render("register", fiber.Map{
			"title": "GoCart - Register",
		})
	})

	app.Post("/register", func(c *fiber.Ctx) error {
		var registerReq RegisterRequest
		err := c.BodyParser(&registerReq)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Error on request",
			})
		}

		if registerReq.FirstName == "" || registerReq.LastName == "" || registerReq.Email == "" || registerReq.Password == "" {
			return c.Render("register", fiber.Map{
				"title": "GoCart - Register",
			})
		}

		err = sm.RegisterUser(registerReq.FirstName, registerReq.LastName, registerReq.Email, registerReq.Password)
		if err != nil {
			return c.Render("register", fiber.Map{
				"title":       "GoCart - Register",
				"errorMessage": fmt.Sprintf("Error registering user. Please try again.\n%s", err),
			})
		}

		return c.Render("register", fiber.Map{
			"title":          "GoCart - Register",
			"successMessage": fmt.Sprintf("User %s successfully registered", registerReq.Email),
		})
	})
}
