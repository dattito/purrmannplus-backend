package api

import (
	"github.com/dattito/purrmannplus-backend/api/providers"
)

var api providers.Provider

func Init() error {
	api = providers.GetProvider()
	return api.Init()
}

func StartListening() error {
	return api.StartListening()
}
