package models

type AccountInfo struct {
	Account     Account
	PhoneNumber string
}

func NewAccountInfo(account Account, phoneNumber string) (*AccountInfo, error) {
	return &AccountInfo{
		Account:     account,
		PhoneNumber: phoneNumber,
	}, nil
}