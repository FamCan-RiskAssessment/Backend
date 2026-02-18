package external

import formdto "github.com/FamCan-RiskAssessment/Backend/internal/application/dto/form"

type VerificationClient interface {
	VerifyPhoneAndSSN(phone string, ssn string) (bool, error)
	GetAddressByPostalCode(postalCode string) (*formdto.PostalCodeInfoResponse, error)
}
