package database

import (
	"github.com/dattito/purrmannplus-backend/config"
	"github.com/dattito/purrmannplus-backend/database/providers"
)

var DB provider.Provider

func Init() error {
	var err error
	DB, err = provider.GetProvider()

	if config.DATABASE_AUTOMIGRATE {
		return DB.CreateTables()
	}

	return err
}
