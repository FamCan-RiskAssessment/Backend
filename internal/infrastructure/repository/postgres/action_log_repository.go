package postgres

import (
	"github.com/FamCan-RiskAssessment/Backend/internal/domain/entity"
	"github.com/FamCan-RiskAssessment/Backend/internal/domain/repository/postgres"
	"github.com/FamCan-RiskAssessment/Backend/internal/infrastructure/database"
	"gorm.io/gorm"
)

type ActionLogRepository struct {
}

func NewActionLogRepository() *ActionLogRepository {
	return &ActionLogRepository{}
}

func ApplyActionLogFilters(query *gorm.DB, filters *postgres.ActionLogFilters) *gorm.DB {
	if filters == nil {
		return query
	}
	if filters.Action != nil {
		query = query.Where("action = ?", *filters.Action)
	}
	if filters.ActorID != nil {
		query = query.Where("actor_id = ?", *filters.ActorID)
	}
	if filters.DateFrom != nil {
		query = query.Where("created_at >= ?", *filters.DateFrom)
	}
	if filters.DateTo != nil {
		query = query.Where("created_at <= ?", *filters.DateTo)
	}
	return query
}

func (r *ActionLogRepository) CreateActionLog(db database.Database, actionLog *entity.ActionLog) error {
	return db.GetDB().Create(actionLog).Error
}

func (r *ActionLogRepository) FindAllActionLogs(db database.Database, options *postgres.QueryOptions, filters *postgres.ActionLogFilters) ([]*entity.ActionLog, error) {
	var logs []*entity.ActionLog

	query := db.GetDB()
	query = ApplyActionLogFilters(query, filters)
	query = applyQueryOptions(query, options)

	err := query.Preload("Actor").Preload("Actor.UserProfile").Preload("Target").Preload("Target.UserProfile").Find(&logs).Error
	return logs, err
}

func (r *ActionLogRepository) CountAllActionLogs(db database.Database, options *postgres.QueryOptions, filters *postgres.ActionLogFilters) (int64, error) {
	var count int64

	query := db.GetDB().Model(&entity.ActionLog{})
	query = ApplyActionLogFilters(query, filters)
	query = applySearchOnly(query, options)

	err := query.Count(&count).Error
	return count, err
}
