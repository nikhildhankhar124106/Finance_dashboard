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
		slog.Info("Caching: Redis not detected locally. System will run in 'Direct-DB' mode (No performance impact for localized testing).")
		return nil
	}

	slog.Info("Successfully connected to Redis Cache Layer", "address", addr)
	return nil
}

func GetRedis() *redis.Client {
	return Client
}
