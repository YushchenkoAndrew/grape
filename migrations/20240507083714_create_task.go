package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upCreateTask, downCreateTask)
}

func upCreateTask(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is applied.
	_, err := tx.Exec(`
	CREATE TABLE IF NOT EXISTS tasks (
		id bigserial PRIMARY KEY NOT NULL,
		uuid character varying NOT NULL,
		created_at timestamp(6) without time zone NOT NULL,
		updated_at timestamp(6) without time zone NOT NULL,
		organization_id bigint NOT NULL REFERENCES organizations(id),
		name character varying NOT NULL,
		description character varying DEFAULT '',
		status integer NOT NULL DEFAULT 1,
		"order" integer NOT NULL DEFAULT 1,
		owner_id bigint NOT NULL REFERENCES users(id),
		stage_id bigint NOT NULL REFERENCES stages(id)
	);
	`)

	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func downCreateTask(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	return nil
}
