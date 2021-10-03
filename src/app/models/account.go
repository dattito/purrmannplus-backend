package models

import (
	"errors"

	"github.com/dattito/purrmannplus-backend/database/models"
)

type Account struct {
	Id     string
	AuthId string
	AuthPw string
}

func NewValidAccount(authId, authPw string) (*Account, error) {
	if authId == "" {
		return nil, errors.New("authId is empty")
	}

	if authPw == "" {
		return nil, errors.New("authPw is empty")
	}

	return &Account{
		AuthId: authId,
		AuthPw: authPw,
	}, nil
}

func AcccountDBModelToAccount(accountDBModel models.AccountDBModel) Account {
	return Account{
		Id:     accountDBModel.Id,
		AuthId: accountDBModel.AuthId,
		AuthPw: accountDBModel.AuthPw,
	}
}
