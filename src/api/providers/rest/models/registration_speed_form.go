package models

type PostRegistrationSpeedFormRequest struct {
	Username    string `form:"username"`
	Password    string `form:"password"`
	PhoneNumber string `form:"phoneNumber"`
}

type PostCustomSubsitutionCredentialsRequest struct {
	AuthId    string `form:"authId"`
	AuthPw string `form:"authPw"`
}

type PostValidateRegistrationSpeedFormRequest struct {
	Code string `form:"code"`
}
