package models

type Substitutions struct {
	AccountId string
	Entries   map[string][]string
}

type SubstitutionInfo struct {
	AuthId          string
	AuthPw          string
	PhoneNumber     string
	AccountId       string
	SubstitutionsId string
	Entries         map[string][]string
	NotSetYet       bool
}
