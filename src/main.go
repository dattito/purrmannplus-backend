package main

import (
	"log"

	"github.com/datti-to/purrmannplus-backend/api"
	"github.com/datti-to/purrmannplus-backend/database"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	api.InitRoutes(app)

	err := database.Init()
	if err != nil {
		log.Fatal(err)
	}

	app.Listen(":3000")
}
