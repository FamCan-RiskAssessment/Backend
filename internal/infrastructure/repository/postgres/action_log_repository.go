package postgres

import (
	"github.com/FamCan-RiskAssessment/Backend/internal/domain/entity"
	"github.com/FamCan-RiskAssessment/Backend/internal/domain/repository/postgres"
	"github.com/FamCan-RiskAssessment/Backend/internal/infrastructure/database"
)

type ActionLogRepository struct {
}

func NewActionLogRepository() *ActionLogRepository {
	return &ActionLogRepository{}
}

func (r *ActionLogRepository) CreateActionLog(db database.Database, actionLog *entity.ActionLog) error {
	return db.GetDB().Create(actionLog).Error
}

func (r *ActionLogRepository) FindAllActionLogs(db database.Database, options *postgres.QueryOptions) ([]*entity.ActionLog, error) {
	var logs []*entity.ActionLog

	query := db.GetDB()
	query = applyQueryOptions(query, options)

	err := query.Find(&logs).Error

	return logs, err
}

func (r *ActionLogRepository) CountAllActionLogs(db database.Database, options *postgres.QueryOptions) (int64, error) {
	var count int64

	query := db.GetDB()
	err := query.Model(&entity.ActionLog{}).Count(&count).Error

	return count, err
}
