package db

import (
	"api/config"
	"api/interfaces"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

func Init(tables []interfaces.Table) (*gorm.DB, *redis.Client) {
	db := ConnectToDB()
	for _, table := range tables {
		if err := table.Migrate(db, config.ENV.ForceMigrate); err != nil {
			panic(err)
		}
	}

	return db, ConnectToRedis()
}
