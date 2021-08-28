package models

import app_models "github.com/datti-to/purrmannplus-backend/app/models"

type PostAccountRequest struct {
	AuthId string `json:"auth_id"`
	AuthPw string `json:"auth_pw"`
}

type PostAccountResponse struct {
	Id string `json:"id"`
}

func PostAccountRequestToAccount(req *PostAccountRequest) (*app_models.Account, error) {
	return app_models.NewAccount(req.AuthId, req.AuthPw)
}
