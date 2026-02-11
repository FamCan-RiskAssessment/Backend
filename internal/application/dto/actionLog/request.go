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

type GetAllActionLogsRequest struct {
	Offset    int
	Limit     int
	SortBy    *string
	SortOrder *string
	Search    *string
	Action    *enum.ActionType
	ActorID   *uint
	DateFrom  *string
	DateTo    *string
}
