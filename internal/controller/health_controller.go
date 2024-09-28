package controller

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rcjeferson/go-api-products/internal/db"
	"github.com/rcjeferson/go-api-products/internal/model"
	"github.com/redis/go-redis/v9"
)

type HealthController struct {
	databaseConn *sql.DB
	rdb          *redis.Client
	useCache     bool
}

func NewHealthController(db *sql.DB, rdb *redis.Client, useCache bool) HealthController {
	return HealthController{
		databaseConn: db,
		rdb:          rdb,
		useCache:     useCache,
	}
}

func (hc *HealthController) Check(ctx *gin.Context) {
	hm := model.HealthMetrics{}
	finalStatus := http.StatusOK

	// Database Check
	dbSM := db.DbPing(hc.databaseConn)
	hm.Database.Status = dbSM.Status
	hm.Database.Error = dbSM.Error
	hm.Database.Latency = dbSM.Latency

	// Redis Check
	if hc.useCache {
		redisSM := db.RedisPing(hc.rdb)
		hm.Redis.Status = redisSM.Status
		hm.Redis.Error = redisSM.Error
		hm.Redis.Latency = redisSM.Latency
	} else {
		hm.Redis.Status = "OK"
		hm.Redis.Error = "Cache is not used!"
		hm.Redis.Latency = "0"
	}

	// Return error only for required resources.
	// Cache is optional.
	if hm.Database.Status == "FAILED" {
		finalStatus = http.StatusServiceUnavailable
	}

	ctx.JSON(finalStatus, hm)
}
