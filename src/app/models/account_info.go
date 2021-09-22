package models

import (
	"github.com/datti-to/purrmannplus-backend/database/models"
	"github.com/datti-to/purrmannplus-backend/utils"
)

type AccountInfo struct {
	Id          string
	Account     Account
	PhoneNumber string
}

func NewAccountInfo(account Account, phoneNumber string) (*AccountInfo, error) {

	formattedPhoneNumber, err := utils.FormatPhoneNumber(phoneNumber)
	if err != nil {
		return nil, err
	}

	return &AccountInfo{
		Account:     account,
		PhoneNumber: formattedPhoneNumber,
	}, nil
}

func AccountInfoDBModelToAccount(accountInfoDBModel models.AccountInfoDBModel) (AccountInfo, error) {

	return AccountInfo{
		Id: accountInfoDBModel.Id,
		Account: Account{
			Id: accountInfoDBModel.AccountId,
		},
		PhoneNumber: accountInfoDBModel.PhoneNumber,
	}, nil
}
