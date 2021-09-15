package models

import (
	provider_models "github.com/datti-to/purrmannplus-backend/database/models"
)

type AccountInfoDB struct {
	Model
	AccountId   string    `gorm:"account_id,uniqueIndex"`
	AccountDB   AccountDB `gorm:"foreignkey:account_id"`
	PhoneNumber string    `gorm:"phone_number"`
}

func (AccountInfoDB) TableName() string {
	return "account_infos"
}

func AccountInfoDBToAccountInfoDBModel(a AccountInfoDB) provider_models.AccountInfoDBModel {
	return provider_models.AccountInfoDBModel{
		AccountId:   a.AccountId,
		PhoneNumber: a.PhoneNumber,
	}
}
