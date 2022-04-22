package client

import (
	"api/config"
	"api/interfaces"
	"api/logs"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnDB(tables []interfaces.Table) *gorm.DB {
	var db, err = gorm.Open(postgres.Open(
		"host="+config.ENV.DBHost+
			" user="+config.ENV.DBUser+
			" password="+config.ENV.DBPass+
			" port="+config.ENV.DBPort+
			" dbname="+config.ENV.DBName), &gorm.Config{})

	if err != nil {
		logs.SendLogs(&logs.Message{
			Stat:    "ERR",
			Name:    "API",
			File:    "/db/db.go",
			Message: "Bruhhh, did you even start the Postgres ???",
			Desc:    err.Error(),
		})
		panic("Failed on db connection")
	}

	for _, table := range tables {
		if err := table.Migrate(db, config.ENV.ForceMigrate); err != nil {
			panic(err)
		}
	}

	return db
}
