package routes

import (
	"fmt"
	//"time"

	"github.com/gofiber/fiber/v2"
	"gocart/v2/auth"
	"github.com/gofiber/session/v2"
	"github.com/nbutton23/zxcvbn-go"
	"net/mail"
	// "github.com/gofiber/storage/memory"
	
)

func alreadyLoggedIn(sessionMiddleware *session.Session, c *fiber.Ctx) error {
	sess := sessionMiddleware.Get(c)
	_, ok := sess.Get("user").(string)
	if ok {
		// If user already logged in
		return c.Redirect("/dash")
	}
	return c.Next()
}

func SetupNoauthRoutes(app fiber.Router, sm *auth.SessionManager, sessionMiddleware *session.Session, registrationEnabled bool) {

	app.Get("/login", func(c *fiber.Ctx) error {
		return alreadyLoggedIn(sessionMiddleware, c)
	}, func(c *fiber.Ctx) error {
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
	if registrationEnabled {

		app.Get("/register", func(c *fiber.Ctx) error {
			return c.Render("register", fiber.Map{
				"title": "GoCart - Register",
			})
		})
	
		// register request, so we can strictly control the data that we work with.
		type RegisterRequest struct {
			FirstName string `json:"first_name"`
			LastName  string `json:"last_name"`
			Email     string `json:"email"`
			Password  string `json:"password"`
		}
	
		app.Post("/register", func(c *fiber.Ctx) error {
			var registerReq RegisterRequest
			if err := c.BodyParser(&registerReq); err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"message": "Bad request",
				})
			}
		
			// Check all the fields are of type string and meet length requirements
			if registerReq.FirstName == "" || len(registerReq.FirstName) > 16 ||
				registerReq.LastName == "" || len(registerReq.LastName) > 16 ||
				registerReq.Email == "" || len(registerReq.Email) > 30 ||
				registerReq.Password == "" || len(registerReq.Password) > 64 {
				return c.Render("register", fiber.Map{
					"title": "GoCart - Register",
					"errorMessage": "Invalid or out of length parameters. Please check your inputs and try again.",
				})
			}
		
			// Check if the email is valid
			_, err := mail.ParseAddress(registerReq.Email)
			if err != nil {
				return c.Render("register", fiber.Map{
					"title": "GoCart - Register",
					"errorMessage": "Invalid email format. Please enter a valid email.",
				})
			}
		
			// Check password strength
			userInputs := []string{registerReq.FirstName, registerReq.LastName, registerReq.Email}
			strength := zxcvbn.PasswordStrength(registerReq.Password, userInputs)
			if strength.Score < 3 {
				return c.Render("register", fiber.Map{
					"title":       "GoCart - Register",
					"errorMessage": "Password too weak. Please use a stronger password.",
				})
			}
		
			err = sm.RegisterUser(registerReq.FirstName, registerReq.LastName, registerReq.Email, registerReq.Password)
			if err != nil {
				return c.Render("register", fiber.Map{
					"title":       "GoCart - Register",
					"errorMessage": fmt.Sprintf("Error registering user. Please try again.\n%s", err),
				})
			}
		
			return c.Redirect("/login")
		})

	}
	
	
	
}
