package api

import (
	"github.com/datti-to/purrmannplus-backend/api/controllers"
	"github.com/gofiber/fiber/v2"
)

func InitPublicRoutes(app *fiber.App) {
	app.Get("/health", controllers.GetHealth)

	app.Post("/login", controllers.AccountLogin)

	app.Post("/accounts", controllers.AddAccount)
	app.Get("/accounts", controllers.GetAccounts)
}

func InitJWTRestrictedRoutes(app *fiber.App) {
	app.Post("/accounts/phone_number_confirmation", controllers.SendPhoneNumberConfirmationLink)
}
