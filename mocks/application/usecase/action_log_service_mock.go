package usecase

import (
	actionlogdto "github.com/FamCan-RiskAssessment/Backend/internal/application/dto/actionLog"
	generaldto "github.com/FamCan-RiskAssessment/Backend/internal/application/dto/general"
	"github.com/stretchr/testify/mock"
)

type ActionLogServiceMock struct {
	mock.Mock
}

func NewActionLogServiceMock() *ActionLogServiceMock {
	return &ActionLogServiceMock{}
}

func (a *ActionLogServiceMock) LogAction(request actionlogdto.LogAction) error {
	args := a.Called(request)
	return args.Error(0)
}

func (a *ActionLogServiceMock) GetAllActionLogs(offset, limit int) ([]actionlogdto.LogResponse, int64, error) {
	args := a.Called(offset, limit)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]actionlogdto.LogResponse), args.Get(1).(int64), args.Error(2)
}

func (a *ActionLogServiceMock) GetAllActionTypes() ([]generaldto.EnumResponse, error) {
	args := a.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]generaldto.EnumResponse), args.Error(1)
}
