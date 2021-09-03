package models

import "github.com/datti-to/purrmannplus-backend/app/models"

type AccountDB struct {
	Model
	ID     string `gorm:"primary_key,size:32"`
	AuthId string `gorm:"auth_id,unique"`
	AuthPw string `gorm:"auth_pw"`
}

func (a AccountDB) TableName() string {
	return "accounts"
}

func AccountToAccountDB(a models.Account) AccountDB {

	return AccountDB{
		ID:     a.Id,
		AuthId: a.AuthId,
		AuthPw: a.AuthPw,
	}
}

func AccountDBToAccount(a AccountDB) models.Account {
	return models.Account{
		Id:     a.ID,
		AuthId: a.AuthId,
		AuthPw: a.AuthPw,
	}
}
