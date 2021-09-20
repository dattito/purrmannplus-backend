package provider

import (
	"github.com/datti-to/purrmannplus-backend/database/models"
	"github.com/datti-to/purrmannplus-backend/database/providers/gorm"
)

type Provider interface {
	CreateTables() error
	CloseDB() error

	AddAccount(authId, authPw string) (models.AccountDBModel, error)
	GetAccount(id string) (models.AccountDBModel, error)
	GetAccountByCredentials(authId, authPw string) (models.AccountDBModel, error)
	GetAccounts() ([]models.AccountDBModel, error)
	DeleteAccount(id string) error
	AddAccountInfo(accountId, phoneNumber string) (models.AccountInfoDBModel, error)
	GetAccountInfo(accountId string) (models.AccountInfoDBModel, error)

	SetSubstitutions(accountId string, substitutions map[string][]string) (models.SubstitutionDBModel, error)
	AddAccountToSubstitutionUpdater(accountId string) error
	RemoveAccountFromSubstitutionUpdater(accountId string) error
	GetSubstitutions(accountId string) (models.SubstitutionDBModel, error)
}

func GetProvider() (Provider, error) {
	// * Add / Change Provider here
	return gorm.NewGormProvider()
}
