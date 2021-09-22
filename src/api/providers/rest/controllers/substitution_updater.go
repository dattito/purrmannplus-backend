package controllers

import (
	"log"

	"github.com/datti-to/purrmannplus-backend/app/commands"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

func RegisterToSubstitutionUpdater(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	accountId := claims["account_id"].(string)

	ok, err := commands.ValidAccountId(accountId)
	if err != nil {
		log.Println(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"error": "Something went wrong",
		})
	}

	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"error": "account not found",
		})
	}

	if err := commands.AddToSubstitutionUpdater(accountId); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"error": err.Error(),
		})
	}

	return c.SendStatus(fiber.StatusCreated)
}

func UnregisterFromSubstitutionUpdater(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	accountId := claims["account_id"].(string)

	err := commands.RemoveFromSubstitutionUpdater(accountId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"error": err.Error(),
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
