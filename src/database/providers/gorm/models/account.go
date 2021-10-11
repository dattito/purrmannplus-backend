package models

import (
	provider_models "github.com/dattito/purrmannplus-backend/database/models"
	"gorm.io/gorm"
)

type AccountDB struct {
	Model
	Username string `gorm:"column:auth_id;uniqueIndex"`
	Password string `gorm:"column:auth_pw"`
}

func (AccountDB) TableName() string {
	return "accounts"
}

func (a *AccountDB) BeforeDelete(tx *gorm.DB) error {
	if err := tx.Where("account_id = ?", a.Id).Delete(&SubstitutionDB{}).Error; err != nil {
		return err
	}

	return tx.Where("account_id = ?", a.Id).Delete(&AccountInfoDB{}).Error
}

func AccountDBToAccountDBModel(a AccountDB) provider_models.AccountDBModel {
	return provider_models.AccountDBModel{
		Id:       a.Id,
		Username: a.Username,
		Password: a.Password,
	}
}
