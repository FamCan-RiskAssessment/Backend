package exception

type ForbiddenType string

const (
	ForbiddenTypeBannedUser            ForbiddenType = "banned_user"
	ForbiddenTypeUnapprovedCorporation ForbiddenType = "unapproved_corporation"
	ForbiddenTypeVerificationFailed    ForbiddenType = "verification_failed"
)

type ForbiddenError struct {
	Type     ForbiddenType
	Resource string
	Message  string
}

func (e ForbiddenError) Error() string {
	if e.Message != "" {
		return "Forbidden: " + e.Message
	}
	return "Forbidden: Access to " + e.Resource + " is not allowed"
}

func NewBannedUserForbiddenError() ForbiddenError {
	return ForbiddenError{
		Type:    ForbiddenTypeBannedUser,
		Message: "Your account is banned and cannot perform this operation.",
	}
}

func NewUnapprovedCorporationForbiddenError() ForbiddenError {
	return ForbiddenError{
		Type:    ForbiddenTypeUnapprovedCorporation,
		Message: "Vendor approval is required to access this resource.",
	}
}

func NewVerificationFailedForbiddenError() ForbiddenError {
	return ForbiddenError{
		Type:    ForbiddenTypeVerificationFailed,
		Message: "Phone and social security number verification failed.",
	}
}
