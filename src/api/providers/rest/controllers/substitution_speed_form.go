package controllers

import (
	"errors"
	"fmt"

	"github.com/dattito/purrmannplus-backend/api/providers/rest/models"
	"github.com/dattito/purrmannplus-backend/api/providers/rest/routes"
	"github.com/dattito/purrmannplus-backend/api/providers/rest/session"
	"github.com/dattito/purrmannplus-backend/app/commands"
	"github.com/dattito/purrmannplus-backend/services/signal_message_sender"
	"github.com/dattito/purrmannplus-backend/utils"
	"github.com/dattito/purrmannplus-backend/utils/logging"
	"github.com/gofiber/fiber/v2"
	"github.com/nyaruka/phonenumbers"
)

func SaveRequestInSession(c *fiber.Ctx, pr models.PostSubstitutionSpeedFormRequest) error {
	session, err := session.SessionStore.Get(c)
	if err != nil {
		return err
	}

	session.Set("username", pr.Username)
	session.Set("password", pr.Password)
	session.Set("phone_number", pr.PhoneNumber)

	return session.Save()
}

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

		validNumber, err := utils.FormatPhoneNumber(pr.PhoneNumber)
		if err != nil {
			if errors.Is(err, phonenumbers.ErrNotANumber) {
				return c.Status(fiber.StatusInternalServerError).Render("substitution_speed_form_full", fiber.Map{
					"FormPostRoute": routes.SubstitutionSpeedFormRoute,
					"ErrorMessage":  "Bitte gebe eine gültige Telefonnummer an",
				}, "layouts/main")
			}
			logging.Errorf("Error formatting number: %v", err)
			return c.Status(fiber.StatusInternalServerError).Render("substitution_speed_form_full", fiber.Map{
				"FormPostRoute": routes.SubstitutionSpeedFormRoute,
				"ErrorMessage":  "Etwas ist schiefgelaufen...",
			}, "layouts/main")
		}

		validationCode := utils.GenerateValidationCode(6)

		if err := signal_message_sender.SignalMessageSender.Send(
			fmt.Sprintf("Willkommen bei PurrmannPlus! Dein Bestätigungscode lautet: %s", validationCode),
			validNumber,
		); err != nil {
			return c.Status(fiber.StatusInternalServerError).Render("substitution_speed_form_full", fiber.Map{
				"FormPostRoute": routes.SubstitutionSpeedFormRoute,
				"ErrorMessage":  "Etwas ist schiefgelaufen...",
			}, "layouts/main")
		}

		err = SaveRequestInSession(c, pr)
		if err != nil {
			logging.Errorf("Error saving request in session: %v", err)
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
		session, err := session.SessionStore.Get(c)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).Render("substitution_speed_pn_validate", fiber.Map{}, "layouts/main")
		}

		username := session.Get("username")
		if username == nil {
			return c.Redirect(routes.SubstitutionSpeedFormRoute)
		}

		return c.Render("substitution_speed_pn_validate", fiber.Map{}, "layouts/main")
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
