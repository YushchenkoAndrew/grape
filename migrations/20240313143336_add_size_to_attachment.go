package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upAddSizeToAttachment, downAddSizeToAttachment)
}

func upAddSizeToAttachment(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is applied.
	if _, err := tx.Exec(`ALTER TABLE attachments ADD COLUMN size bigint NOT NULL;`); err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

func downAddSizeToAttachment(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	return nil
}
