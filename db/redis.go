package db

import (
	"context"
	"fmt"

	"api/config"
	"api/logs"
	"api/models"

	"github.com/go-redis/redis/v8"
)

var Redis *redis.Client

func ConnectToRedis() {
	Redis = redis.NewClient(&redis.Options{
		Addr:     config.ENV.RedisHost + ":" + config.ENV.RedisPort,
		Password: config.ENV.RedisPass,
	})

	ctx := context.Background()
	if _, err := Redis.Ping(ctx).Result(); err != nil {
		logs.SendLogs(&models.LogMessage{
			Stat:    "ERR",
			Name:    "API",
			File:    "/db/redis.go",
			Message: "Bruhhh, did you even start the Redis ???",
			Desc:    err.Error(),
		})
		panic("Failed on Redis connection")
	}
}

func RedisInitDefault() {
	var nInfo int64
	var nWorld int64
	var nFile int64

	DB.Model(&models.Info{}).Count(&nInfo)
	DB.Model(&models.World{}).Count(&nWorld)
	DB.Model(&models.File{}).Count(&nFile)

	// FIXME: ERROR log format
	ctx := context.Background()
	var SetVar = func(ctx *context.Context, param string, value interface{}) {
		if err := Redis.Set(*ctx, param, value, 0).Err(); err != nil {
			fmt.Println("ERROR:")
		}
	}

	SetVar(&ctx, "nInfo", nInfo)
	SetVar(&ctx, "nWorld", nWorld)
	SetVar(&ctx, "nFile", nFile)
	SetVar(&ctx, "Mutex", 1)
}
