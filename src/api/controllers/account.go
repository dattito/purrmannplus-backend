package controllers

import (
	api_models "github.com/datti-to/purrmannplus-backend/api/models"
	"github.com/datti-to/purrmannplus-backend/database"
	"github.com/gofiber/fiber/v2"
)

func AddAccount(c *fiber.Ctx) error {
	accApi := new(api_models.PostAccountRequest)
	if err := c.BodyParser(accApi); err != nil {
		return err
	}

	acc, err := api_models.PostAccountRequestToAccount(accApi)
	if err != nil {
		return c.Status(400).JSON(&fiber.Map{
			"error": err.Error(),
		})
	}

	err = database.DB.AddAccount(acc)
	if err != nil {
		return c.Status(400).JSON(&fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(&fiber.Map{
		"id": acc.ID,
	})
}
