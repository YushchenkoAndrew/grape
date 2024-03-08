package client

import (
	"context"
	"fmt"
	"grape/src/common/config"

	"github.com/go-redis/redis/v8"
)

func ConnRedis(cfg *config.Config) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
		Username: cfg.Redis.User,
		Password: cfg.Redis.Pass,
	})

	ctx := context.Background()
	if _, err := client.Ping(ctx).Result(); err != nil {
		panic("failed on redis connection")
	}

	return client
}
