package models

type PostSubstitutionSpeedFormRequest struct {
	Username    string `form:"username"`
	Password    string `form:"password"`
	PhoneNumber string `form:"phoneNumber"`
}

type PostValidateSubstitutionSpeedFormRequest struct {
	Code string `form:"code"`
}
