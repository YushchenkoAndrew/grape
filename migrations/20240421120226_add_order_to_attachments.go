package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upAddOrderToAttachments, downAddOrderToAttachments)
}

func upAddOrderToAttachments(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is applied.
	if _, err := tx.Exec(`ALTER TABLE attachments ADD COLUMN "order" integer NOT NULL DEFAULT 1;`); err != nil {
		tx.Rollback()
		return err
	}

	_, err := tx.Exec(`
		UPDATE attachments SET "order"=a.rank
		FROM (
		    SELECT RANK() OVER (PARTITION BY attachable_id, attachable_type ORDER BY created_at) rank, id FROM attachments
		) a
		WHERE attachments.id = a.id;
	`)

	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func downAddOrderToAttachments(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	return nil
}
