package models

import (
	app_models "github.com/dattito/purrmannplus-backend/app/models"
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

func (a AccountInfoDB) ToAccountInfo() app_models.AccountInfo {
	return app_models.AccountInfo{
		Account: app_models.Account{
			Id: a.AccountId,
		},
		PhoneNumber: a.PhoneNumber,
	}
}
