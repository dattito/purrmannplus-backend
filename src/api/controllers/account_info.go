package controllers

import "github.com/gofiber/fiber/v2"

func SendPhoneNumberConfirmationLink(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusOK)
}
