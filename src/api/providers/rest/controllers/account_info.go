package controllers

import (
	"fmt"

	api_models "github.com/dattito/purrmannplus-backend/api/providers/rest/models"
	"github.com/dattito/purrmannplus-backend/api/providers/rest/routes"
	"github.com/dattito/purrmannplus-backend/app/commands"
	"github.com/dattito/purrmannplus-backend/app/models"
	"github.com/dattito/purrmannplus-backend/config"
	"github.com/dattito/purrmannplus-backend/services/signal_message_sender"
	utils_jwt "github.com/dattito/purrmannplus-backend/utils/jwt"
	"github.com/dattito/purrmannplus-backend/utils/logging"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

// Sends a message with a link to the user to confirm his phone number
func SendPhoneNumberConfirmationLink(c *fiber.Ctx) error {
	pr := new(api_models.PostSendPhoneNumberConfirmationLinkRequest)
	if err := c.BodyParser(pr); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"error": err.Error(),
		})
	}

	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	accountId := claims["account_id"].(string)

	ok, err := commands.ValidAccountId(accountId)
	if err != nil {
		logging.Errorf("Error while validating account id: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"error": "Something went wrong",
		})
	}

	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"error": "account not found",
		})
	}

	account_info, err := models.NewAccountInfo(models.Account{Id: accountId}, pr.PhoneNumber)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"error": err.Error(),
		})
	}

	has_phone_number, err := commands.HasPhoneNumber(accountId)
	if err != nil {
		logging.Errorf("Error while checking if account has a phone-number: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"error": "Something went wrong",
		})
	}

	if has_phone_number {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"error": "Phone number already added",
		})
	}

	token, err := utils_jwt.NewAccountIdPhoneNumberToken(account_info.Account.Id, account_info.PhoneNumber)
	if err != nil {
		logging.Errorf("Error while creating token: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"error": "couln't create token",
		})
	}

	text := fmt.Sprintf("Willkommen bei PurrmannPlus. Um deine Telefonnummer zu bestätigen, drücke "+
		"auf den nachfolgenden Link. Er ist 10 Minuten lang gültig. Du hast den Link nicht angefordert? Dann kannst du ihn ignorieren. "+
		"%s/v1%s?token=%s",
		config.API_URL, routes.AddPhoneNumberRoute, token)

	err = signal_message_sender.SignalMessageSender.Send(text, account_info.PhoneNumber)

	if err != nil {
		logging.Errorf("Error while sending signal message: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"error": "couln't send signal message",
		})
	}

	return c.SendStatus(fiber.StatusCreated)
}

// Validates the phone number of the user and adds it to the database
func AddPhoneNumber(c *fiber.Ctx) error {
	p := new(api_models.PostAddPhoneNumberRequest)
	if err := c.QueryParser(p); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"error": err.Error(),
		})
	}

	if p.Token == "" {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"error": "Token is required",
		})
	}

	accountId, phoneNumber, err := utils_jwt.ParseAccountIdPhoneNumberToken(p.Token)
	if err != nil {
		logging.Errorf("Error while parsing token: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"error": "Couln't parse token",
		})
	}

	ok, err := commands.ValidAccountId(accountId)
	if err != nil {
		logging.Errorf("Error while validating account id: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"error": "Something went wrong",
		})
	}

	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"error": "account not found",
		})
	}

	_, user_err, internal_error := commands.AddAccountInfo(accountId, phoneNumber)
	if internal_error != nil {
		logging.Errorf("Error while adding account info: %v", internal_error)
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"error": "Something went wrong",
		})
	}

	if user_err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"error": user_err.Error(),
		})
	}

	return c.SendStatus(fiber.StatusCreated)
}
