package service

import (
	"fmt"
	"time"

	"github.com/FamCan-RiskAssessment/Backend/bootstrap"
	actionlogdto "github.com/FamCan-RiskAssessment/Backend/internal/application/dto/actionLog"
	formdto "github.com/FamCan-RiskAssessment/Backend/internal/application/dto/form"
	generaldto "github.com/FamCan-RiskAssessment/Backend/internal/application/dto/general"
	"github.com/FamCan-RiskAssessment/Backend/internal/application/usecase"
	"github.com/FamCan-RiskAssessment/Backend/internal/domain/entity"
	"github.com/FamCan-RiskAssessment/Backend/internal/domain/enum"
	"github.com/FamCan-RiskAssessment/Backend/internal/domain/exception"
	postgres "github.com/FamCan-RiskAssessment/Backend/internal/domain/repository/postgres"
	"github.com/FamCan-RiskAssessment/Backend/internal/domain/storage/s3"
	"github.com/FamCan-RiskAssessment/Backend/internal/infrastructure/database"
)

type FormService struct {
	constants        *bootstrap.Constants
	formRepository   postgres.FormRepository
	userService      usecase.UserService
	actionLogService usecase.ActionLogService
	s3Storage        s3.S3Storage
	db               database.Database
}

func NewFormService(
	constants *bootstrap.Constants,
	formRepository postgres.FormRepository,
	userService usecase.UserService,
	actionLogService usecase.ActionLogService,
	s3Storage s3.S3Storage,
	db database.Database,
) *FormService {
	return &FormService{
		constants:        constants,
		formRepository:   formRepository,
		userService:      userService,
		actionLogService: actionLogService,
		s3Storage:        s3Storage,
		db:               db,
	}
}

func (formService *FormService) canOperatorEditForm(form *entity.Form, operatorID uint) (bool, error) {
	if form.FilledByOperatorID != nil && *form.FilledByOperatorID == operatorID {
		return true, nil
	}

	if form.OperatorID != nil && *form.OperatorID == operatorID {
		return true, nil
	}

	return false, nil
}

func (formService *FormService) isOperator(userID uint) (bool, error) {
	userRoles, err := formService.userService.GetUserRoles(userID)
	if err != nil {
		return false, err
	}

	for _, role := range userRoles {
		if role.Name == enum.Operator.String() {
			return true, nil
		}
	}

	return false, nil
}

func (formService *FormService) isSupervisor(userID uint) (bool, error) {
	userRoles, err := formService.userService.GetUserRoles(userID)
	if err != nil {
		return false, err
	}

	for _, role := range userRoles {
		if role.Name == enum.Supervisor.String() {
			return true, nil
		}
	}

	return false, nil
}

func (formService *FormService) hasPermission(userID uint, permission enum.PermissionType) (bool, error) {
	userRoles, err := formService.userService.GetUserRoles(userID)
	if err != nil {
		return false, err
	}

	permissionName := permission.String()
	allPermissionName := enum.PermissionAll.String()

	for _, role := range userRoles {
		for _, perm := range role.Permissions {
			// PermissionResponse.Name is permission.Type.String()
			if perm.Name == permissionName || perm.Name == allPermissionName {
				return true, nil
			}
		}
	}

	return false, nil
}

func (formService *FormService) logOperatorFormUpdate(operatorID, formID uint, sectionName string) {
	log := actionlogdto.LogAction{
		ActorID:    operatorID,
		Action:     enum.ActionTypeOperatorUpdatedForm,
		ResourceID: &formID,
		Details:    "اپراتور بخش " + sectionName + " فرم را بروزرسانی کرد",
	}
	formService.actionLogService.LogAction(log)
}

func (formService *FormService) CreateBasicInfoForm(request formdto.CreateBasicFormRequest) (formdto.BasicFormResponse, error) {
	user, err := formService.userService.GetUserByID(request.UserID)
	if err != nil {
		return formdto.BasicFormResponse{}, err
	}
	if user == nil {
		notFoundError := exception.NotFoundError{Item: formService.constants.Field.User}
		return formdto.BasicFormResponse{}, notFoundError
	}

	if request.FilledByOperatorID != nil {
		err = formService.userService.ValidateUserForFormCreation(*request.FilledByOperatorID, request.UserID)
		if err != nil {
			return formdto.BasicFormResponse{}, err
		}
		operator, err := formService.userService.GetUserByID(*request.FilledByOperatorID)
		if err != nil {
			return formdto.BasicFormResponse{}, err
		}
		if operator == nil {
			notFoundError := exception.NotFoundError{Item: formService.constants.Field.User}
			return formdto.BasicFormResponse{}, notFoundError
		}
	}

	form := &entity.Form{
		UserID:             request.UserID,
		Status:             enum.FormStatusPending,
		FilledByOperatorID: request.FilledByOperatorID,
	}

	if err = formService.formRepository.CreateForm(formService.db, form); err != nil {
		return formdto.BasicFormResponse{}, err
	}

	basic := &entity.BasicInfo{
		FormID:               form.ID,
		Gender:               enum.Gender(uint(request.Gender)),
		BirthDate:            request.BirthDate,
		IsAtba:               request.IsAtba,
		SocialSecurityNumber: request.SocialSecurityNumber,
		Height:               request.Height,
		Weight:               request.Weight,
	}

	if err = formService.formRepository.CreateBasicInfo(formService.db, basic); err != nil {
		return formdto.BasicFormResponse{}, err
	}

	// Log operator action if form was created by operator
	if request.FilledByOperatorID != nil {
		log := actionlogdto.LogAction{
			ActorID:    *request.FilledByOperatorID,
			TargetID:   &request.UserID,
			Action:     enum.ActionTypeOperatorCreatedForm,
			ResourceID: &form.ID,
			Details:    "اپراتور فرم را برای کاربر ایجاد کرد",
		}
		formService.actionLogService.LogAction(log)
	}

	response := formdto.BasicFormResponse{
		FormID:             form.ID,
		Status:             form.Status.String(),
		UserID:             form.UserID,
		FilledByOperatorID: form.FilledByOperatorID,
		CreatedAt:          form.CreatedAt,
		UpdatedAt:          form.UpdatedAt,
	}

	return response, nil
}

func (formService *FormService) UpsertGeneralHealth(request formdto.UpsertGeneralHealthRequest) error {
	form, err := formService.formRepository.FindFormByID(formService.db, request.FormID)
	if err != nil {
		return err
	}
	if form == nil {
		notFoundError := exception.NotFoundError{Item: formService.constants.Field.Form}
		return notFoundError
	}

	isOp, err := formService.isOperator(request.UserID)
	if err != nil {
		return err
	}
	if isOp {
		canEdit, err := formService.canOperatorEditForm(form, request.UserID)
		if err != nil {
			return err
		}
		if !canEdit {
			forbiddenError := exception.ForbiddenError{Resource: formService.constants.Field.Form}
			return forbiddenError
		}
	}

	// if form.UserID != request.UserID {
	// 	ForbiddenError := exception.ForbiddenError{Message: formService.constants.Field.Form}
	// 	return ForbiddenError
	// }

	info, err := formService.formRepository.FindGeneralHealthByFormID(formService.db, request.FormID)
	if err != nil {
		return err
	}

	if info == nil {
		info = &entity.GeneralHealthInfo{FormID: request.FormID}
	}

	info.DrinksAlcohol = request.DrinksAlcohol
	info.CupsPerWeek = request.CupsPerWeek
	info.LastMonthSabzijatMeal = request.LastMonthSabzijatMeal
	info.LastMonthSabzijatWeight = request.LastMonthSabzijatWeight
	info.MediumActivityMonthInYear = request.MediumActivityMonthInYear
	info.MediumActivityHourInWeek = request.MediumActivityHourInWeek
	info.HardActivityMonthInYear = request.HardActivityMonthInYear
	info.HardActivityHourInWeek = request.HardActivityHourInWeek
	info.SmokeAtLeast100 = request.SmokeAtLeast100
	info.SmokingAge = request.SmokingAge
	info.SmokingNow = request.SmokingNow
	info.LeaveSmokingAge = request.LeaveSmokingAge
	info.CountSmokingDaily = request.CountSmokingDaily
	info.CountGheliandaily = request.CountGheliandaily
	info.CountSmokingDailyPast = request.CountSmokingDailyPast
	info.CountGheliandailyPast = request.CountGheliandailyPast

	if info.ID == 0 {
		err = formService.formRepository.CreateGeneralHealth(formService.db, info)
	} else {
		err = formService.formRepository.UpdateGeneralHealth(formService.db, info)
	}
	if err != nil {
		return err
	}

	// Log operator action if performed by operator
	if isOp {
		formService.logOperatorFormUpdate(request.UserID, request.FormID, "اطلاعات سلامت عمومی")
	}

	return nil
}

func (formService *FormService) UpsertMamography(request formdto.UpsertMamographyRequest) error {
	form, err := formService.formRepository.FindFormByID(formService.db, request.FormID)
	if err != nil {
		return err
	}
	if form == nil {
		notFoundError := exception.NotFoundError{Item: formService.constants.Field.Form}
		return notFoundError
	}

	isOp, err := formService.isOperator(request.UserID)
	if err != nil {
		return err
	}
	if isOp {
		canEdit, err := formService.canOperatorEditForm(form, request.UserID)
		if err != nil {
			return err
		}
		if !canEdit {
			forbiddenError := exception.ForbiddenError{Resource: formService.constants.Field.Form}
			return forbiddenError
		}
	}

	info, err := formService.formRepository.FindMamographyByFormID(formService.db, request.FormID)
	if err != nil {
		return err
	}
	if info == nil {
		info = &entity.MamoGraphyInfo{FormID: request.FormID}
	}

	info.GhaedeAge = request.GhaedeAge
	info.HasChildren = request.HasChildren
	info.NumberOfChildren = request.NumberOfChildren
	info.AgeOfFirstBirth = request.AgeOfFirstBirth
	info.MenopausalStatus = enum.MenopausalStatus(uint(request.MenopausalStatus))
	info.MenopauseAge = request.MenopauseAge
	info.HRT = request.HRT
	info.HRTUseLength = request.HRTUseLength
	info.LastFiveYearsHRTUse = request.LastFiveYearsHRTUse
	info.CurrentHRTUse = request.CurrentHRTUse
	info.IntendedHRTUse = request.IntendedHRTUse
	info.HRTType = request.HRTType
	info.Oral = request.Oral
	info.OralDuration = request.OralDuration
	info.OralTwoLastYears = request.OralTwoLastYears
	info.MamoGraphy = request.MamoGraphy
	info.Falop = request.Falop
	info.Andometrioz = request.Andometrioz
	info.LeavePestan = request.LeavePestan
	info.LeaveTokhmdan = request.LeaveTokhmdan
	info.LaDeColon = request.LaDeColon
	info.LaDePol = request.LaDePol
	info.AspLaMo = request.AspLaMo
	info.NsaiDLaMo = request.NsaiDLaMo
	info.LastFiveYearBloodTestInStool = request.LastFiveYearBloodTestInStool
	info.NumberOfBreastBiopsies = request.NumberOfBreastBiopsies
	info.HyperplasiaInBiopsy = (*enum.HyperplasiaInBiopsyStatus)(request.HyperplasiaInBiopsy)

	var pictureKey string
	var oldPicturePath *string
	if request.MamoGraphy != nil && *request.MamoGraphy && request.MamoGraphyPicture != nil && request.MamoGraphyPicture.Filename != "" {
		// Save the old picture path BEFORE overwriting it
		if info.MamoGraphyPicturePath != nil {
			oldPicturePath = info.MamoGraphyPicturePath
		}

		pictureKey = formService.constants.BucketPath.GetMamoGraphyPath(request.FormID, request.MamoGraphyPicture.Filename)
		info.MamoGraphyPicturePath = &pictureKey
	}

	isCreate := info.ID == 0

	if isCreate {
		err = formService.db.WithTransaction(func(tx database.Database) error {
			if err := formService.formRepository.CreateMamography(tx, info); err != nil {
				return err
			}

			if pictureKey != "" && request.MamoGraphyPicture != nil {
				if err := formService.s3Storage.UploadObject(enum.BucketTypeMamography, pictureKey, request.MamoGraphyPicture); err != nil {
					return err
				}
			}

			return nil
		})
		if err != nil {
			return err
		}
	} else {
		if pictureKey != "" && request.MamoGraphyPicture != nil {
			if oldPicturePath != nil && *oldPicturePath != "" && *oldPicturePath != pictureKey {
				if err := formService.s3Storage.DeleteObject(enum.BucketTypeMamography, *oldPicturePath); err != nil {
					fmt.Printf("Warning: failed to delete old mamography picture %s: %v\n", *oldPicturePath, err)

				}
			}

			if err := formService.s3Storage.UploadObject(enum.BucketTypeMamography, pictureKey, request.MamoGraphyPicture); err != nil {
				return err
			}
		}

		if err := formService.formRepository.UpdateMamography(formService.db, info); err != nil {
			return err
		}
	}

	// Log operator action if performed by operator
	if isOp {
		formService.logOperatorFormUpdate(request.UserID, request.FormID, "اطلاعات ماموگرافی")
	}

	return nil
}

func (formService *FormService) CreateCancer(request formdto.CreateCancerRequest) error {
	form, err := formService.formRepository.FindFormByID(formService.db, request.FormID)
	if err != nil {
		return err
	}
	if form == nil {
		notFoundError := exception.NotFoundError{Item: formService.constants.Field.Form}
		return notFoundError
	}

	isOp, err := formService.isOperator(request.UserID)
	if err != nil {
		return err
	}
	if isOp {
		canEdit, err := formService.canOperatorEditForm(form, request.UserID)
		if err != nil {
			return err
		}
		if !canEdit {
			forbiddenError := exception.ForbiddenError{Resource: formService.constants.Field.Form}
			return forbiddenError
		}
	}

	// Check for duplicate cancer (same type and age)
	existingCancers, err := formService.formRepository.FindCancersByFormID(formService.db, request.FormID)
	if err != nil {
		return err
	}
	for _, existing := range existingCancers {
		if existing.CancerType == enum.CancerType(request.CancerType) && existing.CancerAge == request.CancerAge {
			conflictError := exception.ConflictErrors{}
			conflictError.Add(formService.constants.Field.Cancer, formService.constants.Tag.AlreadyExist)
			return conflictError
		}
	}

	info := &entity.CancerInfo{
		FormID:     request.FormID,
		CancerAge:  request.CancerAge,
		CancerType: enum.CancerType(request.CancerType),
	}

	// Handle picture upload
	if request.Picture != nil && request.Picture.Filename != "" {
		pictureKey := formService.constants.BucketPath.GetCancerPath(
			request.FormID,
			enum.CancerType(request.CancerType),
			request.Picture.Filename,
		)
		info.PicturePath = &pictureKey

		// Upload picture
		if err := formService.s3Storage.UploadObject(enum.BucketTypeCancer, pictureKey, request.Picture); err != nil {
			return err
		}
	}

	if err := formService.formRepository.CreateCancer(formService.db, info); err != nil {
		return err
	}

	// Log operator action if performed by operator
	if isOp {
		formService.logOperatorFormUpdate(request.UserID, request.FormID, "اطلاعات سرطان")
	}

	return nil
}

func (formService *FormService) UpdateCancer(request formdto.UpdateCancerRequest) (formdto.UpdateCancerResponse, error) {
	form, err := formService.formRepository.FindFormByID(formService.db, request.FormID)
	if err != nil {
		return formdto.UpdateCancerResponse{}, err
	}
	if form == nil {
		notFoundError := exception.NotFoundError{Item: formService.constants.Field.Form}
		return formdto.UpdateCancerResponse{}, notFoundError
	}

	isOp, err := formService.isOperator(request.UserID)
	if err != nil {
		return formdto.UpdateCancerResponse{}, err
	}
	if isOp {
		canEdit, err := formService.canOperatorEditForm(form, request.UserID)
		if err != nil {
			return formdto.UpdateCancerResponse{}, err
		}
		if !canEdit {
			forbiddenError := exception.ForbiddenError{Resource: formService.constants.Field.Form}
			return formdto.UpdateCancerResponse{}, forbiddenError
		}
	}

	cancer, err := formService.formRepository.FindCancerByID(formService.db, request.CancerID)
	if err != nil {
		return formdto.UpdateCancerResponse{}, err
	}
	if cancer == nil {
		notFoundError := exception.NotFoundError{Item: "cancer"}
		return formdto.UpdateCancerResponse{}, notFoundError
	}

	// Verify the cancer belongs to the form
	if cancer.FormID != request.FormID {
		forbiddenError := exception.ForbiddenError{Resource: formService.constants.Field.Form}
		return formdto.UpdateCancerResponse{}, forbiddenError
	}

	// Check for duplicate cancer (same type and age) excluding the current cancer being updated
	if cancer.CancerType != enum.CancerType(request.CancerType) || cancer.CancerAge != request.CancerAge {
		existingCancers, err := formService.formRepository.FindCancersByFormID(formService.db, request.FormID)
		if err != nil {
			return formdto.UpdateCancerResponse{}, err
		}
		for _, existing := range existingCancers {
			if existing.ID != request.CancerID &&
				existing.CancerType == enum.CancerType(request.CancerType) &&
				existing.CancerAge == request.CancerAge {
				conflictError := exception.ConflictErrors{}
				conflictError.Add(formService.constants.Field.Cancer, formService.constants.Tag.AlreadyExist)
				return formdto.UpdateCancerResponse{}, conflictError
			}
		}
	}

	// Save old picture path for cleanup
	var oldPicturePath *string
	if cancer.PicturePath != nil {
		oldPicturePath = cancer.PicturePath
	}

	cancer.CancerType = enum.CancerType(request.CancerType)
	cancer.CancerAge = request.CancerAge

	// Handle picture upload
	if request.Picture != nil && request.Picture.Filename != "" {
		pictureKey := formService.constants.BucketPath.GetCancerPath(
			request.FormID,
			enum.CancerType(request.CancerType),
			request.Picture.Filename,
		)
		cancer.PicturePath = &pictureKey

		// Upload new picture
		if err := formService.s3Storage.UploadObject(enum.BucketTypeCancer, pictureKey, request.Picture); err != nil {
			return formdto.UpdateCancerResponse{}, err
		}

		// Delete old picture if it exists and is different
		if oldPicturePath != nil && *oldPicturePath != "" && *oldPicturePath != pictureKey {
			if err := formService.s3Storage.DeleteObject(enum.BucketTypeCancer, *oldPicturePath); err != nil {
				fmt.Printf("Warning: failed to delete old cancer picture %s: %v\n", *oldPicturePath, err)
			}
		}
	}

	if err := formService.formRepository.UpdateCancer(formService.db, cancer); err != nil {
		return formdto.UpdateCancerResponse{}, err
	}

	var pictureURL *string
	if cancer.PicturePath != nil && *cancer.PicturePath != "" {
		presignedURL, err := formService.s3Storage.GetPresignedURL(enum.BucketTypeCancer, *cancer.PicturePath, 8*time.Hour)
		if err != nil {
			return formdto.UpdateCancerResponse{}, err
		}
		pictureURL = &presignedURL
	}

	// Log operator action if performed by operator
	if isOp {
		formService.logOperatorFormUpdate(request.UserID, request.FormID, "اطلاعات سرطان")
	}

	response := formdto.UpdateCancerResponse{
		Cancer: formdto.CancerResponse{
			ID:         cancer.ID,
			CancerType: cancer.CancerType,
			CancerAge:  cancer.CancerAge,
			Picture:    pictureURL,
		},
	}

	return response, nil
}

func (formService *FormService) DeleteCancer(request formdto.DeleteCancerRequest) error {
	form, err := formService.formRepository.FindFormByID(formService.db, request.FormID)
	if err != nil {
		return err
	}
	if form == nil {
		notFoundError := exception.NotFoundError{Item: formService.constants.Field.Form}
		return notFoundError
	}

	isOp, err := formService.isOperator(request.UserID)
	if err != nil {
		return err
	}
	if isOp {
		canEdit, err := formService.canOperatorEditForm(form, request.UserID)
		if err != nil {
			return err
		}
		if !canEdit {
			forbiddenError := exception.ForbiddenError{Resource: formService.constants.Field.Form}
			return forbiddenError
		}
	}

	cancer, err := formService.formRepository.FindCancerByID(formService.db, request.CancerID)
	if err != nil {
		return err
	}
	if cancer == nil {
		notFoundError := exception.NotFoundError{Item: "cancer"}
		return notFoundError
	}

	// Verify the cancer belongs to the form
	if cancer.FormID != request.FormID {
		forbiddenError := exception.ForbiddenError{Resource: formService.constants.Field.Form}
		return forbiddenError
	}

	// Delete picture if it exists
	if cancer.PicturePath != nil && *cancer.PicturePath != "" {
		if err := formService.s3Storage.DeleteObject(enum.BucketTypeCancer, *cancer.PicturePath); err != nil {
			fmt.Printf("Warning: failed to delete cancer picture %s: %v\n", *cancer.PicturePath, err)
		}
	}

	// Delete cancer record
	if err := formService.formRepository.DeleteCancerByID(formService.db, request.CancerID); err != nil {
		return err
	}

	// Log operator action if performed by operator
	if isOp {
		formService.logOperatorFormUpdate(request.UserID, request.FormID, "اطلاعات سرطان")
	}

	return nil
}

func (formService *FormService) CreateFamilyCancer(request formdto.CreateFamilyCancerRequest) (formdto.CreateFamilyCancerResponse, error) {
	form, err := formService.formRepository.FindFormByID(formService.db, request.FormID)
	if err != nil {
		return formdto.CreateFamilyCancerResponse{}, err
	}
	if form == nil {
		notFoundError := exception.NotFoundError{Item: formService.constants.Field.Form}
		return formdto.CreateFamilyCancerResponse{}, notFoundError
	}

	isOp, err := formService.isOperator(request.UserID)
	if err != nil {
		return formdto.CreateFamilyCancerResponse{}, err
	}
	if isOp {
		canEdit, err := formService.canOperatorEditForm(form, request.UserID)
		if err != nil {
			return formdto.CreateFamilyCancerResponse{}, err
		}
		if !canEdit {
			forbiddenError := exception.ForbiddenError{Resource: formService.constants.Field.Form}
			return formdto.CreateFamilyCancerResponse{}, forbiddenError
		}
	}

	info := &entity.FamilyCancerInfo{
		FormID:           request.FormID,
		Relative:         request.Relative,
		RelativeRelation: request.RelativeRelation,
		Name:             request.Name,
		LifeStatus:       request.LifeStatus,
		CancerAge:        request.CancerAge,
		CancerType:       enum.CancerType(request.CancerType),
	}

	// Handle picture upload
	if request.Picture != nil && request.Picture.Filename != "" {
		pictureKey := formService.constants.BucketPath.GetFamilyCancerPath(
			request.FormID,
			enum.CancerType(request.CancerType),
			request.Picture.Filename,
		)
		info.PicturePath = &pictureKey

		// Upload picture
		if err := formService.s3Storage.UploadObject(enum.BucketTypeCancer, pictureKey, request.Picture); err != nil {
			return formdto.CreateFamilyCancerResponse{}, err
		}
	}

	if err := formService.formRepository.CreateFamilyCancer(formService.db, info); err != nil {
		return formdto.CreateFamilyCancerResponse{}, err
	}

	var pictureURL *string
	if info.PicturePath != nil && *info.PicturePath != "" {
		presignedURL, err := formService.s3Storage.GetPresignedURL(enum.BucketTypeCancer, *info.PicturePath, 8*time.Hour)
		if err != nil {
			return formdto.CreateFamilyCancerResponse{}, err
		}
		pictureURL = &presignedURL
	}

	// Log operator action if performed by operator
	if isOp {
		formService.logOperatorFormUpdate(request.UserID, request.FormID, "اطلاعات سرطان خانوادگی")
	}

	response := formdto.CreateFamilyCancerResponse{
		FamilyCancer: formdto.FamilyCancerItemResponse{
			ID:               info.ID,
			Relative:         info.Relative,
			RelativeRelation: info.RelativeRelation,
			Name:             info.Name,
			LifeStatus:       info.LifeStatus,
			CancerType:       info.CancerType,
			CancerAge:        info.CancerAge,
			Picture:          pictureURL,
		},
	}

	return response, nil
}

func (formService *FormService) UpdateFamilyCancer(request formdto.UpdateFamilyCancerRequest) (formdto.UpdateFamilyCancerResponse, error) {
	form, err := formService.formRepository.FindFormByID(formService.db, request.FormID)
	if err != nil {
		return formdto.UpdateFamilyCancerResponse{}, err
	}
	if form == nil {
		notFoundError := exception.NotFoundError{Item: formService.constants.Field.Form}
		return formdto.UpdateFamilyCancerResponse{}, notFoundError
	}

	isOp, err := formService.isOperator(request.UserID)
	if err != nil {
		return formdto.UpdateFamilyCancerResponse{}, err
	}
	if isOp {
		canEdit, err := formService.canOperatorEditForm(form, request.UserID)
		if err != nil {
			return formdto.UpdateFamilyCancerResponse{}, err
		}
		if !canEdit {
			forbiddenError := exception.ForbiddenError{Resource: formService.constants.Field.Form}
			return formdto.UpdateFamilyCancerResponse{}, forbiddenError
		}
	}

	familyCancer, err := formService.formRepository.FindFamilyCancerByID(formService.db, request.FamilyCancerID)
	if err != nil {
		return formdto.UpdateFamilyCancerResponse{}, err
	}
	if familyCancer == nil {
		notFoundError := exception.NotFoundError{Item: "family cancer"}
		return formdto.UpdateFamilyCancerResponse{}, notFoundError
	}

	// Verify the family cancer belongs to the form
	if familyCancer.FormID != request.FormID {
		forbiddenError := exception.ForbiddenError{Resource: formService.constants.Field.Form}
		return formdto.UpdateFamilyCancerResponse{}, forbiddenError
	}

	// Save old picture path for cleanup
	var oldPicturePath *string
	if familyCancer.PicturePath != nil {
		oldPicturePath = familyCancer.PicturePath
	}

	familyCancer.Relative = request.Relative
	familyCancer.RelativeRelation = request.RelativeRelation
	familyCancer.Name = request.Name
	familyCancer.LifeStatus = request.LifeStatus
	familyCancer.CancerType = enum.CancerType(request.CancerType)
	familyCancer.CancerAge = request.CancerAge

	// Handle picture upload
	if request.Picture != nil && request.Picture.Filename != "" {
		pictureKey := formService.constants.BucketPath.GetFamilyCancerPath(
			request.FormID,
			enum.CancerType(request.CancerType),
			request.Picture.Filename,
		)
		familyCancer.PicturePath = &pictureKey

		// Upload new picture
		if err := formService.s3Storage.UploadObject(enum.BucketTypeCancer, pictureKey, request.Picture); err != nil {
			return formdto.UpdateFamilyCancerResponse{}, err
		}

		// Delete old picture if it exists and is different
		if oldPicturePath != nil && *oldPicturePath != "" && *oldPicturePath != pictureKey {
			if err := formService.s3Storage.DeleteObject(enum.BucketTypeCancer, *oldPicturePath); err != nil {
				fmt.Printf("Warning: failed to delete old family cancer picture %s: %v\n", *oldPicturePath, err)
			}
		}
	}

	if err := formService.formRepository.UpdateFamilyCancer(formService.db, familyCancer); err != nil {
		return formdto.UpdateFamilyCancerResponse{}, err
	}

	var pictureURL *string
	if familyCancer.PicturePath != nil && *familyCancer.PicturePath != "" {
		presignedURL, err := formService.s3Storage.GetPresignedURL(enum.BucketTypeCancer, *familyCancer.PicturePath, 8*time.Hour)
		if err != nil {
			return formdto.UpdateFamilyCancerResponse{}, err
		}
		pictureURL = &presignedURL
	}

	// Log operator action if performed by operator
	if isOp {
		formService.logOperatorFormUpdate(request.UserID, request.FormID, "اطلاعات سرطان خانوادگی")
	}

	response := formdto.UpdateFamilyCancerResponse{
		FamilyCancer: formdto.FamilyCancerItemResponse{
			ID:               familyCancer.ID,
			Relative:         familyCancer.Relative,
			RelativeRelation: familyCancer.RelativeRelation,
			Name:             familyCancer.Name,
			LifeStatus:       familyCancer.LifeStatus,
			CancerType:       familyCancer.CancerType,
			CancerAge:        familyCancer.CancerAge,
			Picture:          pictureURL,
		},
	}

	return response, nil
}

func (formService *FormService) DeleteFamilyCancer(request formdto.DeleteFamilyCancerRequest) error {
	form, err := formService.formRepository.FindFormByID(formService.db, request.FormID)
	if err != nil {
		return err
	}
	if form == nil {
		notFoundError := exception.NotFoundError{Item: formService.constants.Field.Form}
		return notFoundError
	}

	isOp, err := formService.isOperator(request.UserID)
	if err != nil {
		return err
	}
	if isOp {
		canEdit, err := formService.canOperatorEditForm(form, request.UserID)
		if err != nil {
			return err
		}
		if !canEdit {
			forbiddenError := exception.ForbiddenError{Resource: formService.constants.Field.Form}
			return forbiddenError
		}
	}

	familyCancer, err := formService.formRepository.FindFamilyCancerByID(formService.db, request.FamilyCancerID)
	if err != nil {
		return err
	}
	if familyCancer == nil {
		notFoundError := exception.NotFoundError{Item: "family cancer"}
		return notFoundError
	}

	// Verify the family cancer belongs to the form
	if familyCancer.FormID != request.FormID {
		forbiddenError := exception.ForbiddenError{Resource: formService.constants.Field.Form}
		return forbiddenError
	}

	// Delete picture if it exists
	if familyCancer.PicturePath != nil && *familyCancer.PicturePath != "" {
		if err := formService.s3Storage.DeleteObject(enum.BucketTypeCancer, *familyCancer.PicturePath); err != nil {
			fmt.Printf("Warning: failed to delete family cancer picture %s: %v\n", *familyCancer.PicturePath, err)
		}
	}

	// Delete family cancer record
	if err := formService.formRepository.DeleteFamilyCancerByID(formService.db, request.FamilyCancerID); err != nil {
		return err
	}

	// Log operator action if performed by operator
	if isOp {
		formService.logOperatorFormUpdate(request.UserID, request.FormID, "اطلاعات سرطان خانوادگی")
	}

	return nil
}

func (formService *FormService) UpsertContact(request formdto.UpsertContactRequest) error {
	form, err := formService.formRepository.FindFormByID(formService.db, request.FormID)
	if err != nil {
		return err
	}
	if form == nil {
		notFoundError := exception.NotFoundError{Item: formService.constants.Field.Form}
		return notFoundError
	}

	isOp, err := formService.isOperator(request.UserID)
	if err != nil {
		return err
	}
	if isOp {
		canEdit, err := formService.canOperatorEditForm(form, request.UserID)
		if err != nil {
			return err
		}
		if !canEdit {
			forbiddenError := exception.ForbiddenError{Resource: formService.constants.Field.Form}
			return forbiddenError
		}
	}

	// if form.UserID != request.UserID {
	// 	ForbiddenError := exception.ForbiddenError{Message: formService.constants.Field.Form}
	// 	return ForbiddenError
	// }

	info, err := formService.formRepository.FindContactByFormID(formService.db, request.FormID)
	if err != nil {
		return err
	}
	if info == nil {
		info = &entity.ContactInfo{FormID: request.FormID}
	}

	// Save old picture paths for cleanup
	var oldTestGenPicturePath *string
	var oldFatherTestGenPicturePath *string
	var oldMotherTestGenPicturePath *string
	if info.TestGenPicturePath != nil {
		oldTestGenPicturePath = info.TestGenPicturePath
	}
	if info.FatherTestGenPicturePath != nil {
		oldFatherTestGenPicturePath = info.FatherTestGenPicturePath
	}
	if info.MotherTestGenPicturePath != nil {
		oldMotherTestGenPicturePath = info.MotherTestGenPicturePath
	}

	info.Name = request.Name
	info.TestGen = request.TestGen
	info.FmTestGen = request.FmTestGen
	info.CallExpert = request.CallExpert
	info.BirthCountry = request.BirthCountry
	info.Province = request.Province
	info.City = request.City
	info.Country = request.Country
	info.Address = request.Address
	info.PostalCode = request.PostalCode

	// Handle TestGen picture upload
	if request.TestGenPicture != nil && request.TestGenPicture.Filename != "" {
		pictureKey := formService.constants.BucketPath.GetGeneticTestPath(
			request.FormID,
			"user",
			request.TestGenPicture.Filename,
		)
		info.TestGenPicturePath = &pictureKey

		// Upload picture
		if err := formService.s3Storage.UploadObject(enum.BucketTypeGeneticTest, pictureKey, request.TestGenPicture); err != nil {
			return err
		}

		// Delete old picture if it exists and is different
		if oldTestGenPicturePath != nil && *oldTestGenPicturePath != "" && *oldTestGenPicturePath != pictureKey {
			if err := formService.s3Storage.DeleteObject(enum.BucketTypeGeneticTest, *oldTestGenPicturePath); err != nil {
				fmt.Printf("Warning: failed to delete old genetic test picture %s: %v\n", *oldTestGenPicturePath, err)
			}
		}
	}

	// Handle FatherTestGen picture upload
	if request.FatherTestGenPicture != nil && request.FatherTestGenPicture.Filename != "" {
		pictureKey := formService.constants.BucketPath.GetFatherGeneticTestPath(
			request.FormID,
			request.FatherTestGenPicture.Filename,
		)
		info.FatherTestGenPicturePath = &pictureKey

		// Upload picture
		if err := formService.s3Storage.UploadObject(enum.BucketTypeGeneticTest, pictureKey, request.FatherTestGenPicture); err != nil {
			return err
		}

		// Delete old picture if it exists and is different
		if oldFatherTestGenPicturePath != nil && *oldFatherTestGenPicturePath != "" && *oldFatherTestGenPicturePath != pictureKey {
			if err := formService.s3Storage.DeleteObject(enum.BucketTypeGeneticTest, *oldFatherTestGenPicturePath); err != nil {
				fmt.Printf("Warning: failed to delete old father genetic test picture %s: %v\n", *oldFatherTestGenPicturePath, err)
			}
		}
	}

	// Handle MotherTestGen picture upload
	if request.MotherTestGenPicture != nil && request.MotherTestGenPicture.Filename != "" {
		pictureKey := formService.constants.BucketPath.GetMotherGeneticTestPath(
			request.FormID,
			request.MotherTestGenPicture.Filename,
		)
		info.MotherTestGenPicturePath = &pictureKey

		// Upload picture
		if err := formService.s3Storage.UploadObject(enum.BucketTypeGeneticTest, pictureKey, request.MotherTestGenPicture); err != nil {
			return err
		}

		// Delete old picture if it exists and is different
		if oldMotherTestGenPicturePath != nil && *oldMotherTestGenPicturePath != "" && *oldMotherTestGenPicturePath != pictureKey {
			if err := formService.s3Storage.DeleteObject(enum.BucketTypeGeneticTest, *oldMotherTestGenPicturePath); err != nil {
				fmt.Printf("Warning: failed to delete old mother genetic test picture %s: %v\n", *oldMotherTestGenPicturePath, err)
			}
		}
	}

	if info.ID == 0 {
		err = formService.formRepository.CreateContact(formService.db, info)
	} else {
		err = formService.formRepository.UpdateContact(formService.db, info)
	}
	if err != nil {
		return err
	}

	// Log operator action if performed by operator
	if isOp {
		formService.logOperatorFormUpdate(request.UserID, request.FormID, "اطلاعات تماس")
	}

	return nil
}

func (formService *FormService) UpsertLungCancer(request formdto.UpsertLungCancerRequest) error {
	form, err := formService.formRepository.FindFormByID(formService.db, request.FormID)
	if err != nil {
		return err
	}
	if form == nil {
		notFoundError := exception.NotFoundError{Item: formService.constants.Field.Form}
		return notFoundError
	}

	isOp, err := formService.isOperator(request.UserID)
	if err != nil {
		return err
	}
	if isOp {
		canEdit, err := formService.canOperatorEditForm(form, request.UserID)
		if err != nil {
			return err
		}
		if !canEdit {
			forbiddenError := exception.ForbiddenError{Resource: formService.constants.Field.Form}
			return forbiddenError
		}
	}

	// if form.UserID != request.UserID {
	// 	ForbiddenError := exception.ForbiddenError{Message: formService.constants.Field.Form}
	// 	return ForbiddenError
	// }

	info, err := formService.formRepository.FindLungCancerByFormID(formService.db, request.FormID)
	if err != nil {
		return err
	}
	if info == nil {
		info = &entity.LungCancerInfo{FormID: request.FormID}
	}

	info.InsuranceStatus = request.InsuranceStatus
	info.SupplementaryInsurances = request.SupplementaryInsurances
	info.Hypertension = request.Hypertension
	info.HypertensionTreatment = request.HypertensionTreatment
	info.HeartDisease = request.HeartDisease
	info.HeartDiseaseTreatment = request.HeartDiseaseTreatment
	info.Diabetes = request.Diabetes
	info.DiabetesTreatment = request.DiabetesTreatment
	info.ChronicLungDisease = request.ChronicLungDisease
	info.ChronicLungDiseaseType = request.ChronicLungDiseaseType
	info.LungCancerHistory = request.LungCancerHistory
	info.OtherCancerHistory = request.OtherCancerHistory
	if request.OtherCancerType != nil {
		info.OtherCancerType = (*enum.CancerType)(request.OtherCancerType)
	}
	info.LungCancerFamily = request.LungCancerFamily
	info.LungCancerFamilyRelation = request.LungCancerFamilyRelation
	info.OtherCancerFamily = request.OtherCancerFamily
	if request.OtherCancerFamilyType != nil {
		info.OtherCancerFamilyType = (*enum.CancerType)(request.OtherCancerFamilyType)
	}
	info.OtherCancerFamilyRelation = request.OtherCancerFamilyRelation
	info.OccupationalExposure = request.OccupationalExposure
	info.CurrentSmoking = request.CurrentSmoking
	info.SmokingStartAgeCurrent = request.SmokingStartAgeCurrent
	info.SmokingTypesCurrent = request.SmokingTypesCurrent
	info.CigarettesPerDayCurrent = request.CigarettesPerDayCurrent
	info.CigarPerDayCurrent = request.CigarPerDayCurrent
	info.ECigPerDayCurrent = request.ECigPerDayCurrent
	info.PipePerDayCurrent = request.PipePerDayCurrent
	info.ChapoghPerDayCurrent = request.ChapoghPerDayCurrent
	info.SmokedOpiumPerDayCurrent = request.SmokedOpiumPerDayCurrent
	info.ChewedOpiumPerDayCurrent = request.ChewedOpiumPerDayCurrent
	info.HookahPerWeekCurrent = request.HookahPerWeekCurrent
	info.PastSmoking = request.PastSmoking
	info.SmokingStartAgePast = request.SmokingStartAgePast
	info.SmokingTypesPast = request.SmokingTypesPast
	info.CigarettesPerDayPast = request.CigarettesPerDayPast
	info.CigarPerDayPast = request.CigarPerDayPast
	info.ECigPerDayPast = request.ECigPerDayPast
	info.PipePerDayPast = request.PipePerDayPast
	info.ChapoghPerDayPast = request.ChapoghPerDayPast
	info.SmokedOpiumPerDayPast = request.SmokedOpiumPerDayPast
	info.ChewedOpiumPerDayPast = request.ChewedOpiumPerDayPast
	info.HookahPerWeekPast = request.HookahPerWeekPast
	info.SecondhandSmoke = request.SecondhandSmoke
	info.SecondhandSmokeLocation = request.SecondhandSmokeLocation

	if info.ID == 0 {
		err = formService.formRepository.CreateLungCancer(formService.db, info)
	} else {
		err = formService.formRepository.UpdateLungCancer(formService.db, info)
	}
	if err != nil {
		return err
	}

	// Log operator action if performed by operator
	if isOp {
		formService.logOperatorFormUpdate(request.UserID, request.FormID, "اطلاعات سرطان ریه")
	}

	return nil
}

func (formService *FormService) ChangeFormStatus(request formdto.ChangeFormStatusRequest) (formdto.ChangeFormStatusResponse, error) {
	form, err := formService.formRepository.FindFormByID(formService.db, request.FormID)
	if err != nil {
		return formdto.ChangeFormStatusResponse{}, err
	}
	if form == nil {
		notFoundError := exception.NotFoundError{Item: formService.constants.Field.Form}
		return formdto.ChangeFormStatusResponse{}, notFoundError
	}

	// if form.UserID != request.UserID {
	// 	ForbiddenError := exception.ForbiddenError{Message: formService.constants.Field.Form}
	// 	return formdto.ChangeFormStatusResponse{}, ForbiddenError
	// }

	form.Status = enum.FormStatusReady
	err = formService.formRepository.UpdateForm(formService.db, form)
	if err != nil {
		return formdto.ChangeFormStatusResponse{}, err
	}

	response := formdto.ChangeFormStatusResponse{
		Form: formdto.BasicFormResponse{
			FormID:     form.ID,
			Status:     form.Status.String(),
			OperatorID: form.OperatorID,
			UserID:     form.UserID,
			CreatedAt:  form.CreatedAt,
			UpdatedAt:  form.UpdatedAt,
		},
	}

	return response, nil
}

func (formService *FormService) GetBasicForm(request formdto.GetPartialFormRequest) (formdto.GetBasicFormResponse, error) {
	form, err := formService.formRepository.FindFormByID(formService.db, request.FormID)
	if err != nil {
		return formdto.GetBasicFormResponse{}, err
	}
	if form == nil {
		notFoundError := exception.NotFoundError{Item: formService.constants.Field.Form}
		return formdto.GetBasicFormResponse{}, notFoundError
	}

	// if form.UserID != request.UserID {
	// 	ForbiddenError := exception.ForbiddenError{Message: formService.constants.Field.Form}
	// 	return formdto.GetBasicFormResponse{}, ForbiddenError
	// }

	basic, err := formService.formRepository.FindBasicInfoByFormID(formService.db, request.FormID)
	if err != nil {
		return formdto.GetBasicFormResponse{}, err
	}
	if basic == nil {
		notFoundError := exception.NotFoundError{Item: formService.constants.Field.Form}
		return formdto.GetBasicFormResponse{}, notFoundError
	}

	return formdto.GetBasicFormResponse{
		ID:                   basic.ID,
		Gender:               basic.Gender,
		BirthDate:            basic.BirthDate,
		IsAtba:               basic.IsAtba,
		SocialSecurityNumber: basic.SocialSecurityNumber,
		Height:               basic.Height,
		Weight:               basic.Weight,
	}, nil
}

func (formService *FormService) GetGeneralHealth(request formdto.GetPartialFormRequest) (formdto.GetGeneralHealthResponse, error) {
	form, err := formService.formRepository.FindFormByID(formService.db, request.FormID)
	if err != nil {
		return formdto.GetGeneralHealthResponse{}, err
	}
	if form == nil {
		notFoundError := exception.NotFoundError{Item: formService.constants.Field.Form}
		return formdto.GetGeneralHealthResponse{}, notFoundError
	}

	// if form.UserID != request.UserID {
	// 	ForbiddenError := exception.ForbiddenError{Message: formService.constants.Field.Form}
	// 	return formdto.GetGeneralHealthResponse{}, ForbiddenError
	// }

	info, err := formService.formRepository.FindGeneralHealthByFormID(formService.db, request.FormID)
	if err != nil {
		return formdto.GetGeneralHealthResponse{}, err
	}
	if info == nil {
		notFoundError := exception.NotFoundError{Item: formService.constants.Field.Form}
		return formdto.GetGeneralHealthResponse{}, notFoundError
	}

	return formdto.GetGeneralHealthResponse{
		ID:                        info.ID,
		DrinksAlcohol:             info.DrinksAlcohol,
		CupsPerWeek:               info.CupsPerWeek,
		LastMonthSabzijatMeal:     info.LastMonthSabzijatMeal,
		LastMonthSabzijatWeight:   info.LastMonthSabzijatWeight,
		MediumActivityMonthInYear: info.MediumActivityMonthInYear,
		MediumActivityHourInWeek:  info.MediumActivityHourInWeek,
		HardActivityMonthInYear:   info.HardActivityMonthInYear,
		HardActivityHourInWeek:    info.HardActivityHourInWeek,
		SmokeAtLeast100:           info.SmokeAtLeast100,
		SmokingAge:                info.SmokingAge,
		SmokingNow:                info.SmokingNow,
		LeaveSmokingAge:           info.LeaveSmokingAge,
		CountSmokingDaily:         info.CountSmokingDaily,
		CountGheliandaily:         info.CountGheliandaily,
		CountSmokingDailyPast:     info.CountSmokingDailyPast,
		CountGheliandailyPast:     info.CountGheliandailyPast,
	}, nil
}

func (formService *FormService) GetMamography(request formdto.GetPartialFormRequest) (formdto.GetMamographyResponse, error) {
	form, err := formService.formRepository.FindFormByID(formService.db, request.FormID)
	if err != nil {
		return formdto.GetMamographyResponse{}, err
	}
	if form == nil {
		notFoundError := exception.NotFoundError{Item: formService.constants.Field.Form}
		return formdto.GetMamographyResponse{}, notFoundError
	}

	// if form.UserID != request.UserID {
	// 	ForbiddenError := exception.ForbiddenError{Message: formService.constants.Field.Form}
	// 	return formdto.GetMamographyResponse{}, ForbiddenError
	// }

	info, err := formService.formRepository.FindMamographyByFormID(formService.db, request.FormID)
	if err != nil {
		return formdto.GetMamographyResponse{}, err
	}
	if info == nil {
		notFoundError := exception.NotFoundError{Item: formService.constants.Field.Form}
		return formdto.GetMamographyResponse{}, notFoundError
	}

	var mamoGraphyPicture *string
	if info.MamoGraphyPicturePath != nil && *info.MamoGraphyPicturePath != "" {
		presignedURL, err := formService.s3Storage.GetPresignedURL(enum.BucketTypeMamography, *info.MamoGraphyPicturePath, 8*time.Hour)
		if err != nil {
			return formdto.GetMamographyResponse{}, err
		}
		mamoGraphyPicture = &presignedURL
	}

	return formdto.GetMamographyResponse{
		ID:                           info.ID,
		GhaedeAge:                    info.GhaedeAge,
		HasChildren:                  info.HasChildren,
		NumberOfChildren:             info.NumberOfChildren,
		AgeOfFirstBirth:              info.AgeOfFirstBirth,
		MenopausalStatus:             info.MenopausalStatus,
		MenopauseAge:                 info.MenopauseAge,
		HRT:                          info.HRT,
		HRTUseLength:                 info.HRTUseLength,
		LastFiveYearsHRTUse:          info.LastFiveYearsHRTUse,
		CurrentHRTUse:                info.CurrentHRTUse,
		IntendedHRTUse:               info.IntendedHRTUse,
		HRTType:                      info.HRTType,
		Oral:                         info.Oral,
		OralDuration:                 info.OralDuration,
		OralTwoLastYears:             info.OralTwoLastYears,
		MamoGraphy:                   info.MamoGraphy,
		MamoGraphyPicture:            mamoGraphyPicture,
		Falop:                        info.Falop,
		Andometrioz:                  info.Andometrioz,
		LeavePestan:                  info.LeavePestan,
		LeaveTokhmdan:                info.LeaveTokhmdan,
		LaDeColon:                    info.LaDeColon,
		LaDePol:                      info.LaDePol,
		AspLaMo:                      info.AspLaMo,
		NsaiDLaMo:                    info.NsaiDLaMo,
		LastFiveYearBloodTestInStool: info.LastFiveYearBloodTestInStool,
	}, nil
}

func (formService *FormService) GetCancers(request formdto.GetPartialFormRequest) (formdto.GetCancersResponse, error) {
	form, err := formService.formRepository.FindFormByID(formService.db, request.FormID)
	if err != nil {
		return formdto.GetCancersResponse{}, err
	}
	if form == nil {
		notFoundError := exception.NotFoundError{Item: formService.constants.Field.Form}
		return formdto.GetCancersResponse{}, notFoundError
	}

	// if form.UserID != request.UserID {
	// 	ForbiddenError := exception.ForbiddenError{Message: formService.constants.Field.Form}
	// 	return formdto.GetCancerResponse{}, ForbiddenError
	// }

	info, err := formService.formRepository.FindCancersByFormID(formService.db, request.FormID)
	if err != nil {
		return formdto.GetCancersResponse{}, err
	}

	cancersResponse := formdto.GetCancersResponse{}
	if len(info) == 0 {
		cancersResponse.Cancer = false
		cancersResponse.Cancers = make([]formdto.CancerResponse, len(info))
		return cancersResponse, nil
	}

	cancersResponse.Cancer = true
	for _, v := range info {
		var pictureURL *string
		if v.PicturePath != nil && *v.PicturePath != "" {
			presignedURL, err := formService.s3Storage.GetPresignedURL(enum.BucketTypeCancer, *v.PicturePath, 8*time.Hour)
			if err != nil {
				return formdto.GetCancersResponse{}, err
			}
			pictureURL = &presignedURL
		}
		cancersResponse.Cancers = append(cancersResponse.Cancers, formdto.CancerResponse{
			ID:         v.ID,
			CancerType: v.CancerType,
			CancerAge:  v.CancerAge,
			Picture:    pictureURL,
		})
	}
	return cancersResponse, nil
}

func (formService *FormService) GetFamilyCancer(request formdto.GetPartialFormRequest) (formdto.GetFamilyCancerResponse, error) {
	form, err := formService.formRepository.FindFormByID(formService.db, request.FormID)
	if err != nil {
		return formdto.GetFamilyCancerResponse{}, err
	}
	if form == nil {
		notFoundError := exception.NotFoundError{Item: formService.constants.Field.Form}
		return formdto.GetFamilyCancerResponse{}, notFoundError
	}

	// if form.UserID != request.UserID {
	// 	ForbiddenError := exception.ForbiddenError{Message: formService.constants.Field.Form}
	// 	return formdto.GetFamilyCancerResponse{}, ForbiddenError
	// }

	info, err := formService.formRepository.FindFamilyCancersByFormID(formService.db, request.FormID)
	if err != nil {
		return formdto.GetFamilyCancerResponse{}, err
	}
	if info == nil {
		notFoundError := exception.NotFoundError{Item: formService.constants.Field.Form}
		return formdto.GetFamilyCancerResponse{}, notFoundError
	}

	FamilyCancersResponse := formdto.GetFamilyCancerResponse{}
	if len(info) == 0 {
		FamilyCancersResponse.FamilyCancers = make([]formdto.FamilyCancerResponse, 0)
		return FamilyCancersResponse, nil
	}

	for i, v := range info {
		if i == 0 || FamilyCancersResponse.FamilyCancers[len(FamilyCancersResponse.FamilyCancers)-1].Relative != v.Relative ||
			!((v.Name != nil && FamilyCancersResponse.FamilyCancers[len(FamilyCancersResponse.FamilyCancers)-1].Name != nil && *v.Name == *FamilyCancersResponse.FamilyCancers[len(FamilyCancersResponse.FamilyCancers)-1].Name) || (v.Name == nil && FamilyCancersResponse.FamilyCancers[len(FamilyCancersResponse.FamilyCancers)-1].Name == nil)) ||
			!((v.RelativeRelation != nil && FamilyCancersResponse.FamilyCancers[len(FamilyCancersResponse.FamilyCancers)-1].RelativeRelation != nil && *v.RelativeRelation == *FamilyCancersResponse.FamilyCancers[len(FamilyCancersResponse.FamilyCancers)-1].RelativeRelation) || (v.RelativeRelation == nil && FamilyCancersResponse.FamilyCancers[len(FamilyCancersResponse.FamilyCancers)-1].RelativeRelation == nil)) {
			familyInfo := formdto.FamilyCancerResponse{Relative: v.Relative, RelativeRelation: v.RelativeRelation, Name: v.Name, LifeStatus: v.LifeStatus}
			FamilyCancersResponse.FamilyCancers = append(FamilyCancersResponse.FamilyCancers, familyInfo)
		}
		var pictureURL *string
		if v.PicturePath != nil && *v.PicturePath != "" {
			presignedURL, err := formService.s3Storage.GetPresignedURL(enum.BucketTypeCancer, *v.PicturePath, 8*time.Hour)
			if err != nil {
				return formdto.GetFamilyCancerResponse{}, err
			}
			pictureURL = &presignedURL
		}
		FamilyCancersResponse.FamilyCancers[len(FamilyCancersResponse.FamilyCancers)-1].Cancers = append(FamilyCancersResponse.FamilyCancers[len(FamilyCancersResponse.FamilyCancers)-1].Cancers, formdto.CancerResponse{
			ID:         v.ID,
			CancerType: v.CancerType,
			CancerAge:  v.CancerAge,
			Picture:    pictureURL,
		})
	}

	return FamilyCancersResponse, nil
}

func (formService *FormService) GetContact(request formdto.GetPartialFormRequest) (formdto.GetContactResponse, error) {
	form, err := formService.formRepository.FindFormByID(formService.db, request.FormID)
	if err != nil {
		return formdto.GetContactResponse{}, err
	}
	if form == nil {
		notFoundError := exception.NotFoundError{Item: formService.constants.Field.Form}
		return formdto.GetContactResponse{}, notFoundError
	}

	// if form.UserID != request.UserID {
	// 	ForbiddenError := exception.ForbiddenError{Message: formService.constants.Field.Form}
	// 	return formdto.GetContactResponse{}, ForbiddenError
	// }

	info, err := formService.formRepository.FindContactByFormID(formService.db, request.FormID)
	if err != nil {
		return formdto.GetContactResponse{}, err
	}
	if info == nil {
		notFoundError := exception.NotFoundError{Item: formService.constants.Field.Form}
		return formdto.GetContactResponse{}, notFoundError
	}

	// Generate presigned URLs for pictures
	var testGenPictureURL *string
	if info.TestGenPicturePath != nil && *info.TestGenPicturePath != "" {
		presignedURL, err := formService.s3Storage.GetPresignedURL(enum.BucketTypeGeneticTest, *info.TestGenPicturePath, 8*time.Hour)
		if err != nil {
			return formdto.GetContactResponse{}, err
		}
		testGenPictureURL = &presignedURL
	}

	var fatherTestGenPictureURL *string
	if info.FatherTestGenPicturePath != nil && *info.FatherTestGenPicturePath != "" {
		presignedURL, err := formService.s3Storage.GetPresignedURL(enum.BucketTypeGeneticTest, *info.FatherTestGenPicturePath, 8*time.Hour)
		if err != nil {
			return formdto.GetContactResponse{}, err
		}
		fatherTestGenPictureURL = &presignedURL
	}

	var motherTestGenPictureURL *string
	if info.MotherTestGenPicturePath != nil && *info.MotherTestGenPicturePath != "" {
		presignedURL, err := formService.s3Storage.GetPresignedURL(enum.BucketTypeGeneticTest, *info.MotherTestGenPicturePath, 8*time.Hour)
		if err != nil {
			return formdto.GetContactResponse{}, err
		}
		motherTestGenPictureURL = &presignedURL
	}

	return formdto.GetContactResponse{
		ID:                   info.ID,
		Name:                 info.Name,
		TestGen:              info.TestGen,
		TestGenPicture:       testGenPictureURL,
		FmTestGen:            info.FmTestGen,
		FatherTestGenPicture: fatherTestGenPictureURL,
		MotherTestGenPicture: motherTestGenPictureURL,
		CallExpert:           info.CallExpert,
		BirthCountry:         info.BirthCountry,
		Province:             info.Province,
		City:                 info.City,
		Country:              info.Country,
		Address:              info.Address,
		PostalCode:           info.PostalCode,
	}, nil
}

func (formService *FormService) GetLungCancer(request formdto.GetPartialFormRequest) (formdto.GetLungCancerResponse, error) {
	form, err := formService.formRepository.FindFormByID(formService.db, request.FormID)
	if err != nil {
		return formdto.GetLungCancerResponse{}, err
	}
	if form == nil {
		notFoundError := exception.NotFoundError{Item: formService.constants.Field.Form}
		return formdto.GetLungCancerResponse{}, notFoundError
	}

	// if form.UserID != request.UserID {
	// 	ForbiddenError := exception.ForbiddenError{Message: formService.constants.Field.Form}
	// 	return formdto.GetLungCancerResponse{}, ForbiddenError
	// }

	info, err := formService.formRepository.FindLungCancerByFormID(formService.db, request.FormID)
	if err != nil {
		return formdto.GetLungCancerResponse{}, err
	}
	if info == nil {
		notFoundError := exception.NotFoundError{Item: formService.constants.Field.Form}
		return formdto.GetLungCancerResponse{}, notFoundError
	}

	return formdto.GetLungCancerResponse{
		ID:                        info.ID,
		InsuranceStatus:           info.InsuranceStatus,
		SupplementaryInsurances:   info.SupplementaryInsurances,
		Hypertension:              info.Hypertension,
		HypertensionTreatment:     info.HypertensionTreatment,
		HeartDisease:              info.HeartDisease,
		HeartDiseaseTreatment:     info.HeartDiseaseTreatment,
		Diabetes:                  info.Diabetes,
		DiabetesTreatment:         info.DiabetesTreatment,
		ChronicLungDisease:        info.ChronicLungDisease,
		ChronicLungDiseaseType:    info.ChronicLungDiseaseType,
		LungCancerHistory:         info.LungCancerHistory,
		OtherCancerHistory:        info.OtherCancerHistory,
		OtherCancerType:           info.OtherCancerType,
		LungCancerFamily:          info.LungCancerFamily,
		LungCancerFamilyRelation:  info.LungCancerFamilyRelation,
		OtherCancerFamily:         info.OtherCancerFamily,
		OtherCancerFamilyType:     info.OtherCancerFamilyType,
		OtherCancerFamilyRelation: info.OtherCancerFamilyRelation,
		OccupationalExposure:      info.OccupationalExposure,
		CurrentSmoking:            info.CurrentSmoking,
		SmokingStartAgeCurrent:    info.SmokingStartAgeCurrent,
		SmokingTypesCurrent:       info.SmokingTypesCurrent,
		CigarettesPerDayCurrent:   info.CigarettesPerDayCurrent,
		CigarPerDayCurrent:        info.CigarPerDayCurrent,
		ECigPerDayCurrent:         info.ECigPerDayCurrent,
		PipePerDayCurrent:         info.PipePerDayCurrent,
		ChapoghPerDayCurrent:      info.ChapoghPerDayCurrent,
		SmokedOpiumPerDayCurrent:  info.SmokedOpiumPerDayCurrent,
		ChewedOpiumPerDayCurrent:  info.ChewedOpiumPerDayCurrent,
		HookahPerWeekCurrent:      info.HookahPerWeekCurrent,
		PastSmoking:               info.PastSmoking,
		SmokingStartAgePast:       info.SmokingStartAgePast,
		SmokingTypesPast:          info.SmokingTypesPast,
		CigarettesPerDayPast:      info.CigarettesPerDayPast,
		CigarPerDayPast:           info.CigarPerDayPast,
		ECigPerDayPast:            info.ECigPerDayPast,
		PipePerDayPast:            info.PipePerDayPast,
		ChapoghPerDayPast:         info.ChapoghPerDayPast,
		SmokedOpiumPerDayPast:     info.SmokedOpiumPerDayPast,
		ChewedOpiumPerDayPast:     info.ChewedOpiumPerDayPast,
		HookahPerWeekPast:         info.HookahPerWeekPast,
		SecondhandSmoke:           info.SecondhandSmoke,
		SecondhandSmokeLocation:   info.SecondhandSmokeLocation,
	}, nil
}

func (formService *FormService) GetUserForms(request formdto.GetUserFormsRequest) ([]formdto.BasicFormResponse, int64, error) {
	user, err := formService.userService.GetUserByID(request.UserID)
	if err != nil {
		return nil, 0, err
	}
	if user == nil {
		notFoundError := exception.NotFoundError{Item: formService.constants.Field.User}
		return nil, 0, notFoundError
	}

	options := postgres.NewQueryOptions().
		WithPagination(request.Limit, request.Offset)

	forms, err := formService.formRepository.FindFormsByUserID(formService.db, request.UserID, options)
	if err != nil {
		return nil, 0, err
	}

	count, err := formService.formRepository.CountFormsByUserID(formService.db, request.UserID)
	if err != nil {
		return nil, 0, err
	}

	formResponses := make([]formdto.BasicFormResponse, len(forms))
	for i, form := range forms {
		formResponses[i] = formdto.BasicFormResponse{
			FormID:     form.ID,
			Status:     form.Status.String(),
			UserID:     form.UserID,
			OperatorID: form.OperatorID,
			CreatedAt:  form.CreatedAt,
			UpdatedAt:  form.UpdatedAt,
		}
	}

	return formResponses, count, nil
}

func (formService *FormService) UpdateForm(request formdto.UpdateBasicFormRequest) error {
	form, err := formService.formRepository.FindFormByID(formService.db, request.FormID)
	if err != nil {
		return err
	}
	if form == nil {
		notFoundError := exception.NotFoundError{Item: formService.constants.Field.Form}
		return notFoundError
	}

	err = formService.formRepository.UpdateForm(formService.db, form)
	if err != nil {
		return err
	}

	return nil
}

func (formService *FormService) UpdateBasicInfo(request formdto.UpdateBasicFormRequest) error {
	form, err := formService.formRepository.FindFormByID(formService.db, request.FormID)
	if err != nil {
		return err
	}
	if form == nil {
		notFoundError := exception.NotFoundError{Item: formService.constants.Field.Form}
		return notFoundError
	}

	isOp, err := formService.isOperator(request.UserID)
	if err != nil {
		return err
	}
	if isOp {
		canEdit, err := formService.canOperatorEditForm(form, request.UserID)
		if err != nil {
			return err
		}
		if !canEdit {
			forbiddenError := exception.ForbiddenError{Resource: formService.constants.Field.Form}
			return forbiddenError
		}
	}

	info, err := formService.formRepository.FindBasicInfoByFormID(formService.db, request.FormID)
	if err != nil {
		return err
	}
	if info == nil {
		notFoundError := exception.NotFoundError{Item: formService.constants.Field.Form}
		return notFoundError
	}

	if request.BirthDate != nil {
		info.BirthDate = *request.BirthDate
	}
	if request.SocialSecurityNumber != nil {
		info.SocialSecurityNumber = *request.SocialSecurityNumber
	}
	if request.Gender != nil {
		info.Gender = enum.Gender(uint(*request.Gender))
	}
	if request.IsAtba != nil {
		info.IsAtba = *request.IsAtba
	}
	if request.Height != nil {
		info.Height = *request.Height
	}
	if request.Weight != nil {
		info.Weight = *request.Weight
	}

	err = formService.formRepository.UpdateBasicInfo(formService.db, info)
	if err != nil {
		return err
	}

	// Log operator action if performed by operator
	if isOp {
		formService.logOperatorFormUpdate(request.UserID, request.FormID, "اطلاعات پایه")
	}

	return nil
}

func (formService *FormService) DeleteForm(request formdto.DeleteFormRequest) error {
	form, err := formService.formRepository.FindFormByID(formService.db, request.FormID)
	if err != nil {
		return err
	}
	if form == nil {
		notFoundError := exception.NotFoundError{Item: formService.constants.Field.Form}
		return notFoundError
	}

	// if form.UserID != request.UserID {
	// 	ForbiddenError := exception.ForbiddenError{Message: formService.constants.Field.Form}
	// 	return ForbiddenError
	// }

	err = formService.formRepository.DeleteForm(formService.db, request.FormID)
	if err != nil {
		return err
	}

	return nil
}

func (formService *FormService) GetAllForms(offset, limit int, filters *postgres.FormFilters) ([]formdto.BasicFormResponse, int64, error) {
	forms, err := formService.formRepository.FindAllForms(formService.db, offset, limit, filters)
	if err != nil {
		return nil, 0, err
	}

	count, err := formService.formRepository.CountAllForms(formService.db, filters)
	if err != nil {
		return nil, 0, err
	}

	formResponses := make([]formdto.BasicFormResponse, len(forms))
	for i, form := range forms {
		formResponses[i] = formdto.BasicFormResponse{
			FormID:     form.ID,
			Status:     form.Status.String(),
			UserID:     form.UserID,
			OperatorID: form.OperatorID,
			CreatedAt:  form.CreatedAt,
			UpdatedAt:  form.UpdatedAt,
		}
	}

	return formResponses, count, nil
}

func (formService *FormService) GetAllOperatorForms(offset, limit int, filters *postgres.OperatorFormFilters) ([]formdto.BasicFormResponse, int64, error) {
	operator, err := formService.userService.GetUserByID(filters.OperatorID)
	if err != nil {
		return nil, 0, err
	}
	if operator == nil {
		notFoundError := exception.NotFoundError{Item: formService.constants.Field.User}
		return nil, 0, notFoundError
	}
	userRoles, err := formService.userService.GetUserRoles(filters.OperatorID)
	if err != nil {
		return nil, 0, err
	}
	hasOperatorRole := false
	for _, role := range userRoles {
		if role.Name == enum.Operator.String() {
			hasOperatorRole = true
			break
		}
	}
	if !hasOperatorRole {
		forbiddenError := exception.ForbiddenError{Resource: formService.constants.Field.Role}
		return nil, 0, forbiddenError
	}

	forms, err := formService.formRepository.FindAllOperatorForms(formService.db, offset, limit, filters)
	if err != nil {
		return nil, 0, err
	}

	count, err := formService.formRepository.CountAllOperatorForms(formService.db, filters)
	if err != nil {
		return nil, 0, err
	}

	formResponses := make([]formdto.BasicFormResponse, len(forms))
	for i, form := range forms {
		formResponses[i] = formdto.BasicFormResponse{
			FormID:     form.ID,
			Status:     form.Status.String(),
			UserID:     form.UserID,
			OperatorID: form.OperatorID,
			CreatedAt:  form.CreatedAt,
			UpdatedAt:  form.UpdatedAt,
		}
	}

	return formResponses, count, nil
}

func (formService *FormService) AcceptForm(formID uint, userID uint) error {
	form, err := formService.formRepository.FindFormByID(formService.db, formID)
	if err != nil {
		return err
	}

	if form == nil {
		notFoundError := exception.NotFoundError{Item: formService.constants.Field.Form}
		return notFoundError
	}

	// Check if user has PermissionHandleOperators (supervisor)
	hasSupervisorPermission, err := formService.hasPermission(userID, enum.PermissionHandleOperators)
	if err != nil {
		return err
	}

	// Check if user is the assigned operator
	isAssignedOperator := form.OperatorID != nil && *form.OperatorID == userID

	// Authorization: Supervisor (has PermissionHandleOperators) or assigned operator can accept forms
	if !hasSupervisorPermission && !isAssignedOperator {
		forbiddenError := exception.ForbiddenError{Resource: formService.constants.Field.Form, Message: "Only supervisor or assigned operator can accept forms"}
		return forbiddenError
	}

	form.Status = enum.FormStatusApproved
	err = formService.formRepository.UpdateForm(formService.db, form)
	if err != nil {
		return err
	}

	return nil
}

func (formService *FormService) RejectForm(formID uint, userID uint) error {
	form, err := formService.formRepository.FindFormByID(formService.db, formID)
	if err != nil {
		return err
	}

	if form == nil {
		notFoundError := exception.NotFoundError{Item: formService.constants.Field.Form}
		return notFoundError
	}

	// Check if user has PermissionHandleOperators (supervisor)
	hasSupervisorPermission, err := formService.hasPermission(userID, enum.PermissionHandleOperators)
	if err != nil {
		return err
	}

	// Check if user is the assigned operator
	isAssignedOperator := form.OperatorID != nil && *form.OperatorID == userID

	// Authorization: Supervisor (has PermissionHandleOperators) or assigned operator can reject forms
	if !hasSupervisorPermission && !isAssignedOperator {
		forbiddenError := exception.ForbiddenError{Resource: formService.constants.Field.Form, Message: "Only supervisor or assigned operator can reject forms"}
		return forbiddenError
	}

	form.Status = enum.FormStatusRejected
	err = formService.formRepository.UpdateForm(formService.db, form)
	if err != nil {
		return err
	}

	return nil
}

func (formService *FormService) UpdateGeneralHealth(request formdto.UpdateGeneralHealthRequest) error {
	form, err := formService.formRepository.FindFormByID(formService.db, request.FormID)
	if err != nil {
		return err
	}
	if form == nil {
		notFoundError := exception.NotFoundError{Item: formService.constants.Field.Form}
		return notFoundError
	}

	isOp, err := formService.isOperator(request.UserID)
	if err != nil {
		return err
	}
	if isOp {
		canEdit, err := formService.canOperatorEditForm(form, request.UserID)
		if err != nil {
			return err
		}
		if !canEdit {
			forbiddenError := exception.ForbiddenError{Resource: formService.constants.Field.Form}
			return forbiddenError
		}
	}

	// if form.UserID != request.UserID {
	// 	ForbiddenError := exception.ForbiddenError{Message: formService.constants.Field.Form}
	// 	return ForbiddenError
	// }

	info, err := formService.formRepository.FindGeneralHealthByFormID(formService.db, request.FormID)
	if err != nil {
		return err
	}

	if info == nil {
		ForbiddenError := exception.ForbiddenError{Message: formService.constants.Field.Form}
		return ForbiddenError
	}

	info.DrinksAlcohol = request.DrinksAlcohol
	info.CupsPerWeek = request.CupsPerWeek
	if request.LastMonthSabzijatMeal != nil {
		info.LastMonthSabzijatMeal = *request.LastMonthSabzijatMeal
	}
	if request.LastMonthSabzijatWeight != nil {
		info.LastMonthSabzijatWeight = *request.LastMonthSabzijatWeight
	}
	if request.MediumActivityMonthInYear != nil {
		info.MediumActivityMonthInYear = *request.MediumActivityMonthInYear
	}
	if request.MediumActivityHourInWeek != nil {
		info.MediumActivityHourInWeek = *request.MediumActivityHourInWeek
	}
	if request.HardActivityMonthInYear != nil {
		info.HardActivityMonthInYear = *request.HardActivityMonthInYear
	}
	if request.HardActivityHourInWeek != nil {
		info.HardActivityHourInWeek = *request.HardActivityHourInWeek
	}
	info.SmokeAtLeast100 = request.SmokeAtLeast100
	info.SmokingAge = request.SmokingAge
	if request.SmokingNow != nil {
		info.SmokingNow = *request.SmokingNow
	}
	info.LeaveSmokingAge = request.LeaveSmokingAge
	info.CountSmokingDaily = request.CountSmokingDaily
	info.CountGheliandaily = request.CountGheliandaily
	info.CountSmokingDailyPast = request.CountSmokingDailyPast
	info.CountGheliandailyPast = request.CountGheliandailyPast

	if info.ID == 0 {
		err = formService.formRepository.CreateGeneralHealth(formService.db, info)
	} else {
		err = formService.formRepository.UpdateGeneralHealth(formService.db, info)
	}
	if err != nil {
		return err
	}

	// Log operator action if performed by operator
	if isOp {
		formService.logOperatorFormUpdate(request.UserID, request.FormID, "اطلاعات سلامت عمومی")
	}

	return nil
}

func (formService *FormService) UpdateMamography(request formdto.UpdateMamographyRequest) error {
	form, err := formService.formRepository.FindFormByID(formService.db, request.FormID)
	if err != nil {
		return err
	}
	if form == nil {
		notFoundError := exception.NotFoundError{Item: formService.constants.Field.Form}
		return notFoundError
	}

	isOp, err := formService.isOperator(request.UserID)
	if err != nil {
		return err
	}
	if isOp {
		canEdit, err := formService.canOperatorEditForm(form, request.UserID)
		if err != nil {
			return err
		}
		if !canEdit {
			forbiddenError := exception.ForbiddenError{Resource: formService.constants.Field.Form}
			return forbiddenError
		}
	}

	// if form.UserID != request.UserID {
	// 	ForbiddenError := exception.ForbiddenError{Message: formService.constants.Field.Form}
	// 	return ForbiddenError
	// }

	info, err := formService.formRepository.FindMamographyByFormID(formService.db, request.FormID)
	if err != nil {
		return err
	}
	if info == nil {
		ForbiddenError := exception.ForbiddenError{Message: formService.constants.Field.Form}
		return ForbiddenError
	}

	if request.GhaedeAge != nil {
		info.GhaedeAge = *request.GhaedeAge
	}
	if request.HasChildren != nil {
		info.HasChildren = *request.HasChildren
	}
	info.NumberOfChildren = request.NumberOfChildren
	info.AgeOfFirstBirth = request.AgeOfFirstBirth
	if request.MenopausalStatus != nil {
		info.MenopausalStatus = enum.MenopausalStatus(uint(*request.MenopausalStatus))
	}
	info.MenopauseAge = request.MenopauseAge
	info.HRT = request.HRT
	info.HRTUseLength = request.HRTUseLength
	if request.LastFiveYearsHRTUse != nil {
		info.LastFiveYearsHRTUse = *request.LastFiveYearsHRTUse
	}
	info.CurrentHRTUse = request.CurrentHRTUse
	info.IntendedHRTUse = request.IntendedHRTUse
	info.HRTType = request.HRTType
	info.Oral = request.Oral
	info.OralDuration = request.OralDuration
	info.OralTwoLastYears = request.OralTwoLastYears
	info.MamoGraphy = request.MamoGraphy
	info.Falop = request.Falop
	info.Andometrioz = request.Andometrioz
	if request.LeavePestan != nil {
		info.LeavePestan = *request.LeavePestan
	}
	if request.LeaveTokhmdan != nil {
		info.LeaveTokhmdan = *request.LeaveTokhmdan
	}
	info.LaDeColon = request.LaDeColon
	info.LaDePol = request.LaDePol
	info.AspLaMo = request.AspLaMo
	info.NsaiDLaMo = request.NsaiDLaMo
	info.LastFiveYearBloodTestInStool = request.LastFiveYearBloodTestInStool
	if request.NumberOfBreastBiopsies != nil {
		info.NumberOfBreastBiopsies = request.NumberOfBreastBiopsies
	}
	if request.HyperplasiaInBiopsy != nil {
		info.HyperplasiaInBiopsy = (*enum.HyperplasiaInBiopsyStatus)(request.HyperplasiaInBiopsy)
	}

	if info.ID == 0 {
		err = formService.formRepository.CreateMamography(formService.db, info)
	} else {
		err = formService.formRepository.UpdateMamography(formService.db, info)
	}
	if err != nil {
		return err
	}

	// Log operator action if performed by operator
	if isOp {
		formService.logOperatorFormUpdate(request.UserID, request.FormID, "اطلاعات ماموگرافی")
	}

	return nil
}

func (formService *FormService) UpdateContact(request formdto.UpdateContactRequest) error {
	form, err := formService.formRepository.FindFormByID(formService.db, request.FormID)
	if err != nil {
		return err
	}
	if form == nil {
		notFoundError := exception.NotFoundError{Item: formService.constants.Field.Form}
		return notFoundError
	}

	isOp, err := formService.isOperator(request.UserID)
	if err != nil {
		return err
	}
	if isOp {
		canEdit, err := formService.canOperatorEditForm(form, request.UserID)
		if err != nil {
			return err
		}
		if !canEdit {
			forbiddenError := exception.ForbiddenError{Resource: formService.constants.Field.Form}
			return forbiddenError
		}
	}

	// if form.UserID != request.UserID {
	// 	ForbiddenError := exception.ForbiddenError{Message: formService.constants.Field.Form}
	// 	return ForbiddenError
	// }

	info, err := formService.formRepository.FindContactByFormID(formService.db, request.FormID)
	if err != nil {
		return err
	}
	if info == nil {
		info = &entity.ContactInfo{FormID: request.FormID}
	}

	// Save old picture paths for cleanup
	var oldTestGenPicturePath *string
	var oldFatherTestGenPicturePath *string
	var oldMotherTestGenPicturePath *string
	if info.TestGenPicturePath != nil {
		oldTestGenPicturePath = info.TestGenPicturePath
	}
	if info.FatherTestGenPicturePath != nil {
		oldFatherTestGenPicturePath = info.FatherTestGenPicturePath
	}
	if info.MotherTestGenPicturePath != nil {
		oldMotherTestGenPicturePath = info.MotherTestGenPicturePath
	}

	if request.Name != nil {
		info.Name = *request.Name
	}
	info.TestGen = request.TestGen
	info.FmTestGen = request.FmTestGen
	if request.CallExpert != nil {
		info.CallExpert = *request.CallExpert
	}
	info.BirthCountry = request.BirthCountry
	info.Province = request.Province
	info.City = request.City
	info.Country = request.Country
	if request.Address != nil {
		info.Address = *request.Address
	}
	if request.PostalCode != nil {
		info.PostalCode = *request.PostalCode
	}

	// Handle TestGen picture upload
	if request.TestGenPicture != nil && request.TestGenPicture.Filename != "" {
		pictureKey := formService.constants.BucketPath.GetGeneticTestPath(
			request.FormID,
			"user",
			request.TestGenPicture.Filename,
		)
		info.TestGenPicturePath = &pictureKey

		// Upload picture
		if err := formService.s3Storage.UploadObject(enum.BucketTypeGeneticTest, pictureKey, request.TestGenPicture); err != nil {
			return err
		}

		// Delete old picture if it exists and is different
		if oldTestGenPicturePath != nil && *oldTestGenPicturePath != "" && *oldTestGenPicturePath != pictureKey {
			if err := formService.s3Storage.DeleteObject(enum.BucketTypeGeneticTest, *oldTestGenPicturePath); err != nil {
				fmt.Printf("Warning: failed to delete old genetic test picture %s: %v\n", *oldTestGenPicturePath, err)
			}
		}
	}

	// Handle FatherTestGen picture upload
	if request.FatherTestGenPicture != nil && request.FatherTestGenPicture.Filename != "" {
		pictureKey := formService.constants.BucketPath.GetFatherGeneticTestPath(
			request.FormID,
			request.FatherTestGenPicture.Filename,
		)
		info.FatherTestGenPicturePath = &pictureKey

		// Upload picture
		if err := formService.s3Storage.UploadObject(enum.BucketTypeGeneticTest, pictureKey, request.FatherTestGenPicture); err != nil {
			return err
		}

		// Delete old picture if it exists and is different
		if oldFatherTestGenPicturePath != nil && *oldFatherTestGenPicturePath != "" && *oldFatherTestGenPicturePath != pictureKey {
			if err := formService.s3Storage.DeleteObject(enum.BucketTypeGeneticTest, *oldFatherTestGenPicturePath); err != nil {
				fmt.Printf("Warning: failed to delete old father genetic test picture %s: %v\n", *oldFatherTestGenPicturePath, err)
			}
		}
	}

	// Handle MotherTestGen picture upload
	if request.MotherTestGenPicture != nil && request.MotherTestGenPicture.Filename != "" {
		pictureKey := formService.constants.BucketPath.GetMotherGeneticTestPath(
			request.FormID,
			request.MotherTestGenPicture.Filename,
		)
		info.MotherTestGenPicturePath = &pictureKey

		// Upload picture
		if err := formService.s3Storage.UploadObject(enum.BucketTypeGeneticTest, pictureKey, request.MotherTestGenPicture); err != nil {
			return err
		}

		// Delete old picture if it exists and is different
		if oldMotherTestGenPicturePath != nil && *oldMotherTestGenPicturePath != "" && *oldMotherTestGenPicturePath != pictureKey {
			if err := formService.s3Storage.DeleteObject(enum.BucketTypeGeneticTest, *oldMotherTestGenPicturePath); err != nil {
				fmt.Printf("Warning: failed to delete old mother genetic test picture %s: %v\n", *oldMotherTestGenPicturePath, err)
			}
		}
	}

	if info.ID == 0 {
		err = formService.formRepository.CreateContact(formService.db, info)
	} else {
		err = formService.formRepository.UpdateContact(formService.db, info)
	}
	if err != nil {
		return err
	}

	// Log operator action if performed by operator
	if isOp {
		formService.logOperatorFormUpdate(request.UserID, request.FormID, "اطلاعات تماس")
	}

	return nil
}

func (formService *FormService) UpdateLungCancer(request formdto.UpdateLungCancerRequest) error {
	form, err := formService.formRepository.FindFormByID(formService.db, request.FormID)
	if err != nil {
		return err
	}
	if form == nil {
		notFoundError := exception.NotFoundError{Item: formService.constants.Field.Form}
		return notFoundError
	}

	isOp, err := formService.isOperator(request.UserID)
	if err != nil {
		return err
	}
	if isOp {
		canEdit, err := formService.canOperatorEditForm(form, request.UserID)
		if err != nil {
			return err
		}
		if !canEdit {
			forbiddenError := exception.ForbiddenError{Resource: formService.constants.Field.Form}
			return forbiddenError
		}
	}

	// if form.UserID != request.UserID {
	// 	ForbiddenError := exception.ForbiddenError{Message: formService.constants.Field.Form}
	// 	return ForbiddenError
	// }

	info, err := formService.formRepository.FindLungCancerByFormID(formService.db, request.FormID)
	if err != nil {
		return err
	}
	if info == nil {
		info = &entity.LungCancerInfo{FormID: request.FormID}
	}

	info.InsuranceStatus = request.InsuranceStatus
	info.SupplementaryInsurances = request.SupplementaryInsurances
	if request.Hypertension != nil {
		info.Hypertension = *request.Hypertension
	}
	info.HypertensionTreatment = request.HypertensionTreatment
	if request.HeartDisease != nil {
		info.HeartDisease = *request.HeartDisease
	}
	info.HeartDiseaseTreatment = request.HeartDiseaseTreatment
	if request.Diabetes != nil {
		info.Diabetes = *request.Diabetes
	}
	info.DiabetesTreatment = request.DiabetesTreatment
	info.ChronicLungDisease = request.ChronicLungDisease
	info.ChronicLungDiseaseType = request.ChronicLungDiseaseType
	if request.LungCancerHistory != nil {
		info.LungCancerHistory = *request.LungCancerHistory
	}
	if request.OtherCancerHistory != nil {
		info.OtherCancerHistory = *request.OtherCancerHistory
	}
	if request.OtherCancerType != nil {
		info.OtherCancerType = (*enum.CancerType)(request.OtherCancerType)
	}
	info.LungCancerFamily = request.LungCancerFamily
	info.LungCancerFamilyRelation = request.LungCancerFamilyRelation
	info.OtherCancerFamily = request.OtherCancerFamily
	if request.OtherCancerFamilyType != nil {
		info.OtherCancerFamilyType = (*enum.CancerType)(request.OtherCancerFamilyType)
	}
	info.OtherCancerFamilyRelation = request.OtherCancerFamilyRelation
	info.OccupationalExposure = request.OccupationalExposure
	if request.CurrentSmoking != nil {
		info.CurrentSmoking = *request.CurrentSmoking
	}
	info.SmokingStartAgeCurrent = request.SmokingStartAgeCurrent
	info.SmokingTypesCurrent = request.SmokingTypesCurrent
	info.CigarettesPerDayCurrent = request.CigarettesPerDayCurrent
	info.CigarPerDayCurrent = request.CigarPerDayCurrent
	info.ECigPerDayCurrent = request.ECigPerDayCurrent
	info.PipePerDayCurrent = request.PipePerDayCurrent
	info.ChapoghPerDayCurrent = request.ChapoghPerDayCurrent
	info.SmokedOpiumPerDayCurrent = request.SmokedOpiumPerDayCurrent
	info.ChewedOpiumPerDayCurrent = request.ChewedOpiumPerDayCurrent
	info.HookahPerWeekCurrent = request.HookahPerWeekCurrent
	info.PastSmoking = request.PastSmoking
	info.SmokingStartAgePast = request.SmokingStartAgePast
	info.SmokingTypesPast = request.SmokingTypesPast
	info.CigarettesPerDayPast = request.CigarettesPerDayPast
	info.CigarPerDayPast = request.CigarPerDayPast
	info.ECigPerDayPast = request.ECigPerDayPast
	info.PipePerDayPast = request.PipePerDayPast
	info.ChapoghPerDayPast = request.ChapoghPerDayPast
	info.SmokedOpiumPerDayPast = request.SmokedOpiumPerDayPast
	info.ChewedOpiumPerDayPast = request.ChewedOpiumPerDayPast
	info.HookahPerWeekPast = request.HookahPerWeekPast
	if request.SecondhandSmoke != nil {
		info.SecondhandSmoke = *request.SecondhandSmoke
	}
	info.SecondhandSmokeLocation = request.SecondhandSmokeLocation

	if info.ID == 0 {
		err = formService.formRepository.CreateLungCancer(formService.db, info)
	} else {
		err = formService.formRepository.UpdateLungCancer(formService.db, info)
	}
	if err != nil {
		return err
	}

	// Log operator action if performed by operator
	if isOp {
		formService.logOperatorFormUpdate(request.UserID, request.FormID, "اطلاعات سرطان ریه")
	}

	return nil
}

func (formService *FormService) AssignOperator(request formdto.AssignOperatorRequest) error {
	form, err := formService.formRepository.FindFormByID(formService.db, request.FormID)
	if err != nil {
		return err
	}
	if form == nil {
		notFoundError := exception.NotFoundError{Item: formService.constants.Field.Form}
		return notFoundError
	}

	operator, err := formService.userService.GetUserByID(request.OperatorID)
	if err != nil {
		return err
	}
	if operator == nil {
		notFoundError := exception.NotFoundError{Item: formService.constants.Field.User}
		return notFoundError
	}
	userRoles, err := formService.userService.GetUserRoles(request.OperatorID)
	if err != nil {
		return err
	}
	hasOperatorRole := false
	for _, role := range userRoles {
		if role.Name == enum.Operator.String() {
			hasOperatorRole = true
			break
		}
	}
	if !hasOperatorRole {
		forbiddenError := exception.ForbiddenError{Resource: formService.constants.Field.Role}
		return forbiddenError
	}

	form.OperatorID = &request.OperatorID
	form.Status = enum.FormStatusAssigned
	err = formService.formRepository.UpdateForm(formService.db, form)
	if err != nil {
		return err
	}

	log := actionlogdto.LogAction{
		ActorID:    request.UserID,
		TargetID:   &request.OperatorID,
		Action:     enum.ActionTypeFormAssigned,
		ResourceID: &request.FormID,
		Details:    "فرم به اپراتور اساین شد",
	}
	formService.actionLogService.LogAction(log)

	return nil
}
func (formService *FormService) UnassignOperator(request formdto.UnassignOperatorRequest) error {
	form, err := formService.formRepository.FindFormByID(formService.db, request.FormID)
	if err != nil {
		return err
	}
	if form == nil {
		notFoundError := exception.NotFoundError{Item: formService.constants.Field.Form}
		return notFoundError
	}

	form.OperatorID = nil
	err = formService.formRepository.UpdateForm(formService.db, form)
	if err != nil {
		return err
	}

	log := actionlogdto.LogAction{
		ActorID:    request.UserID,
		ResourceID: &request.FormID,
		Action:     enum.ActionTypeFormUnAssigned,
		Details:    "فرم از اپراتور گرفته شد",
	}
	formService.actionLogService.LogAction(log)

	return nil
}

func (formService *FormService) GetAllCancerTypes() ([]generaldto.EnumResponse, error) {
	cancerTypes := enum.GetAllCancerTypes()
	response := make([]generaldto.EnumResponse, len(cancerTypes))
	for i, cancerType := range cancerTypes {
		response[i] = generaldto.EnumResponse{
			ID:   uint(cancerType),
			Name: cancerType.String(),
		}
	}
	return response, nil
}

func (formService *FormService) GetAllGenders() ([]generaldto.EnumResponse, error) {
	genders := enum.GetAllGenders()
	response := make([]generaldto.EnumResponse, len(genders))
	for i, gender := range genders {
		response[i] = generaldto.EnumResponse{
			ID:   uint(gender),
			Name: gender.String(),
		}
	}
	return response, nil
}

func (formService *FormService) GetAllMenopausalStatuses() ([]generaldto.EnumResponse, error) {
	statuses := enum.GetAllMenopausalStatuses()
	response := make([]generaldto.EnumResponse, len(statuses))
	for i, status := range statuses {
		response[i] = generaldto.EnumResponse{
			ID:   uint(status),
			Name: status.String(),
		}
	}
	return response, nil
}

func (formService *FormService) GetAllFormStatuses() ([]generaldto.EnumResponse, error) {
	statuses := enum.GetAllFormStatuses()
	response := make([]generaldto.EnumResponse, len(statuses))
	for i, status := range statuses {
		response[i] = generaldto.EnumResponse{
			ID:   uint(status),
			Name: status.String(),
		}
	}
	return response, nil
}

func (formService *FormService) GetAllHyperplasiaInBiopsyStatuses() ([]generaldto.EnumResponse, error) {
	statuses := enum.GetAllHyperplasiaInBiopsyStatuses()
	response := make([]generaldto.EnumResponse, len(statuses))
	for i, status := range statuses {
		response[i] = generaldto.EnumResponse{
			ID:   uint(status),
			Name: status.String(),
		}
	}
	return response, nil
}

func (formService *FormService) GetAllLifeStatuses() ([]generaldto.EnumResponse, error) {
	statuses := enum.GetAllLifeStatus()
	response := make([]generaldto.EnumResponse, len(statuses))
	for i, status := range statuses {
		response[i] = generaldto.EnumResponse{
			ID:   uint(status),
			Name: status.String(),
		}
	}
	return response, nil
}

func (formService *FormService) GetAllRelativeTypes() ([]generaldto.EnumResponse, error) {
	statuses := enum.GetAllRelatives()
	response := make([]generaldto.EnumResponse, len(statuses))
	for i, status := range statuses {
		response[i] = generaldto.EnumResponse{
			ID:   uint(status),
			Name: status.String(),
		}
	}
	return response, nil
}
