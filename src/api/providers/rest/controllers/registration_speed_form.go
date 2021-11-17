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
	"github.com/dattito/purrmannplus-backend/services/substitutions"
	"github.com/dattito/purrmannplus-backend/utils"
	"github.com/dattito/purrmannplus-backend/utils/logging"
	"github.com/gofiber/fiber/v2"
	"github.com/nyaruka/phonenumbers"
)

func SaveCustomSubstitutionCredentials(c *fiber.Ctx, customSubstitutionAuthId, customSubstitutionAuthPw string) error {
	session, err := session.SessionStore.Get(c)
	if err != nil {
		return err
	}

	session.Set("custom_substitution_auth_id", customSubstitutionAuthId)
	session.Set("custom_substitution_auth_pw", customSubstitutionAuthPw)

	return session.Save()
}

func SaveNeedsCustomSubstitutionCredentials(c *fiber.Ctx) error {
	session, err := session.SessionStore.Get(c)
	if err != nil {
		return err
	}

	session.Set("needs_custom_substitution_credentials", true)

	return session.Save()
}

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

func sendConfirmationCode(c *fiber.Ctx) error {
	session, err := session.SessionStore.Get(c)
	if err != nil {
		return err
	}

	phoneNumber := session.Get("phone_number")
	code := session.Get("code")

	if phoneNumber == nil || code == nil {
		return errors.New("phone number or code not found in session")
	}

	return signal_message_sender.SignalMessageSender.Send(
		fmt.Sprintf("Willkommen bei PurrmannPlus! Dein Bestätigungscode lautet: %s", code),
		phoneNumber.(string),
	)
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
		if _, err := commands.GetAccountByCredentials(pr.Username, pr.Password); err != nil {
			if !errors.Is(err, &db_errors.ErrRecordNotFound) {
				logging.Errorf("Error getting account by credentials: %v", err)
				return internalServerErrorResponse
			}
		} else {
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

		err = SaveRequestInSession(c, pr.Username, pr.Password, validNumber, code)
		if err != nil {
			logging.Errorf("Error saving request in session: %v", err)
			return internalServerErrorResponse
		}

		ok, err := substitutions.CheckCredentials(pr.Username, pr.Password)
		if err != nil {
			logging.Errorf("Error checking credentials: %v", err)
			return internalServerErrorResponse
		}

		if !ok {
			if err := SaveNeedsCustomSubstitutionCredentials(c); err != nil {
				logging.Errorf("Error saving needs custom substitution credentials: %v", err)
				return internalServerErrorResponse
			}

			return c.Redirect(routes.RegistrationSpeedFormSubstitutionCredentialsRoute)
		}

		if err := sendConfirmationCode(c); err != nil {
			logging.Errorf("Error sending confirmation code: %v", err)
			session, err := session.SessionStore.Get(c)
			if err != nil {
				return internalServerErrorResponse
			}
			session.Destroy()
			return internalServerErrorResponse
		}

		return c.Redirect(routes.RegistrationSpeedFormValidationRoute)
	} else {
		return fiber.ErrMethodNotAllowed
	}
}

func SubstitutionCredentialsSpeedForm(c *fiber.Ctx) error {
	internalServerErrorResponse := c.Status(fiber.StatusInternalServerError).Render("registration_speed_form_substitution_credentials", fiber.Map{
		"FormPostRoute": routes.RegistrationSpeedFormSubstitutionCredentialsRoute,
		"ErrorMessage":  "Etwas ist schiefgelaufen...",
	}, "layouts/main")

	session, err := session.SessionStore.Get(c)
	if err != nil {
		session.Destroy()
		logging.Errorf("Error getting session: %v", err)
		return internalServerErrorResponse
	}

	if c.Method() == fiber.MethodGet {
		needsCustomSubstitutionCredentials := session.Get("needs_custom_substitution_credentials")
		if needsCustomSubstitutionCredentials == nil || needsCustomSubstitutionCredentials == false {
			return c.Redirect(routes.RegistrationSpeedFormRoute)
		}
		return c.Render("registration_speed_form_substitution_credentials", fiber.Map{
			"FormPostRoute": routes.RegistrationSpeedFormSubstitutionCredentialsRoute,
		}, "layouts/main")
	} else if c.Method() == fiber.MethodPost {
		var pr models.PostCustomSubsitutionCredentialsRequest
		if err := c.BodyParser(&pr); err != nil {
			session.Destroy()
			logging.Errorf("Error parsing body: %v", err)
			return internalServerErrorResponse
		}

		pr.AuthId = strings.ToLower(pr.AuthId)
		ok, err := substitutions.CheckCredentials(pr.AuthId, pr.AuthPw)
		if err != nil {
			session.Destroy()
			logging.Errorf("Error checking substitution credentials: %v", err)
			return internalServerErrorResponse
		}
		if !ok {
			return c.Status(fiber.StatusUnauthorized).Render("registration_speed_form_substitution_credentials", fiber.Map{
				"FormPostRoute": routes.RegistrationSpeedFormSubstitutionCredentialsRoute,
				"ErrorMessage":  "Falsche Anmeldedaten",
			}, "layouts/main")
		}
		if err := SaveCustomSubstitutionCredentials(c, pr.AuthId, pr.AuthPw); err != nil {
			logging.Errorf("Error saving custom substitution credentials: %v", err)
			session.Destroy()
			return internalServerErrorResponse
		}

		if err := sendConfirmationCode(c); err != nil {
			logging.Errorf("Error sending confirmation code: %v", err)
			session.Destroy()
			return internalServerErrorResponse
		}

		return c.Redirect(routes.RegistrationSpeedFormValidationRoute)
	}
	return fiber.ErrMethodNotAllowed
}

func ValidateRegistrationSpeedForm(c *fiber.Ctx) error {
	internalServerErrorResponse := c.Status(fiber.StatusInternalServerError).Render("registration_speed_form_pn_validate", fiber.Map{
		"FormPostRoute": routes.RegistrationSpeedFormValidationRoute,
		"ErrorMessage":  "Etwas ist schiefgelaufen...",
	}, "layouts/main")

	session, err := session.SessionStore.Get(c)
	if err != nil {
		return internalServerErrorResponse
	}

	if c.Method() == fiber.MethodGet {
		username := session.Get("username")
		if username == nil {
			return c.Redirect(routes.RegistrationSpeedFormRoute)
		}

		return c.Render("registration_speed_form_pn_validate", fiber.Map{
			"FormPostRoute": routes.RegistrationSpeedFormValidationRoute,
		}, "layouts/main")
	}

	if c.Method() == fiber.MethodPost {
		username := session.Get("username")
		if username == nil {
			return c.Redirect(routes.RegistrationSpeedFormRoute)
		}
		needsCustomSubstitutionCredentials := session.Get("needs_custom_substitution_credentials")
		customSubstitionAuthId := session.Get("custom_substitution_auth_id")
		customSubstitionAuthPw := session.Get("custom_substitution_auth_pw")
		if needsCustomSubstitutionCredentials != nil &&
			needsCustomSubstitutionCredentials == true &&
			(customSubstitionAuthId == nil || customSubstitionAuthPw == nil) {
			return c.Redirect(routes.RegistrationSpeedFormSubstitutionCredentialsRoute)
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

		if _, err := commands.AddAccountToMoodleAssignmentUpdater(acc.Id); err != nil {
			session.Destroy()
			return internalServerErrorResponse
		}

		if needsCustomSubstitutionCredentials != nil && needsCustomSubstitutionCredentials == true {
			if _, err := commands.AddAccountToSubstitutionUpdaterWithCustomCredentials(acc.Id, session.Get("custom_substitution_auth_id").(string), session.Get("custom_substitution_auth_pw").(string)); err != nil {
				session.Destroy()
				logging.Errorf("Error adding account to substitution updater with custom credentials: %v", err)
				return internalServerErrorResponse
			}
		} else {
			if _, err := commands.AddAccountToSubstitutionUpdater(acc.Id); err != nil {
				session.Destroy()
				logging.Errorf("Error adding account to substitution updater: %v", err)
				return internalServerErrorResponse
			}
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
