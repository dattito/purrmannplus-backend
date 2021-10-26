package controllers

import (
	"fmt"

	"github.com/dattito/purrmannplus-backend/api/providers/rest/models"
	"github.com/dattito/purrmannplus-backend/api/providers/rest/routes"
	"github.com/dattito/purrmannplus-backend/api/providers/rest/session"
	"github.com/dattito/purrmannplus-backend/app/commands"
	"github.com/dattito/purrmannplus-backend/utils"
	"github.com/dattito/purrmannplus-backend/utils/logging"
	"github.com/gofiber/fiber/v2"
)

func SubstitutionSpeedForm(c *fiber.Ctx) error {
	if c.Method() == fiber.MethodGet {
		return c.Render("substitution_speed_form_full", fiber.Map{
			"FormPostRoute": routes.SubstitutionSpeedFormRoute,
		}, "layouts/main")
	} else if c.Method() == fiber.MethodPost {
		var pr models.PostSubstitutionSpeedFormRequest
		if err := c.BodyParser(&pr); err != nil {
			logging.Errorf("Error parsing body: %v", err)
			return c.Status(fiber.StatusInternalServerError).Render("substitution_speed_form_full", fiber.Map{
				"FormPostRoute": routes.SubstitutionSpeedFormRoute,
				"ErrorMessage":  "Etwas ist schiefgelaufen...",
			}, "layouts/main")
		}

		fmt.Printf("username: %s", pr.Username)

		correct, err := commands.CheckCredentials(pr.Username, pr.Password)
		if err != nil {
			logging.Errorf("Error checking credentials: %v", err)
			return c.Status(fiber.StatusInternalServerError).Render("substitution_speed_form_full", fiber.Map{
				"FormPostRoute": routes.SubstitutionSpeedFormRoute,
				"ErrorMessage":  "Etwas ist schiefgelaufen...",
			}, "layouts/main")
		}

		if !correct {
			return c.Status(fiber.StatusUnauthorized).Render("substitution_speed_form_full", fiber.Map{
				"FormPostRoute": routes.SubstitutionSpeedFormRoute,
				"ErrorMessage":  "Falsche Anmeldedaten",
			}, "layouts/main")
		}

		_, err = utils.FormatPhoneNumber(pr.PhoneNumber)
		if err != nil {
			logging.Errorf("Error formatting number: %v", err)
			return c.Status(fiber.StatusInternalServerError).Render("substitution_speed_form_full", fiber.Map{
				"FormPostRoute": routes.SubstitutionSpeedFormRoute,
				"ErrorMessage":  "Etwas ist schiefgelaufen...",
			}, "layouts/main")
		}
		return c.Redirect(routes.ValidateSubstitutionSpeedFormRoute)
	} else {
		return fiber.ErrMethodNotAllowed
	}
}

func ValidateSubstitutionSpeedForm(c *fiber.Ctx) error {
	if c.Method() == fiber.MethodGet {
		return c.Render("substitution_speed_form_pn_validate", fiber.Map{}, "layouts/main")
	}

	if c.Method() == fiber.MethodPost {
		_, err := session.SessionStore.Get(c)
		if err != nil {
			logging.Errorf("Can't get session from session store: %v", err.Error())
			return c.Status(500).JSON(&fiber.Map{
				"error": "Something went wrong",
			})
		}

	}

	return fiber.ErrMethodNotAllowed
}
