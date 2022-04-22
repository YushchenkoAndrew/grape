package client

import (
	"context"

	"api/config"
	"api/logs"

	"github.com/go-redis/redis/v8"
)

func ConnRedis() *redis.Client {
	var client = redis.NewClient(&redis.Options{
		Addr:     config.ENV.RedisHost + ":" + config.ENV.RedisPort,
		Password: config.ENV.RedisPass,
	})

	ctx := context.Background()
	if _, err := client.Ping(ctx).Result(); err != nil {
		logs.SendLogs(&logs.Message{
			Stat:    "ERR",
			Name:    "API",
			File:    "/db/redis.go",
			Message: "Bruhhh, did you even start the Redis ???",
			Desc:    err.Error(),
		})
		panic("Failed on Redis connection")
	}

	client.Set(ctx, "Mutex", 1, 0)
	return client
}
