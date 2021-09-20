package controllers

import (
	"github.com/datti-to/purrmannplus-backend/app/commands"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

func RegisterToSubstitutionUpdater(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	accountId := claims["account_id"].(string)

	err := commands.AddToSubstitutionUpdater(accountId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"error": err.Error(),
		})
	}

	return c.SendStatus(fiber.StatusCreated)
}
