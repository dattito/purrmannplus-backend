package models

type PostRegistrationSpeedFormRequest struct {
	Username    string `form:"username"`
	Password    string `form:"password"`
	PhoneNumber string `form:"phoneNumber"`
}

type PostValidateRegistrationSpeedFormRequest struct {
	Code string `form:"code"`
}
