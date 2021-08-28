package database

import (
	"github.com/datti-to/purrmannplus-backend/database/provider"
)

var DB provider.Provider

func Init() error {
	var err error
	DB, err = provider.GetProvider()

	return err
}
