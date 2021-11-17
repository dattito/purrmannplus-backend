package commands

import (
	"errors"

	"github.com/dattito/purrmannplus-backend/app/models"
	"github.com/dattito/purrmannplus-backend/database"
	db_errors "github.com/dattito/purrmannplus-backend/database/errors"
)

// Return the account info for the given account id; error produced by user; error not produced by user
func AddAccountInfo(accountId, phoneNumber string) (models.AccountInfo, error, error) {
	_, err := models.NewAccountInfo(models.Account{Id: accountId}, phoneNumber)
	if err != nil {
		return models.AccountInfo{}, err, nil
	}

	ai, err := database.DB.AddAccountInfo(accountId, phoneNumber)
	if err != nil {
		return models.AccountInfo{}, nil, err
	}

	return ai, nil, err
}

// Returns true if an phone number was added to this user
func HasPhoneNumber(account_id string) (bool, error) {
	ai, err := database.DB.GetAccountInfo(account_id)
	if err != nil {
		if errors.Is(err, &db_errors.ErrRecordNotFound) {
			return false, nil
		}

		return false, err
	}

	return ai.PhoneNumber != "", nil
}
