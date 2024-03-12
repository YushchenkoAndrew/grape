package migrations

import (
	"bufio"
	"context"
	"database/sql"
	"fmt"
	"grape/src/common/config"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upCreateLocations, downCreateLocations)
}

func upCreateLocations(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is applied.
	_, err := tx.Exec(`
	CREATE TABLE IF NOT EXISTS locations (
		id bigserial PRIMARY KEY NOT NULL,
		uuid character varying NOT NULL,
		created_at timestamp(6) without time zone NOT NULL,
		updated_at timestamp(6) without time zone NOT NULL,
		locale_code character varying NOT NULL,
		continent_code character varying(2) NOT NULL,
		continent_name character varying NOT NULL,
		country_iso_code character varying(2) NOT NULL,
		country_name character varying NOT NULL,
		geoname_id bigint NOT NULL
	);
	`)

	if err != nil {
		tx.Rollback()
		return err
	}

	cfg := config.NewConfig("configs/", "config", "yaml")

	var path string
	if value, ok := os.LookupEnv(config.CONFIG_ARG); ok {
		path = filepath.Join(value, cfg.Server.Migration)
	} else {
		path = filepath.Join(value, cfg.Server.Migration)
	}

	file, err := os.Open(filepath.Join(path, "tmp", "GeoLite2-Country-Locations-en.csv"))
	if err != nil {
		tx.Rollback()
		return err
	}

	defer file.Close()

	const (
		GEONAME_ID       = 0
		LOCALE_CODE      = 1
		CONTINENT_CODE   = 2
		CONTINENT_NAME   = 3
		COUNTRY_ISO_CODE = 4
		COUNTRY_NAME     = 5
	)

	header := true
	scanner := bufio.NewScanner(file)
	insert := func(chunk []string) error {
		_, err := tx.Exec(fmt.Sprintf(`
			INSERT INTO locations(uuid, created_at, updated_at, locale_code, continent_code, continent_name, country_iso_code, country_name, geoname_id)
				VALUES %s;
			`, strings.Join(chunk, ", ")))

		return err
	}

	var chunk []string
	for scanner.Scan() {
		if header {
			scanner.Text()
			header = false
			continue
		}

		s := strings.Split(scanner.Text(), ",")
		id, _ := strconv.ParseInt(s[GEONAME_ID], 10, 64)
		chunk = append(chunk, fmt.Sprintf(`('%s', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, '%s', '%s', '%s', '%s', '%s', %d)`, uuid.New().String(), s[LOCALE_CODE], s[CONTINENT_CODE], s[CONTINENT_NAME], s[COUNTRY_ISO_CODE], s[COUNTRY_NAME], id))

		if len(chunk) == 500 {
			if err := insert(chunk); err != nil {
				tx.Rollback()
				return err
			}

			chunk = []string{}
		}
	}

	if err := insert(chunk); err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func downCreateLocations(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	return nil
}
