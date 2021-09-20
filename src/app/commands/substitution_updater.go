package commands

import (
	"errors"
	"log"

	"github.com/datti-to/purrmannplus-backend/config"
	"github.com/datti-to/purrmannplus-backend/database"
	db_errors "github.com/datti-to/purrmannplus-backend/database/errors"
	"github.com/datti-to/purrmannplus-backend/services/scheduler"
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

func UpdateAllSubstitutions() error {
	return nil
}

func EnableSubstitutionUpdater() {
	scheduler.AddJob(config.SUBSTITUTIONS_UPDATECRON, func() {
		if err := UpdateAllSubstitutions(); err != nil {
			log.Println(err.Error())
		}
	})
}
