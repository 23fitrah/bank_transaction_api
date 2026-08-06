package config

import (
	"context"
	"fmt"
	"time"

	"github.com/bytedance/sonic"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/cast"
)

var ctx = context.Background()

func ConnectRedis() (*redis.Client, error) {

	cfg, err := Load()

	if err != nil {
		return nil, fmt.Errorf("Failed to load config: %v", err)
	}

	addr := cfg.RedisDB.RedisHost + ":" + cfg.RedisDB.RedisPort
	password := cfg.RedisDB.RedisPassword
	db := cast.ToInt(cfg.RedisDB.RedisDatabase)

	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	if _, err := rdb.Ping(ctx).Result(); err != nil {
		logEntry := map[string]interface{}{
			"level":     "error",
			"timestamp": time.Now().Format(time.RFC3339),
			"message":   "Failed to connect to redis: ",
			"error":     err.Error(),
		}
		jsonBytes, _ := sonic.Marshal(logEntry)

		fmt.Println(string(jsonBytes))
	} else {
		logEntry := map[string]interface{}{
			"level":     "info",
			"timestamp": time.Now().Format(time.RFC3339),
			"message":   "redis connected successfully",
		}
		jsonBytes, _ := sonic.Marshal(logEntry)
		fmt.Println(string(jsonBytes))

	}
	return rdb, err
}
