package postgres

import (
	"github.com/FamCan-RiskAssessment/Backend/internal/domain/entity"
	"github.com/FamCan-RiskAssessment/Backend/internal/infrastructure/database"
)

type ActionLogRepository interface {
	CreateActionLog(db database.Database, actionLog *entity.ActionLog) error
	FindAllActionLogs(db database.Database, options *QueryOptions) ([]*entity.ActionLog, error)
	CountAllActionLogs(db database.Database, options *QueryOptions) (int64, error)
}
