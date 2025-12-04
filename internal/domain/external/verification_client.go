package external

type VerificationClient interface {
	VerifyPhoneAndSSN(phone string, ssn string) (bool, error)
}
