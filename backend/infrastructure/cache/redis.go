package cache

import (
	"context"
	"fmt"
	"log/slog"

	"backend/config"

	"github.com/redis/go-redis/v9"
)

var Client *redis.Client

func ConnectRedis(cfg *config.Config) error {
	addr := fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort)

	Client = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "", // Set securely in active deployment environments
		DB:       0,
	})

	// Verify connection successfully
	ctx := context.Background()
	_, err := Client.Ping(ctx).Result()
	if err != nil {
		slog.Error("Failed to initialize remote Redis Cache Store", "error", err)
		return err
	}

	slog.Info("Successfully connected to Redis Cache Layer", "address", addr)
	return nil
}

func GetRedis() *redis.Client {
	return Client
}
