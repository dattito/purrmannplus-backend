package controllers

import (
	api_models "github.com/dattito/purrmannplus-backend/api/providers/rest/models"
	"github.com/dattito/purrmannplus-backend/app/commands"
	"github.com/dattito/purrmannplus-backend/utils/logging"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

// Creates a new account and returns the account id
func AddAccount(c *fiber.Ctx) error {
	accApi := new(api_models.PostAccountRequest)
	if err := c.BodyParser(accApi); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"error": err.Error(),
		})
	}

	acc, user_err, db_err := commands.CreateAccount(accApi.Username, accApi.Password)

	if user_err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"error": user_err.Error(),
		})
	}

	if db_err != nil {
		logging.Errorf("Error while creating account: %v", db_err.Error())

		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"error": "Something went wrong",
		})
	}

	return c.JSON(api_models.AccountToPostAccountResponse(&acc))
}

// Retruns the account_id and the credentials of all accounts
func GetAccounts(c *fiber.Ctx) error {
	accs, err := commands.GetAllAccounts()
	if err != nil {
		logging.Errorf("Error while getting accounts: %v", err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"error": "Something went wrong",
		})
	}

	return c.JSON(api_models.AccountsToGetAccountResponses(accs))
}

// Delets an account
func DeleteAccount(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	accountId := claims["account_id"].(string)

	if err := commands.DeleteAccount(accountId); err != nil {
		logging.Errorf("Error while deleting account: %v", err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"error": err.Error(),
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
