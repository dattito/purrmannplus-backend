package controllers

import "github.com/gofiber/fiber/v2"

func GetSubstitutionSpeedForm(c *fiber.Ctx) error {
	return c.Render("index", fiber.Map{}, "layouts/main")
}
