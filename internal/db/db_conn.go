package db

import (
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"os"
	"time"

	_ "github.com/lib/pq"
	"github.com/rcjeferson/go-api-products/internal/model"
)

func ConnectDB() (*sql.DB, error) {
	var (
		host     = os.Getenv("DB_HOST")
		port     = os.Getenv("DB_PORT")
		user     = os.Getenv("DB_USER")
		password = os.Getenv("DB_PASSWORD")
		dbname   = os.Getenv("DB_DATABASE")
	)

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
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

	slog.Info(fmt.Sprintf("Successfully connected on database %s on host %s:%s!", dbname, host, port))

	return db, nil
}

func DbPing(db *sql.DB) model.ServiceMetrics {
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
