package migrations

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upInsertDefaultStage, downInsertDefaultStage)
}

func upInsertDefaultStage(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is applied.
	var count int
	if err := tx.QueryRow("SELECT COUNT(*) FROM stages").Scan(&count); err != nil {
		tx.Rollback()
		return err
	}

	if count > 0 {
		return nil
	}

	_, err := tx.Exec(`
	INSERT INTO stages(uuid, created_at, updated_at, organization_id, name)
		VALUES ($1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 1, 'default')
	`, uuid.New().String())

	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec(`
	INSERT INTO tasks(uuid, created_at, updated_at, organization_id, name, owner_id, stage_id)
		VALUES ($1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 1, 'default', 1, 1)
	`, uuid.New().String())

	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func downInsertDefaultStage(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	return nil
}
