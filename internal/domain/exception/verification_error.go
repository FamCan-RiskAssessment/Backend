package exception

type VerificationError struct {
	Field   string
	Message string
}

func (e VerificationError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	return "verification failed for " + e.Field
}
