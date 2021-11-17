package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	app_models "github.com/dattito/purrmannplus-backend/app/models"
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
	AccountId string    `gorm:"column:account_id;uniqueIndex"`
	AuthId    string    `gorm:"column:auth_id;uniqueIndex"`
	AuthPw    string    `gorm:"column:auth_pw"`
	AccountDB AccountDB `gorm:"foreignKey:account_id"`
	Entries   *Entries  `gorm:"entries;default:{}"`
	NotSetYet bool      `gorm:"column:not_set_yet"`
}

func (SubstitutionDB) TableName() string {
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

func (s SubstitutionDB) ToSubstitutions() app_models.Substitutions {
	return app_models.Substitutions{
		AccountId: s.AccountId,
		Entries:   *s.Entries,
	}
}

type SubstitutionInfoDB struct {
	AuthId          string   `gorm:"column:auth_id"`
	AuthPw          string   `gorm:"column:auth_pw"`
	PhoneNumber     string   `gorm:"column:phone_number"`
	AccountId       string   `gorm:"column:account_id"`
	SubstitutionsId string   `gorm:"column:substitutions_id"`
	Entries         *Entries `gorm:"column:entries"`
	NotSetYet       bool     `gorm:"column:not_set_yet"`
}

func (a SubstitutionInfoDB) ToSubstitutionInfo() app_models.SubstitutionInfo {
	return app_models.SubstitutionInfo{
		AuthId:          a.AuthId,
		AuthPw:          a.AuthPw,
		PhoneNumber:     a.PhoneNumber,
		AccountId:       a.AccountId,
		SubstitutionsId: a.SubstitutionsId,
		Entries:         *a.Entries,
		NotSetYet:       a.NotSetYet,
	}
}
