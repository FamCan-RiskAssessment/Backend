package service

import (
	"github.com/FamCan-RiskAssessment/Backend/bootstrap"
	actionlogdto "github.com/FamCan-RiskAssessment/Backend/internal/application/dto/actionLog"
	generaldto "github.com/FamCan-RiskAssessment/Backend/internal/application/dto/general"
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

var actionLogAllowedSortColumns = []string{"id", "created_at", "action"}

func (als *ActionLogService) GetAllActionLogs(request actionlogdto.GetAllActionLogsRequest) ([]actionlogdto.LogResponse, int64, error) {
	options := postgres.NewQueryOptions().
		WithPagination(request.Limit, request.Offset)

	sortCol := "created_at"
	if request.SortBy != nil {
		sortCol = postgres.ValidateSortColumn(*request.SortBy, actionLogAllowedSortColumns, "created_at")
	}
	asc := false
	if request.SortOrder != nil && *request.SortOrder == "asc" {
		asc = true
	}
	options.WithSorting(sortCol, asc)

	if request.Search != nil && *request.Search != "" {
		options.WithSearch(*request.Search, []string{"details", "resource"})
	}

	filters := &postgres.ActionLogFilters{
		Action:   request.Action,
		ActorID:  request.ActorID,
		DateFrom: request.DateFrom,
		DateTo:   request.DateTo,
	}

	actionLogs, err := als.actionLogRepository.FindAllActionLogs(als.db, options, filters)
	if err != nil {
		return nil, 0, err
	}

	count, err := als.actionLogRepository.CountAllActionLogs(als.db, options, filters)
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

func (als *ActionLogService) GetAllActionTypes() ([]generaldto.EnumResponse, error) {
	actionTypes := enum.GetAllActionTypes()
	response := make([]generaldto.EnumResponse, len(actionTypes))
	for i, actionType := range actionTypes {
		response[i] = generaldto.EnumResponse{
			ID:   uint(actionType),
			Name: actionType.String(),
		}
	}
	return response, nil
}
