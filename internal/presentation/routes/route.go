package routes

import (
	"github.com/FamCan-RiskAssessment/Backend/internal/domain/enum"
	"github.com/FamCan-RiskAssessment/Backend/wire"
	"github.com/gin-gonic/gin"
)

func Run(ginEngine *gin.Engine, app *wire.Application) {
	ginEngine.Use(gin.Logger())
	ginEngine.Use(app.Middlewares.Recovery.Recovery)
	ginEngine.Use(app.Middlewares.Localization.Localization)
	ginEngine.Use(app.Middlewares.Cors.CORS())
	ginEngine.Use(app.Middlewares.RateLimit.RateLimit())
	registerGeneralRoutes(ginEngine.Group("/"), app)
	registerAdminRoutes(ginEngine.Group("/admin"), app)
	registerCustomerRoutes(ginEngine.Group("/"), app)
}

func registerAdminRoutes(ginEngine *gin.RouterGroup, app *wire.Application) {
	ginEngine.Use(app.Middlewares.Auth.AuthRequired)
	ginEngine.Use(app.Middlewares.Auth.RequiredWithPermission([]enum.PermissionType{enum.PermissionSetPassword}))
	SetupAdminRoutes(ginEngine, app)
}

func registerCustomerRoutes(ginEngine *gin.RouterGroup, app *wire.Application) {
	ginEngine.Use(app.Middlewares.Auth.AuthRequired)
	SetupCustomerRoutes(ginEngine, app)
}

func registerGeneralRoutes(ginEngine *gin.RouterGroup, app *wire.Application) {
	SetupGeneralRoutes(ginEngine, app)
}
