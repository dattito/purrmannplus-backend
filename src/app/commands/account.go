package commands

import (
	"errors"

	"github.com/dattito/purrmannplus-backend/app/models"
	"github.com/dattito/purrmannplus-backend/database"
	db_errors "github.com/dattito/purrmannplus-backend/database/errors"
	"github.com/dattito/purrmannplus-backend/services/moodle"
	"github.com/dattito/purrmannplus-backend/utils/logging"
)

// Returns the accountId of the new account; error produced by user; error not produced by user
func CreateAccount(username, password string) (models.Account, error, error) {
	if _, err := models.NewValidAccount(username, password); err != nil {
		return models.Account{}, err, nil
	}

	correct, err := moodle.CheckCredentials(username, password)
	if err != nil {
		return models.Account{}, nil, err
	}

	if !correct {
		return models.Account{}, errors.New("incorrect credentials"), nil
	}

	a, err := database.DB.AddAccount(username, password)
	if err == nil {
		logging.Infof("Created account %s", a.Username)
	}

	return a, nil, err
}

// Returns the id and the credentials of all accounts
func GetAllAccounts() ([]models.Account, error) {

	return database.DB.GetAccounts()
}

// Returns the accountId and the credentials of the account
func GetAccount(accountId string) (models.Account, error) {
	a, err := database.DB.GetAccount(accountId)
	if err != nil {
		return models.Account{}, err
	}

	return a, nil
}

// Returns true if the accountId was found in the database
func ValidAccountId(accountId string) (bool, error) {
	_, err := GetAccount(accountId)
	if err != nil {
		if errors.Is(err, &db_errors.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

// Returns the accountId and the credentials of the account matching the credentials
func GetAccountByCredentials(username, password string) (models.Account, error) {
	if username == "" {
		return models.Account{}, errors.New("missing authId")
	}

	if password == "" {
		return models.Account{}, errors.New("missing authPw")
	}

	a, err := database.DB.GetAccountByCredentials(username, password)
	if err != nil {
		return models.Account{}, err
	}

	return a, nil
}

// Deleting an account
func DeleteAccount(accountId string) error {
	err := database.DB.DeleteAccount(accountId)
	if err == nil {
		logging.Infof("Deleted account %s", accountId)
	}

	return err
}

// Checks the credentials of an account, should be the same as mooodle.CheckCredentials(username, password)
func CheckCredentials(username, password string) (bool, error) {
	return moodle.CheckCredentials(username, password)
}
