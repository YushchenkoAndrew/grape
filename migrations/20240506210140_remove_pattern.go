package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upRemovePattern, downRemovePattern)
}

func upRemovePattern(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is applied.
	if _, err := tx.Exec(`ALTER TABLE projects DROP COLUMN pattern_id;`); err != nil {
		tx.Rollback()
		return err
	}

	if _, err := tx.Exec(`DROP TABLE patterns;`); err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func downRemovePattern(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	return nil
}
