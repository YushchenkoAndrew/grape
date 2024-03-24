package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upCreateAttachments, downCreateAttachments)
}

func upCreateAttachments(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is applied.

	_, err := tx.Exec(`
	CREATE TABLE IF NOT EXISTS attachments (
		id bigserial PRIMARY KEY NOT NULL,
		uuid character varying NOT NULL,
		created_at timestamp(6) without time zone NOT NULL,
		updated_at timestamp(6) without time zone NOT NULL,
		name character varying NOT NULL,
		path character varying NOT NULL DEFAULT '/',
		type character varying NOT NULL,
		attachable_id bigint NOT NULL,
		attachable_type character varying  NOT NULL
	);
	`)

	if err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

func downCreateAttachments(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	return nil
}
