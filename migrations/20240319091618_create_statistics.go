package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upCreateStatistics, downCreateStatistics)
}

func upCreateStatistics(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is applied.
	_, err := tx.Exec(`
	CREATE TABLE IF NOT EXISTS statistics (
		id bigserial PRIMARY KEY NOT NULL,
		uuid character varying NOT NULL,
		created_at timestamp(6) without time zone NOT NULL,
		updated_at timestamp(6) without time zone NOT NULL,
		views integer NOT NULL DEFAULT 0,
		clicks integer NOT NULL DEFAULT 0,
		media integer NOT NULL DEFAULT 0
	);
	`)

	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func downCreateStatistics(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	return nil
}
