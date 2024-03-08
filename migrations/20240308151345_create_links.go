package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upCreateLinks, downCreateLinks)
}

func upCreateLinks(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is applied.

	_, err := tx.Exec(`
	CREATE TABLE IF NOT EXISTS links (
		id bigserial PRIMARY KEY NOT NULL,
		uuid character varying NOT NULL,
		created_at timestamp(6) without time zone NOT NULL,
		updated_at timestamp(6) without time zone NOT NULL,
		name character varying NOT NULL,
		link character varying NOT NULL,
		type character varying NOT NULL,
		project_id bigint REFERENCES projects(id)
	);
	`)

	if err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

func downCreateLinks(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	return nil
}
