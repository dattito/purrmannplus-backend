package routes

const (
	HealthRoute                          = "/health"
	AboutRoute                           = "/about"
	AccountLoginRoute                    = "/login"
	AccountLogoutRoute                   = "/logout"
	IsLoggedInRoute                      = "/login_check"
	AddAccountRoute                      = "/accounts"
	GetAccountsRoute                     = "/accounts"
	DeleteAccountRoute                   = "/accounts"
	SendPhoneNumberConfirmationLinkRoute = "/accounts/phone_number"
	AddPhoneNumberRoute                  = "/accounts/phone_number/validate"

	AddAccountToSubstitutionUpdaterRoute      = "/substitution_updater"
	RemoveAccountFromSubstitutionUpdaterRoute = "/substitution_updater"

	AddAccountToMoodleAssignmentUpdaterRoute      = "/moodle_assignment_updater"
	RemoveAccountFromMoodleAssignmentUpdaterRoute = "/moodle_assignment_updater"

	RegistrationSpeedFormRoute                        = "/registration_speed_form"
	RegistrationSpeedFormSubstitutionCredentialsRoute = "/registration_speed_form/substitution-credentials"
	RegistrationSpeedFormValidationRoute              = "/registration_speed_form/validate"
	RegistrationSpeedFormFinishRoute                  = "/registration_speed_form/finish"
	RegistrationSpeedFormInfoRoute                    = "/registration_speed_form/info"
)
