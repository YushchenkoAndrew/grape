package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upCreateContext, downCreateContext)
}

func upCreateContext(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is applied.
	_, err := tx.Exec(`
	CREATE TABLE IF NOT EXISTS contexts (
		id bigserial PRIMARY KEY NOT NULL,
		uuid character varying NOT NULL,
		created_at timestamp(6) without time zone NOT NULL,
		updated_at timestamp(6) without time zone NOT NULL,
		name character varying NOT NULL,
		"order" integer NOT NULL DEFAULT 1,
		contextable_id bigint NOT NULL,
		contextable_type character varying  NOT NULL
	);
	`)

	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func downCreateContext(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	return nil
}
