package db

import (
	"context"

	"api/config"
	"api/logs"
	"api/models"

	"github.com/go-redis/redis/v8"
)

func ConnectToRedis() *redis.Client {
	var client = redis.NewClient(&redis.Options{
		Addr:     config.ENV.RedisHost + ":" + config.ENV.RedisPort,
		Password: config.ENV.RedisPass,
	})

	ctx := context.Background()
	if _, err := client.Ping(ctx).Result(); err != nil {
		logs.SendLogs(&models.LogMessage{
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

// func FlushValue(client *redis.Client, key string) {
// 	ctx := context.Background()
// 	iter := client.Scan(ctx, 0, fmt.Sprintf("%s:*", key), 0).Iterator()

// 	for iter.Next(ctx) {
// 		go client.Del(ctx, iter.Val())
// 	}

// 	if err := iter.Err(); err != nil {
// 		fmt.Println("[Redis] Error happed while setting interating through keys")
// 	}
// }
