package database

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
)

func RunMigrations(db *sql.DB, migrationsPath, direction string, args ...string) error {
	instance, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("create postgres driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(migrationsPath, "postgres", instance)
	if err != nil {
		return fmt.Errorf("create migrate instance: %w", err)
	}

	switch direction {
	case "up":
		err = m.Up()
	case "down":
		err = m.Steps(-1)
	case "force":
		if len(os.Args) < 3 {
			return fmt.Errorf("force requires a version argument")
		}
		version, err := strconv.Atoi(os.Args[2])
		if err != nil {
			return fmt.Errorf("invalid version: %w", err)
		}
		err = m.Force(version)
	default:
		return fmt.Errorf("invalid direction %q: must be 'up' or 'down'", direction)
	}

	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("migration %s: %w", direction, err)
	}

	return nil
}
