package routes

import (
	"github.com/FamCan-RiskAssessment/Backend/wire"
	"github.com/gin-gonic/gin"
)

func Run(ginEngine *gin.Engine, app *wire.Application) {
	ginEngine.Use(app.Middlewares.Localization.Localization)
	ginEngine.Use(app.Middlewares.Recovery.Recovery)
	SetupGeneralRoutes(ginEngine.Group("/"), app)
	SetupAdminRoutes(ginEngine.Group("/admin"), app)
}
