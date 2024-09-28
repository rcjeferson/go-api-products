package db

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"time"

	"github.com/rcjeferson/go-api-products/internal/model"
	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

func ConnectRedis() (*redis.Client, error) {
	var (
		host     = os.Getenv("REDIS_HOST")
		password = os.Getenv("REDIS_PASSWORD")
		dbEnv    = os.Getenv("REDIS_DATABASE")
	)

	db, err := strconv.Atoi(dbEnv)
	if err != nil {
		slog.Error("Failed to convert dbEnv to int.")

		return nil, err
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     host,
		Password: password,
		DB:       db,
	})

	return rdb, nil
}

func GetCache(key string, rdb *redis.Client) (string, error) {
	cachedItem, err := rdb.Get(ctx, key).Result()

	if err == redis.Nil {
		slog.Info(fmt.Sprintf("Cache not found for key '%s'!", key))
		return "", redis.Nil

	} else if err != nil {
		slog.Warn(fmt.Sprintf("Error while getting key '%s' from cache: ", key), err)

		return "", err

	} else {
		slog.Info(fmt.Sprintf("Returning '%s' from cache!", key))

		return cachedItem, nil
	}
}

func SetCache(key string, value []byte, ttl time.Duration, rdb *redis.Client) {
	ttlDuration := time.Second * ttl

	err := rdb.Set(ctx, key, value, ttlDuration).Err()
	if err != nil {
		slog.Error(fmt.Sprintf("Error while setting cache for key '%s'", key), err)
	}

	slog.Info(fmt.Sprintf("Cache created for key '%s'", key), err)
}

func RedisPing(rdb *redis.Client) model.ServiceMetrics {
	sm := model.ServiceMetrics{}

	start := time.Now()
	result := rdb.Ping(ctx)
	elapsed := time.Since(start)

	sm.Status = "OK"
	sm.Error = ""
	sm.Latency = elapsed.String()

	if result.Err() != nil {
		sm.Status = "FAILED"
		sm.Error = result.String()

		slog.Error("Failed to Ping Redis: " + result.String())
	}

	return sm
}
