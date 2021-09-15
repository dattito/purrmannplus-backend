package models

import "github.com/datti-to/purrmannplus-backend/app/models"

type AccountDB struct {
	Model
	AuthId string `gorm:"auth_id,uniqueIndex"`
	AuthPw string `gorm:"auth_pw"`
}

func (a AccountDB) TableName() string {
	return "accounts"
}

func AccountToAccountDB(a models.Account) AccountDB {
	ab := AccountDB{
		AuthId: a.AuthId,
		AuthPw: a.AuthPw,
	}
	ab.ID = a.Id

	return ab
}

func AccountDBToAccount(a AccountDB) models.Account {
	return models.Account{
		Id:     a.ID,
		AuthId: a.AuthId,
		AuthPw: a.AuthPw,
	}
}
