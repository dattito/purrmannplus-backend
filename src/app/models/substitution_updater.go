package models

import "github.com/datti-to/purrmannplus-backend/database/models"

type Substitutions struct {
	Id      string
	Account Account
	Entries map[string][]string
}

func NewSubstitutions(authId, authPw string, entries map[string][]string) (*Substitutions, error) {
	a, err := NewValidAccount(authId, authPw)
	if err != nil {
		return nil, err
	}
	return &Substitutions{
		Account: *a,
		Entries: entries,
	}, nil
}

func NewEmptySubstitutions(authId, authPw string) (*Substitutions, error) {
	return NewSubstitutions(authId, authPw, map[string][]string{})
}

type SubstitutionUpdateInfos struct {
	AuthId          string
	AuthPw          string
	PhoneNumber     string
	AccountId       string
	SubstitutionsId string
	Entries         map[string][]string
}

func AccountCredentialsAndPhoneNumberAndSubstitutionsDBModelToSubstitutionUpdateInfos(m *models.AccountCredentialsAndPhoneNumberAndSubstitutionsDBModel) SubstitutionUpdateInfos {
	return SubstitutionUpdateInfos{
		AuthId:          m.AuthId,
		AuthPw:          m.AuthPw,
		PhoneNumber:     m.PhoneNumber,
		AccountId:       m.AccountId,
		SubstitutionsId: m.SubstitutionsId,
		Entries:         m.Entries,
	}
}
