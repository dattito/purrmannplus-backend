package models

type SubstitutionDBModel struct {
	Id        string
	AccountId string
	Entries   map[string][]string
}
