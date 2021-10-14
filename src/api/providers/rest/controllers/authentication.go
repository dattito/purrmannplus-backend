package controllers

import (
	"errors"
	"time"

	"github.com/dattito/purrmannplus-backend/api/providers/rest/models"
	"github.com/dattito/purrmannplus-backend/app/commands"
	"github.com/dattito/purrmannplus-backend/config"
	db_errors "github.com/dattito/purrmannplus-backend/database/errors"
	"github.com/dattito/purrmannplus-backend/utils/jwt"
	"github.com/dattito/purrmannplus-backend/utils/logging"
	"github.com/gofiber/fiber/v2"
)

// AccountLogin logs in the user and returns a JWT token or sets a cookie
func AccountLogin(c *fiber.Ctx) error {
	a := new(models.PostLoginRequest)
	if err := c.BodyParser(a); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"error": err.Error(),
		})
	}

	dbAcc, err := commands.GetAccountByCredentials(a.Username, a.Password)
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

		if a.StayLoggedIn {
			cookie.Expires = time.Now().Add(time.Hour * 24 * 30)
		}
		cookie.HTTPOnly = true

		if config.AUTHORIZATION_COOKIE_DOMAIN != "" {
			cookie.Domain = config.AUTHORIZATION_COOKIE_DOMAIN
		}

		c.Cookie(cookie)

		return c.Status(fiber.StatusCreated).JSON(&fiber.Map{
			"ok": true,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(&fiber.Map{
		"token": token,
	})
}

// AccountLogout deletes the authorization cookie (logs out the user)
func AccountLogout(c *fiber.Ctx) error {
	cookie := new(fiber.Cookie)
	cookie.Name = "Authorization"
	cookie.Value = ""
	cookie.Expires = time.Now().Add(-1 * time.Hour)
	cookie.HTTPOnly = config.AUTHORIZATION_COOKIE_HTTPONLY
	
	cookie.Secure = config.AUTHORIZATION_COOKIE_SECURE

	if config.AUTHORIZATION_COOKIE_DOMAIN != "" {
		cookie.Domain = config.AUTHORIZATION_COOKIE_DOMAIN
	}

	c.Cookie(cookie)

	return c.SendStatus(fiber.StatusOK)
}

func IsLoggedIn(c *fiber.Ctx) error {
	return c.JSON(&fiber.Map{
		"loggedIn": true,
	})
}
