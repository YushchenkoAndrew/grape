package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upMakeLinksPolymorphic, downMakeLinksPolymorphic)
}

func upMakeLinksPolymorphic(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is applied.
	if _, err := tx.Exec(`ALTER TABLE links DROP COLUMN project_id;`); err != nil {
		tx.Rollback()
		return err
	}

	if _, err := tx.Exec(`ALTER TABLE links ADD COLUMN linkable_id bigint NOT NULL;`); err != nil {
		tx.Rollback()
		return err
	}

	if _, err := tx.Exec(`ALTER TABLE links ADD COLUMN linkable_type character varying NOT NULL;`); err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func downMakeLinksPolymorphic(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	return nil
}
