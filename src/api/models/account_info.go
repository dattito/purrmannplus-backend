package models

type PostAddPhoneNumberRequest struct {
	PhoneNumber string `json:"phone_number"`
}

type PostSendPhoneNumberConfirmationLinkRequest struct {
	PhoneNumber string `json:"phone_number"`
}