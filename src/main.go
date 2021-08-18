package main

import (
	"github.com/datti-to/purrmannplus-backend/api"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	api.InitRoutes(app)

	app.Listen(":3000")
}
