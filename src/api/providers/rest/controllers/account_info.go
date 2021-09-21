package controllers

import (
	"fmt"
	"log"

	api_models "github.com/datti-to/purrmannplus-backend/api/providers/rest/models"
	"github.com/datti-to/purrmannplus-backend/app/commands"
	"github.com/datti-to/purrmannplus-backend/app/models"
	"github.com/datti-to/purrmannplus-backend/config"
	"github.com/datti-to/purrmannplus-backend/services/signal_message_sender"
	utils_jwt "github.com/datti-to/purrmannplus-backend/utils/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

func SendPhoneNumberConfirmationLink(c *fiber.Ctx) error {
	pr := new(api_models.PostSendPhoneNumberConfirmationLinkRequest)
	if err := c.BodyParser(pr); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"error": err.Error(),
		})
	}

	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	account_id := claims["account_id"].(string)

	account_info, err := models.NewAccountInfo(models.Account{Id: account_id}, pr.PhoneNumber)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"error": err.Error(),
		})
	}

	token, err := utils_jwt.NewAccountIdPhoneNumberToken(account_info.Account.Id, account_info.PhoneNumber)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"error": "couln't create token",
		})
	}

	text := fmt.Sprintf("Willkommen bei PurrmannPlus. Um deine Telefonnummer zu bestätigen, drücke "+
		"auf den nachfolgenden Link. Er ist 10 Minuten lang gültig. Du hast den Link nicht angefordert? Dann kannst du ihn ignorieren. "+
		"%s/v1/accounts/phone_number/validate?token=%s",
		config.API_URL, token)

	err = signal_message_sender.SignalMessageSender.Send(text, account_info.PhoneNumber)

	if err != nil {
		log.Println(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"error": "couln't send signal message",
		})
	}

	return c.SendStatus(fiber.StatusCreated)
}

func AddPhoneNumber(c *fiber.Ctx) error {
	p := new(api_models.PostAddPhoneNumberRequest)
	if err := c.QueryParser(p); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"error": err.Error(),
		})
	}

	if p.Token == "" {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"error": "token is required",
		})
	}

	accountId, phoneNumber, err := utils_jwt.ParseAccountIdPhoneNumberToken(p.Token)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"error": "couln't parse token",
		})
	}

	if _, err = commands.AddAccountInfo(accountId, phoneNumber); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"error": "something went wrong",
		})
	}

	return c.SendStatus(fiber.StatusCreated)
}
