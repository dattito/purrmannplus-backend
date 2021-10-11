package rest

import (
	"fmt"

	"github.com/dattito/purrmannplus-backend/api/providers/rest/controllers"
	"github.com/dattito/purrmannplus-backend/config"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
)

func getJWTConfig() jwtware.Config {
	return jwtware.Config{
		SigningKey: []byte(config.JWT_SECRET),
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			IsMissingOrMalformedJWT := err.Error() == "Missing or malformed JWT"
			var StatusCode int
			if IsMissingOrMalformedJWT {
				StatusCode = fiber.StatusBadRequest
			} else {
				StatusCode = fiber.StatusUnauthorized
			}

			var ErrorText string
			if IsMissingOrMalformedJWT {
				ErrorText = "Missing or malformed JWT"
			} else {
				ErrorText = "Invalid or expired JWT"
			}

			return c.Status(StatusCode).JSON(fiber.Map{
				"error": ErrorText,
			})
		},
		TokenLookup: "header:Authorization,cookie:Authorization",
	}
}

func Protected() fiber.Handler {
	return jwtware.New(getJWTConfig())
}

type RestProvider struct {
	app *fiber.App
}

func (r *RestProvider) Init() error {
	r.app = fiber.New()

	r.app.Get(HealthRoute, controllers.GetHealth)
	r.app.Get(AboutRoute, controllers.About)

	v1 := r.app.Group("/v1")

	v1.Post(AccountLoginRoute, controllers.AccountLogin)
	v1.Get(AccountLogoutRoute, controllers.AccountLogout)

	v1.Post(AddAccountRoute, controllers.AddAccount)
	v1.Delete(DeleteAccountRoute, Protected(), controllers.DeleteAccount)
	//v1.Get(GetAccountsRoute, controllers.GetAccounts)
	v1.Post(SendPhoneNumberConfirmationLinkRoute, Protected(), controllers.SendPhoneNumberConfirmationLink)
	v1.Get(AddPhoneNumberRoute, controllers.AddPhoneNumber)

	v1.Post(RegisterToSubstitutionUpdaterRoute, Protected(), controllers.RegisterToSubstitutionUpdater)
	v1.Delete(UnregisterFromSubstitutionUpdaterRoute, Protected(), controllers.UnregisterFromSubstitutionUpdater)

	return nil
}

func (r *RestProvider) StartListening() error {
	return r.app.Listen(fmt.Sprintf(":%d", config.LISTENING_PORT))
}
