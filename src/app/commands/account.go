package commands

import (
	"errors"

	"github.com/datti-to/purrmannplus-backend/app/models"
	"github.com/datti-to/purrmannplus-backend/database"
	"github.com/datti-to/purrmannplus-backend/services/hpg"
)

func CreateAccount(authId, authPw string) (models.Account, error) {
	if _, err := models.NewValidAccount(authId, authPw); err != nil {
		return models.Account{}, err
	}

	correct, err := hpg.CheckCredentials(authId, authPw)
	if err != nil {
		return models.Account{}, err
	}

	if !correct {
		return models.Account{}, errors.New("incorrect credentials")
	}

	a, err := database.DB.AddAccount(authId, authPw)

	return models.AcccountDBModelToAccount(a), err
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
