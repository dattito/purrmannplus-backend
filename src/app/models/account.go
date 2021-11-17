package models

import (
	"errors"
	"strings"
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
