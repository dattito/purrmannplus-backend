package models

import "github.com/dattito/purrmannplus-backend/app/models"

type PostLoginRequest struct {
	StoreInCookie bool   `json:"store_in_cookie"`
	AuthId        string `json:"auth_id"`
	AuthPw        string `json:"auth_pw"`
}

type PostLoginResponse struct {
	Token string `json:"token"`
}

func PostLoginRequestToAccount(p PostLoginRequest) (*models.Account, error) {
	return models.NewValidAccount(p.AuthId, p.AuthPw)
}
