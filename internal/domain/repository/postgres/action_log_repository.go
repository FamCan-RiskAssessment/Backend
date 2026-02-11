package postgres

import (
	"github.com/FamCan-RiskAssessment/Backend/internal/domain/entity"
	"github.com/FamCan-RiskAssessment/Backend/internal/domain/enum"
	"github.com/FamCan-RiskAssessment/Backend/internal/infrastructure/database"
)

type ActionLogFilters struct {
	Action   *enum.ActionType
	ActorID  *uint
	DateFrom *string
	DateTo   *string
}

type ActionLogRepository interface {
	CreateActionLog(db database.Database, actionLog *entity.ActionLog) error
	FindAllActionLogs(db database.Database, options *QueryOptions, filters *ActionLogFilters) ([]*entity.ActionLog, error)
	CountAllActionLogs(db database.Database, options *QueryOptions, filters *ActionLogFilters) (int64, error)
}
