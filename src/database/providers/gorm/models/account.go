package models

import (
	provider_models "github.com/datti-to/purrmannplus-backend/database/models"
)

type AccountDB struct {
	Model
	AuthId string `gorm:"auth_id,uniqueIndex"`
	AuthPw string `gorm:"auth_pw"`
}

func (AccountDB) TableName() string {
	return "accounts"
}

func AccountDBToAccountDBModel(a AccountDB) provider_models.AccountDBModel {
	return provider_models.AccountDBModel{
		Id:     a.Id,
		AuthId: a.AuthId,
		AuthPw: a.AuthPw,
	}
}
