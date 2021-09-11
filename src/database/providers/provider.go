package provider

import (
	"github.com/datti-to/purrmannplus-backend/app/models"
	"github.com/datti-to/purrmannplus-backend/database/providers/gorm"
)

type Provider interface {
	CreateTables() error
	CloseDB() error

	AddAccount(account *models.Account) error
	GetAccount(id string) (models.Account, error)
	GetAccountByCredentials(a models.Account) (models.Account, error)
	GetAccounts() ([]models.Account, error)
	UpdateAccount(account models.Account) error
	DeleteAccount(id string) error
	GetAccountInfo(accountId string) (models.AccountInfo, error)

	AddSubstitution(substitutions *models.Substitutions) error
}

func GetProvider() (Provider, error) {
	// * Add / Change Provider here
	return gorm.NewGormProvider()
}
