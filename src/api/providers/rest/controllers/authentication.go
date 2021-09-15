package controllers

import (
	"errors"

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

	return c.JSON(&fiber.Map{
		"token": token,
	})
}
