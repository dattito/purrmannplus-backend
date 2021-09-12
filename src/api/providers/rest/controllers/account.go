package controllers

import (
	api_models "github.com/datti-to/purrmannplus-backend/api/providers/rest/models"
	"github.com/datti-to/purrmannplus-backend/database"
	"github.com/datti-to/purrmannplus-backend/services/hpg"
	"github.com/gofiber/fiber/v2"
)

func AddAccount(c *fiber.Ctx) error {
	accApi := new(api_models.PostAccountRequest)
	if err := c.BodyParser(accApi); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"error": err.Error(),
		})
	}

	acc, err := api_models.PostAccountRequestToAccount(accApi)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"error": err.Error(),
		})
	}

	correct, err := hpg.CheckCredentials(acc.AuthId, acc.AuthPw)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"error": err.Error(),
		})
	}

	if !correct {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"error": "invalid credentials",
		})
	}

	err = database.DB.AddAccount(acc)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(api_models.AccountToPostAccountResponse(acc))
}

func GetAccounts(c *fiber.Ctx) error {
	accs, err := database.DB.GetAccounts()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(api_models.AccountsToGetAccountResponses(accs))
}
