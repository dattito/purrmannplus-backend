package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	provider_models "github.com/datti-to/purrmannplus-backend/database/models"
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

func SubstitutionDBToSubstitutionDBModel(s SubstitutionDB) provider_models.SubstitutionDBModel {
	return provider_models.SubstitutionDBModel{
		Id:        s.Id,
		AccountId: s.AccountId,
		Entries:   *s.Entries,
		NotSetYet: s.NotSetYet,
	}
}

type AccountCredentialsAndPhoneNumberAndSubstitutionsDB struct {
	AuthId          string   `gorm:"column:auth_id"`
	AuthPw          string   `gorm:"column:auth_pw"`
	PhoneNumber     string   `gorm:"column:phone_number"`
	AccountId       string   `gorm:"column:account_id"`
	SubstitutionsId string   `gorm:"column:substitutions_id"`
	Entries         *Entries `gorm:"column:entries"`
	NotSetYet       bool     `gorm:"column:not_set_yet"`
}

func ACPSDBtoACPDSDBM(s AccountCredentialsAndPhoneNumberAndSubstitutionsDB) provider_models.AccountCredentialsAndPhoneNumberAndSubstitutionsDBModel {
	return provider_models.AccountCredentialsAndPhoneNumberAndSubstitutionsDBModel{
		AuthId:          s.AuthId,
		AuthPw:          s.AuthPw,
		PhoneNumber:     s.PhoneNumber,
		AccountId:       s.AccountId,
		SubstitutionsId: s.SubstitutionsId,
		Entries:         *s.Entries,
		NotSetYet:       s.NotSetYet,
	}
}
