package controllers

import (
	"errors"

	"github.com/datti-to/purrmannplus-backend/api/models"
	"github.com/datti-to/purrmannplus-backend/database"
	db_errors "github.com/datti-to/purrmannplus-backend/database/errors"
	"github.com/datti-to/purrmannplus-backend/utils/jwt"
	"github.com/gofiber/fiber/v2"
)

func AccountLogin(c *fiber.Ctx) error {
	a := new(models.PostLoginRequest)
	if err := c.BodyParser(a); err != nil {
		return err
	}

	acc, err := models.PostLoginRequestToAccount(*a)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{
			"error": err.Error(),
		})
	}

	dbAcc, err := database.DB.GetAccountByCredentials(*acc)
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

	token, err := jwt.NewAccountLoginToken(dbAcc.Id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(&fiber.Map{
		"token": token,
	})
}
