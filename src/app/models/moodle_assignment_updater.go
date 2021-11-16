package models

import "github.com/dattito/purrmannplus-backend/database/models"

type Assignments struct {
	Courses []struct {
		FullName    string `json:"fullname"`
		Assignments []struct {
			ID int `json:"id"`
		} `json:"assignments"`
	} `json:"courses"`
}

type MoodleAssignmentUpdateInfos struct {
	AuthId                  string
	AuthPw                  string
	PhoneNumber             string
	AccountId               string
	MoodleUserAssignmentsId string
	AssignmentIds           []int
	NotSetYet               bool
}

func AccountCredentialsAndPhoneNumberAndMoodleUserAssignmentsDBModelToMoodleAssignmentUpdateInfos(m *models.AccountCredentialsAndPhoneNumberAndMoodleUserAssignmentsDBModel) MoodleAssignmentUpdateInfos {
	return MoodleAssignmentUpdateInfos{
		AuthId:                  m.AuthId,
		AuthPw:                  m.AuthPw,
		PhoneNumber:             m.PhoneNumber,
		AccountId:               m.AccountId,
		MoodleUserAssignmentsId: m.MoodleUserAssignmentsId,
		AssignmentIds:           m.AssignmentIds,
		NotSetYet:               m.NotSetYet,
	}
}
