package controllers

import (
	"errors"
	"fmt"
	"strings"

	"github.com/dattito/purrmannplus-backend/api/providers/rest/models"
	"github.com/dattito/purrmannplus-backend/api/providers/rest/routes"
	"github.com/dattito/purrmannplus-backend/api/providers/rest/session"
	"github.com/dattito/purrmannplus-backend/app/commands"
	db_errors "github.com/dattito/purrmannplus-backend/database/errors"
	"github.com/dattito/purrmannplus-backend/services/signal_message_sender"
	"github.com/dattito/purrmannplus-backend/utils"
	"github.com/dattito/purrmannplus-backend/utils/logging"
	"github.com/gofiber/fiber/v2"
	"github.com/nyaruka/phonenumbers"
)

func SaveRequestInSession(c *fiber.Ctx, username, password, phoneNumber, code string) error {
	session, err := session.SessionStore.Get(c)
	if err != nil {
		return err
	}

	session.Set("username", username)
	session.Set("password", password)
	session.Set("phone_number", phoneNumber)
	session.Set("code", code)

	return session.Save()
}

func SubstitutionSpeedForm(c *fiber.Ctx) error {
	if c.Method() == fiber.MethodGet {
		return c.Render("substitution_speed_form_full", fiber.Map{
			"InfoRoute":     routes.SubstitutionSpeedFormInfoRoute,
			"FormPostRoute": routes.SubstitutionSpeedFormRoute,
		}, "layouts/main")
	} else if c.Method() == fiber.MethodPost {
		internalServerErrorResponse := c.Status(fiber.StatusInternalServerError).Render("substitution_speed_form_full", fiber.Map{
			"InfoRoute":     routes.SubstitutionSpeedFormInfoRoute,
			"FormPostRoute": routes.SubstitutionSpeedFormRoute,
			"ErrorMessage":  "Etwas ist schiefgelaufen...",
		}, "layouts/main")

		var pr models.PostSubstitutionSpeedFormRequest
		if err := c.BodyParser(&pr); err != nil {
			logging.Errorf("Error parsing body: %v", err)
			return internalServerErrorResponse
		}

		pr.Username = strings.ToLower(pr.Username)
		if len(pr.Username) < 4 && utils.NumberInString(pr.Username) {
			return c.Status(fiber.StatusBadRequest).Render("substitution_speed_form_full", fiber.Map{
				"InfoRoute":     routes.SubstitutionSpeedFormInfoRoute,
				"FormPostRoute": routes.SubstitutionSpeedFormRoute,
				"ErrorMessage":  "Momentan können sich nur Schüler der Oberstufe anmelden...",
			}, "layouts/main")
		}

		correct, err := commands.CheckCredentials(pr.Username, pr.Password)
		if err != nil {
			logging.Errorf("Error checking credentials: %v", err)
			return internalServerErrorResponse
		}

		if !correct {
			return c.Status(fiber.StatusUnauthorized).Render("substitution_speed_form_full", fiber.Map{
				"InfoRoute":     routes.SubstitutionSpeedFormInfoRoute,
				"FormPostRoute": routes.SubstitutionSpeedFormRoute,
				"ErrorMessage":  "Falsche Anmeldedaten",
			}, "layouts/main")
		}

		// Check if accounts already exist
		account, err := commands.GetAccountByCredentials(pr.Username, pr.Password)
		if err != &db_errors.ErrRecordNotFound && err != nil {
			logging.Errorf("Error getting account by credentials: %v", err)
			return internalServerErrorResponse
		}

		if account.Username != "" {
			return c.Status(fiber.StatusUnauthorized).Render("substitution_speed_form_full", fiber.Map{
				"InfoRoute":     routes.SubstitutionSpeedFormInfoRoute,
				"FormPostRoute": routes.SubstitutionSpeedFormRoute,
				"ErrorMessage":  "Das Konto exestiert bereits",
			}, "layouts/main")
		}

		validNumber, err := utils.FormatPhoneNumber(pr.PhoneNumber)
		if err != nil {
			if errors.Is(err, phonenumbers.ErrNotANumber) {
				return c.Status(fiber.StatusInternalServerError).Render("substitution_speed_form_full", fiber.Map{
					"InfoRoute":     routes.SubstitutionSpeedFormInfoRoute,
					"FormPostRoute": routes.SubstitutionSpeedFormRoute,
					"ErrorMessage":  "Bitte gebe eine gültige Telefonnummer an",
				}, "layouts/main")
			}
			logging.Errorf("Error formatting number: %v", err)
			return internalServerErrorResponse
		}

		code := utils.GenerateValidationCode(6)

		if err := signal_message_sender.SignalMessageSender.Send(
			fmt.Sprintf("Willkommen bei PurrmannPlus! Dein Bestätigungscode lautet: %s", code),
			validNumber,
		); err != nil {
			return internalServerErrorResponse
		}

		err = SaveRequestInSession(c, pr.Username, pr.Password, validNumber, code)
		if err != nil {
			logging.Errorf("Error saving request in session: %v", err)
			return internalServerErrorResponse
		}

		return c.Redirect(routes.SubstitutionSpeedFormValidationRoute)
	} else {
		return fiber.ErrMethodNotAllowed
	}
}

func ValidateSubstitutionSpeedForm(c *fiber.Ctx) error {

	internalServerErrorResponse := c.Status(fiber.StatusInternalServerError).Render("substitution_speed_pn_validate", fiber.Map{
		"FormPostRoute": routes.SubstitutionSpeedFormValidationRoute,
		"ErrorMessage":  "Etwas ist schiefgelaufen...",
	}, "layouts/main")

	if c.Method() == fiber.MethodGet {
		session, err := session.SessionStore.Get(c)
		if err != nil {
			return c.Render("substitution_speed_pn_validate", fiber.Map{
				"FormPostRoute": routes.SubstitutionSpeedFormValidationRoute,
			}, "layouts/main")
		}

		username := session.Get("username")
		if username == nil {
			return c.Redirect(routes.SubstitutionSpeedFormRoute)
		}

		return c.Render("substitution_speed_pn_validate", fiber.Map{
			"FormPostRoute": routes.SubstitutionSpeedFormValidationRoute,
		}, "layouts/main")
	}

	if c.Method() == fiber.MethodPost {
		session, err := session.SessionStore.Get(c)
		if err != nil {
			return internalServerErrorResponse
		}

		username := session.Get("username")
		if username == nil {
			return c.Redirect(routes.SubstitutionSpeedFormRoute)
		}

		var pr models.PostValidateSubstitutionSpeedFormRequest
		if err := c.BodyParser(&pr); err != nil {
			logging.Errorf("Error parsing body: %v", err)
			session.Destroy()
			return internalServerErrorResponse
		}

		if pr.Code != session.Get("code") {
			return c.Status(fiber.StatusUnauthorized).Render("substitution_speed_pn_validate", fiber.Map{
				"FormPostRoute": routes.SubstitutionSpeedFormValidationRoute,
				"ErrorMessage":  "Falscher Code",
			}, "layouts/main")
		}

		acc, userErr, internalErr := commands.CreateAccount(session.Get("username").(string), session.Get("password").(string))
		if internalErr != nil {
			session.Destroy()
			return internalServerErrorResponse
		}

		if userErr != nil {
			return internalServerErrorResponse
		}

		_, userErr, internalErr = commands.AddAccountInfo(acc.Id, session.Get("phone_number").(string))
		if internalErr != nil {
			session.Destroy()
			return internalServerErrorResponse
		}
		if userErr != nil {
			return internalServerErrorResponse
		}

		_, err = commands.AddToSubstitutionUpdater(acc.Id)
		if err != nil {
			session.Destroy()
			return internalServerErrorResponse
		}

		session.Destroy()
		return c.Redirect(routes.SubstitutionSpeedFormFinishRoute)
	}
	return fiber.ErrMethodNotAllowed
}

func FinishSubstitutionSpeedForm(c *fiber.Ctx) error {
	return c.Render("substitution_speed_finish", fiber.Map{}, "layouts/main")
}

func InfoSubstitutionSpeedForm(c *fiber.Ctx) error {
	return c.Render("substitution_speed_info", fiber.Map{
		"FormRoute": routes.SubstitutionSpeedFormRoute,
	}, "layouts/main")
}
