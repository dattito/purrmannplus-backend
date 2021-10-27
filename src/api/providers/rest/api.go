package rest

import (
	"fmt"

	"github.com/dattito/purrmannplus-backend/api/providers/rest/controllers"
	"github.com/dattito/purrmannplus-backend/api/providers/rest/routes"
	"github.com/dattito/purrmannplus-backend/api/providers/rest/session"
	"github.com/dattito/purrmannplus-backend/config"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/gofiber/template/amber"
)

// Get the JWT configuration for the api
func getJWTConfig() jwtware.Config {
	return jwtware.Config{
		SigningKey: []byte(config.JWT_SECRET),
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			IsMissingOrMalformedJWT := err.Error() == "Missing or malformed JWT"

			var ErrorText string
			if IsMissingOrMalformedJWT {
				ErrorText = "Missing or malformed JWT"
			} else {
				ErrorText = "Invalid or expired JWT"
			}

			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": ErrorText,
			})
		},
		TokenLookup: "header:Authorization,cookie:Authorization",
	}
}

// Protected is a middleware that checks if the user is logged in
func Protected() fiber.Handler {
	return jwtware.New(getJWTConfig())
}

type RestProvider struct {
	app *fiber.App
}

// Initialize the fiber app and sets the routes and middlewares
func (r *RestProvider) Init() error {
	r.app = fiber.New(fiber.Config{
		Views: amber.New(config.PATH_TO_API_VIEWS, ".amber"),
	})

	r.app.Static("/static", config.PATH_TO_API_STATIC)

	if config.CORS_ALLOWED_ORIGINS != "" {
		r.app.Use(cors.New(cors.Config{
			AllowOrigins: config.CORS_ALLOWED_ORIGINS,
			AllowHeaders: "Origin, Content-Type, Accept",
		}))
	}

	r.app.Get(routes.HealthRoute, controllers.GetHealth)
	r.app.Get(routes.AboutRoute, controllers.About)

	v1 := r.app.Group("/v1")

	v1.Post(routes.AccountLoginRoute, controllers.AccountLogin)
	v1.Get(routes.AccountLogoutRoute, controllers.AccountLogout)
	v1.Get(routes.IsLoggedInRoute, Protected(), controllers.IsLoggedIn)

	v1.Post(routes.AddAccountRoute, controllers.AddAccount)
	v1.Delete(routes.DeleteAccountRoute, Protected(), controllers.DeleteAccount)
	//v1.Get(GetAccountsRoute, controllers.GetAccounts)
	v1.Post(routes.SendPhoneNumberConfirmationLinkRoute, Protected(), controllers.SendPhoneNumberConfirmationLink)
	v1.Get(routes.AddPhoneNumberRoute, controllers.AddPhoneNumber)

	v1.Post(routes.RegisterToSubstitutionUpdaterRoute, Protected(), controllers.RegisterToSubstitutionUpdater)
	v1.Delete(routes.UnregisterFromSubstitutionUpdaterRoute, Protected(), controllers.UnregisterFromSubstitutionUpdater)

	r.app.All(routes.SubstitutionSpeedFormRoute, controllers.SubstitutionSpeedForm)
	r.app.All(routes.ValidateSubstitutionSpeedFormRoute, controllers.ValidateSubstitutionSpeedForm)
	r.app.Get(routes.FinishSubstitutionSpeedFormRoute, controllers.FinishSubstitutionSpeedForm)

	session.Init()

	return nil
}

// Start the fiber app and listen on the specified port
func (r *RestProvider) StartListening() error {
	return r.app.Listen(fmt.Sprintf(":%d", config.LISTENING_PORT))
}
