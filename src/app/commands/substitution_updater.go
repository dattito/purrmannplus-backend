package commands

import (
	"errors"
	"log"

	"github.com/datti-to/purrmannplus-backend/app/models"
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

// Updates the substitutions for a given account and it's relevant data and sends a message via signal
func UpdateSubstitutions(m models.SubstitutionUpdateInfos) error {
	return nil
}

func UpdateAllSubstitutions() error {
	mdbs, err := database.DB.GetAllAccountCredentialsAndPhoneNumberAndSubstitutions()
	if err != nil {
		return err
	}

	errCount := 0

	for _, mdb := range mdbs {
		err := UpdateSubstitutions(models.AccountCredentialsAndPhoneNumberAndSubstitutionsDBModelToSubstitutionUpdateInfos(&mdb))
		if err != nil {
			log.Printf("Error updating substitutions for account %s: %s", mdb.AccountId, err.Error())
			errCount++
			if errCount > config.MAX_ERROS_TO_STOP_UPDATING_SUBSTITUTIONS {
				return errors.New("got too many errors updating substitutions, stopping")
			}
		}
	}
	return nil
}

func EnableSubstitutionUpdater() {
	scheduler.AddJob(config.SUBSTITUTIONS_UPDATECRON, func() {
		if err := UpdateAllSubstitutions(); err != nil {
			log.Println(err.Error())
		}
	})
}
