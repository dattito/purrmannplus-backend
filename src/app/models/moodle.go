package models

type Assignments struct {
	Courses []struct {
		FullName    string `json:"fullname"`
		Assignments []struct {
			ID int `json:"id"`
		} `json:"assignments"`
	} `json:"courses"`
}
