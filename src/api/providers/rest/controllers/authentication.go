package controllers

import (
	"errors"
	"time"

	"github.com/dattito/purrmannplus-backend/api/providers/rest/models"
	"github.com/dattito/purrmannplus-backend/app/commands"
	"github.com/dattito/purrmannplus-backend/config"
	db_errors "github.com/dattito/purrmannplus-backend/database/errors"
	"github.com/dattito/purrmannplus-backend/utils/logging"
	"github.com/dattito/purrmannplus-backend/utils/jwt"
	"github.com/gofiber/fiber/v2"
)

func AccountLogin(c *fiber.Ctx) error {
	a := new(models.PostLoginRequest)
	if err := c.BodyParser(a); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"error": err.Error(),
		})
	}

	dbAcc, err := commands.GetAccountByCredentials(a.AuthId, a.AuthPw)
	if err != nil {
		if errors.Is(err, &db_errors.ErrRecordNotFound) {
			return c.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{
				"error": "wrong credentials",
			})
		}
		logging.Errorf("Error while getting account by credentials: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"error": "Something went wrong",
		})
	}

	token, err := jwt.NewAccountIdToken(dbAcc.Id)
	if err != nil {
		logging.Errorf("Error while creating token: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"error": err.Error(),
		})
	}

	if a.StoreInCookie {
		cookie := new(fiber.Cookie)
		cookie.Name = "Authorization"
		cookie.Value = token
		cookie.Expires = time.Now().Add(24 * 30 * time.Hour)
		cookie.HTTPOnly = true

		if config.AUTHORIZATION_COOKIE_DOMAIN != "" {
			cookie.Domain = config.AUTHORIZATION_COOKIE_DOMAIN
		}

		c.Cookie(cookie)

		return c.SendStatus(fiber.StatusCreated)
	}

	return c.JSON(&fiber.Map{
		"token": token,
	})
}
