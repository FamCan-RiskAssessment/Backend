package validation

import (
	"github.com/FamCan-RiskAssessment/Backend/internal/domain/enum"
	"github.com/FamCan-RiskAssessment/Backend/internal/domain/exception"
)

type FormStatusTransition struct {
	FromStatus enum.FormStatus
	ToStatus   enum.FormStatus
}

// validTransitions defines all allowed status transitions
var validTransitions = map[FormStatusTransition]bool{
	// Draft can transition to Submitted
	{enum.FormStatusDraft, enum.FormStatusSubmitted}: true,

	// Submitted can transition to:
	{enum.FormStatusSubmitted, enum.FormStatusWaitingForPatientResponse}: true,
	{enum.FormStatusSubmitted, enum.FormStatusReadyForCalculation}:       true,
	{enum.FormStatusSubmitted, enum.FormStatusRejected}:                  true,

	// WaitingForPatientResponse can transition to:
	{enum.FormStatusWaitingForPatientResponse, enum.FormStatusWaitingForDocuments}: true,

	// WaitingForDocuments can transition to:
	{enum.FormStatusWaitingForDocuments, enum.FormStatusSubmitted}: true,

	// Rejected can transition to:
	{enum.FormStatusRejected, enum.FormStatusDraft}: true,

	// ReadyForCalculation can transition to:
	{enum.FormStatusReadyForCalculation, enum.FormStatusCalculated}: true,
}

// ValidateStatusTransition checks if a status transition is allowed
// Returns an error if the transition is invalid
func ValidateStatusTransition(from, to enum.FormStatus) error {
	transition := FormStatusTransition{FromStatus: from, ToStatus: to}

	if !validTransitions[transition] {
		return exception.FieldError{
			Field: "status",
			Tag:   "invalid_transition",
		}
	}

	return nil
}

// CanUserEditFormInStatus checks if users can edit forms in the given status
// Patients can only edit in Draft and WaitingForDocuments
// Operators/admins can edit in all statuses except Calculated
func CanUserEditFormInStatus(status enum.FormStatus, isPatient bool) bool {
	if isPatient {
		// Patients can only edit in Draft and WaitingForDocuments
		return status == enum.FormStatusDraft || status == enum.FormStatusWaitingForDocuments
	}

	// Operators/admins can edit in all statuses except Calculated
	return status != enum.FormStatusCalculated
}

// GetNextPossibleStatuses returns all valid next statuses for the given current status
func GetNextPossibleStatuses(currentStatus enum.FormStatus) []enum.FormStatus {
	var nextStatuses []enum.FormStatus

	for transition := range validTransitions {
		if transition.FromStatus == currentStatus {
			nextStatuses = append(nextStatuses, transition.ToStatus)
		}
	}

	return nextStatuses
}
