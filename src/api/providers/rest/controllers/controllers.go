package controllers

import (
	"github.com/dattito/purrmannplus-backend/config"
	"github.com/gofiber/fiber/v2"
)

// Sends an empty response to check if the server is up and running
func GetHealth(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNoContent)
}

// Sends generel information about the server
func About(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"name":    "purrmannplus-backend",
		"source":  "https://github.com/Dattito/purrmannplus-backend",
		"LICENSE": "AGPL-3.0",
		"version": config.DNT_VERSION,
	})
}
