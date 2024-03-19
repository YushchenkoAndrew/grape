package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upCreateProjects, downCreateProjects)
}

func upCreateProjects(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is applied.
	_, err := tx.Exec(`
	CREATE TABLE IF NOT EXISTS projects (
		id bigserial PRIMARY KEY NOT NULL,
		uuid character varying NOT NULL,
		created_at timestamp(6) without time zone NOT NULL,
		updated_at timestamp(6) without time zone NOT NULL,
		organization_id bigint NOT NULL REFERENCES organizations(id),
		name character varying NOT NULL,
		description character varying DEFAULT '',
		type integer NOT NULL DEFAULT 0,
		status integer NOT NULL DEFAULT 0,
		footer character varying,
		"order" integer NOT NULL DEFAULT 0,
		owner_id bigint NOT NULL REFERENCES users(id)
	);
	`)

	if err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

func downCreateProjects(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	return nil
}
