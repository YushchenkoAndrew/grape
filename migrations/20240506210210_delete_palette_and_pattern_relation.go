package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upDeletePaletteRelation, downDeletePaletteRelation)
}

func upDeletePaletteRelation(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is applied.
	if _, err := tx.Exec(`ALTER TABLE projects DROP COLUMN palette_id;`); err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func downDeletePaletteRelation(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	return nil
}
