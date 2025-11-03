package main

import (
	"fmt"

	"github.com/FamCan-RiskAssessment/Backend/bootstrap"
	"github.com/FamCan-RiskAssessment/Backend/internal/domain/entity"
	"github.com/FamCan-RiskAssessment/Backend/internal/presentation/routes"
	"github.com/FamCan-RiskAssessment/Backend/wire"
	"github.com/gin-gonic/gin"
)

func main() {
	ginEngine := gin.New()

	config := bootstrap.Run()

	app, err := wire.InitializeApplication(config)
	if err != nil {
		panic(err)
	}

	app.Database.DB.GetDB().AutoMigrate(
		&entity.User{},
		&entity.Role{},
		&entity.Permission{},
		&entity.Form{},
		&entity.BasicInfo{},
		&entity.GeneralHealthInfo{},
		&entity.MamoGraphyInfo{},
		&entity.CancerInfo{},
		&entity.FamilyCancerInfo{},
		&entity.ContactInfo{},
		&entity.LungCancerInfo{},
		&entity.ActionLog{},
		&entity.Premm5Result{},
		&entity.BCRAResult{},
	)

	app.Seeds.RoleSeeder.SeedRoles()
	app.Seeds.DummySeeder.SeedDummy()

	routes.Run(ginEngine, app)

	ginEngine.Run(fmt.Sprintf(":%s", config.Env.Server.Port))
}
