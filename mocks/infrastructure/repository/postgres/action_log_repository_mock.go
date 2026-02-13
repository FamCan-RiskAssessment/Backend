package postgres

import (
	"github.com/FamCan-RiskAssessment/Backend/internal/domain/entity"
	"github.com/FamCan-RiskAssessment/Backend/internal/domain/repository/postgres"
	"github.com/FamCan-RiskAssessment/Backend/internal/infrastructure/database"
	"github.com/stretchr/testify/mock"
)

type ActionLogRepositoryMock struct {
	mock.Mock
}

func NewActionLogRepositoryMock() *ActionLogRepositoryMock {
	return &ActionLogRepositoryMock{}
}

func (a *ActionLogRepositoryMock) CreateActionLog(db database.Database, actionLog *entity.ActionLog) error {
	args := a.Called(db, actionLog)
	return args.Error(0)
}

func (a *ActionLogRepositoryMock) FindAllActionLogs(db database.Database, options *postgres.QueryOptions, filters *postgres.ActionLogFilters) ([]*entity.ActionLog, error) {
	args := a.Called(db, options, filters)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.ActionLog), args.Error(1)
}

func (a *ActionLogRepositoryMock) CountAllActionLogs(db database.Database, options *postgres.QueryOptions, filters *postgres.ActionLogFilters) (int64, error) {
	args := a.Called(db, options, filters)
	return args.Get(0).(int64), args.Error(1)
}
