package migrations

import (
	"context"
	"database/sql"
	"grape/src/common/config"

	"github.com/pressly/goose/v3"
)

var cfg *config.DatabaseConfig

func init() {
	cfg = config.NewDatabaseConfig("configs/", "database", "yaml")
	goose.AddMigrationContext(upCreateOrganizations, downCreateOrganizations)
}

func upCreateOrganizations(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is applied.
	_, err := tx.Exec(`
	CREATE TABLE IF NOT EXISTS organizations (
		id bigserial PRIMARY KEY NOT NULL,
		uuid character varying NOT NULL,
		created_at timestamp(6) without time zone NOT NULL,
		updated_at timestamp(6) without time zone NOT NULL,
		name character varying NOT NULL
	);
	`)

	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func downCreateOrganizations(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	return nil
}
