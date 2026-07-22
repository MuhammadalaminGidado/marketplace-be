package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// RunMigrations runs migrations based on environment configuration
func RunMigrations(db *sql.DB, migrationsPath, direction string, args ...string) error {
	// Check if migrations should run
	if !shouldRunMigrations() {
		log.Println("⏭️ Skipping migrations (disabled by configuration)")
		return nil
	}

	log.Println("🔄 Starting database migrations...")

	instance, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("create postgres driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(migrationsPath, "postgres", instance)
	if err != nil {
		return fmt.Errorf("create migrate instance: %w", err)
	}

	// Determine direction - use passed direction or get from environment
	finalDirection := direction
	if finalDirection == "" {
		finalDirection = getMigrationDirection()
	}

	switch finalDirection {
	case "up":
		err = m.Up()
		if err != nil && !errors.Is(err, migrate.ErrNoChange) {
			return fmt.Errorf("migration up: %w", err)
		}
		if errors.Is(err, migrate.ErrNoChange) {
			log.Println("✅ No new migrations to apply")
		} else {
			log.Println("✅ UP migrations completed successfully")
		}

	case "down":
		steps := 1
		if len(args) > 0 {
			if s, err := strconv.Atoi(args[0]); err == nil && s > 0 {
				steps = s
			}
		}
		// Also check environment for steps
		if stepsEnv := os.Getenv("MIGRATION_STEPS"); stepsEnv != "" {
			if s, err := strconv.Atoi(stepsEnv); err == nil && s > 0 {
				steps = s
			}
		}

		log.Printf("⬇️ Running DOWN migrations (%d steps)...", steps)
		err = m.Steps(-steps)
		if err != nil && !errors.Is(err, migrate.ErrNoChange) {
			return fmt.Errorf("migration down: %w", err)
		}
		log.Println("✅ DOWN migrations completed successfully")

	case "force":
		version := 0
		if len(args) > 0 {
			if v, err := strconv.Atoi(args[0]); err == nil {
				version = v
			}
		}
		// Also check environment for force version
		if forceEnv := os.Getenv("MIGRATION_FORCE_VERSION"); forceEnv != "" {
			if v, err := strconv.Atoi(forceEnv); err == nil {
				version = v
			}
		}

		if version == 0 && len(args) == 0 && os.Getenv("MIGRATION_FORCE_VERSION") == "" {
			return fmt.Errorf("force requires a version argument or MIGRATION_FORCE_VERSION env var")
		}

		log.Printf("💪 Forcing migration version to %d...", version)
		err = m.Force(version)
		if err != nil {
			return fmt.Errorf("force version: %w", err)
		}
		log.Printf("✅ Version forced to %d", version)

	default:
		return fmt.Errorf("invalid direction %q: must be 'up', 'down', or 'force'", finalDirection)
	}

	// Log current version
	version, dirty, _ := m.Version()
	log.Printf("📊 Current migration version: %d (dirty: %v)", version, dirty)

	return nil
}

// shouldRunMigrations checks if migrations should be executed
func shouldRunMigrations() bool {
	// Check environment variable
	runMigrations := os.Getenv("RUN_MIGRATIONS")

	// If RUN_MIGRATIONS is not set, use environment-based default
	if runMigrations == "" {
		// Default: Run migrations in production/staging, skip in development
		env := os.Getenv("ENV")
		return env == "production" || env == "staging"
	}

	// Parse boolean value
	val, err := strconv.ParseBool(runMigrations)
	if err != nil {
		log.Printf("Warning: Invalid RUN_MIGRATIONS value '%s', defaulting to false", runMigrations)
		return false
	}

	return val
}

// getMigrationDirection determines which direction to migrate
func getMigrationDirection() string {
	direction := os.Getenv("MIGRATION_DIRECTION")
	if direction == "" {
		return "up"
	}
	return direction
}

// MigrateUp is a convenience function to run up migrations
func MigrateUp(db *sql.DB, migrationsPath string) error {
	return RunMigrations(db, migrationsPath, "up")
}

// MigrateDown is a convenience function to run down migrations
func MigrateDown(db *sql.DB, migrationsPath string, steps int) error {
	return RunMigrations(db, migrationsPath, "down", strconv.Itoa(steps))
}

// MigrateForce is a convenience function to force a version
func MigrateForce(db *sql.DB, migrationsPath string, version int) error {
	return RunMigrations(db, migrationsPath, "force", strconv.Itoa(version))
}
