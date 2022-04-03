package helper

import (
	"api/config"
	"reflect"
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
		client.Del(context.Background(), fmt.Sprintf("%s:%s", prefix, suffix))
		return
	}

	// Encode json to strO
	if str, err := json.Marshal(model); err == nil {
		client.Set(context.Background(), fmt.Sprintf("%s:%s", prefix, suffix), str, time.Duration(config.ENV.LiveTime)*time.Second)
	}
}

func Getcache(db *gorm.DB, client *redis.Client, prefix, suffix string, model interface{}) (error, int64) {
	ctx := context.Background()

	// Check if cache have requested data
	var key = fmt.Sprintf("%s:%s", prefix, suffix)
	if data, err := client.Get(ctx, key).Result(); err == nil {
		if !strings.HasPrefix(data, "[") && reflect.ValueOf(model).Elem().Kind() == reflect.Slice {
			data = fmt.Sprintf("[%s]", data)
		} else if strings.HasPrefix(data, "[") && reflect.ValueOf(model).Elem().Kind() != reflect.Slice {
			// NOTE: Not the BEST solution
			data = strings.Trim(data, "[]")
		}

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

func Popcache(client *redis.Client, prefix, suffix string, models interface{}) error {
	ctx := context.Background()

	data, err := client.Get(ctx, fmt.Sprintf("%s:%s", prefix, suffix)).Result()
	if err != nil {
		return err
	}

	json.Unmarshal([]byte(data), &models)
	if reflect.ValueOf(models).Elem().Kind() != reflect.Slice {
		return fmt.Errorf("Slice was expected")
	}

	items := reflect.ValueOf(models).Elem()
	if items.Len() > 1 {
		Precache(client, prefix, suffix, items.Slice(1, items.Len()))
	} else {
		Precache(client, prefix, suffix, nil)
	}

	return nil

}

func Recache(client *redis.Client, prefix, suffix string, revalue func(string, string) interface{}) error {
	ctx := context.Background()
	iter := client.Scan(ctx, 0, fmt.Sprintf("%s:%s", prefix, suffix), 0).Iterator()

	for iter.Next(ctx) {
		var key = strings.Replace(iter.Val(), prefix+":", "", 1)

		data, _ := client.Get(ctx, iter.Val()).Result()
		Precache(client, prefix, key, revalue(data, key))
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
		client.Del(ctx, iter.Val())
	}

	if err := iter.Err(); err != nil {
		return fmt.Errorf("Error happed while setting interating through keys: %v", err)
	}

	return nil
}
