package hpg

type wrongCredentialsError struct{}

func (*wrongCredentialsError) Error() string {
	return "HPG: wrong credentials"
}

var WrongCredentialsError error = &wrongCredentialsError{}
