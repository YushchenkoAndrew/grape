package helper

import (
	"api/config"
	"crypto/md5"

	// "api/db"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

func Precache(client *redis.Client, prefix, suffix string, model interface{}) {
	hasher := md5.New()
	hasher.Write([]byte(suffix))
	var key = fmt.Sprintf("%s:%s", prefix, hex.EncodeToString(hasher.Sum(nil)))

	// Encode json to strO
	if str, err := json.Marshal(model); err == nil {
		go client.Set(context.Background(), key, str, time.Duration(config.ENV.LiveTime)*time.Second)
	}
}

func Getcache(db *gorm.DB, client *redis.Client, prefix, suffix string, model interface{}) (error, int64) {
	hasher := md5.New()
	hasher.Write([]byte(suffix))
	var key = fmt.Sprintf("%s:%s", prefix, hex.EncodeToString(hasher.Sum(nil)))

	ctx := context.Background()

	// Check if cache have requested data
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

func Delcache(client *redis.Client, prefix, suffix string, model interface{}) {
	hasher := md5.New()
	hasher.Write([]byte(suffix))

	go client.Del(context.Background(), fmt.Sprintf("%s:%s", prefix, hex.EncodeToString(hasher.Sum(nil))))
}
