package api

import (
	"fmt"

	"github.com/datti-to/purrmannplus-backend/config"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
)

var App *fiber.App

func Init() {
	App = fiber.New()

	App.Use(jwtware.New(jwtware.Config{
		SigningKey: []byte(config.JWT_SECRET),
	}))

	InitRoutes(App)
}

func StartListening() error {
	return App.Listen(fmt.Sprintf(":%d", config.LISTENING_PORT))
}
