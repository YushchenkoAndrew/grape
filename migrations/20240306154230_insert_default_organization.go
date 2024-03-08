package migrations

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upInsertDefaultOrganizations, downInsertDefaultOrganizations)
}

func upInsertDefaultOrganizations(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is applied.
	_, err := tx.Exec(`
	INSERT INTO organizations(uuid, created_at, updated_at, name)
		VALUES($1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, $2);
	`, uuid.New().String(), cfg.Organization.Name)

	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func downInsertDefaultOrganizations(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	return nil
}
