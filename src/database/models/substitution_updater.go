package models

type SubstitutionDBModel struct {
	Id        string
	AccountId string
	Entries   map[string][]string
	NotSetYet bool
}

type AccountCredentialsAndPhoneNumberAndSubstitutionsDBModel struct {
	AuthId          string
	AuthPw          string
	PhoneNumber     string
	AccountId       string
	SubstitutionsId string
	Entries         map[string][]string
	NotSetYet       bool
}
