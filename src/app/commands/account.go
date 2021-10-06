package commands

import (
	"errors"

	"github.com/dattito/purrmannplus-backend/app/models"
	"github.com/dattito/purrmannplus-backend/database"
	db_errors "github.com/dattito/purrmannplus-backend/database/errors"
	"github.com/dattito/purrmannplus-backend/services/hpg"
	"github.com/dattito/purrmannplus-backend/utils/logging"
)

// Returns the accountId of the new account; error produced by user; error not produced by user
func CreateAccount(authId, authPw string) (models.Account, error, error) {
	if _, err := models.NewValidAccount(authId, authPw); err != nil {
		return models.Account{}, nil, err
	}

	correct, err := hpg.CheckCredentials(authId, authPw)
	if err != nil {
		return models.Account{}, nil, err
	}

	if !correct {
		return models.Account{}, errors.New("incorrect credentials"), nil
	}

	a, err := database.DB.AddAccount(authId, authPw)
	if err == nil {
		logging.Infof("Created account %s", a.AuthId)
	}

	return models.AcccountDBModelToAccount(a), nil, err
}

func GetAllAccounts() ([]models.Account, error) {

	accounts, err := database.DB.GetAccounts()
	if err != nil {
		return nil, err
	}

	var accs []models.Account
	for _, a := range accounts {
		accs = append(accs, models.AcccountDBModelToAccount(a))
	}

	return accs, nil
}

func GetAccount(accountId string) (models.Account, error) {
	a, err := database.DB.GetAccount(accountId)
	if err != nil {
		return models.Account{}, err
	}

	return models.AcccountDBModelToAccount(a), nil
}

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

func GetAccountByCredentials(authId, authPw string) (models.Account, error) {
	if authId == "" {
		return models.Account{}, errors.New("missing authId")
	}

	if authPw == "" {
		return models.Account{}, errors.New("missing authPw")
	}

	a, err := database.DB.GetAccountByCredentials(authId, authPw)
	if err != nil {
		return models.Account{}, err
	}

	return models.AcccountDBModelToAccount(a), nil
}

func DeleteAccount(accountId string) error {
	err := database.DB.DeleteAccount(accountId)
	if err == nil {
		logging.Infof("Deleted account %s", accountId)
	}

	return err
}
