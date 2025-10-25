package service

import (
	"github.com/FamCan-RiskAssessment/Backend/bootstrap"
	actionlogdto "github.com/FamCan-RiskAssessment/Backend/internal/application/dto/actionLog"
	"github.com/FamCan-RiskAssessment/Backend/internal/domain/entity"
	"github.com/FamCan-RiskAssessment/Backend/internal/domain/enum"
	"github.com/FamCan-RiskAssessment/Backend/internal/domain/repository/postgres"
	"github.com/FamCan-RiskAssessment/Backend/internal/infrastructure/database"
)

type ActionLogService struct {
	constants           *bootstrap.Constants
	actionLogRepository postgres.ActionLogRepository
	// userService         usecase.UserService
	db database.Database
}

func NewActionLogService(
	constants *bootstrap.Constants,
	actionLogRepository postgres.ActionLogRepository,
	// userService usecase.UserService,
	db database.Database,
) *ActionLogService {
	return &ActionLogService{
		constants:           constants,
		actionLogRepository: actionLogRepository,
		// userService:         userService,
		db: db,
	}
}

func (als *ActionLogService) LogAction(request actionlogdto.LogAction) error {
	// user, err := als.userService.GetUserByID(request.ActorID)
	// if err != nil {
	// 	return err
	// }
	// if user == nil {
	// 	notFoundError := exception.NotFoundError{Item: als.constants.Field.User}
	// 	return notFoundError
	// }

	// if request.TargetID != nil {
	// 	user, err := als.userService.GetUserByID(*request.TargetID)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	if user == nil {
	// 		notFoundError := exception.NotFoundError{Item: als.constants.Field.User}
	// 		return notFoundError
	// 	}
	// }

	actionLog := &entity.ActionLog{
		ActorID:    request.ActorID,
		TargetID:   request.TargetID,
		Action:     request.Action,
		Resource:   request.Resource,
		ResourceID: request.ResourceID,
		Details:    request.Details,
	}

	err := als.actionLogRepository.CreateActionLog(als.db, actionLog)
	if err != nil {
		return err
	}
	return nil
}

func (als *ActionLogService) GetAllActionLogs(offset, limit int) ([]actionlogdto.LogResponse, int64, error) {
	options := postgres.NewQueryOptions().
		WithPagination(limit, offset)

	actionLogs, err := als.actionLogRepository.FindAllActionLogs(als.db, options)
	if err != nil {
		return nil, 0, err
	}

	count, err := als.actionLogRepository.CountAllActionLogs(als.db, options)
	if err != nil {
		return nil, 0, err
	}

	response := make([]actionlogdto.LogResponse, len(actionLogs))
	for i, actionLog := range actionLogs {
		response[i] = actionlogdto.LogResponse{
			ID:        actionLog.ID,
			Action:    actionLog.Action.String(),
			Resource:  actionLog.Resource,
			Details:   actionLog.Details,
			CreatedAt: actionLog.CreatedAt,
		}
	}

	return response, count, nil
}

func (als *ActionLogService) GetAllActionTypes() ([]actionlogdto.ActionType, error) {
	actionTypes := enum.GetAllActionTypes()
	response := make([]actionlogdto.ActionType, len(actionTypes))
	for i, actionType := range actionTypes {
		response[i] = actionlogdto.ActionType{
			ID:   uint(actionType),
			Name: actionType.String(),
		}
	}
	return response, nil
}
