package external

import (
	formdto "github.com/FamCan-RiskAssessment/Backend/internal/application/dto/form"
	"github.com/stretchr/testify/mock"
)

type VerificationClientMock struct {
	mock.Mock
}

func NewVerificationClientMock() *VerificationClientMock {
	return &VerificationClientMock{}
}

func (v *VerificationClientMock) VerifyPhoneAndSSN(phone string, ssn string) (bool, error) {
	args := v.Called(phone, ssn)
	return args.Bool(0), args.Error(1)
}

func (v *VerificationClientMock) GetAddressByPostalCode(postalCode string) (*formdto.PostalCodeInfoResponse, error) {
	args := v.Called(postalCode)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*formdto.PostalCodeInfoResponse), args.Error(1)
}
