package exception

import "fmt"

const (
	ServiceSMS = "sms"

	ReasonNetwork  = "network"
	ReasonUpstream = "upstream"
)

type ServiceUnavailableError struct {
	Service     string
	Reason      string
	Message     string
	OriginalErr error
}

func (e ServiceUnavailableError) Error() string {
	if e.OriginalErr != nil {
		return fmt.Sprintf("service unavailable (%s/%s): %s - %v", e.Service, e.Reason, e.Message, e.OriginalErr)
	}
	return fmt.Sprintf("service unavailable (%s/%s): %s", e.Service, e.Reason, e.Message)
}

func NewServiceUnavailableError(service, reason, message string, originalErr error) *ServiceUnavailableError {
	if message == "" {
		message = "service unavailable"
	}
	return &ServiceUnavailableError{
		Service:     service,
		Reason:      reason,
		Message:     message,
		OriginalErr: originalErr,
	}
}
