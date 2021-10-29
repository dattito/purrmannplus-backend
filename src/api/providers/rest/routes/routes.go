package routes

const (
	HealthRoute                            = "/health"
	AboutRoute                             = "/about"
	AccountLoginRoute                      = "/login"
	AccountLogoutRoute                     = "/logout"
	IsLoggedInRoute                        = "/login_check"
	AddAccountRoute                        = "/accounts"
	GetAccountsRoute                       = "/accounts"
	DeleteAccountRoute                     = "/accounts"
	SendPhoneNumberConfirmationLinkRoute   = "/accounts/phone_number"
	AddPhoneNumberRoute                    = "/accounts/phone_number/validate"
	RegisterToSubstitutionUpdaterRoute     = "/substitution_updater"
	UnregisterFromSubstitutionUpdaterRoute = "/substitution_updater"

	SubstitutionSpeedFormRoute           = "/substitution_speed_form"
	SubstitutionSpeedFormValidationRoute = "/substitution_speed_form/validate"
	SubstitutionSpeedFormFinishRoute     = "/substitution_speed_form/finish"
	SubstitutionSpeedFormInfoRoute       = "/substitution_speed_form/info"
)
