package helper

import (
	"api/config"
	"strings"

	// "api/db"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

func Precache(client *redis.Client, prefix, suffix string, model interface{}) {
	if model == nil {
		go client.Del(context.Background(), fmt.Sprintf("%s:%s", prefix, suffix))
		return
	}

	// Encode json to strO
	if str, err := json.Marshal(model); err == nil {
		go client.Set(context.Background(), fmt.Sprintf("%s:%s", prefix, suffix), str, time.Duration(config.ENV.LiveTime)*time.Second)
	}
}

func Getcache(db *gorm.DB, client *redis.Client, prefix, suffix string, model interface{}) (error, int64) {
	ctx := context.Background()

	// Check if cache have requested data
	var key = fmt.Sprintf("%s:%s", prefix, suffix)
	if data, err := client.Get(ctx, key).Result(); err == nil {
		json.Unmarshal([]byte(data), model)
		go client.Expire(ctx, key, time.Duration(config.ENV.LiveTime)*time.Second)
	} else {
		if result := db.Find(model); result.Error != nil || result.RowsAffected == 0 {
			return result.Error, result.RowsAffected
		}

		Precache(client, prefix, suffix, model)
	}

	return nil, -1
}

// func GetcacheAll(db *gorm.DB, client *redis.Client, prefix, suffix string, model interface{}) error {
// 	hasher := md5.New()
// 	hasher.Write([]byte(suffix))
// 	var key = fmt.Sprintf("%s:%s", prefix, hex.EncodeToString(hasher.Sum(nil)))

// 	ctx := context.Background()

// 	// Check if cache have requested data
// 	if data, err := client.Get(ctx, key).Result(); err == nil {
// 		json.Unmarshal([]byte(data), model)
// 		go client.Expire(ctx, key, time.Duration(config.ENV.LiveTime)*time.Second)
// 	} else {
// 		if result := db.Find(model); result.Error != nil {
// 			return fmt.Errorf("Server side error: Something went wrong - %v", result.Error)
// 		}

// 		Precache(client, prefix, suffix, model)
// 	}

// 	return nil
// }

func Recache(client *redis.Client, prefix, suffix string, revalue func(string) interface{}) error {
	ctx := context.Background()
	iter := client.Scan(ctx, 0, fmt.Sprintf("%s:%s", prefix, suffix), 0).Iterator()

	for iter.Next(ctx) {
		data, _ := client.Get(ctx, iter.Val()).Result()
		go Precache(client, prefix, strings.TrimLeft(iter.Val(), prefix+":"), revalue(data))
	}

	if err := iter.Err(); err != nil {
		return fmt.Errorf("Error happed while setting interating through keys: %v", err)
	}

	return nil
}

func Delcache(client *redis.Client, prefix, suffix string) error {
	ctx := context.Background()
	iter := client.Scan(ctx, 0, fmt.Sprintf("%s:%s", prefix, suffix), 0).Iterator()

	for iter.Next(ctx) {
		go client.Del(ctx, iter.Val())
	}

	if err := iter.Err(); err != nil {
		return fmt.Errorf("Error happed while setting interating through keys: %v", err)
	}

	return nil
}
