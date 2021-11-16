package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	provider_models "github.com/dattito/purrmannplus-backend/database/models"
)

type AssignmentIds []int

func (s *AssignmentIds) Scan(val interface{}) error {
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

func (s *AssignmentIds) Value() (driver.Value, error) {
	if s != nil {
		b, err := json.Marshal(s)
		if err != nil {
			return nil, err
		}
		return string(b), nil
	}
	return nil, nil
}

type MoodleUserAssignmentsDB struct {
	Model
	AccountId     string         `gorm:"column:account_id;uniqueIndex"`
	AccountDB     AccountDB      `gorm:"foreignKey:account_id"`
	AssignmentIds *AssignmentIds `gorm:"assignment_ids;default:[]"`
	NotSetYet     bool           `gorm:"column:not_set_yet"`
}

func (MoodleUserAssignmentsDB) TableName() string {
	return "moodle_user_assignments"
}

func (s *MoodleUserAssignmentsDB) Scan(val interface{}) error {
	switch v := val.(type) {
	case []byte:
		json.Unmarshal(v, &s)
		return nil
	default:
		return fmt.Errorf("unsupported type: %T", v)
	}
}

func (s *MoodleUserAssignmentsDB) Value() (driver.Value, error) {
	return json.Marshal(s)
}

func MoodleUserAssignmentsDBTMoodleUserAssignmentsDBModel(s MoodleUserAssignmentsDB) provider_models.MoodleUserAssignmentsDBModel {
	// thats a really long method name, but everybody knows what if means ...
	return provider_models.MoodleUserAssignmentsDBModel{
		Id:            s.Id,
		AccountId:     s.AccountId,
		AssignmentIds: *s.AssignmentIds,
		NotSetYet:     s.NotSetYet,
	}
}

type AccountCredentialsAndPhoneNumberAndMoodleUserAssignmentsDB struct {
	AuthId                  string         `gorm:"column:auth_id"`
	AuthPw                  string         `gorm:"column:auth_pw"`
	PhoneNumber             string         `gorm:"column:phone_number"`
	AccountId               string         `gorm:"column:account_id"`
	MoodleUserAssignmentsId string         `gorm:"column:moodle_user_assignment_id"`
	AssignmentIds           *AssignmentIds `gorm:"column:assignment_ids"`
	NotSetYet               bool           `gorm:"column:not_set_yet"`
}

func ACPADBtoACPADBM(s AccountCredentialsAndPhoneNumberAndMoodleUserAssignmentsDB) provider_models.AccountCredentialsAndPhoneNumberAndMoodleUserAssignmentsDBModel {
	return provider_models.AccountCredentialsAndPhoneNumberAndMoodleUserAssignmentsDBModel{
		AuthId:                  s.AuthId,
		AuthPw:                  s.AuthPw,
		PhoneNumber:             s.PhoneNumber,
		AccountId:               s.AccountId,
		MoodleUserAssignmentsId: s.MoodleUserAssignmentsId,
		AssignmentIds:           *s.AssignmentIds,
		NotSetYet:               s.NotSetYet,
	}
}
