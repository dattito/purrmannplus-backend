package api

import (
	"fmt"

	"github.com/datti-to/purrmannplus-backend/api/controllers"
	"github.com/datti-to/purrmannplus-backend/config"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
)

var App *fiber.App

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

func Init() {
	App = fiber.New()

	App.Get(GetHealthRoute, controllers.GetHealth)

	App.Post(AccountLoginRoute, controllers.AccountLogin)

	App.Post(AddAccountRoute, controllers.AddAccount)
	App.Get(GetAccountsRoute, controllers.GetAccounts)

	App.Get(AddPhoneNumberRoute, controllers.AddPhoneNumber)

	App.Use(jwtware.New(getJWTConfig()))

	App.Post(SendPhoneNumberConfirmationLinkRoute, controllers.SendPhoneNumberConfirmationLink)
	App.Post(RegisterToSubstitutionUpdaterRoute, controllers.RegisterToSubstitutionUpdater)

}

func StartListening() error {
	return App.Listen(fmt.Sprintf(":%d", config.LISTENING_PORT))
}
