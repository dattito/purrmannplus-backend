package controllers

import (
	"github.com/gofiber/fiber/v2"
)

func GetHealth(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNoContent)
}
