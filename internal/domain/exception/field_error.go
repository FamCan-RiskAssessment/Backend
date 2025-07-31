package exception

import "fmt"

const (
	fieldErrMsg = "Error:Field validation for '%s' failed on the '%s' tag"
)

type FieldError struct {
	Field string
	Tag   string
}

func (e FieldError) Error() string {
	return fmt.Sprintf(fieldErrMsg, e.Field, e.Tag)
}
