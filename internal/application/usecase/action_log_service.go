package usecase

import (
	actionlogdto "github.com/FamCan-RiskAssessment/Backend/internal/application/dto/actionLog"
	generaldto "github.com/FamCan-RiskAssessment/Backend/internal/application/dto/general"
)

type ActionLogService interface {
	LogAction(request actionlogdto.LogAction) error
	GetAllActionLogs(offset, limit int) ([]actionlogdto.LogResponse, int64, error)
	GetAllActionTypes() ([]generaldto.EnumResponse, error)
}
