package actionlogdto

import "github.com/FamCan-RiskAssessment/Backend/internal/domain/enum"

type LogAction struct {
	ActorID    uint
	TargetID   *uint
	Action     enum.ActionType
	Resource   *string
	ResourceID *uint
	Details    string
}
