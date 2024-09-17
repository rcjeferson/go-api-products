package controller

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rcjeferson/go-api-products/internal/db"
	"github.com/rcjeferson/go-api-products/internal/model"
)

type HealthController struct {
	DatabaseConn *sql.DB
}

func NewHealthController(db *sql.DB) HealthController {
	return HealthController{
		DatabaseConn: db,
	}
}

func (hc *HealthController) Check(ctx *gin.Context) {
	hm := model.HealthMetrics{}
	finalStatus := http.StatusOK

	// Database Check
	dbSM := db.Ping(hc.DatabaseConn)
	hm.Database.Status = dbSM.Status
	hm.Database.Error = dbSM.Error
	hm.Database.Latency = dbSM.Latency

	if hm.Database.Status == "FAILED" {
		finalStatus = http.StatusServiceUnavailable
	}

	ctx.JSON(finalStatus, hm)
}
