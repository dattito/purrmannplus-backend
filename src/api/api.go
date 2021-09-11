package api

import (
	"fmt"

	"github.com/datti-to/purrmannplus-backend/config"
	"github.com/gofiber/fiber/v2"
)

var App *fiber.App

func Init() {
	App = fiber.New()

	InitPublicRoutes(App)

	InitJWTMiddlewares(App)

	InitJWTRestrictedRoutes(App)
}

func StartListening() error {
	return App.Listen(fmt.Sprintf(":%d", config.LISTENING_PORT))
}
