package middleware

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"reflect"
	"strconv"

	"github.com/FamCan-RiskAssessment/Backend/bootstrap"
	"github.com/FamCan-RiskAssessment/Backend/internal/domain/exception"
	"github.com/FamCan-RiskAssessment/Backend/internal/presentation/controller"
	"github.com/gin-gonic/gin"
)

const (
	genericError = "errors.generic"
)

type RecoveryMiddleware struct {
	constants *bootstrap.Constants
}

func NewRecoveryMiddleware(constants *bootstrap.Constants) *RecoveryMiddleware {
	return &RecoveryMiddleware{
		constants: constants,
	}
}

func (recovery RecoveryMiddleware) Recovery(ctx *gin.Context) {
	defer func() {
		if rec := recover(); rec != nil {
			if err, ok := rec.(error); ok {
				recovery.handleRecoveredError(ctx, err)
				ctx.Abort()
			}
		}
	}()
	ctx.Next()
}

func (recovery RecoveryMiddleware) handleRecoveredError(ctx *gin.Context, err error) {
	slog.Warn("request failed",
		"method", ctx.Request.Method,
		"path", ctx.Request.URL.Path,
		"error", err.Error(),
	)

	if validationErrors, ok := err.(exception.ValidationErrors); ok {
		handleValidationError(ctx, validationErrors, recovery.constants.Context.Translator)
	} else if bindingError, ok := err.(exception.BindingError); ok {
		handleBindingError(ctx, bindingError, recovery.constants.Context.Translator)
	} else if fileValidationError, ok := err.(exception.FileValidationError); ok {
		handleFileValidationError(ctx, fileValidationError)
	} else if rateLimitError, ok := err.(*exception.RateLimitError); ok {
		handleRateLimitError(ctx, *rateLimitError, recovery.constants.Context.Translator)
	} else if conflictErrors, ok := err.(exception.ConflictErrors); ok {
		handleConflictError(ctx, conflictErrors, recovery.constants.Context.Translator)
	} else if authError, ok := err.(*exception.AuthError); ok {
		handleAuthError(ctx, *authError, recovery.constants.Context.Translator)
	} else if calcError, ok := err.(*exception.CalcError); ok {
		handleCalcError(ctx, *calcError, recovery.constants.Context.Translator)
	} else if serviceUnavailableError, ok := err.(*exception.ServiceUnavailableError); ok {
		handleServiceUnavailableError(ctx, *serviceUnavailableError, recovery.constants.Context.Translator)
	} else if notFoundError, ok := err.(exception.NotFoundError); ok {
		handleNotFoundError(ctx, notFoundError, recovery.constants.Context.Translator)
	} else if forbiddenError, ok := err.(exception.ForbiddenError); ok {
		handleForbiddenError(ctx, forbiddenError, recovery.constants.Context.Translator)
	} else if verificationError, ok := err.(exception.VerificationError); ok {
		handleVerificationError(ctx, verificationError)
	} else if fieldError, ok := err.(exception.FieldError); ok {
		handleFieldError(ctx, fieldError, recovery.constants.Context.Translator)
	} else {
		unhandledErrors(ctx, recovery.constants.Context.Translator)
	}
}

func handleValidationError(ctx *gin.Context, validationErrors exception.ValidationErrors, transKey string) {
	trans := controller.GetTranslator(ctx, transKey)
	errorMessages := make(map[string]map[string]string)

	for _, validationError := range validationErrors.Errors {
		if _, ok := errorMessages[validationError.Field]; !ok {
			errorMessages[validationError.Field] = make(map[string]string)
		}
		fieldName, _ := trans.Translate(validationError.Field)
		message, _ := trans.Translate(fmt.Sprintf("errors.%s", validationError.Tag), fieldName)
		errorMessages[validationError.Field][validationError.Tag] = message
	}

	controller.Response(ctx, 422, errorMessages, nil)
}

func handleBindingError(ctx *gin.Context, bindingError exception.BindingError, transKey string) {
	trans := controller.GetTranslator(ctx, transKey)
	message, _ := trans.Translate(genericError)

	var numError *strconv.NumError
	if errors.As(bindingError.Err, &numError) {
		message, _ = trans.Translate("errors.numeric", numError.Num)
	} else if errors.Is(bindingError.Err, http.ErrMissingFile) {
		message, _ = trans.Translate("errors.fileRequired")
	} else {
		var typeError *json.UnmarshalTypeError
		if errors.As(bindingError.Err, &typeError) {
			errorKey := "errors.invalidType"
			if isNumericKind(typeError.Type.Kind()) {
				errorKey = "errors.numeric"
			}
			message, _ = trans.Translate(errorKey, typeError.Field)
		}
	}

	controller.Response(ctx, 400, message, nil)
}

func isNumericKind(k reflect.Kind) bool {
	switch k {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
		reflect.Float32, reflect.Float64:
		return true
	default:
		return false
	}
}

func handleFileValidationError(ctx *gin.Context, fileValidationError exception.FileValidationError) {
	controller.Response(ctx, 400, fileValidationError.Message, nil)
}

func handleRateLimitError(ctx *gin.Context, rateLimitError exception.RateLimitError, transKey string) {
	trans := controller.GetTranslator(ctx, transKey)

	message, _ := trans.Translate(genericError)
	switch rateLimitError.Type {
	case exception.ErrorTypeRequestRateLimit:
		message, _ = trans.Translate("errors.rateLimit")
	case exception.ErrorTypeConcurrentInstallLimit:
		message, _ = trans.Translate("errors.installRateLimit")
	}

	controller.Response(ctx, 429, message, nil)
}

func handleConflictError(ctx *gin.Context, conflictErrors exception.ConflictErrors, transKey string) {
	trans := controller.GetTranslator(ctx, transKey)
	errorMessages := make(map[string]map[string]string)

	for _, conflictError := range conflictErrors.Errors {
		if _, ok := errorMessages[conflictError.Field]; !ok {
			errorMessages[conflictError.Field] = make(map[string]string)
		}
		fieldName, _ := trans.Translate(conflictError.Field)
		message, _ := trans.Translate(fmt.Sprintf("errors.%s", conflictError.Tag), fieldName)
		errorMessages[conflictError.Field][conflictError.Tag] = message
	}

	controller.Response(ctx, 409, errorMessages, nil)
}

func handleAuthError(ctx *gin.Context, authError exception.AuthError, transKey string) {
	trans := controller.GetTranslator(ctx, transKey)

	message, _ := trans.Translate(genericError)
	switch authError.Type {
	case exception.ErrorTypeInvalidCredentials:
		message, _ = trans.Translate("errors.invalidAuthCredentials")
	case exception.ErrorTypeExpiredToken:
		message, _ = trans.Translate("errors.expiredAuthToken")
	case exception.ErrorTypeInvalidToken:
		message, _ = trans.Translate("errors.invalidAuthToken")
	case exception.ErrorTypeUnauthorized:
		message, _ = trans.Translate("errors.unauthorized")
	}

	controller.Response(ctx, 401, message, nil)
}

func handleNotFoundError(ctx *gin.Context, notFoundError exception.NotFoundError, transKey string) {
	trans := controller.GetTranslator(ctx, transKey)
	itemName, _ := trans.Translate(notFoundError.Item)
	message, _ := trans.Translate("errors.notFound", itemName)
	controller.Response(ctx, 404, message, nil)
}

func handleForbiddenError(ctx *gin.Context, forbiddenError exception.ForbiddenError, transKey string) {
	trans := controller.GetTranslator(ctx, transKey)
	var message string
	switch forbiddenError.Type {
	case exception.ForbiddenTypeVerificationFailed:
		message, _ = trans.Translate("errors.verificationFailed")
	default:
		ResourceName, _ := trans.Translate(forbiddenError.Resource)
		message, _ = trans.Translate("errors.forbiddenError", ResourceName)
	}
	controller.Response(ctx, 403, message, nil)
}

func handleFieldError(ctx *gin.Context, fieldError exception.FieldError, transKey string) {
	trans := controller.GetTranslator(ctx, transKey)
	fieldName, _ := trans.Translate(fieldError.Field)
	message, _ := trans.Translate(fmt.Sprintf("errors.%s", fieldError.Tag), fieldName)
	controller.Response(ctx, 400, message, nil)
}

func handleCalcError(ctx *gin.Context, calcError exception.CalcError, transKey string) {
	trans := controller.GetTranslator(ctx, transKey)
	errorMessage, _ := trans.Translate(genericError)

	statusCode := 500
	switch calcError.Type {
	case exception.ErrorTypeAPIFailure:
		errorMessage, _ = trans.Translate("errors.calcAPIFailure", calcError.Model)
		statusCode = 503
	case exception.ErrorTypeMissingData:
		errorMessage, _ = trans.Translate("errors.calcMissingData", calcError.Model)
		statusCode = 422
	case exception.ErrorTypeInvalidData:
		errorMessage, _ = trans.Translate("errors.calcInvalidData", calcError.Model)
		statusCode = 400
	case exception.ErrorTypeDatabaseError:
		errorMessage, _ = trans.Translate("errors.calcDatabaseError")
		statusCode = 500
	}

	// Include specific error details in response data
	responseData := map[string]interface{}{
		"error": calcError.Message,
		"model": calcError.Model,
		"type":  string(calcError.Type),
	}

	controller.Response(ctx, statusCode, errorMessage, responseData)
}

func handleServiceUnavailableError(ctx *gin.Context, serviceUnavailableError exception.ServiceUnavailableError, transKey string) {
	trans := controller.GetTranslator(ctx, transKey)
	errorMessage, _ := trans.Translate(genericError)

	switch serviceUnavailableError.Service {
	case exception.ServiceSMS:
		switch serviceUnavailableError.Reason {
		case exception.ReasonNetwork:
			errorMessage, _ = trans.Translate("errors.smsUnavailableVPN")
		default:
			errorMessage, _ = trans.Translate("errors.serviceUnavailable")
		}
	default:
		errorMessage, _ = trans.Translate("errors.serviceUnavailable")
	}

	controller.Response(ctx, 503, errorMessage, nil)
}

func handleVerificationError(ctx *gin.Context, verificationError exception.VerificationError) {
	errorMessages := map[string]map[string]string{
		verificationError.Field: {
			"mismatch": verificationError.Message,
		},
	}
	controller.Response(ctx, 422, errorMessages, nil)
}

func unhandledErrors(ctx *gin.Context, transKey string) {
	slog.Error("unhandled request error",
		"method", ctx.Request.Method,
		"path", ctx.Request.URL.Path,
	)

	trans := controller.GetTranslator(ctx, transKey)
	errorMessage, _ := trans.Translate(genericError)

	controller.Response(ctx, 500, errorMessage, nil)
}
