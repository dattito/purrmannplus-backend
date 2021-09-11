package api

import (
	"github.com/datti-to/purrmannplus-backend/config"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
)

func InitMiddlewares(app *fiber.App) {
	App.Use(jwtware.New(jwtware.Config{
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
	}))
}
