package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upAddOrderToLink, downAddOrderToLink)
}

func upAddOrderToLink(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is applied.
	if _, err := tx.Exec(`ALTER TABLE links ADD COLUMN "order" integer NOT NULL DEFAULT 1;`); err != nil {
		tx.Rollback()
		return err
	}

	_, err := tx.Exec(`
		UPDATE links SET "order"=a.rank
		FROM (
		    SELECT RANK() OVER (PARTITION BY linkable_id, linkable_type ORDER BY created_at) rank, id FROM links
		) a
		WHERE links.id = a.id;
	`)

	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func downAddOrderToLink(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	return nil
}
