package db

import (
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"time"

	_ "github.com/lib/pq"
	"github.com/rcjeferson/go-api-products/internal/model"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "go-api-products"
)

func ConnectDB() (*sql.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Panic(err)
	}

	err = db.Ping()
	if err != nil {
		log.Panic(err)
	}

	slog.Info(fmt.Sprintf("Successfully connected on database %s on host %s:%d!", dbname, host, port))

	return db, nil
}

func Ping(db *sql.DB) model.ServiceMetrics {
	sm := model.ServiceMetrics{}

	start := time.Now()
	err := db.Ping()
	elapsed := time.Since(start)

	sm.Status = "OK"
	sm.Error = ""
	sm.Latency = elapsed.String()

	if err != nil {
		sm.Status = "FAILED"
		sm.Error = err.Error()

		slog.Error("Failed to Ping Database: " + string(err.Error()))
	}

	return sm
}
