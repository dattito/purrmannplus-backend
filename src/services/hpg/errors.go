package hpg

type wrongCredentialsError struct{}

func (e *wrongCredentialsError) Error() string {
	return "HPG: wrong credentials"
}

var WrongCredentialsError error = &wrongCredentialsError{}
