package models

import (
	"errors"
	"strings"

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
		Username: strings.ToLower(username),
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
