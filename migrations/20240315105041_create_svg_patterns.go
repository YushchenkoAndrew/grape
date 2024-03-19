package migrations

import (
	"bufio"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"grape/src/common/config"
	"grape/src/style/types"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upCreateSvgPatterns, downCreateSvgPatterns)
}

func upCreateSvgPatterns(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is applied.
	_, err := tx.Exec(`
	CREATE TABLE IF NOT EXISTS svg_patterns (
		id bigserial PRIMARY KEY NOT NULL,
		uuid character varying NOT NULL,
		created_at timestamp(6) without time zone NOT NULL,
		updated_at timestamp(6) without time zone NOT NULL,
		organization_id bigint REFERENCES organizations(id),
		mode integer NOT NULL DEFAULT 0,
		colors integer NOT NULL,
		options jsonb NOT NULL DEFAULT '{"max_stroke": 0, "max_scale": 0, "max_spacing_x": 0, "max_spacing_y": 0}'::jsonb,
		width float NOT NULL DEFAULT 0,
		height float NOT NULL DEFAULT 0,
		path text NOT NULL DEFAULT ''
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

	file, err := os.Open(filepath.Join(path, "tmp", "Patterns.csv"))
	if err != nil {
		tx.Rollback()
		return err
	}

	defer file.Close()

	const (
		MODE          = 1
		COLORS        = 2
		MAX_STROKE    = 3
		MAX_SCALE     = 4
		MAX_SPACING_X = 5
		MAX_SPACING_Y = 6
		WIDTH         = 7
		HEIGHT        = 8
		PATH          = 9
	)

	header := true
	scanner := bufio.NewScanner(file)
	insert := func(chunk []string, args []interface{}) error {
		_, err := tx.Exec(fmt.Sprintf(`
			INSERT INTO svg_patterns(uuid, created_at, updated_at, organization_id, mode, colors, options, width, height, path)
				VALUES %s;
			`, strings.Join(chunk, ", ")), args...)

		return err
	}

	var args []interface{}
	var chunk []string
	for scanner.Scan() {
		if header {
			scanner.Text()
			header = false
			continue
		}

		s := strings.Split(scanner.Text(), ",")

		max_stroke, _ := strconv.ParseFloat(s[MAX_STROKE], 32)
		max_scale, _ := strconv.ParseInt(s[MAX_SCALE], 10, 0)
		max_spacing_x, _ := strconv.ParseFloat(s[MAX_SPACING_X], 32)
		max_spacing_y, _ := strconv.ParseFloat(s[MAX_SPACING_Y], 32)

		data := &types.ColorPaletteOptionsType{MaxStroke: float32(max_stroke), MaxScale: int(max_scale), MaxSpacingX: float32(max_spacing_x), MaxSpacingY: float32(max_spacing_y)}
		json, _ := json.Marshal(data)

		chunk = append(chunk, fmt.Sprintf(`('%s', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 1, %d, %s, '%s'::jsonb, %s, %s, $%d)`, uuid.New().String(), int(types.Fill.Value(s[MODE])), s[COLORS], string(json), s[WIDTH], s[HEIGHT], len(args)+1))
		args = append(args, s[PATH])

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

func downCreateSvgPatterns(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	return nil
}
