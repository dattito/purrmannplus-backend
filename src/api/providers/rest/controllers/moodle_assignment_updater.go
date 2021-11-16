package controllers

import (
	"github.com/dattito/purrmannplus-backend/app/commands"
	"github.com/dattito/purrmannplus-backend/utils/logging"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

func AddAccountToMoodleAssignmentUpdater(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	accountId := claims["account_id"].(string)

	ok, err := commands.ValidAccountId(accountId)
	if err != nil {
		logging.Errorf("Error validating account id: %s", err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"error": "Something went wrong",
		})
	}

	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"error": "account not found",
		})
	}

	user_err, db_err := commands.AddAccountToMoodleAssignmentUpdater(accountId)

	if db_err != nil {
		logging.Errorf("Error while adding account to moodle assignment updater: %s", db_err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"error": "Something went wrong",
		})
	}

	if user_err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"error": user_err.Error(),
		})
	}

	return c.SendStatus(fiber.StatusCreated)
}

func RemoveAccountFromMoodleAssignmentUpdater(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	accountId := claims["account_id"].(string)

	err := commands.RemoveAccountFromMoodleAssignmentUpdater(accountId)

	if err != nil {
		logging.Errorf("Error while removing account from moodle assignment updater: %s", err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"error": "Something went wrong",
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
