package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upDeleteRedundantGeonameId, downDeleteRedundantGeonameId)
}

func upDeleteRedundantGeonameId(ctx context.Context, tx *sql.Tx) error {
	if _, err := tx.Exec(`ALTER TABLE networks DROP COLUMN geoname_id;`); err != nil {
		tx.Rollback()
		return err
	}

	if _, err := tx.Exec(`ALTER TABLE locations DROP COLUMN geoname_id;`); err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

func downDeleteRedundantGeonameId(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	return nil
}
