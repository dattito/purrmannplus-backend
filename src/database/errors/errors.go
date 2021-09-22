package errors

type DatabaseError struct {
	message string
}

func (e *DatabaseError) Error() string {
	return e.message
}

var ErrRecordNotFound DatabaseError = DatabaseError{
	message: "record not found",
}
