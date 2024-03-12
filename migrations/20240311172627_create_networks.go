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

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upCreateNetworks, downCreateNetworks)
}

func upCreateNetworks(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is applied.
	_, err := tx.Exec(`
	CREATE TABLE IF NOT EXISTS networks (
		id bigserial PRIMARY KEY NOT NULL,
		network cidr NOT NULL,
		location_id bigint NOT NULL,
		geoname_id bigint NOT NULL
	);
	`)

	if err != nil {
		tx.Rollback()
		return err
	}

	if _, err = tx.Exec(`CREATE INDEX ON networks USING gist (network inet_ops);`); err != nil {
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

	file, err := os.Open(filepath.Join(path, "tmp", "GeoLite2-Country-Blocks.csv"))
	if err != nil {
		tx.Rollback()
		return err
	}

	defer file.Close()

	const (
		NETWORK    = 0
		GEONAME_ID = 1
	)

	header := true
	scanner := bufio.NewScanner(file)
	insert := func(chunk []string) error {
		_, err := tx.Exec(fmt.Sprintf(`
			INSERT INTO networks(network, location_id, geoname_id)
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
		geoname_id, _ := strconv.ParseInt(s[GEONAME_ID], 10, 64)
		chunk = append(chunk, fmt.Sprintf(`('%s', %d, %d)`, s[NETWORK], 0, geoname_id))

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

	if _, err = tx.Exec(`UPDATE networks n SET location_id = l.id FROM locations l WHERE l.geoname_id = n.geoname_id;`); err != nil {
		tx.Rollback()
		return err
	}

	if _, err = tx.Exec(`DELETE FROM networks WHERE location_id = 0;`); err != nil {
		tx.Rollback()
		return err
	}

	if _, err = tx.Exec(`ALTER TABLE networks ADD CONSTRAINT networks_location_id_fkey FOREIGN KEY (location_id) REFERENCES locations(id);`); err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func downCreateNetworks(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	return nil
}
