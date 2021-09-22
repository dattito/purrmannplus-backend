package api

import (
	"github.com/datti-to/purrmannplus-backend/api/providers"
)

var api providers.Provider

func Init() error {
	api = providers.GetProvider()
	return api.Init()
}

func StartListening() error {
	return api.StartListening()
}
