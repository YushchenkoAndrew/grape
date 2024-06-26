package client

import (
	"fmt"
	_ "grape/migrations"
	"grape/src/common/config"
	c "grape/src/common/config"
	"os"
	"path/filepath"
	"time"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func ConnPsql(cfg *config.Config) *gorm.DB {
	dsn := fmt.Sprintf("host=%s user=%s password=%s port=%d dbname=%s sslmode=disable", cfg.Psql.Host, cfg.Psql.Username, cfg.Psql.Password, cfg.Psql.Port, cfg.Psql.Name)

	var config gorm.Config
	if cfg.Psql.Logger {
		config.Logger = logger.Default.LogMode(logger.Info)
	}

	timeout, _ := time.ParseDuration(cfg.Psql.Options.HealthTimeout)

	var retry func(count int) (*gorm.DB, error)
	retry = func(count int) (*gorm.DB, error) {
		db, err := gorm.Open(postgres.Open(dsn), &config)
		if err != nil && count < cfg.Psql.Options.HealthRetries {
			time.Sleep(timeout)
			return retry(count + 1)
		}

		return db, err
	}

	db, err := retry(0)

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
		if value, ok := os.LookupEnv(c.CONFIG_ARG); ok {
			path = filepath.Join(value, cfg.Server.Migrations)
		} else {
			path = filepath.Join(value, cfg.Server.Migrations)
		}

		if err = goose.Up(db, path); err != nil {
			panic(fmt.Errorf("goose: failed to migrate: %w", err))
		}
	}()

	return db
}
