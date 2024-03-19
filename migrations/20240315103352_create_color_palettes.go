package migrations

import (
	"bufio"
	"context"
	"database/sql"
	"fmt"
	"grape/src/common/config"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"github.com/samber/lo"
)

func init() {
	goose.AddMigrationContext(upCreateColorPalettes, downCreateColorPalettes)
}

func upCreateColorPalettes(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is applied.
	_, err := tx.Exec(`
	CREATE TABLE IF NOT EXISTS color_palettes (
		id bigserial PRIMARY KEY NOT NULL,
		uuid character varying NOT NULL,
		created_at timestamp(6) without time zone NOT NULL,
		updated_at timestamp(6) without time zone NOT NULL,
		organization_id bigint NOT NULL REFERENCES organizations(id),
		colors character varying[] NOT NULL DEFAULT array[]::character varying[]
	);
	`)

	if err != nil {
		tx.Rollback()
		return err
	}

	cfg := config.NewConfig("configs/", "config", "yaml")

	var path string
	if value, ok := os.LookupEnv(config.CONFIG_ARG); ok {
		path = filepath.Join(value, cfg.Server.Migrations)
	} else {
		path = filepath.Join(value, cfg.Server.Migrations)
	}

	file, err := os.Open(filepath.Join(path, "tmp", "Colors.csv"))
	if err != nil {
		tx.Rollback()
		return err
	}

	defer file.Close()

	const (
		COLORS = 1
	)

	header := true
	scanner := bufio.NewScanner(file)
	insert := func(chunk []string, args []interface{}) error {
		_, err := tx.Exec(fmt.Sprintf(`
			INSERT INTO color_palettes(uuid, created_at, updated_at, organization_id, colors)
				VALUES %s;
			`, strings.Join(chunk, ", ")), args...)

		return err
	}

	var chunk []string
	var args []interface{}
	for scanner.Scan() {
		if header {
			scanner.Text()
			header = false
			continue
		}

		s := strings.Split(scanner.Text(), ",")
		colors := strings.Split(s[COLORS], "-")
		chunk = append(chunk, fmt.Sprintf(`('%s', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 1, $%d)`, uuid.New().String(), len(args)+1))
		args = append(args,
			pq.StringArray(lo.Map(colors, func(color string, _ int) string {
				return "#" + color
			}),
			))

		if len(chunk) == 500 {
			if err := insert(chunk, args); err != nil {
				tx.Rollback()
				return err
			}

			chunk = []string{}
			args = []interface{}{}
		}
	}

	if err := insert(chunk, args); err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func downCreateColorPalettes(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	return nil
}
