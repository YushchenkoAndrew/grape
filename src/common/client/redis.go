package client

import (
	"context"
	"grape/src/common/config"

	"github.com/go-redis/redis/v8"
)

func ConnRedis(cfg *config.Config) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Address,
		DB:       cfg.Redis.DB,
		Username: cfg.Redis.Username,
		Password: cfg.Redis.Password,
	})

	ctx := context.Background()
	if _, err := client.Ping(ctx).Result(); err != nil {
		panic("failed on redis connection")
	}

	return client
}
