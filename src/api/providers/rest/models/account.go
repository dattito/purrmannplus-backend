package models

import app_models "github.com/datti-to/purrmannplus-backend/app/models"

type PostAccountRequest struct {
	AuthId string `json:"auth_id" form:"auth_id"`
	AuthPw string `json:"auth_pw" form:"auth_pw"`
}

func PostAccountRequestToAccount(req *PostAccountRequest) (*app_models.Account, error) {
	return app_models.NewAccount(req.AuthId, req.AuthPw)
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
	Id     string `json:"id"`
	AuthId string `json:"auth_id" form:"auth_id"`
	AuthPw string `json:"auth_pw" form:"auth_pw"`
}

func AccountToGetAccountResponse(account *app_models.Account) *GetAccountResponse {
	return &GetAccountResponse{
		Id:     account.Id,
		AuthId: account.AuthId,
		AuthPw: account.AuthPw,
	}
}

func AccountsToGetAccountResponses(accounts []app_models.Account) []*GetAccountResponse {
	var getAccountResponses []*GetAccountResponse
	for _, account := range accounts {
		getAccountResponses = append(getAccountResponses, AccountToGetAccountResponse(&account))
	}
	return getAccountResponses
}
