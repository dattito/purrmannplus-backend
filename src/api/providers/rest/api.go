package rest

import (
	"fmt"

	"github.com/datti-to/purrmannplus-backend/api/providers/rest/controllers"
	"github.com/datti-to/purrmannplus-backend/config"
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
	}
}

type RestProvider struct {
	app *fiber.App
}

func (r *RestProvider)Init() error {
	r.app = fiber.New()

	r.app.Get(GetHealthRoute, controllers.GetHealth)

	r.app.Post(AccountLoginRoute, controllers.AccountLogin)

	r.app.Post(AddAccountRoute, controllers.AddAccount)
	r.app.Get(GetAccountsRoute, controllers.GetAccounts)

	r.app.Get(AddPhoneNumberRoute, controllers.AddPhoneNumber)

	r.app.Use(jwtware.New(getJWTConfig()))

	r.app.Post(SendPhoneNumberConfirmationLinkRoute, controllers.SendPhoneNumberConfirmationLink)
	r.app.Post(RegisterToSubstitutionUpdaterRoute, controllers.RegisterToSubstitutionUpdater)

	return nil
}

func (r *RestProvider) StartListening() error {
	return r.app.Listen(fmt.Sprintf(":%d", config.LISTENING_PORT))
}
