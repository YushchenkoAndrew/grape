package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upAddRelationToSvgPatternsAndColorPalettesInProjects, downAddRelationToSvgPatternsAndColorPalettesInProjects)
}

func upAddRelationToSvgPatternsAndColorPalettesInProjects(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is applied.
	if _, err := tx.Exec(`DELETE FROM projects WHERE true;`); err != nil {
		tx.Rollback()
		return err
	}

	if _, err := tx.Exec(`ALTER TABLE projects ADD COLUMN palette_id bigint NOT NULL REFERENCES palettes(id);`); err != nil {
		tx.Rollback()
		return err
	}

	if _, err := tx.Exec(`ALTER TABLE projects ADD COLUMN pattern_id bigint NOT NULL REFERENCES patterns(id);`); err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

func downAddRelationToSvgPatternsAndColorPalettesInProjects(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	return nil
}
