package models

import "github.com/datti-to/purrmannplus-backend/app/models"

type AccountInfoDB struct {
	Model
	AccountDB AccountDB `gorm:"foreignkey:AccountID"`
	PhoneNumber string `gorm:"phone_number"`
}

func (AccountInfoDB) TableName() string {
	return "account_infos"
}

func AccountInfoToAccountInfoDB(account models.AccountInfo) AccountInfoDB {
	aib := AccountInfoDB{
		AccountDB: AccountToAccountDB(account.Account),
		PhoneNumber: account.PhoneNumber,
	}
	aib.ID = account.Id

	return aib
}

func AccountInfoDBToAccountInfo(accountInfo AccountInfoDB) models.AccountInfo {
	account := AccountDBToAccount(accountInfo.AccountDB)

	return models.AccountInfo{
		Id: accountInfo.ID,
		Account: account,
		PhoneNumber: accountInfo.PhoneNumber,
	}
}