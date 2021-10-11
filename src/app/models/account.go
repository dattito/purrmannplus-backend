package models

import (
	"errors"

	"github.com/dattito/purrmannplus-backend/database/models"
)

type Account struct {
	Id       string
	Username string
	Password string
}

func NewValidAccount(username, password string) (*Account, error) {
	if username == "" {
		return nil, errors.New("username is empty")
	}

	if password == "" {
		return nil, errors.New("password is empty")
	}

	return &Account{
		Username: username,
		Password: password,
	}, nil
}

func AcccountDBModelToAccount(accountDBModel models.AccountDBModel) Account {
	return Account{
		Id:       accountDBModel.Id,
		Username: accountDBModel.Username,
		Password: accountDBModel.Password,
	}
}
