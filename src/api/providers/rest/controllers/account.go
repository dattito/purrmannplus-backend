package controllers

import (
	api_models "github.com/dattito/purrmannplus-backend/api/providers/rest/models"
	"github.com/dattito/purrmannplus-backend/app/commands"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

func AddAccount(c *fiber.Ctx) error {
	accApi := new(api_models.PostAccountRequest)
	if err := c.BodyParser(accApi); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"error": err.Error(),
		})
	}

	acc, err := commands.CreateAccount(accApi.AuthId, accApi.AuthPw)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(api_models.AccountToPostAccountResponse(&acc))
}

func GetAccounts(c *fiber.Ctx) error {
	accs, err := commands.GetAllAccounts()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(api_models.AccountsToGetAccountResponses(accs))
}

func DeleteAccount(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	accountId := claims["account_id"].(string)

	if err := commands.DeleteAccount(accountId); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"error": err.Error(),
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
