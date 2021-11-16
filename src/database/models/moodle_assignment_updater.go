package models

type MoodleUserAssignmentsDBModel struct {
	Id            string
	AccountId     string
	AssignmentIds []int
	NotSetYet     bool
}

type AccountCredentialsAndPhoneNumberAndMoodleUserAssignmentsDBModel struct {
	AuthId                  string
	AuthPw                  string
	PhoneNumber             string
	AccountId               string
	MoodleUserAssignmentsId string
	AssignmentIds           []int
	NotSetYet               bool
}
