package provider

import (
	"github.com/dattito/purrmannplus-backend/app/models"
	"github.com/dattito/purrmannplus-backend/database/providers/gorm"
)

type Provider interface {
	CreateTables() error
	CloseDB() error

	AddAccount(username, password string) (models.Account, error)
	GetAccount(id string) (models.Account, error)
	GetAccountByCredentials(username, password string) (models.Account, error)
	GetAccounts() ([]models.Account, error)
	DeleteAccount(id string) error
	AddAccountInfo(accountId, phoneNumber string) (models.AccountInfo, error)
	GetAccountInfo(accountId string) (models.AccountInfo, error)

	AddAccountToSubstitution(accountId, authId, authPw string) error
	SetSubstitutions(accountId string, substitutions map[string][]string, notSetYet bool) error
	RemoveAccountFromSubstitutionUpdater(accountId string) error
	GetSubstitutions(accountId string) (models.Substitutions, error)
	GetAllSubstitutionInfos() ([]models.SubstitutionInfo, error)
	GetSubstitutionInfos(accountId string) (models.SubstitutionInfo, error)

	SetMoodleAssignments(accountId string, assignmentIds []int, notSetYet bool) error
	RemoveAccountFromMoodleAssignmentUpdater(accountId string) error
	GetMoodleAssignments(accountId string) (models.MoodleAssignments, error)
	GetAllMoodleAssignmentInfos() ([]models.MoodleAssignmentInfo, error)
	GetMoodleAssignmentInfos(accountId string) (models.MoodleAssignmentInfo, error)
}

func GetProvider() (Provider, error) {
	// * Add / Change Provider here
	return gorm.NewGormProvider()
}
