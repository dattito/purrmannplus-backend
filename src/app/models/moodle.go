package models

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
