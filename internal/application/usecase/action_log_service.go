package usecase

import actionlogdto "github.com/FamCan-RiskAssessment/Backend/internal/application/dto/actionLog"

type ActionLogService interface {
	LogAction(request actionlogdto.LogAction) error
	GetAllActionLogs(offset, limit int) ([]actionlogdto.LogResponse, int64, error)
	GetAllActionTypes() ([]actionlogdto.ActionType, error)
}
