package db

import (
	"api/config"
	"api/logs"
	"api/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectToDB() *gorm.DB {
	var db, err = gorm.Open(postgres.Open(
		"host="+config.ENV.DBHost+
			" user="+config.ENV.DBUser+
			" password="+config.ENV.DBPass+
			" port="+config.ENV.DBPort+
			" dbname="+config.ENV.DBName), &gorm.Config{})

	if err != nil {
		logs.SendLogs(&models.LogMessage{
			Stat:    "ERR",
			Name:    "API",
			File:    "/db/db.go",
			Message: "Bruhhh, did you even start the Postgres ???",
			Desc:    err.Error(),
		})
		panic("Failed on db connection")
	}

	return db
}
