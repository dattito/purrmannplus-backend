package models

type MoodleCourse struct {
	Courses []struct {
		FullName    string `json:"fullname"`
		Assignments []struct {
			ID int `json:"id"`
		} `json:"assignments"`
	} `json:"courses"`
}

type MoodleAssignments struct {
	AccountId   string
	Assignments []int
}

type MoodleAssignmentInfo struct {
	AuthId                  string
	AuthPw                  string
	PhoneNumber             string
	AccountId               string
	MoodleUserAssignmentsId string
	AssignmentIds           []int
	NotSetYet               bool
}
