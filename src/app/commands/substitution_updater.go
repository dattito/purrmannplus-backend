package commands

import (
	"errors"

	"github.com/datti-to/purrmannplus-backend/database"
	db_errors "github.com/datti-to/purrmannplus-backend/database/errors"
)

func AddToSubstitutionUpdater(accountId string) error {
	ai, err := database.DB.GetAccountInfo(accountId)
	if err != nil {
		if errors.Is(err, &db_errors.ErrRecordNotFound) {
			return errors.New("phone number has to be added first")
		}
		return err
	}

	if ai.PhoneNumber == "" {
		return errors.New("phone number has to be added first")
	}

	return database.DB.AddAccountToSubstitutionUpdater(accountId)
}

func RemoveFromSubstitutionUpdater(accountId string) error {
	return database.DB.RemoveAccountFromSubstitutionUpdater(accountId)
}
