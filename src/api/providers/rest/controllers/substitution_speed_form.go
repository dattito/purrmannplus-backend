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

func SaveRequestInSession(c *fiber.Ctx, pr models.PostSubstitutionSpeedFormRequest, code string) error {
	session, err := session.SessionStore.Get(c)
	if err != nil {
		return err
	}

	session.Set("username", pr.Username)
	session.Set("password", pr.Password)
	session.Set("phone_number", pr.PhoneNumber)
	session.Set("code", code)

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

		code := utils.GenerateValidationCode(6)

		if err := signal_message_sender.SignalMessageSender.Send(
			fmt.Sprintf("Willkommen bei PurrmannPlus! Dein Bestätigungscode lautet: %s", code),
			validNumber,
		); err != nil {
			return c.Status(fiber.StatusInternalServerError).Render("substitution_speed_form_full", fiber.Map{
				"FormPostRoute": routes.SubstitutionSpeedFormRoute,
				"ErrorMessage":  "Etwas ist schiefgelaufen...",
			}, "layouts/main")
		}

		err = SaveRequestInSession(c, pr, code)
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
			return c.Status(fiber.StatusInternalServerError).Render("substitution_speed_pn_validate", fiber.Map{
				"FormPostRoute": routes.ValidateSubstitutionSpeedFormRoute,
			}, "layouts/main")
		}

		username := session.Get("username")
		if username == nil {
			return c.Redirect(routes.SubstitutionSpeedFormRoute)
		}

		return c.Render("substitution_speed_pn_validate", fiber.Map{
			"FormPostRoute": routes.ValidateSubstitutionSpeedFormRoute,
		}, "layouts/main")
	}

	if c.Method() == fiber.MethodPost {
		session, err := session.SessionStore.Get(c)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).Render("substitution_speed_pn_validate", fiber.Map{
				"FormPostRoute": routes.ValidateSubstitutionSpeedFormRoute,
			}, "layouts/main")
		}

		username := session.Get("username")
		if username == nil {
			return c.Redirect(routes.SubstitutionSpeedFormRoute)
		}

		var pr models.PostValidateSubstitutionSpeedFormRequest
		if err := c.BodyParser(&pr); err != nil {
			logging.Errorf("Error parsing body: %v", err)
			session.Destroy()
			return c.Status(fiber.StatusInternalServerError).Render("substitution_speed_pn_validate", fiber.Map{
				"FormPostRoute": routes.ValidateSubstitutionSpeedFormRoute,
				"ErrorMessage":  "Etwas ist schiefgelaufen...",
			}, "layouts/main")
		}

		if pr.Code != session.Get("code") {
			return c.Status(fiber.StatusUnauthorized).Render("substitution_speed_pn_validate", fiber.Map{
				"FormPostRoute": routes.ValidateSubstitutionSpeedFormRoute,
				"ErrorMessage":  "Falscher Code",
			}, "layouts/main")
		}

		acc, userErr, internalErr := commands.CreateAccount(session.Get("username").(string), session.Get("password").(string))
		if internalErr != nil {
			session.Destroy()
			return c.Status(fiber.StatusInternalServerError).Render("substitution_speed_pn_validate", fiber.Map{
				"FormPostRoute": routes.ValidateSubstitutionSpeedFormRoute,
				"ErrorMessage":  "Etwas ist schiefgelaufen...",
			}, "layouts/main")
		}
		if userErr != nil {
			return c.Status(fiber.StatusInternalServerError).Render("substitution_speed_pn_validate", fiber.Map{
				"FormPostRoute": routes.ValidateSubstitutionSpeedFormRoute,
				"ErrorMessage":  "Etwas ist schiefgelaufen...",
			}, "layouts/main")
		}

		_, userErr, internalErr = commands.AddAccountInfo(acc.Id, session.Get("phoneNumber").(string))
		if internalErr != nil {
			session.Destroy()
			return c.Status(fiber.StatusInternalServerError).Render("substitution_speed_pn_validate", fiber.Map{
				"FormPostRoute": routes.ValidateSubstitutionSpeedFormRoute,
				"ErrorMessage":  "Etwas ist schiefgelaufen...",
			}, "layouts/main")
		}
		if userErr != nil {
			return c.Status(fiber.StatusInternalServerError).Render("substitution_speed_pn_validate", fiber.Map{
				"FormPostRoute": routes.ValidateSubstitutionSpeedFormRoute,
				"ErrorMessage":  "Etwas ist schiefgelaufen...",
			}, "layouts/main")
		}

		_, err = commands.AddToSubstitutionUpdater(acc.Id)
		if err != nil {
			session.Destroy()
			return c.Status(fiber.StatusInternalServerError).Render("substitution_speed_pn_validate", fiber.Map{
				"FormPostRoute": routes.ValidateSubstitutionSpeedFormRoute,
				"ErrorMessage":  "Etwas ist schiefgelaufen...",
			}, "layouts/main")
		}

		session.Destroy()
		return c.Redirect(routes.FinishSubstitutionSpeedFormRoute)
	}
	return fiber.ErrMethodNotAllowed
}

func FinishSubstitutionSpeedForm(c *fiber.Ctx) error {
	return c.Render("substitution_speed_finish", fiber.Map{}, "layouts/main")
}
