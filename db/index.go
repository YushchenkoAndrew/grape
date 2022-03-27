package db

import (
	"api/config"
	"api/interfaces"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

func Init(tables []interfaces.Table) (*gorm.DB, *redis.Client) {
	db, client := ConnectToDB(), ConnectToRedis()

	for _, table := range tables {
		table.Migrate(db, config.ENV.ForceMigrate)

		if err := table.Redis(db, client); err != nil {
			panic(err)
		}
	}

	return db, client
}
