package migrations

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/pressly/goose/v3"
	"golang.org/x/crypto/bcrypt"
)

func init() {
	goose.AddMigrationContext(upInsertDefaultUser, downInsertDefaultUser)
}

func upInsertDefaultUser(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is applied.
	hash, err := bcrypt.GenerateFromPassword([]byte(cfg.User.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
	INSERT INTO users(uuid, created_at, updated_at, organization_id, name, password)
		SELECT $1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, o.id, $2, $3
		FROM organizations o LIMIT 1;
	`, uuid.New().String(), cfg.User.Name, string(hash))

	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func downInsertDefaultUser(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	return nil
}
