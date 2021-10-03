package providers

import "github.com/dattito/purrmannplus-backend/api/providers/rest"

type Provider interface {
	Init() error
	StartListening() error
}

func GetProvider() Provider {
	return &rest.RestProvider{}
}
