package migrations

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upInsertDefaultProject, downInsertDefaultProject)
}

func upInsertDefaultProject(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is applied.
	var count int
	if err := tx.QueryRow("SELECT COUNT(*) FROM projects").Scan(&count); err != nil {
		tx.Rollback()
		return err
	}

	if count > 0 {
		return nil
	}

	_, err := tx.Exec(`
	INSERT INTO statistics(uuid, created_at, updated_at)
		VALUES ($1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	`, uuid.New().String())

	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec(`
	INSERT INTO projects(uuid, created_at, updated_at, organization_id, name, type, status, "order", owner_id, statistic_id)
		VALUES ($1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 1, 'default', 1, 1, 1, 1, 1)
	`, uuid.New().String())

	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func downInsertDefaultProject(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	return nil
}
