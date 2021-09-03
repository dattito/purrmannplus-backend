package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/datti-to/purrmannplus-backend/app/models"
)

type Entries map[string][]string

func (s *Entries) Scan(val interface{}) error {
	switch v := val.(type) {
	case []byte:
		json.Unmarshal(v, &s)
		return nil
	case string:
		json.Unmarshal([]byte(v), &s)
		return nil
	default:
		return fmt.Errorf("unsupported type: %T", v)
	}
}

func (s *Entries) Value() (driver.Value, error) {
	if s != nil {
		b, err := json.Marshal(s)
		if err != nil {
			return nil, err
		}
		return string(b), nil
	}
	return nil, nil
}

type SubstitutionDB struct {
	Model
	ID        string    `gorm:"primary_key,size:32"`
	AccountDB AccountDB `gorm:"foreignkey:AccountID"`
	Entries   Entries   `gorm:"entries"`
}

func (s SubstitutionDB) TableName() string {
	return "substitutions"
}

func (s *SubstitutionDB) Scan(val interface{}) error {
	switch v := val.(type) {
	case []byte:
		json.Unmarshal(v, &s)
		return nil
	default:
		return fmt.Errorf("unsupported type: %T", v)
	}
}

func (s *SubstitutionDB) Value() (driver.Value, error) {
	return json.Marshal(s)
}

func SubstitutionDBToSubstitution(s SubstitutionDB) *models.Substitutions {
	return &models.Substitutions{
		Id:      s.ID,
		Account: AccountDBToAccount(s.AccountDB),
		Entries: s.Entries,
	}
}

func SubstitutionsToSubstitutionDB(s *models.Substitutions) *SubstitutionDB {
	return &SubstitutionDB{
		ID:        s.Id,
		AccountDB: AccountToAccountDB(s.Account),
		Entries:   s.Entries,
	}
}
