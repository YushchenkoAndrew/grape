package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upAddHomeToAttachments, downAddHomeToAttachments)
}

func upAddHomeToAttachments(ctx context.Context, tx *sql.Tx) error {
	if _, err := tx.Exec(`DELETE FROM attachments WHERE true;`); err != nil {
		tx.Rollback()
		return err
	}

	// This code is executed when the migration is applied.
	if _, err := tx.Exec(`ALTER TABLE attachments ADD COLUMN home character varying NOT NULL;`); err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

func downAddHomeToAttachments(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	return nil
}
