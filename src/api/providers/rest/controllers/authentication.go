package controllers

import (
	"errors"
	"time"

	"github.com/datti-to/purrmannplus-backend/api/providers/rest/models"
	"github.com/datti-to/purrmannplus-backend/app/commands"
	db_errors "github.com/datti-to/purrmannplus-backend/database/errors"
	"github.com/datti-to/purrmannplus-backend/utils/jwt"
	"github.com/gofiber/fiber/v2"
)

func AccountLogin(c *fiber.Ctx) error {
	a := new(models.PostLoginRequest)
	if err := c.BodyParser(a); err != nil {
		return err
	}

	dbAcc, err := commands.GetAccountByCredentials(a.AuthId, a.AuthPw)
	if err != nil {
		if errors.Is(err, &db_errors.ErrRecordNotFound) {
			return c.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{
				"error": "wrong credentials",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"error": err.Error(),
		})
	}

	token, err := jwt.NewAccountIdToken(dbAcc.Id)
	if err != nil {
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

		c.Cookie(cookie)

		return c.SendStatus(fiber.StatusCreated)
	}

	return c.JSON(&fiber.Map{
		"token": token,
	})
}
