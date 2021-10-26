package models

type PostSubstitutionSpeedFormRequest struct {
	Username    string `form:"username"`
	Password    string `form:"password"`
	PhoneNumber string `form:"phone_number"`
}

type PostValidateSubstitutionSpeedFormRequest struct {
	Code string `form:"code"`
}
