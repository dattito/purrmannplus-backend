package providers

import "github.com/datti-to/purrmannplus-backend/api/providers/rest"

type Provider interface {
	Init() error
	StartListening() error
}

func GetProvider() Provider {
	return &rest.RestProvider{}
}
