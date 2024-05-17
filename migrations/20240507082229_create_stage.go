package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upCreateStage, downCreateStage)
}

func upCreateStage(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is applied.
	_, err := tx.Exec(`
	CREATE TABLE IF NOT EXISTS stages (
		id bigserial PRIMARY KEY NOT NULL,
		uuid character varying NOT NULL,
		created_at timestamp(6) without time zone NOT NULL,
		updated_at timestamp(6) without time zone NOT NULL,
		organization_id bigint NOT NULL REFERENCES organizations(id),
		name character varying NOT NULL,
		status integer NOT NULL DEFAULT 1,
		"order" integer NOT NULL DEFAULT 1
	);
	`)

	if err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

func downCreateStage(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	return nil
}
