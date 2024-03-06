package client

import (
	"grape/src/common/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnDB() *gorm.DB {
	var db, err = gorm.Open(postgres.Open(
		"host="+config.ENV.DBHost+
			" user="+config.ENV.DBUser+
			" password="+config.ENV.DBPass+
			" port="+config.ENV.DBPort+
			" dbname="+config.ENV.DBName), &gorm.Config{})

	if err != nil {
		SendLogs(&Message{
			Stat:    "ERR",
			Name:    "grape",
			File:    "/db/db.go",
			Message: "Bruhhh, did you even start the Postgres ???",
			Desc:    err.Error(),
		})
		panic("Failed on db connection")
	}

	return db
}
