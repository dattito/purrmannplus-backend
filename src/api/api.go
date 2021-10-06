package api

import (
	"github.com/dattito/purrmannplus-backend/api/providers"
	"github.com/dattito/purrmannplus-backend/config"
	"github.com/dattito/purrmannplus-backend/utils/logging"
)

var api providers.Provider

func Init() error {
	api = providers.GetProvider()
	return api.Init()
}

func StartListening() error {
	logging.Infof("Starting listening on port %d", config.LISTENING_PORT)
	if config.ENABLE_API {
		return api.StartListening()
	}
	return nil
}
