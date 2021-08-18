package api

import (
	"github.com/datti-to/purrmannplus-backend/api/controllers"
	"github.com/gofiber/fiber/v2"
)

func InitRoutes(app *fiber.App) {
	app.Get("/health", controllers.GetHealth)

	app.Post("/signal-code-message", nil)
}
