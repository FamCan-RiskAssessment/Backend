package health

import (
	"context"
	"net/http"
	"time"

	"github.com/FamCan-RiskAssessment/Backend/internal/infrastructure/database"
	"github.com/gin-gonic/gin"
)

const probeTimeout = 3 * time.Second

type HealthController struct {
	db  database.Database
	rdb database.Cache
}

func NewHealthController(db database.Database, rdb database.Cache) *HealthController {
	return &HealthController{
		db:  db,
		rdb: rdb,
	}
}

type readinessResponse struct {
	Status string            `json:"status"`
	Checks map[string]string `json:"checks"`
}

func (c *HealthController) Readiness(ctx *gin.Context) {
	checks := map[string]string{
		"postgres": c.checkPostgres(ctx),
		"redis":    c.checkRedis(ctx),
	}

	status := "ready"
	statusCode := http.StatusOK
	for _, result := range checks {
		if result != "ok" {
			status = "not_ready"
			statusCode = http.StatusServiceUnavailable
			break
		}
	}

	ctx.JSON(statusCode, readinessResponse{
		Status: status,
		Checks: checks,
	})
}

func (c *HealthController) checkPostgres(ctx *gin.Context) string {
	sqlDB, err := c.db.GetDB().DB()
	if err != nil {
		return err.Error()
	}

	probeCtx, cancel := context.WithTimeout(ctx.Request.Context(), probeTimeout)
	defer cancel()

	if err := sqlDB.PingContext(probeCtx); err != nil {
		return err.Error()
	}

	return "ok"
}

func (c *HealthController) checkRedis(ctx *gin.Context) string {
	probeCtx, cancel := context.WithTimeout(ctx.Request.Context(), probeTimeout)
	defer cancel()

	if err := c.rdb.GetRDB().Ping(probeCtx).Err(); err != nil {
		return err.Error()
	}

	return "ok"
}
