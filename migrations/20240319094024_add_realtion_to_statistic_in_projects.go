package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upAddRealtionToStatisticInProjects, downAddRealtionToStatisticInProjects)
}

func upAddRealtionToStatisticInProjects(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is applied.
	if _, err := tx.Exec(`DELETE FROM projects WHERE true;`); err != nil {
		tx.Rollback()
		return err
	}

	if _, err := tx.Exec(`ALTER TABLE projects ADD COLUMN statistic_id bigint NOT NULL REFERENCES statistics(id);`); err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

func downAddRealtionToStatisticInProjects(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	return nil
}
