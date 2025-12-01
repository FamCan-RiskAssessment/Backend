package exception

import "fmt"

// CalcErrorType represents the type of calculation error
type CalcErrorType string

const (
	ErrorTypeAPIFailure     CalcErrorType = "api_failure"
	ErrorTypeMissingData    CalcErrorType = "missing_data"
	ErrorTypeInvalidData    CalcErrorType = "invalid_data"
	ErrorTypeDatabaseError  CalcErrorType = "database_error"
)

type CalcError struct {
	Type      CalcErrorType
	Model     string // e.g., "PREMM5", "BCRA", "GAIL", "PLCO"
	Message   string
	OrigErr   error // Original error for debugging
}

func (ce CalcError) Error() string {
	if ce.OrigErr != nil {
		return fmt.Sprintf("calculation error (%s): %s - %v", ce.Type, ce.Message, ce.OrigErr)
	}
	return fmt.Sprintf("calculation error (%s): %s", ce.Type, ce.Message)
}
