package controllers

import (
	"errors"
	"fmt"
	"strings"

	"github.com/dattito/purrmannplus-backend/api/providers/rest/models"
	"github.com/dattito/purrmannplus-backend/api/providers/rest/routes"
	"github.com/dattito/purrmannplus-backend/api/providers/rest/session"
	"github.com/dattito/purrmannplus-backend/app/commands"
	"github.com/dattito/purrmannplus-backend/config"
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

func RegistrationSpeedForm(c *fiber.Ctx) error {
	if c.Method() == fiber.MethodGet {
		return c.Render("registration_speed_form", fiber.Map{
			"InfoRoute":        routes.RegistrationSpeedFormInfoRoute,
			"FormPostRoute":    routes.RegistrationSpeedFormRoute,
			"ContactEmail":     config.CONTACT_EMAIL,
			"ContactInstagram": config.CONTACT_INSTAGRAM,
		}, "layouts/main")
	} else if c.Method() == fiber.MethodPost {
		internalServerErrorResponse := c.Status(fiber.StatusInternalServerError).Render("registration_speed_form", fiber.Map{
			"InfoRoute":        routes.RegistrationSpeedFormInfoRoute,
			"FormPostRoute":    routes.RegistrationSpeedFormRoute,
			"ErrorMessage":     "Etwas ist schiefgelaufen...",
			"ContactEmail":     config.CONTACT_EMAIL,
			"ContactInstagram": config.CONTACT_INSTAGRAM,
		}, "layouts/main")

		var pr models.PostRegistrationSpeedFormRequest
		if err := c.BodyParser(&pr); err != nil {
			logging.Errorf("Error parsing body: %v", err)
			return internalServerErrorResponse
		}

		pr.Username = strings.ToLower(pr.Username)
		if len(pr.Username) < 4 && utils.NumberInString(pr.Username) {
			return c.Status(fiber.StatusBadRequest).Render("registration_speed_form", fiber.Map{
				"InfoRoute":        routes.RegistrationSpeedFormInfoRoute,
				"FormPostRoute":    routes.RegistrationSpeedFormRoute,
				"ErrorMessage":     "Momentan können sich nur Schüler der Oberstufe anmelden...",
				"ContactEmail":     config.CONTACT_EMAIL,
				"ContactInstagram": config.CONTACT_INSTAGRAM,
			}, "layouts/main")
		}

		correct, err := commands.CheckCredentials(pr.Username, pr.Password)
		if err != nil {
			logging.Errorf("Error checking credentials: %v", err)
			return internalServerErrorResponse
		}

		if !correct {
			return c.Status(fiber.StatusUnauthorized).Render("registration_speed_form", fiber.Map{
				"InfoRoute":        routes.RegistrationSpeedFormInfoRoute,
				"FormPostRoute":    routes.RegistrationSpeedFormRoute,
				"ErrorMessage":     "Falsche Anmeldedaten",
				"ContactEmail":     config.CONTACT_EMAIL,
				"ContactInstagram": config.CONTACT_INSTAGRAM,
			}, "layouts/main")
		}

		// Check if accounts already exist
		account, err := commands.GetAccountByCredentials(pr.Username, pr.Password)
		if err != &db_errors.ErrRecordNotFound && err != nil {
			logging.Errorf("Error getting account by credentials: %v", err)
			return internalServerErrorResponse
		}

		if account.Username != "" {
			return c.Status(fiber.StatusUnauthorized).Render("registration_speed_form", fiber.Map{
				"InfoRoute":        routes.RegistrationSpeedFormInfoRoute,
				"FormPostRoute":    routes.RegistrationSpeedFormRoute,
				"ErrorMessage":     "Das Konto exestiert bereits",
				"ContactEmail":     config.CONTACT_EMAIL,
				"ContactInstagram": config.CONTACT_INSTAGRAM,
			}, "layouts/main")
		}

		validNumber, err := utils.FormatPhoneNumber(pr.PhoneNumber)
		if err != nil {
			if errors.Is(err, phonenumbers.ErrNotANumber) {
				return c.Status(fiber.StatusInternalServerError).Render("registration_speed_form", fiber.Map{
					"InfoRoute":        routes.RegistrationSpeedFormInfoRoute,
					"FormPostRoute":    routes.RegistrationSpeedFormRoute,
					"ErrorMessage":     "Bitte gebe eine gültige Telefonnummer an",
					"ContactEmail":     config.CONTACT_EMAIL,
					"ContactInstagram": config.CONTACT_INSTAGRAM,
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

		return c.Redirect(routes.RegistrationSpeedFormValidationRoute)
	} else {
		return fiber.ErrMethodNotAllowed
	}
}

func ValidateRegistrationSpeedForm(c *fiber.Ctx) error {

	internalServerErrorResponse := c.Status(fiber.StatusInternalServerError).Render("registration_speed_form_pn_validate", fiber.Map{
		"FormPostRoute": routes.RegistrationSpeedFormValidationRoute,
		"ErrorMessage":  "Etwas ist schiefgelaufen...",
	}, "layouts/main")

	if c.Method() == fiber.MethodGet {
		session, err := session.SessionStore.Get(c)
		if err != nil {
			return c.Render("registration_speed_form_pn_validate", fiber.Map{
				"FormPostRoute": routes.RegistrationSpeedFormValidationRoute,
			}, "layouts/main")
		}

		username := session.Get("username")
		if username == nil {
			return c.Redirect(routes.RegistrationSpeedFormRoute)
		}

		return c.Render("registration_speed_form_pn_validate", fiber.Map{
			"FormPostRoute": routes.RegistrationSpeedFormValidationRoute,
		}, "layouts/main")
	}

	if c.Method() == fiber.MethodPost {
		session, err := session.SessionStore.Get(c)
		if err != nil {
			return internalServerErrorResponse
		}

		username := session.Get("username")
		if username == nil {
			return c.Redirect(routes.RegistrationSpeedFormRoute)
		}

		var pr models.PostValidateRegistrationSpeedFormRequest
		if err := c.BodyParser(&pr); err != nil {
			logging.Errorf("Error parsing body: %v", err)
			session.Destroy()
			return internalServerErrorResponse
		}

		if pr.Code != session.Get("code") {
			return c.Status(fiber.StatusUnauthorized).Render("registration_speed_form_pn_validate", fiber.Map{
				"FormPostRoute": routes.RegistrationSpeedFormValidationRoute,
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

		if _, err := commands.AddAccountToSubstitutionUpdater(acc.Id); err != nil {
			session.Destroy()
			return internalServerErrorResponse
		}

		if _, err := commands.AddAccountToMoodleAssignmentUpdater(acc.Id); err != nil {
			session.Destroy()
			return internalServerErrorResponse
		}

		if err := signal_message_sender.SignalMessageSender.Send(
			fmt.Sprintf("Dein Account '%s' wurde mit dieser Telefonnummer verbunden. Ab jetzt erhälst du über diesen Chat neue Infos über Vertretungen und Moodle-Aufgaben!",
				acc.Username),
			session.Get("phone_number").(string),
		); err != nil {
			return internalServerErrorResponse
		}

		session.Destroy()
		return c.Redirect(routes.RegistrationSpeedFormFinishRoute)
	}
	return fiber.ErrMethodNotAllowed
}

func FinishRegistrationSpeedForm(c *fiber.Ctx) error {
	return c.Render("registration_speed_form_finish", fiber.Map{}, "layouts/main")
}

func InfoRegsitrationSpeedForm(c *fiber.Ctx) error {
	return c.Render("registration_speed_form_info", fiber.Map{
		"FormRoute": routes.RegistrationSpeedFormRoute,
	}, "layouts/main")
}
