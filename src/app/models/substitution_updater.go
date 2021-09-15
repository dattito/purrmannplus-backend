package models

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
