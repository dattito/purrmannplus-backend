package models

type PostAddPhoneNumberRequest struct {
	Token string `query:"token"`
}

type PostSendPhoneNumberConfirmationLinkRequest struct {
	PhoneNumber string `json:"phone_number"`
}
