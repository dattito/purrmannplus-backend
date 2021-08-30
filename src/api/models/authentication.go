package models

import "github.com/datti-to/purrmannplus-backend/app/models"

type PostLoginRequest struct {
	AuthId string `json:"auth_id"`
	AuthPw string `json:"auth_pw"`
}

type PostLoginResponse struct {
	Token string `json:"token"`
}

func PostLoginRequestToAccount(p PostLoginRequest) (*models.Account, error) {
	return models.NewAccount(p.AuthId, p.AuthPw)
}