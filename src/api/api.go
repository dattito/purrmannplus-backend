package api

import (
	"github.com/dattito/purrmannplus-backend/api/providers"
	"github.com/dattito/purrmannplus-backend/config"
	"github.com/dattito/purrmannplus-backend/utils/logging"
)

var api providers.Provider

// Initialize the api object
func Init() error {
	api = providers.GetProvider()
	return api.Init()
}

// Start the api to listen on the configured port
func StartListening() error {
	logging.Infof("Starting listening on port %d", config.LISTENING_PORT)
	return api.StartListening()
}
