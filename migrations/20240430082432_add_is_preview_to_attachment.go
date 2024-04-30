package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upAddIsPreviewToAttachment, downAddIsPreviewToAttachment)
}

func upAddIsPreviewToAttachment(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is applied.
	if _, err := tx.Exec(`ALTER TABLE attachments ADD COLUMN "preview" boolean NOT NULL DEFAULT false;`); err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

func downAddIsPreviewToAttachment(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	return nil
}
