package models

import "github.com/dattito/purrmannplus-backend/app/models"

type PostLoginRequest struct {
	StoreInCookie bool   `json:"store_in_cookie"`
	Username      string `json:"username"`
	Password      string `json:"password"`
}

type PostLoginResponse struct {
	Token string `json:"token"`
}

func PostLoginRequestToAccount(p PostLoginRequest) (*models.Account, error) {
	return models.NewValidAccount(p.Username, p.Password)
}
