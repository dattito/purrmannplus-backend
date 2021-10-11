package provider

import (
	"github.com/dattito/purrmannplus-backend/database/models"
	"github.com/dattito/purrmannplus-backend/database/providers/gorm"
)

type Provider interface {
	CreateTables() error
	CloseDB() error

	AddAccount(username, password string) (models.AccountDBModel, error)
	GetAccount(id string) (models.AccountDBModel, error)
	GetAccountByCredentials(username, password string) (models.AccountDBModel, error)
	GetAccounts() ([]models.AccountDBModel, error)
	DeleteAccount(id string) error
	AddAccountInfo(accountId, phoneNumber string) (models.AccountInfoDBModel, error)
	GetAccountInfo(accountId string) (models.AccountInfoDBModel, error)

	SetSubstitutions(accountId string, substitutions map[string][]string, notSetYet bool) (models.SubstitutionDBModel, error)
	RemoveAccountFromSubstitutionUpdater(accountId string) error
	GetSubstitutions(accountId string) (models.SubstitutionDBModel, error)
	GetAllAccountCredentialsAndPhoneNumberAndSubstitutions() ([]models.AccountCredentialsAndPhoneNumberAndSubstitutionsDBModel, error)
	GetAccountCredentialsAndPhoneNumberAndSubstitutions(accountId string) (models.AccountCredentialsAndPhoneNumberAndSubstitutionsDBModel, error)
}

func GetProvider() (Provider, error) {
	// * Add / Change Provider here
	return gorm.NewGormProvider()
}
