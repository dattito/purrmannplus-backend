package models

import (
	provider_models "github.com/dattito/purrmannplus-backend/database/models"
)

type AccountInfoDB struct {
	Model
	AccountId   string    `gorm:"column:account_id;uniqueIndex"`
	AccountDB   AccountDB `gorm:"foreignkey:account_id"`
	PhoneNumber string    `gorm:"column:phone_number"`
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
