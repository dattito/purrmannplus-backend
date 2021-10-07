package models

import app_models "github.com/dattito/purrmannplus-backend/app/models"

type PostAccountRequest struct {
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
}

func PostAccountRequestToAccount(req *PostAccountRequest) (*app_models.Account, error) {
	return app_models.NewValidAccount(req.Username, req.Password)
}

type PostAccountResponse struct {
	Id string `json:"id"`
}

func AccountToPostAccountResponse(account *app_models.Account) *PostAccountResponse {
	return &PostAccountResponse{
		Id: account.Id,
	}
}

type GetAccountResponse struct {
	Id       string `json:"id"`
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
}

func AccountToGetAccountResponse(account *app_models.Account) *GetAccountResponse {
	return &GetAccountResponse{
		Id:       account.Id,
		Username: account.Username,
		Password: account.Password,
	}
}

func AccountsToGetAccountResponses(accounts []app_models.Account) []*GetAccountResponse {
	var getAccountResponses []*GetAccountResponse
	for _, account := range accounts {
		getAccountResponses = append(getAccountResponses, AccountToGetAccountResponse(&account))
	}
	return getAccountResponses
}
