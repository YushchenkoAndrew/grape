package client

import (
	"fmt"
	_ "grape/migrations"
	"grape/src/common/config"
	"os"
	"path/filepath"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnPsql(cfg *config.Config) *gorm.DB {
	dsn := fmt.Sprintf("host=%s user=%s password=%s port=%d dbname=%s sslmode=disable", cfg.Psql.Host, cfg.Psql.User, cfg.Psql.Pass, cfg.Psql.Port, cfg.Psql.Name)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		// TODO: Think about better solution for that !
		// SendLogs(&Message{
		// 	Stat:    "ERR",
		// 	Name:    "grape",
		// 	File:    "/psql/psql.go",
		// 	Message: "Bruhhh, did you even start the Postgres ???",
		// 	Desc:    err.Error(),
		// })
		panic("failed on psql connection")
	}

	defer func() {
		db, err := goose.OpenDBWithDriver("postgres", dsn)
		if err != nil {
			panic(fmt.Errorf("goose: failed to open psql: %w", err))
		}

		if err = db.Ping(); err != nil {
			panic(fmt.Errorf("goose: failed to connect to psql: %w", err))
		}

		defer func() {
			if err := db.Close(); err != nil {
				panic(fmt.Errorf("goose: failed to close psql: %w", err))
			}
		}()

		var path string
		if value, ok := os.LookupEnv(config.CONFIG_ARG); ok {
			path = filepath.Join(value, cfg.Server.Migration)
		} else {
			path = filepath.Join(value, cfg.Server.Migration)
		}

		if err = goose.Up(db, path); err != nil {
			panic(fmt.Errorf("goose: failed to migrate: %w", err))
		}
	}()

	return db
}
