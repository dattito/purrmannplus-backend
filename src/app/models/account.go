package models

import "errors"

type Account struct {
	ID     string
	AuthId string
	AuthPw string
}

func NewAccount(authId, authPw string) (*Account, error) {
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
