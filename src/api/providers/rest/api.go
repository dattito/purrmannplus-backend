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

func (r *RestProvider) Init() error {
	r.app = fiber.New()

	v1 := r.app.Group("/v1")

	v1.Get(GetHealthRoute, controllers.GetHealth)

	v1.Post(AccountLoginRoute, controllers.AccountLogin)

	v1.Post(AddAccountRoute, controllers.AddAccount)
	//v1.Get(GetAccountsRoute, controllers.GetAccounts)

	v1.Get(AddPhoneNumberRoute, controllers.AddPhoneNumber)

	v1.Use(jwtware.New(getJWTConfig()))

	v1.Post(SendPhoneNumberConfirmationLinkRoute, controllers.SendPhoneNumberConfirmationLink)
	v1.Post(RegisterToSubstitutionUpdaterRoute, controllers.RegisterToSubstitutionUpdater)
	v1.Delete(UnregisterFromSubstitutionUpdaterRoute, controllers.UnregisterFromSubstitutionUpdater)

	return nil
}

func (r *RestProvider) StartListening() error {
	return r.app.Listen(fmt.Sprintf(":%d", config.LISTENING_PORT))
}
