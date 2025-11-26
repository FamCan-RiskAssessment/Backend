package exception

type FileValidationError struct {
	Message string
}

func (e FileValidationError) Error() string {
	return e.Message
}
