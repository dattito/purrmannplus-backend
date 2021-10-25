package controllers

import (
	"github.com/dattito/purrmannplus-backend/api/providers/rest/session"
	"github.com/dattito/purrmannplus-backend/utils/logging"
	"github.com/gofiber/fiber/v2"
)

func GetSubstitutionSpeedForm(c *fiber.Ctx) error {
	if c.Method() == fiber.MethodGet {
		return c.Render("substitution_speed_form_full", fiber.Map{}, "layouts/main")
	}

	if c.Method() == fiber.MethodPost {
		sess, err := session.SessionStore.Get(c)
		if err != nil {
			logging.Errorf("Can't get session from session store: %v", err.Error())
			return c.Status(500).JSON(&fiber.Map{
				"error": "Something went wrong",
			})
		}

		if sess.Get("") == nil {
			// TODO
		}
	}

	return fiber.ErrMethodNotAllowed

	// sess, err := session.SessionStore.Get(c)
	// if err != nil {
	// logging.Errorf("Can't get session from session store: %v", err.Error())
	// 	return c.Status(500).JSON(&fiber.Map{
	// 		"error": "Something went wrong",
	// 	})
	// }
}
